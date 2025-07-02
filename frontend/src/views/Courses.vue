<template>
  <div class="courses-container">
    <div class="header">
      <h1>课程列表</h1>
      <div class="filters">
        <input v-model="searchQuery" placeholder="搜索课程..." class="search-input" />
        <select v-model="selectedSubject" class="filter-select">
          <option value="">所有学科</option>
          <option value="math">数学</option>
          <option value="physics">物理</option>
          <option value="chemistry">化学</option>
          <option value="computer">计算机</option>
        </select>
      </div>
    </div>
    
    <div class="courses-grid">
      <div v-for="course in filteredCourses" :key="course.id" class="course-card">
        <div class="course-image">
          <img :src="course.cover_image || '/default-course.jpg'" :alt="course.name" />
        </div>
        <div class="course-content">
          <h3>{{ course.name }}</h3>
          <p class="course-description">{{ course.description }}</p>
          <div class="course-meta">
            <span class="subject">{{ course.subject }}</span>
            <span class="grade">{{ course.grade }}</span>
          </div>
          <div class="course-stats">
            <span>📚 {{ course.chapterCount || 0 }} 章节</span>
            <span>👥 {{ course.studentCount || 0 }} 学生</span>
          </div>
          <div class="course-teacher">
            <span>👨‍🏫 {{ course.teacher?.real_name || '未知教师' }}</span>
          </div>
          <router-link :to="`/dashboard/courses/${course.id}`" class="view-btn">
            查看详情
          </router-link>
        </div>
      </div>
    </div>
    
    <div v-if="loading" class="loading">
      加载中...
    </div>
    
    <div v-if="!loading && filteredCourses.length === 0" class="empty-state">
      <p>暂无课程</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { courseApi } from '@/api/course'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const searchQuery = ref('')
const selectedSubject = ref('')
const courses = ref([])

const filteredCourses = computed(() => {
  return courses.value.filter(course => {
    const matchesSearch = course.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         course.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesSubject = !selectedSubject.value || course.subject === selectedSubject.value
    return matchesSearch && matchesSubject
  })
})

onMounted(async () => {
  loading.value = true
  try {
    const response = await courseApi.getCourses()
    courses.value = response.data.data || []
  } catch (error) {
    console.error('获取课程失败:', error)
    ElMessage.error('获取课程列表失败')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.courses-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 20px;
}

.header h1 {
  color: #333;
  margin: 0;
}

.filters {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.search-input, .filter-select {
  padding: 10px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.search-input {
  min-width: 200px;
}

.courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 25px;
}

.course-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.course-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.15);
}

.course-image {
  height: 200px;
  overflow: hidden;
}

.course-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.course-content {
  padding: 20px;
}

.course-content h3 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1.2rem;
}

.course-description {
  color: #666;
  margin: 0 0 15px 0;
  line-height: 1.5;
}

.course-meta {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.subject, .grade {
  background: #f0f0f0;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.9rem;
  color: #666;
}

.course-stats {
  display: flex;
  gap: 15px;
  margin-bottom: 10px;
  font-size: 0.9rem;
  color: #666;
}

.course-teacher {
  margin-bottom: 20px;
  font-size: 0.9rem;
  color: #666;
}

.view-btn {
  display: inline-block;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 10px 20px;
  border-radius: 6px;
  text-decoration: none;
  text-align: center;
  transition: opacity 0.3s;
}

.view-btn:hover {
  opacity: 0.9;
}

.loading, .empty-state {
  text-align: center;
  padding: 40px;
  color: #666;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filters {
    justify-content: center;
  }
  
  .courses-grid {
    grid-template-columns: 1fr;
  }
}
</style> 