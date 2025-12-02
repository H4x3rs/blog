import request from '../utils/request'

// 发送消息（AI对话需要更长的超时时间）
export function sendChatMessage(data) {
  return request({
    url: '/chat/message',
    method: 'post',
    data,
    timeout: 120000 // 120秒超时，AI对话API可能需要更长时间
  })
}

// 流式发送消息
export async function sendChatMessageStream(data, onChunk, onDone) {
  const token = localStorage.getItem('token')
  const baseURL = '/api'
  
  const response = await fetch(`${baseURL}/chat/message`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'text/event-stream',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {})
    },
    body: JSON.stringify(data)
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder('utf-8')
  let buffer = ''
  let pendingEvent = null

  try {
    while (true) {
      const { value, done } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || '' // 保留最后一个不完整的行

      for (const line of lines) {
        const trimmedLine = line.trim()
        if (!trimmedLine) {
          // 空行表示一个完整的事件结束
          pendingEvent = null
          continue
        }

        // 调试日志
        console.log('收到行:', trimmedLine.substring(0, 100))

        if (trimmedLine.startsWith('event: ')) {
          // 保存事件类型
          pendingEvent = trimmedLine.substring(7)
          console.log('设置事件类型:', pendingEvent)
          continue
        }

        if (trimmedLine.startsWith('data: ')) {
          const dataStr = trimmedLine.substring(6)
          
          if (dataStr === '[DONE]') {
            continue
          }

          // 如果是done事件，解析sessionId
          if (pendingEvent === 'done') {
            try {
              const jsonData = JSON.parse(dataStr)
              if (jsonData.sessionId) {
                console.log('会话完成，sessionId:', jsonData.sessionId)
                onDone(jsonData.sessionId)
              }
            } catch (e) {
              console.error('解析完成事件失败:', e, 'data:', dataStr)
            }
            pendingEvent = null
            continue
          }

          // 处理内容数据（默认或没有event的情况）
          try {
            // 尝试解析JSON（可能是转义的字符串）
            const jsonData = JSON.parse(dataStr)
            if (typeof jsonData === 'string') {
              // 如果是字符串，直接作为内容块
              console.log('收到内容块:', jsonData)
              onChunk(jsonData)
            } else {
              console.warn('未知的JSON格式:', jsonData)
            }
          } catch (e) {
            // 如果不是JSON，尝试作为JSON转义的字符串解析
            if (dataStr.startsWith('"') && dataStr.endsWith('"')) {
              try {
                const content = JSON.parse(dataStr)
                console.log('解析JSON转义字符串:', content)
                onChunk(content)
              } catch (e2) {
                // 解析失败，直接使用原始内容
                console.log('使用原始内容:', dataStr)
                onChunk(dataStr)
              }
            } else {
              // 直接作为文本内容
              console.log('直接使用文本:', dataStr)
              onChunk(dataStr)
            }
          }
          pendingEvent = null
        }
      }
    }
  } finally {
    reader.releaseLock()
  }
}

// 清除会话
export function clearChatSession(data) {
  return request({
    url: '/chat/clearSession',
    method: 'post',
    data
  })
}
