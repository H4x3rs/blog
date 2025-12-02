package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"blog/internal/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sChat struct{}

var Chat = sChat{}

// 消息结构
type ChatMessage struct {
	Role    string `json:"role"`           // "user", "assistant", "function"
	Content string `json:"content"`        // 消息内容
	Name    string `json:"name,omitempty"` // function名称（当role为function时）
}

// 会话历史
type ChatHistory struct {
	Messages  []ChatMessage `json:"messages"`
	UpdatedAt int64         `json:"updatedAt"`
}

// 通义千问API响应
type QwenResponse struct {
	Output struct {
		Choices []struct {
			Message struct {
				Role      string                   `json:"role"`
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

// 发送消息
func (s *sChat) SendMessage(ctx context.Context, userId int, message string, sessionId string) (reply string, newSessionId string, err error) {
	// 如果没有sessionId，创建一个新的
	if sessionId == "" {
		sessionId = fmt.Sprintf("chat_%d_%d", userId, time.Now().Unix())
	}

	// 获取或创建会话历史
	history, err := s.getOrCreateHistory(ctx, userId, sessionId)
	if err != nil {
		return "", "", err
	}

	// 添加用户消息
	history.Messages = append(history.Messages, ChatMessage{
		Role:    "user",
		Content: message,
	})

	// 循环处理function calling（最多5轮，避免无限循环）
	maxIterations := 5
	for i := 0; i < maxIterations; i++ {
		// 调用通义千问API
		reply, functionCall, err := s.callQwenAPI(ctx, history.Messages, userId)
		if err != nil {
			return "", "", err
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
			Content: reply,
		})
		break
	}

	history.UpdatedAt = time.Now().Unix()

	// 保存会话历史（限制历史记录数量，避免过长）
	maxHistory := 20 // 最多保留20轮对话
	if len(history.Messages) > maxHistory*2 {
		// 保留最近的对话，但保留第一条系统消息（如果有）
		keepCount := maxHistory * 2
		history.Messages = history.Messages[len(history.Messages)-keepCount:]
	}

	// 保存到Redis
	err = s.saveHistory(ctx, userId, sessionId, history)
	if err != nil {
		g.Log().Warning(ctx, "保存会话历史失败:", err)
		// 不返回错误，因为对话已经成功
	}

	return reply, sessionId, nil
}

// 调用通义千问API
// 返回：完整回复内容、function call（如果有）、错误
func (s *sChat) callQwenAPI(ctx context.Context, messages []ChatMessage, userId int) (string, *FunctionCall, error) {
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
	g.Log().Infof(ctx, "准备调用通义千问API（非流式），模型: %s, 函数数量: %d", model, len(tools))
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

	// 构建请求体（启用function calling）
	requestBody := map[string]interface{}{
		"model": model,
		"input": map[string]interface{}{
			"messages": finalMessages,
		},
		"parameters": map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  2000,
		},
		"tools": tools, // 添加function calling支持
	}

	// 记录请求体（不记录敏感信息）
	requestBodyForLog := map[string]interface{}{
		"model":          model,
		"messages_count": len(finalMessages),
		"tools_count":    len(tools),
	}
	g.Log().Debugf(ctx, "请求体摘要: %+v", requestBodyForLog)

	// 发送HTTP请求
	client := g.Client()
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client.SetHeader("Content-Type", "application/json")
	client.SetTimeout(120 * time.Second) // 设置120秒超时，AI API可能需要更长时间

	response, err := client.Post(ctx, apiEndpoint, requestBody)
	if err != nil {
		return "", nil, gerror.Wrap(err, "调用通义千问API失败")
	}

	defer response.Close()

	// 读取响应内容
	responseBody := response.ReadAll()

	// 检查HTTP状态码
	if response.StatusCode != 200 {
		return "", nil, gerror.Newf("通义千问API返回错误状态码: %d, 响应: %s", response.StatusCode, string(responseBody))
	}

	// 解析响应
	var qwenResp QwenResponse
	if err := json.Unmarshal(responseBody, &qwenResp); err != nil {
		// 如果解析失败，尝试查看原始响应
		g.Log().Errorf(ctx, "通义千问API响应解析失败，原始响应: %s", string(responseBody))
		return "", nil, gerror.Wrap(err, "解析通义千问API响应失败")
	}

	// 检查错误
	if qwenResp.Code != "" && qwenResp.Code != "Success" {
		return "", nil, gerror.Newf("通义千问API返回错误: %s - %s", qwenResp.Code, qwenResp.Message)
	}

	// 提取回复内容
	if len(qwenResp.Output.Choices) == 0 {
		g.Log().Errorf(ctx, "通义千问API返回空回复，完整响应: %+v", qwenResp)
		return "", nil, gerror.New("通义千问API返回空回复")
	}

	choice := qwenResp.Output.Choices[0]
	reply := choice.Message.Content

	// 调试日志：记录响应信息
	g.Log().Infof(ctx, "API响应 - finish_reason: %s, tool_calls数量: %d, content长度: %d",
		choice.FinishReason, len(choice.Message.ToolCalls), len(reply))
	if len(choice.Message.ToolCalls) > 0 {
		g.Log().Infof(ctx, "收到tool_calls: %+v", choice.Message.ToolCalls)
	}

	// 检查是否有function call
	if choice.FinishReason == "tool_calls" && len(choice.Message.ToolCalls) > 0 {
		g.Log().Infof(ctx, "开始处理function call")
		// 处理第一个tool call
		for _, toolCall := range choice.Message.ToolCalls {
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
					return reply, &FunctionCall{
						Name:      functionName,
						Arguments: arguments,
					}, nil
				}
			}
		}
	}

	// Fallback: 尝试从回复内容中解析函数调用信息
	// AI可能在content中返回了函数调用信息，格式如：{"function": "query_categories", "arguments": {}}
	if functionCall := s.parseFunctionCallFromContent(ctx, reply); functionCall != nil {
		g.Log().Infof(ctx, "从回复内容中解析到function call: %s, arguments: %+v", functionCall.Name, functionCall.Arguments)
		return reply, functionCall, nil
	}

	if reply == "" {
		return "", nil, gerror.New("通义千问API返回空内容")
	}

	return reply, nil, nil
}

