<template>
  <div class="teacher-courses">
    <div class="page-header">
      <h2>我的课程</h2>
      <el-button type="primary" @click="$router.push('/teacher/courses/create')">
        <el-icon><Plus /></el-icon>
        创建新课程
      </el-button>
    </div>

    <div class="courses-grid">
      <el-card v-for="course in courses" :key="course.id" class="course-card">
        <div class="course-header">
          <h3>{{ course.name }}</h3>
          <el-tag :type="course.status === 'active' ? 'success' : 'info'" size="small">
            {{ course.status === 'active' ? '进行中' : '已结束' }}
          </el-tag>
        </div>
        
        <p class="course-description">{{ course.description }}</p>
        
        <div class="course-stats">
          <div class="stat-item">
            <el-icon><User /></el-icon>
            <span>{{ course.student_count }}名学生</span>
          </div>
          <div class="stat-item">
            <el-icon><Document /></el-icon>
            <span>{{ course.chapter_count }}个章节</span>
          </div>
          <div class="stat-item">
            <el-icon><Edit /></el-icon>
            <span>{{ course.exercise_count }}个练习</span>
          </div>
        </div>
        
        <div class="course-actions">
          <el-button type="primary" size="small" @click="$router.push(`/teacher/courses/${course.id}`)">
            管理课程
          </el-button>
          <el-button type="success" size="small" @click="$router.push(`/teacher/courses/${course.id}/materials`)">
            资料管理
          </el-button>
          <el-button type="warning" size="small" @click="$router.push(`/teacher/courses/${course.id}/exercises`)">
            练习管理
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 空状态 -->
    <div v-if="courses.length === 0" class="empty-state">
      <el-empty description="暂无课程">
        <el-button type="primary" @click="$router.push('/teacher/courses/create')">
          创建第一个课程
        </el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const courses = ref([])

const loadCourses = async () => {
  // 这里应该调用API获取教师的课程列表
  // 暂时使用模拟数据
  courses.value = [
    {
      id: 1,
      name: 'Web开发基础',
      description: '学习HTML、CSS、JavaScript基础知识，掌握现代Web开发技术栈',
      status: 'active',
      student_count: 45,
      chapter_count: 8,
      exercise_count: 12
    },
    {
      id: 2,
      name: '数据结构与算法',
      description: '掌握基本数据结构和算法思想，提高编程能力和问题解决能力',
      status: 'active',
      student_count: 38,
      chapter_count: 10,
      exercise_count: 15
    },
    {
      id: 3,
      name: '数据库原理',
      description: '学习数据库设计和SQL语言，理解数据存储和查询原理',
      status: 'active',
      student_count: 42,
      chapter_count: 6,
      exercise_count: 8
    },
    {
      id: 4,
      name: '软件工程',
      description: '了解软件开发流程和项目管理方法，培养工程化思维',
      status: 'active',
      student_count: 35,
      chapter_count: 7,
      exercise_count: 10
    }
  ]
}

onMounted(() => {
  loadCourses()
})
</script>

<style scoped>
.teacher-courses {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.course-card {
  border: none;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  transition: all 0.3s;
  height: fit-content;
}

.course-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 20px rgba(0,0,0,0.15);
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.course-header h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 18px;
  font-weight: 600;
  flex: 1;
  margin-right: 12px;
}

.course-description {
  color: #7f8c8d;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 16px;
  min-height: 40px;
}

.course-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #7f8c8d;
  font-size: 12px;
}

.course-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .courses-grid {
    grid-template-columns: 1fr;
  }
  
  .course-stats {
    flex-direction: column;
    gap: 8px;
  }
  
  .course-actions {
    flex-direction: column;
  }
}
</style> 