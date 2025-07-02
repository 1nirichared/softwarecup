<template>
  <div class="chat-container">
    <h1>AI 助手</h1>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="chat-sessions">
      <div v-for="session in sessions" :key="session.id" class="session-card">
        <h3>{{ session.title || '未命名对话' }}</h3>
        <p>{{ session.lastMessage || '暂无消息' }}</p>
        <div class="session-meta">
          <span>{{ formatDate(session.updated_at) }}</span>
          <span>{{ session.type }}</span>
        </div>
        <div class="session-actions">
          <router-link :to="`/dashboard/chat/${session.id}`" class="open-btn">
            打开对话
          </router-link>
          <button @click="deleteSession(session.id)" class="delete-btn">
            删除
          </button>
        </div>
      </div>
    </div>
    <button class="new-chat-btn" @click="createNewChat" :disabled="creating">
      {{ creating ? '创建中...' : '新建对话' }}
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { chatApi } from '@/api'

const router = useRouter()
const sessions = ref([])
const loading = ref(true)
const error = ref('')
const creating = ref(false)

const loadSessions = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await chatApi.getChatSessions()
    if (response.data && response.data.code === 200) {
      sessions.value = response.data.data || []
    } else {
      error.value = response.data?.message || '获取会话列表失败'
    }
  } catch (err) {
    console.error('获取会话列表失败:', err)
    error.value = err.response?.data?.message || '网络错误，请重试'
  } finally {
    loading.value = false
  }
}

const createNewChat = async () => {
  try {
    creating.value = true
    const response = await chatApi.createChatSession({
      title: '新对话',
      type: 'general'
    })
    if (response.data && response.data.code === 200) {
      const sessionId = response.data.data.id
      router.push(`/dashboard/chat/${sessionId}`)
    } else {
      alert(response.data?.message || '创建会话失败')
    }
  } catch (err) {
    console.error('创建会话失败:', err)
    alert(err.response?.data?.message || '创建会话失败，请重试')
  } finally {
    creating.value = false
  }
}

const deleteSession = async (sessionId) => {
  if (!confirm('确定要删除这个对话吗？')) return
  
  try {
    const response = await chatApi.deleteChatSession(sessionId)
    if (response.data && response.data.code === 200) {
      await loadSessions() // 重新加载列表
    } else {
      alert(response.data?.message || '删除会话失败')
    }
  } catch (err) {
    console.error('删除会话失败:', err)
    alert(err.response?.data?.message || '删除会话失败，请重试')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  loadSessions()
})
</script>

<style scoped>
.chat-container {
  padding: 20px;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  color: #f56c6c;
}

.chat-sessions {
  display: grid;
  gap: 15px;
  margin: 20px 0;
}

.session-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.session-meta {
  display: flex;
  gap: 15px;
  margin: 10px 0;
  font-size: 0.9rem;
  color: #666;
}

.session-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

.open-btn {
  background: #667eea;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  text-decoration: none;
  font-size: 0.9rem;
}

.delete-btn {
  background: #f56c6c;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.new-chat-btn {
  background: #28a745;
  color: white;
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  transition: opacity 0.3s;
}

.new-chat-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.new-chat-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style> 