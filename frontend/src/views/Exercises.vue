<template>
  <div class="exercises-container">
    <h1>练习列表</h1>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="exercises-grid">
      <div v-for="exercise in exercises" :key="exercise.id" class="exercise-card">
        <h3>{{ exercise.title }}</h3>
        <p>{{ exercise.description }}</p>
        <div class="exercise-meta">
          <span>时长: {{ exercise.duration }}分钟</span>
          <span>总分: {{ exercise.total_score }}分</span>
          <span>类型: {{ exercise.type }}</span>
        </div>
        <router-link :to="`/dashboard/exercises/${exercise.id}`" class="start-btn">
          开始练习
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { exerciseApi } from '@/api'

const exercises = ref([])
const loading = ref(true)
const error = ref('')

const loadExercises = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await exerciseApi.getExercises()
    if (response.data && response.data.code === 200) {
      exercises.value = response.data.data || []
    } else {
      error.value = response.data?.message || '获取练习列表失败'
    }
  } catch (err) {
    console.error('获取练习列表失败:', err)
    error.value = err.response?.data?.message || '网络错误，请重试'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadExercises()
})
</script>

<style scoped>
.exercises-container {
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

.exercises-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.exercise-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.exercise-meta {
  display: flex;
  gap: 15px;
  margin: 15px 0;
  font-size: 0.9rem;
  color: #666;
  flex-wrap: wrap;
}

.start-btn {
  background: #667eea;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  text-decoration: none;
  display: inline-block;
}
</style> 