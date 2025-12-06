# Function Call 故障排查指南

## 问题：查询时没有调用函数

如果输入"查询所有已发布的文章"时没有调用函数，请按以下步骤排查：

## 1. 检查模型配置

确保使用的是支持 function calling 的模型：

- ✅ **支持**：`qwen-max`, `qwen-plus`, `qwen3-max`, `qwen3-plus`
- ❌ **不支持**：`qwen-turbo`（基础版不支持 function calling）

在 `backend/manifest/config/config.yaml` 中检查：

```yaml
qwen:
  model: "qwen3-max"  # 确保使用支持function calling的模型
```

## 2. 查看后端日志

启动后端服务后，在日志中查找以下信息：

### 发送请求时的日志：
```
准备调用通义千问API，模型: qwen3-max, 函数数量: 15
可用函数列表: [...]
请求体摘要: {...}
```

### 收到响应时的日志：
```
收到finish_reason: tool_calls
收到tool_calls: [...]
检测到function call: query_articles, arguments: {...}
```

### 如果没有看到 function call 相关日志：
- 检查 `finish_reason` 是否为 `tool_calls`
- 检查 `tool_calls` 是否为空数组
- 查看完整的API响应内容

## 3. 检查 API 响应格式

在日志中查找完整的API响应，检查：

1. **finish_reason** 字段：
   - 如果是 `stop`：AI 直接回答了，没有调用函数
   - 如果是 `tool_calls`：应该会调用函数

2. **tool_calls** 字段：
   - 应该包含函数调用的信息
   - 格式应该是：
     ```json
     {
       "type": "function",
       "function": {
         "name": "query_articles",
         "arguments": "{\"status\":\"published\"}"
       }
     }
     ```

## 4. 可能的原因和解决方案

### 原因1：模型不支持 function calling
**解决方案**：更换为支持 function calling 的模型（qwen-max, qwen-plus, qwen3-max等）

### 原因2：AI 没有识别需要调用函数
**解决方案**：
- 尝试更明确的提问，例如："请使用函数查询所有已发布的文章"
- 检查 system message 是否正确添加（代码中已自动添加）

### 原因3：API 响应格式解析错误
**解决方案**：
- 查看日志中的完整响应内容
- 检查 `tool_calls` 字段的格式是否正确

### 原因4：流式响应中 tool_calls 在最后才返回
**解决方案**：代码已优化，会在流式响应结束后检查是否有 tool_calls

## 5. 测试步骤

1. **启动后端服务**，确保日志级别设置为 `debug` 或 `info`

2. **发送测试消息**："查询所有已发布的文章"

3. **查看日志输出**，查找：
   - `准备调用通义千问API` - 确认请求已发送
   - `收到finish_reason` - 查看响应类型
   - `收到tool_calls` - 查看是否有函数调用
   - `检测到function call` - 确认函数调用成功

4. **如果没有 function call**：
   - 查看完整的API响应内容
   - 检查 `finish_reason` 的值
   - 检查模型配置是否正确

## 6. 调试技巧

### 添加更详细的日志

如果问题仍然存在，可以在代码中添加更详细的日志：

```go
// 在 chat_stream.go 的 callQwenAPIStream 函数中
g.Log().Infof(ctx, "完整响应内容: %s", dataStr)
```

### 测试不同的提问方式

尝试以下提问方式：

1. "查询所有已发布的文章"（自然语言）
2. "请使用函数查询所有已发布的文章"（明确要求）
3. "帮我查询状态为published的文章"（更具体）

### 检查 API Key 和权限

确保：
- API Key 有效且有足够的额度
- API Key 对应的账户有权限使用 function calling 功能

## 7. 常见问题

### Q: finish_reason 是 "stop" 而不是 "tool_calls"
**A**: AI 可能认为不需要调用函数，或者模型不支持 function calling。尝试：
- 更换为支持 function calling 的模型
- 使用更明确的提问方式
- 检查 system message 是否正确

### Q: tool_calls 是空数组
**A**: API 返回了 tool_calls 字段，但内容为空。可能原因：
- 函数定义格式不正确
- 模型无法识别需要调用的函数
- API 版本问题

### Q: 解析 function arguments 失败
**A**: arguments 字段可能是字符串格式的 JSON，需要解析。代码中已处理，如果仍然失败，检查日志中的 arguments 内容。

## 8. 联系支持

如果以上步骤都无法解决问题，请：
1. 收集完整的日志输出
2. 记录 API 请求和响应的完整内容
3. 检查通义千问官方文档是否有更新

## 9. 验证 Function Calling 是否工作

成功调用函数时，日志应该显示：

```
检测到function call: query_articles, arguments: map[status:published]
执行函数调用: query_articles
查询文章列表成功，返回 X 条记录
```

如果看到这些日志，说明 function calling 正常工作。


