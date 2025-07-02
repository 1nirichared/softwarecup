<template>
  <div class="chat-session-container">
    <div class="chat-header">
      <h1>{{ session?.title || '对话详情' }}</h1>
      <div class="header-actions">
        <button @click="$router.push('/dashboard/chat')" class="back-btn">
          返回列表
        </button>
      </div>
    </div>
    
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="chat-content">
      <div class="messages" ref="messagesContainer">
        <div v-for="message in messages" :key="message.id" class="message">
          <div :class="['message-content', message.role]">
            <div class="message-text">{{ message.content }}</div>
            <div class="message-time">{{ formatTime(message.created_at) }}</div>
          </div>
        </div>
      </div>
      
      <div class="input-area">
        <textarea 
          v-model="newMessage" 
          placeholder="输入消息..." 
          @keydown.ctrl.enter="sendMessage"
          :disabled="sending"
          rows="3"
        ></textarea>
        <button @click="sendMessage" :disabled="sending || !newMessage.trim()">
          {{ sending ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { chatApi } from '@/api'

const route = useRoute()
const session = ref(null)
const messages = ref([])
const newMessage = ref('')
const loading = ref(true)
const error = ref('')
const sending = ref(false)
const messagesContainer = ref(null)

const loadSession = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await chatApi.getChatSession(route.params.id)
    if (response.data && response.data.code === 200) {
      session.value = response.data.data
      messages.value = response.data.data.messages || []
      scrollToBottom()
    } else {
      error.value = response.data?.message || '获取会话详情失败'
    }
  } catch (err) {
    console.error('获取会话详情失败:', err)
    error.value = err.response?.data?.message || '网络错误，请重试'
  } finally {
    loading.value = false
  }
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || sending.value) return
  
  const messageText = newMessage.value.trim()
  newMessage.value = ''
  
  try {
    sending.value = true
    
    // 立即显示用户消息
    const userMessage = {
      id: Date.now(),
      role: 'user',
      content: messageText,
      created_at: new Date().toISOString()
    }
    messages.value.push(userMessage)
    scrollToBottom()
    
    // 添加AI消息占位
    const aiMessage = {
      id: Date.now() + 1,
      role: 'assistant',
      content: '',
      created_at: new Date().toISOString()
    }
    messages.value.push(aiMessage)
    scrollToBottom()
    
    // 流式请求AI回复
    let fullReply = ''
    const source = chatApi.streamAIChat(
      messageText,
      (chunk) => {
        fullReply += chunk
        aiMessage.content = fullReply
        scrollToBottom()
      },
      () => {
        sending.value = false
      },
      (err) => {
        aiMessage.content += '\n[AI回复失败]'
        sending.value = false
      }
    )
  } catch (err) {
    console.error('发送消息失败:', err)
    alert(err.response?.data?.message || '发送消息失败，请重试')
    sending.value = false
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const formatTime = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleTimeString('zh-CN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

// 监听消息变化，自动滚动到底部
watch(messages, () => {
  scrollToBottom()
}, { deep: true })

onMounted(() => {
  loadSession()
})
</script>

<style scoped>
.chat-session-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 20px;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.back-btn {
  background: #666;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.chat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  overflow: hidden;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  color: #f56c6c;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.message {
  margin-bottom: 20px;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  position: relative;
}

.message-content.user {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  margin-left: auto;
}

.message-content.assistant {
  background: #f0f0f0;
  color: #333;
}

.message-text {
  margin-bottom: 5px;
  line-height: 1.5;
  white-space: pre-wrap;
}

.message-time {
  font-size: 0.8rem;
  opacity: 0.7;
}

.input-area {
  display: flex;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #eee;
  background: #f8f9fa;
}

.input-area textarea {
  flex: 1;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  resize: none;
  font-family: inherit;
}

.input-area textarea:focus {
  outline: none;
  border-color: #667eea;
}

.input-area button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: opacity 0.3s;
}

.input-area button:hover:not(:disabled) {
  opacity: 0.9;
}

.input-area button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style> 