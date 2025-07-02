<template>
  <div class="home-container">
    <div class="welcome-section">
      <h1>欢迎回来，{{ user?.real_name || user?.username }}！</h1>
      <p class="subtitle">智能教学平台 - 让学习更高效</p>
    </div>
    
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📚</div>
        <div class="stat-content">
          <h3>{{ stats.courseCount || 0 }}</h3>
          <p>我的课程</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">📝</div>
        <div class="stat-content">
          <h3>{{ stats.exerciseCount || 0 }}</h3>
          <p>练习数量</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">🎯</div>
        <div class="stat-content">
          <h3>{{ stats.avgScore || 0 }}%</h3>
          <p>平均成绩</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">⏱️</div>
        <div class="stat-content">
          <h3>{{ stats.studyTime || 0 }}h</h3>
          <p>学习时长</p>
        </div>
      </div>
    </div>
    
    <div class="quick-actions">
      <h2>快速操作</h2>
      <div class="action-grid">
        <router-link to="/dashboard/courses" class="action-card">
          <div class="action-icon">📖</div>
          <h3>浏览课程</h3>
          <p>查看所有可用课程</p>
        </router-link>
        
        <router-link to="/dashboard/exercises" class="action-card">
          <div class="action-icon">✏️</div>
          <h3>开始练习</h3>
          <p>进行在线练习</p>
        </router-link>
        
        <router-link to="/dashboard/chat" class="action-card">
          <div class="action-icon">💬</div>
          <h3>AI 助手</h3>
          <p>与AI助手对话</p>
        </router-link>
        
        <router-link to="/dashboard/profile" class="action-card">
          <div class="action-icon">👤</div>
          <h3>个人资料</h3>
          <p>管理个人信息</p>
        </router-link>
      </div>
    </div>
    
    <div class="recent-activities" v-if="recentActivities.length > 0">
      <h2>最近活动</h2>
      <div class="activity-list">
        <div v-for="activity in recentActivities" :key="activity.id" class="activity-item">
          <div class="activity-icon">{{ activity.icon }}</div>
          <div class="activity-content">
            <h4>{{ activity.title }}</h4>
            <p>{{ activity.description }}</p>
            <span class="activity-time">{{ activity.time }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

const stats = ref({
  courseCount: 0,
  exerciseCount: 0,
  avgScore: 0,
  studyTime: 0
})

const recentActivities = ref([
  {
    id: 1,
    icon: '📚',
    title: '完成了课程学习',
    description: '完成了《JavaScript基础》第三章的学习',
    time: '2小时前'
  },
  {
    id: 2,
    icon: '✏️',
    title: '提交了练习答案',
    description: '完成了《数据结构》练习，得分85分',
    time: '1天前'
  },
  {
    id: 3,
    icon: '💬',
    title: '与AI助手对话',
    description: '询问了关于算法复杂度的问题',
    time: '2天前'
  }
])

onMounted(async () => {
  // 这里可以加载用户统计数据
  // await loadUserStats()
})
</script>

<style scoped>
.home-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-section {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 15px;
}

.welcome-section h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-card {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 20px;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-icon {
  font-size: 2.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-content h3 {
  font-size: 2rem;
  color: #333;
  margin: 0 0 5px 0;
}

.stat-content p {
  color: #666;
  margin: 0;
}

.quick-actions {
  margin-bottom: 40px;
}

.quick-actions h2 {
  margin-bottom: 20px;
  color: #333;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.action-card {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-decoration: none;
  color: inherit;
  transition: all 0.3s ease;
  text-align: center;
}

.action-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.15);
}

.action-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.action-card h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.action-card p {
  color: #666;
  margin: 0;
}

.recent-activities {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.recent-activities h2 {
  margin-bottom: 20px;
  color: #333;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-radius: 8px;
  background: #f8f9fa;
}

.activity-icon {
  font-size: 1.5rem;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 50%;
}

.activity-content h4 {
  margin: 0 0 5px 0;
  color: #333;
}

.activity-content p {
  margin: 0 0 5px 0;
  color: #666;
}

.activity-time {
  font-size: 0.9rem;
  color: #999;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .action-grid {
    grid-template-columns: 1fr;
  }
  
  .welcome-section h1 {
    font-size: 2rem;
  }
}
</style> 