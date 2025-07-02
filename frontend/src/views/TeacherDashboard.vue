<template>
  <el-container class="teacher-dashboard">
    <!-- 侧边栏 -->
    <el-aside width="280px" class="sidebar">
      <div class="logo">
        <h3>教师管理平台</h3>
        <p class="subtitle">Teaching Management System</p>
      </div>
      
      <el-menu
        :default-active="$route.path"
        class="sidebar-menu"
        router
        background-color="#2c3e50"
        text-color="#ecf0f1"
        active-text-color="#3498db"
      >
        <el-menu-item index="/teacher/dashboard">
          <el-icon><House /></el-icon>
          <span>工作台</span>
        </el-menu-item>
        
        <el-sub-menu index="courses">
          <template #title>
            <el-icon><Reading /></el-icon>
            <span>课程管理</span>
          </template>
          <el-menu-item index="/teacher/courses">
            <el-icon><Document /></el-icon>
            <span>我的课程</span>
          </el-menu-item>
          <el-menu-item index="/teacher/courses/create">
            <el-icon><Plus /></el-icon>
            <span>创建课程</span>
          </el-menu-item>
          <el-menu-item index="/teacher/materials">
            <el-icon><Folder /></el-icon>
            <span>资料管理</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="exercises">
          <template #title>
            <el-icon><Edit /></el-icon>
            <span>练习管理</span>
          </template>
          <el-menu-item index="/teacher/exercises">
            <el-icon><List /></el-icon>
            <span>练习列表</span>
          </el-menu-item>
          <el-menu-item index="/teacher/exercises/create">
            <el-icon><Plus /></el-icon>
            <span>创建练习</span>
          </el-menu-item>
          <el-menu-item index="/teacher/exercises/ai-generate">
            <el-icon><MagicStick /></el-icon>
            <span>AI生成练习</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="students">
          <template #title>
            <el-icon><User /></el-icon>
            <span>学生管理</span>
          </template>
          <el-menu-item index="/teacher/students">
            <el-icon><Avatar /></el-icon>
            <span>学生列表</span>
          </el-menu-item>
          <el-menu-item index="/teacher/students/progress">
            <el-icon><TrendCharts /></el-icon>
            <span>学习进度</span>
          </el-menu-item>
          <el-menu-item index="/teacher/students/performance">
            <el-icon><DataAnalysis /></el-icon>
            <span>成绩分析</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-menu-item index="/teacher/chat">
          <el-icon><ChatDotRound /></el-icon>
          <span>AI助手</span>
        </el-menu-item>
        
        <el-menu-item index="/teacher/lesson-plan">
          <el-icon><Calendar /></el-icon>
          <span>备课助手</span>
        </el-menu-item>
        
        <el-menu-item index="/teacher/analytics">
          <el-icon><PieChart /></el-icon>
          <span>数据分析</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/teacher/dashboard' }">教师工作台</el-breadcrumb-item>
            <el-breadcrumb-item v-if="$route.name !== 'TeacherHome'">
              {{ getPageTitle() }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-center">
          <el-tag type="success" size="small">教师端</el-tag>
        </div>
        
        <div class="header-right">
          <el-button type="primary" size="small" @click="switchToStudent">
            <el-icon><Switch /></el-icon>
            切换到学生端
          </el-button>
          
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="36" :src="user?.avatar">
                {{ user?.real_name?.charAt(0) }}
              </el-avatar>
              <span class="username">{{ user?.real_name || user?.username }}</span>
              <el-tag type="warning" size="small">教师</el-tag>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人资料
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>
                  系统设置
                </el-dropdown-item>
                <el-dropdown-item command="help">
                  <el-icon><QuestionFilled /></el-icon>
                  帮助文档
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const user = computed(() => authStore.user)

const getPageTitle = () => {
  const routeMap = {
    'TeacherCourses': '课程管理',
    'TeacherCourseCreate': '创建课程',
    'TeacherCourseDetail': '课程详情',
    'TeacherMaterials': '资料管理',
    'TeacherExercises': '练习管理',
    'TeacherExerciseCreate': '创建练习',
    'TeacherExerciseGenerate': 'AI生成练习',
    'TeacherStudents': '学生管理',
    'TeacherStudentProgress': '学习进度',
    'TeacherStudentPerformance': '成绩分析',
    'TeacherChat': 'AI助手',
    'TeacherLessonPlan': '备课助手',
    'TeacherAnalytics': '数据分析',
    'TeacherProfile': '个人资料'
  }
  return routeMap[route.name] || ''
}

const switchToStudent = () => {
  router.push('/dashboard')
}

const handleCommand = async (command) => {
  switch (command) {
    case 'profile':
      router.push('/teacher/profile')
      break
    case 'settings':
      router.push('/teacher/settings')
      break
    case 'help':
      window.open('/help', '_blank')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        authStore.logout()
        router.push('/login')
      } catch {
        // 用户取消
      }
      break
  }
}
</script>

<style scoped>
.teacher-dashboard {
  height: 100vh;
}

.sidebar {
  background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
  color: white;
  box-shadow: 2px 0 8px rgba(0,0,0,0.1);
}

.logo {
  height: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #34495e;
  background: rgba(255,255,255,0.05);
}

.logo h3 {
  color: white;
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
}

.subtitle {
  color: #bdc3c7;
  margin: 0;
  font-size: 12px;
}

.sidebar-menu {
  border: none;
  margin-top: 10px;
}

.sidebar-menu .el-menu-item,
.sidebar-menu .el-sub-menu__title {
  height: 50px;
  line-height: 50px;
  margin: 2px 0;
  border-radius: 0 25px 25px 0;
  margin-right: 10px;
}

.sidebar-menu .el-menu-item:hover,
.sidebar-menu .el-sub-menu__title:hover {
  background-color: rgba(52, 152, 219, 0.2) !important;
}

.sidebar-menu .el-menu-item.is-active {
  background-color: #3498db !important;
  color: white !important;
}

.header {
  background: white;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.header-left {
  flex: 1;
}

.header-center {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.3s;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
}

.user-info:hover {
  background-color: #e9ecef;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.username {
  margin: 0 8px;
  color: #333;
  font-weight: 500;
}

.main {
  background: #f5f7fa;
  padding: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 200px !important;
  }
  
  .header {
    padding: 0 10px;
  }
  
  .header-right {
    gap: 10px;
  }
  
  .user-info {
    padding: 6px 8px;
  }
  
  .username {
    display: none;
  }
}
</style> 