package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"blog/internal/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 流式调用通义千问API
func (s *sChat) SendMessageStream(ctx context.Context, writer *ghttp.Response, userId int, message string, sessionId string) (newSessionId string, err error) {
	// 如果没有sessionId，创建一个新的
	if sessionId == "" {
		sessionId = fmt.Sprintf("chat_%d_%d", userId, time.Now().Unix())
	}

	// 获取或创建会话历史
	history, err := s.getOrCreateHistory(ctx, userId, sessionId)
	if err != nil {
		return "", err
	}

	// 添加用户消息
	history.Messages = append(history.Messages, ChatMessage{
		Role:    "user",
		Content: message,
	})

	// 循环处理function calling（最多5轮，避免无限循环）
	maxIterations := 5
	for i := 0; i < maxIterations; i++ {
		// 调用流式通义千问API
		fullReply, functionCall, err := s.callQwenAPIStream(ctx, writer, history.Messages, userId)
		if err != nil {
			return "", err
		}

		// 如果有function call，执行函数并继续
		if functionCall != nil {
			// 执行函数
			result, err := s.ExecuteFunction(ctx, userId, functionCall)
			if err != nil {
				// 函数执行失败，发送错误信息
				history.Messages = append(history.Messages, ChatMessage{
					Role:    "function",
					Content: fmt.Sprintf(`{"error": "%s"}`, err.Error()),
					Name:    functionCall.Name,
				})
				continue
			}

			// 将函数结果添加到消息历史
			resultJSON, _ := json.Marshal(result.Content)
			history.Messages = append(history.Messages, ChatMessage{
				Role:    "function",
				Content: string(resultJSON),
				Name:    functionCall.Name,
			})

			// 继续下一轮对话
			continue
		}

		// 没有function call，添加AI回复并结束
		history.Messages = append(history.Messages, ChatMessage{
			Role:    "assistant",
			Content: fullReply,
		})
		break
	}

	history.UpdatedAt = time.Now().Unix()

	// 保存会话历史（限制历史记录数量，避免过长）
	maxHistory := 20 // 最多保留20轮对话
	if len(history.Messages) > maxHistory*2 {
		keepCount := maxHistory * 2
		history.Messages = history.Messages[len(history.Messages)-keepCount:]
	}

	// 保存到Redis
	err = s.saveHistory(ctx, userId, sessionId, history)
	if err != nil {
		g.Log().Warning(ctx, "保存会话历史失败:", err)
		// 不返回错误，因为对话已经成功
	}

	return sessionId, nil
}

