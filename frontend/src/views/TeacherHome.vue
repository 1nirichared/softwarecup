<template>
  <div class="teacher-home">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1>欢迎回来，{{ user?.real_name || user?.username }}老师！</h1>
        <p>今天是 {{ currentDate }}，祝您教学愉快！</p>
      </div>
      <div class="welcome-actions">
        <el-button type="primary" @click="$router.push('/teacher/courses/create')">
          <el-icon><Plus /></el-icon>
          创建新课程
        </el-button>
        <el-button type="success" @click="$router.push('/teacher/exercises/create')">
          <el-icon><Edit /></el-icon>
          创建练习
        </el-button>
        <el-button type="warning" @click="$router.push('/teacher/lesson-plan')">
          <el-icon><Calendar /></el-icon>
          备课助手
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon courses">
            <el-icon><Reading /></el-icon>
          </div>
          <div class="stat-info">
            <h3>{{ stats.courseCount }}</h3>
            <p>我的课程</p>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon students">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <h3>{{ stats.studentCount }}</h3>
            <p>学生总数</p>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon exercises">
            <el-icon><Edit /></el-icon>
          </div>
          <div class="stat-info">
            <h3>{{ stats.exerciseCount }}</h3>
            <p>练习数量</p>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon materials">
            <el-icon><Folder /></el-icon>
          </div>
          <div class="stat-info">
            <h3>{{ stats.materialCount }}</h3>
            <p>教学资料</p>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <div class="content-left">
        <!-- 最近课程 -->
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span>最近课程</span>
              <el-button type="text" @click="$router.push('/teacher/courses')">
                查看全部
              </el-button>
            </div>
          </template>
          <div class="recent-courses">
            <div v-for="course in recentCourses" :key="course.id" class="course-item">
              <div class="course-info">
                <h4>{{ course.name }}</h4>
                <p>{{ course.description }}</p>
                <div class="course-meta">
                  <el-tag size="small" :type="course.status === 'active' ? 'success' : 'info'">
                    {{ course.status === 'active' ? '进行中' : '已结束' }}
                  </el-tag>
                  <span class="student-count">{{ course.student_count }}名学生</span>
                </div>
              </div>
              <el-button type="primary" size="small" @click="$router.push(`/teacher/courses/${course.id}`)">
                管理
              </el-button>
            </div>
          </div>
        </el-card>

        <!-- 学生进度概览 -->
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span>学生进度概览</span>
              <el-button type="text" @click="$router.push('/teacher/students/progress')">
                查看详情
              </el-button>
            </div>
          </template>
          <div class="progress-overview">
            <div class="progress-item">
              <div class="progress-label">平均完成率</div>
              <el-progress :percentage="stats.avgCompletion" :color="getProgressColor(stats.avgCompletion)" />
            </div>
            <div class="progress-item">
              <div class="progress-label">平均成绩</div>
              <el-progress :percentage="stats.avgScore" :color="getProgressColor(stats.avgScore)" />
            </div>
          </div>
        </el-card>
      </div>

      <div class="content-right">
        <!-- 快捷操作 -->
        <el-card class="content-card">
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" plain block @click="$router.push('/teacher/exercises/ai-generate')">
              <el-icon><MagicStick /></el-icon>
              AI生成练习
            </el-button>
            <el-button type="success" plain block @click="$router.push('/teacher/materials')">
              <el-icon><Upload /></el-icon>
              上传资料
            </el-button>
            <el-button type="warning" plain block @click="$router.push('/teacher/chat')">
              <el-icon><ChatDotRound /></el-icon>
              AI教学助手
            </el-button>
            <el-button type="info" plain block @click="$router.push('/teacher/analytics')">
              <el-icon><DataAnalysis /></el-icon>
              数据分析
            </el-button>
          </div>
        </el-card>

        <!-- 系统通知 -->
        <el-card class="content-card">
          <template #header>
            <span>系统通知</span>
          </template>
          <div class="notifications">
            <div v-for="notification in notifications" :key="notification.id" class="notification-item">
              <div class="notification-icon">
                <el-icon :color="notification.type === 'success' ? '#67c23a' : '#e6a23c'">
                  <component :is="notification.icon" />
                </el-icon>
              </div>
              <div class="notification-content">
                <p>{{ notification.message }}</p>
                <span class="notification-time">{{ notification.time }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

// 当前日期
const currentDate = computed(() => {
  return new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})

// 统计数据
const stats = ref({
  courseCount: 0,
  studentCount: 0,
  exerciseCount: 0,
  materialCount: 0,
  avgCompletion: 0,
  avgScore: 0
})

// 最近课程
const recentCourses = ref([])

// 系统通知
const notifications = ref([
  {
    id: 1,
    type: 'success',
    icon: 'Check',
    message: '新课程《Web开发基础》创建成功',
    time: '2小时前'
  },
  {
    id: 2,
    type: 'warning',
    icon: 'Bell',
    message: '有5名学生提交了新的练习答案',
    time: '4小时前'
  },
  {
    id: 3,
    type: 'success',
    icon: 'Check',
    message: 'AI助手已为您生成10道练习题',
    time: '1天前'
  }
])

// 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage >= 80) return '#67c23a'
  if (percentage >= 60) return '#e6a23c'
  return '#f56c6c'
}

// 加载数据
const loadData = async () => {
  // 这里应该调用API获取真实数据
  // 暂时使用模拟数据
  stats.value = {
    courseCount: 8,
    studentCount: 156,
    exerciseCount: 45,
    materialCount: 23,
    avgCompletion: 78,
    avgScore: 85
  }

  recentCourses.value = [
    {
      id: 1,
      name: 'Web开发基础',
      description: '学习HTML、CSS、JavaScript基础知识',
      status: 'active',
      student_count: 45
    },
    {
      id: 2,
      name: '数据结构与算法',
      description: '掌握基本数据结构和算法思想',
      status: 'active',
      student_count: 38
    },
    {
      id: 3,
      name: '数据库原理',
      description: '学习数据库设计和SQL语言',
      status: 'active',
      student_count: 42
    }
  ]
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.teacher-home {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 30px;
  border-radius: 12px;
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-content h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.welcome-content p {
  margin: 0;
  opacity: 0.9;
}

.welcome-actions {
  display: flex;
  gap: 12px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  border: none;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.courses { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.stat-icon.students { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.stat-icon.exercises { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.stat-icon.materials { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }

.stat-info h3 {
  margin: 0 0 4px 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
}

.stat-info p {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
}

.main-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

.content-card {
  margin-bottom: 24px;
  border: none;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.recent-courses {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.course-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  transition: all 0.3s;
}

.course-item:hover {
  border-color: #3498db;
  box-shadow: 0 2px 8px rgba(52, 152, 219, 0.2);
}

.course-info h4 {
  margin: 0 0 4px 0;
  color: #2c3e50;
}

.course-info p {
  margin: 0 0 8px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.course-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.student-count {
  color: #7f8c8d;
  font-size: 12px;
}

.progress-overview {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.progress-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.progress-label {
  font-size: 14px;
  color: #2c3e50;
  font-weight: 500;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notifications {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  background: #f8f9fa;
  transition: background-color 0.3s;
}

.notification-item:hover {
  background: #e9ecef;
}

.notification-icon {
  margin-top: 2px;
}

.notification-content p {
  margin: 0 0 4px 0;
  color: #2c3e50;
  font-size: 14px;
}

.notification-time {
  color: #7f8c8d;
  font-size: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .welcome-section {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .welcome-actions {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .main-content {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style> 