// 获取或创建会话历史
func (s *sChat) getOrCreateHistory(ctx context.Context, userId int, sessionId string) (*ChatHistory, error) {
	key := s.getHistoryKey(userId, sessionId)

	// 从Redis获取
	redis := g.Redis()
	value, err := redis.Do(ctx, "GET", key)
	if err != nil {
		// Redis错误，记录日志但不影响功能，返回新的历史
		g.Log().Warning(ctx, "从Redis获取会话历史失败:", err)
		return &ChatHistory{
			Messages:  []ChatMessage{},
			UpdatedAt: time.Now().Unix(),
		}, nil
	}
	if value == nil {
		// 不存在，创建新的
		return &ChatHistory{
			Messages:  []ChatMessage{},
			UpdatedAt: time.Now().Unix(),
		}, nil
	}

	// 解析历史记录
	var history ChatHistory
	valueStr := value.String()
	if valueStr == "" {
		return &ChatHistory{
			Messages:  []ChatMessage{},
			UpdatedAt: time.Now().Unix(),
		}, nil
	}

	if err := json.Unmarshal([]byte(valueStr), &history); err != nil {
		// 解析失败，创建新的
		return &ChatHistory{
			Messages:  []ChatMessage{},
			UpdatedAt: time.Now().Unix(),
		}, nil
	}

	return &history, nil
}

// 保存会话历史
func (s *sChat) saveHistory(ctx context.Context, userId int, sessionId string, history *ChatHistory) error {
	key := s.getHistoryKey(userId, sessionId)
	redis := g.Redis()

	// 序列化历史记录
	data, err := json.Marshal(history)
	if err != nil {
		return err
	}

	// 保存到Redis，设置7天过期
	_, err = redis.Do(ctx, "SETEX", key, 7*24*3600, string(data))
	return err
}

// 获取Redis key
func (s *sChat) getHistoryKey(userId int, sessionId string) string {
	return fmt.Sprintf("chat:history:%d:%s", userId, sessionId)
}

// 从回复内容中解析函数调用信息（Fallback机制）
// 支持格式：
// 1. {"function": "query_categories", "arguments": {}}
// 2. {"function": "query_categories", "arguments": {"page": 1}}
// 3. 包含JSON格式的文本内容
func (s *sChat) parseFunctionCallFromContent(ctx context.Context, content string) *FunctionCall {
	// 尝试直接解析JSON（如果整个内容就是JSON）
	var directCall struct {
		Function  string                 `json:"function"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal([]byte(content), &directCall); err == nil && directCall.Function != "" {
		g.Log().Infof(ctx, "从内容中直接解析到函数调用: %s", directCall.Function)
		return &FunctionCall{
			Name:      directCall.Function,
			Arguments: directCall.Arguments,
		}
	}

	// 尝试从文本中提取JSON格式的函数调用
	// 查找类似 {"function": "...", "arguments": {...}} 的模式
	re := regexp.MustCompile(`\{[^{}]*"function"\s*:\s*"([^"]+)"[^{}]*"arguments"\s*:\s*(\{[^}]*\})[^{}]*\}`)
	matches := re.FindStringSubmatch(content)
	if len(matches) >= 3 {
		functionName := matches[1]
		argumentsStr := matches[2]

		var arguments map[string]interface{}
		if err := json.Unmarshal([]byte(argumentsStr), &arguments); err == nil {
			g.Log().Infof(ctx, "从内容中正则解析到函数调用: %s", functionName)
			return &FunctionCall{
				Name:      functionName,
				Arguments: arguments,
			}
		}
	}

	// 尝试查找包含函数名的JSON块（逐行检查）
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}

		// 检查是否包含function和arguments字段
		if strings.Contains(line, `"function"`) && strings.Contains(line, `"arguments"`) {
			var call struct {
				Function  string                 `json:"function"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			if err := json.Unmarshal([]byte(line), &call); err == nil && call.Function != "" {
				g.Log().Infof(ctx, "从行中解析到函数调用: %s", call.Function)
				return &FunctionCall{
					Name:      call.Function,
					Arguments: call.Arguments,
				}
			}
		}
	}

	// 尝试查找代码块中的JSON（markdown代码块）
	codeBlockRe := regexp.MustCompile("```(?:json)?\\s*\\n?\\s*(\\{[^`]*\\})\\s*```")
	codeMatches := codeBlockRe.FindStringSubmatch(content)
	if len(codeMatches) >= 2 {
		var call struct {
			Function  string                 `json:"function"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal([]byte(codeMatches[1]), &call); err == nil && call.Function != "" {
			g.Log().Infof(ctx, "从代码块中解析到函数调用: %s", call.Function)
			return &FunctionCall{
				Name:      call.Function,
				Arguments: call.Arguments,
			}
		}
	}

	return nil
}

// 清除会话
func (s *sChat) ClearSession(ctx context.Context, userId int, sessionId string) error {
	key := s.getHistoryKey(userId, sessionId)
	redis := g.Redis()
	_, err := redis.Do(ctx, "DEL", key)
	return err
}
