import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue')
      },
      {
        path: 'courses',
        name: 'Courses',
        component: () => import('@/views/Courses.vue')
      },
      {
        path: 'courses/:id',
        name: 'CourseDetail',
        component: () => import('@/views/CourseDetail.vue')
      },
      {
        path: 'exercises',
        name: 'Exercises',
        component: () => import('@/views/Exercises.vue')
      },
      {
        path: 'exercises/:id',
        name: 'ExerciseDetail',
        component: () => import('@/views/ExerciseDetail.vue')
      },
      {
        path: 'chat',
        name: 'Chat',
        component: () => import('@/views/Chat.vue')
      },
      {
        path: 'chat/:id',
        name: 'ChatSession',
        component: () => import('@/views/ChatSession.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue')
      },
      {
        path: 'admin',
        name: 'Admin',
        component: () => import('@/views/Admin.vue'),
        meta: { requiresAdmin: true }
      }
    ]
  },
  {
    path: '/teacher',
    name: 'TeacherDashboard',
    component: () => import('@/views/TeacherDashboard.vue'),
    meta: { requiresAuth: true, requiresTeacher: true },
    children: [
      {
        path: '',
        redirect: '/teacher/dashboard'
      },
      {
        path: 'dashboard',
        name: 'TeacherHome',
        component: () => import('@/views/TeacherHome.vue')
      },
      {
        path: 'courses',
        name: 'TeacherCourses',
        component: () => import('@/views/TeacherCourses.vue')
      },
      {
        path: 'courses/create',
        name: 'TeacherCourseCreate',
        component: () => import('@/views/TeacherCourseCreate.vue')
      },
      {
        path: 'courses/:id',
        name: 'TeacherCourseDetail',
        component: () => import('@/views/TeacherCourseDetail.vue')
      },
      {
        path: 'materials',
        name: 'TeacherMaterials',
        component: () => import('@/views/TeacherMaterials.vue')
      },
      {
        path: 'exercises',
        name: 'TeacherExercises',
        component: () => import('@/views/TeacherExercises.vue')
      },
      {
        path: 'exercises/create',
        name: 'TeacherExerciseCreate',
        component: () => import('@/views/TeacherExerciseCreate.vue')
      },
      {
        path: 'exercises/ai-generate',
        name: 'TeacherExerciseGenerate',
        component: () => import('@/views/TeacherExerciseGenerate.vue')
      },
      {
        path: 'students',
        name: 'TeacherStudents',
        component: () => import('@/views/TeacherStudents.vue')
      },
      {
        path: 'students/progress',
        name: 'TeacherStudentProgress',
        component: () => import('@/views/TeacherStudentProgress.vue')
      },
      {
        path: 'students/performance',
        name: 'TeacherStudentPerformance',
        component: () => import('@/views/TeacherStudentPerformance.vue')
      },
      {
        path: 'chat',
        name: 'TeacherChat',
        component: () => import('@/views/TeacherChat.vue')
      },
      {
        path: 'lesson-plan',
        name: 'TeacherLessonPlan',
        component: () => import('@/views/TeacherLessonPlan.vue')
      },
      {
        path: 'analytics',
        name: 'TeacherAnalytics',
        component: () => import('@/views/TeacherAnalytics.vue')
      },
      {
        path: 'profile',
        name: 'TeacherProfile',
        component: () => import('@/views/TeacherProfile.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresAdmin && authStore.user?.role !== 'admin') {
    next('/dashboard')
  } else if (to.meta.requiresTeacher && authStore.user?.role !== 'teacher') {
    next('/dashboard')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    // 根据用户角色重定向到对应的工作台
    if (authStore.user?.role === 'teacher') {
      next('/teacher/dashboard')
    } else {
      next('/dashboard')
    }
  } else {
    next()
  }
})

export default router 