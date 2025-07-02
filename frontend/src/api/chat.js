import api from './index'

export const chatApi = {
  // 获取聊天会话列表
  getChatSessions: () => {
    return api.get('/chat/sessions')
  },

  // 创建聊天会话
  createChatSession: (data) => {
    return api.post('/chat/sessions', data)
  },

  // 获取聊天会话详情
  getChatSession: (id) => {
    return api.get(`/chat/sessions/${id}`)
  },

  // 发送消息
  sendMessage: (sessionId, data) => {
    return api.post(`/chat/sessions/${sessionId}/messages`, data)
  },

  // 删除聊天会话
  deleteChatSession: (id) => {
    return api.delete(`/chat/sessions/${id}`)
  },

  // 获取学习建议
  getLearningAdvice: () => {
    return api.get('/chat/advice')
  },

  // 流式AI回复（SSE）
  streamAIChat: (message, onMessage, onEnd, onError) => {
    const token = localStorage.getItem('token')
    const source = new EventSource(`/api/v1/chat/stream?message=${encodeURIComponent(message)}&token=${token}`)
    
    source.addEventListener('message', (e) => {
      if (onMessage) onMessage(e.data)
    })
    
    source.addEventListener('error', (e) => {
      if (onError) onError(e)
      source.close()
    })
    
    source.addEventListener('end', (e) => {
      if (onEnd) onEnd(e.data)
      source.close()
    })
    
    return source
  }
} 