// 流式调用通义千问API
// 返回：完整回复内容、function call（如果有）、错误
func (s *sChat) callQwenAPIStream(ctx context.Context, writer *ghttp.Response, messages []ChatMessage, userId int) (string, *FunctionCall, error) {
	// 获取配置
	apiKey := utils.GetConfigString(ctx, "QWEN_API_KEY", "qwen.apiKey", "")
	apiEndpoint := utils.GetConfigString(ctx, "QWEN_API_ENDPOINT", "qwen.apiEndpoint", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation")
	model := utils.GetConfigString(ctx, "QWEN_MODEL", "qwen.model", "qwen-turbo")

	if apiKey == "" {
		return "", nil, gerror.New("通义千问API密钥未配置，请在环境变量QWEN_API_KEY或配置文件中设置qwen.apiKey")
	}

	// 获取可用函数列表
	tools := s.GetAvailableFunctions()

	// 添加调试日志
	g.Log().Infof(ctx, "准备调用通义千问API，模型: %s, 函数数量: %d", model, len(tools))
	if len(tools) > 0 {
		g.Log().Debugf(ctx, "可用函数列表: %+v", tools)
	}

	// 检查是否有system message，如果没有则添加
	hasSystemMessage := false
	for _, msg := range messages {
		if msg.Role == "system" {
			hasSystemMessage = true
			break
		}
	}

	// 构建消息列表，确保有system message引导AI使用function calling
	finalMessages := make([]ChatMessage, 0, len(messages)+1)
	if !hasSystemMessage && len(tools) > 0 {
		// 添加system message，引导AI使用function calling
		finalMessages = append(finalMessages, ChatMessage{
			Role:    "system",
			Content: `你是一个智能助手，可以帮助用户查询数据库中的信息。当用户询问文章、用户、分类、标签、菜单、权限、角色、专题或设置等信息时，你应该使用提供的函数来查询数据库，而不是直接回答。只有在查询到实际数据后，才能基于数据给出回答。`,
		})
	}
	finalMessages = append(finalMessages, messages...)

	// 构建请求体（启用流式输出和function calling）
	requestBody := map[string]interface{}{
		"model": model,
		"input": map[string]interface{}{
			"messages": finalMessages,
		},
		"parameters": map[string]interface{}{
			"temperature":        0.7,
			"max_tokens":         2000,
			"incremental_output": true, // 启用流式输出
		},
		"tools": tools, // 添加function calling支持
	}

	// 记录请求体（不记录敏感信息）
	requestBodyForLog := map[string]interface{}{
		"model":          model,
		"messages_count": len(messages),
		"tools_count":    len(tools),
	}
	g.Log().Debugf(ctx, "请求体摘要: %+v", requestBodyForLog)

	// 发送HTTP请求
	client := g.Client()
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "text/event-stream")
	client.SetTimeout(120 * time.Second)

	// 使用流式请求
	response, err := client.Post(ctx, apiEndpoint, requestBody)
	if err != nil {
		return "", nil, gerror.Wrap(err, "调用通义千问API失败")
	}
	defer response.Close()

	// 检查HTTP状态码
	if response.StatusCode != 200 {
		responseBody := response.ReadAll()
		return "", nil, gerror.Newf("通义千问API返回错误状态码: %d, 响应: %s", response.StatusCode, string(responseBody))
	}

	// 读取流式响应
	var fullReply strings.Builder
	reader := bufio.NewReader(response.Body)
	lineCount := 0
	var lastToolCalls []map[string]interface{} // 累积tool_calls

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", nil, gerror.Wrap(err, "读取流式响应失败")
		}

		lineCount++
		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			// 空行表示事件结束
			continue
		}

		// 解析SSE格式
		if strings.HasPrefix(line, "id:") {
			// 忽略id行
			continue
		} else if strings.HasPrefix(line, "event:") {
			// 忽略event行
			continue
		} else if strings.HasPrefix(line, ":") {
			// 忽略注释行（如 :HTTP_STATUS/200）
			continue
		} else if strings.HasPrefix(line, "data:") {
			// 解析data行（注意：data:后面可能没有空格）
			dataStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

			if dataStr == "[DONE]" {
				break
			}

			// 解析JSON数据
			var qwenResp struct {
				Output struct {
					Choices []struct {
						Message struct {
							Content   string                   `json:"content"`
							ToolCalls []map[string]interface{} `json:"tool_calls"`
						} `json:"message"`
						FinishReason string `json:"finish_reason"`
					} `json:"choices"`
				} `json:"output"`
				RequestId string `json:"request_id"`
				Code      string `json:"code"`
				Message   string `json:"message"`
			}

			if err := json.Unmarshal([]byte(dataStr), &qwenResp); err != nil {
				// 如果不是JSON格式，可能是纯文本数据，跳过
				if !strings.HasPrefix(dataStr, "{") {
					continue
				}
				g.Log().Warningf(ctx, "解析流式响应数据失败: %v, data: %s", err, dataStr)
				continue
			}

			// 详细日志：记录每次收到的响应
			if len(qwenResp.Output.Choices) > 0 {
				choice := qwenResp.Output.Choices[0]
				g.Log().Debugf(ctx, "收到响应块 - finish_reason: %s, content长度: %d, tool_calls数量: %d",
					choice.FinishReason, len(choice.Message.Content), len(choice.Message.ToolCalls))
			}

			// 检查错误
			if qwenResp.Code != "" && qwenResp.Code != "Success" {
				return "", nil, gerror.Newf("通义千问API返回错误: %s - %s", qwenResp.Code, qwenResp.Message)
			}

			// 提取增量内容和function call
			if len(qwenResp.Output.Choices) > 0 {
				choice := qwenResp.Output.Choices[0]
				content := choice.Message.Content

				// 调试日志：记录finish_reason和tool_calls
				if choice.FinishReason != "" {
					g.Log().Debugf(ctx, "收到finish_reason: %s", choice.FinishReason)
				}
				if len(choice.Message.ToolCalls) > 0 {
					g.Log().Infof(ctx, "收到tool_calls: %+v", choice.Message.ToolCalls)
				}

				// 累积tool_calls（流式响应中tool_calls可能是分片返回的）
				if len(choice.Message.ToolCalls) > 0 {
					lastToolCalls = choice.Message.ToolCalls
					g.Log().Infof(ctx, "累积tool_calls，当前数量: %d", len(lastToolCalls))
				}

				// 发送文本内容
				if content != "" {
					fullReply.WriteString(content)

					// 发送增量内容到前端（JSON转义）
					writer.Writefln("data: %s\n\n", jsonEscape(content))
					writer.Flush()
				}

				// 检查是否完成
				if choice.FinishReason == "stop" || choice.FinishReason == "tool_calls" {
					g.Log().Infof(ctx, "流式响应完成，原因: %s, lastToolCalls数量: %d", choice.FinishReason, len(lastToolCalls))

					// 如果是tool_calls，处理function call
					if choice.FinishReason == "tool_calls" && len(lastToolCalls) > 0 {
						g.Log().Infof(ctx, "开始处理function call，tool_calls: %+v", lastToolCalls)
						// 处理第一个tool call
						for _, toolCall := range lastToolCalls {
							if functionType, ok := toolCall["type"].(string); ok && functionType == "function" {
								if function, ok := toolCall["function"].(map[string]interface{}); ok {
									functionName, _ := function["name"].(string)
									argumentsStr, _ := function["arguments"].(string)

									var arguments map[string]interface{}
									if err := json.Unmarshal([]byte(argumentsStr), &arguments); err != nil {
										g.Log().Warningf(ctx, "解析function arguments失败: %v", err)
										continue
									}

									g.Log().Infof(ctx, "检测到function call: %s, arguments: %+v", functionName, arguments)
									return fullReply.String(), &FunctionCall{
										Name:      functionName,
										Arguments: arguments,
									}, nil
								}
							}
						}
					}
					break
				}
			}
		}
	}

	g.Log().Infof(ctx, "流式响应完成，总共处理%d行，完整回复长度: %d, 最终lastToolCalls数量: %d",
		lineCount, fullReply.Len(), len(lastToolCalls))

	// 如果流式响应结束但还有tool_calls，尝试处理
	if len(lastToolCalls) > 0 {
		g.Log().Infof(ctx, "流式响应结束后发现tool_calls，尝试处理: %+v", lastToolCalls)
		for _, toolCall := range lastToolCalls {
			if functionType, ok := toolCall["type"].(string); ok && functionType == "function" {
				if function, ok := toolCall["function"].(map[string]interface{}); ok {
					functionName, _ := function["name"].(string)
					argumentsStr, _ := function["arguments"].(string)

					var arguments map[string]interface{}
					if err := json.Unmarshal([]byte(argumentsStr), &arguments); err != nil {
						g.Log().Warningf(ctx, "解析function arguments失败: %v", err)
						continue
					}

					g.Log().Infof(ctx, "检测到function call: %s, arguments: %+v", functionName, arguments)
					return fullReply.String(), &FunctionCall{
						Name:      functionName,
						Arguments: arguments,
					}, nil
				}
			}
		}
	}

	// Fallback: 尝试从回复内容中解析函数调用信息
	// AI可能在content中返回了函数调用信息，格式如：{"function": "query_categories", "arguments": {}}
	fullReplyStr := fullReply.String()
	if functionCall := s.parseFunctionCallFromContent(ctx, fullReplyStr); functionCall != nil {
		g.Log().Infof(ctx, "从回复内容中解析到function call: %s, arguments: %+v", functionCall.Name, functionCall.Arguments)
		return fullReplyStr, functionCall, nil
	}

	return fullReply.String(), nil, nil
}

// JSON转义
func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
