<template>
  <div class="exercise-detail-container">
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="exercise" class="exercise-content">
      <h1>{{ exercise.title }}</h1>
      <p>{{ exercise.description }}</p>
      <div class="exercise-info">
        <span>时长: {{ exercise.duration }}分钟</span>
        <span>总分: {{ exercise.total_score }}分</span>
        <span>类型: {{ exercise.type }}</span>
      </div>
      
      <div v-if="!isStarted" class="start-section">
        <button class="start-btn" @click="startExercise" :disabled="starting">
          {{ starting ? '开始中...' : '开始练习' }}
        </button>
      </div>
      
      <div v-else class="exercise-questions">
        <h2>练习题目</h2>
        <div v-for="(question, index) in questions" :key="question.id" class="question-item">
          <h3>第{{ index + 1 }}题</h3>
          <p>{{ question.title }}</p>
          <div v-if="question.content" class="question-content">
            {{ question.content }}
          </div>
          
          <div v-if="question.type === 'single' || question.type === 'multiple'" class="options">
            <div v-for="option in parseOptions(question.options)" :key="option.key" class="option">
              <input 
                :type="question.type === 'single' ? 'radio' : 'checkbox'"
                :name="`question-${question.id}`"
                :value="option.key"
                v-model="answers[question.id]"
              />
              <label>{{ option.key }}. {{ option.value }}</label>
            </div>
          </div>
          
          <div v-else class="answer-input">
            <textarea 
              v-model="answers[question.id]"
              placeholder="请输入你的答案..."
              rows="4"
            ></textarea>
          </div>
        </div>
        
        <div class="submit-section">
          <button class="submit-btn" @click="submitAnswers" :disabled="submitting">
            {{ submitting ? '提交中...' : '提交答案' }}
          </button>
        </div>
      </div>
    </div>
    <div v-else class="error">练习不存在</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { exerciseApi } from '@/api'

const route = useRoute()
const router = useRouter()

const exercise = ref(null)
const questions = ref([])
const answers = ref({})
const loading = ref(true)
const error = ref('')
const isStarted = ref(false)
const starting = ref(false)
const submitting = ref(false)
const recordId = ref(null)

const loadExercise = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await exerciseApi.getExercise(route.params.id)
    if (response.data && response.data.code === 200) {
      exercise.value = response.data.data
      questions.value = response.data.data.questions || []
    } else {
      error.value = response.data?.message || '获取练习详情失败'
    }
  } catch (err) {
    console.error('获取练习详情失败:', err)
    error.value = err.response?.data?.message || '网络错误，请重试'
  } finally {
    loading.value = false
  }
}

const startExercise = async () => {
  try {
    starting.value = true
    const response = await exerciseApi.startExercise(route.params.id)
    if (response.data && response.data.code === 200) {
      isStarted.value = true
      recordId.value = response.data.data.record_id
      // 初始化答案对象
      questions.value.forEach(q => {
        answers.value[q.id] = q.type === 'multiple' ? [] : ''
      })
    } else {
      alert(response.data?.message || '开始练习失败')
    }
  } catch (err) {
    console.error('开始练习失败:', err)
    alert(err.response?.data?.message || '开始练习失败，请重试')
  } finally {
    starting.value = false
  }
}

const submitAnswers = async () => {
  try {
    submitting.value = true
    
    // 提交每个答案
    for (const [questionId, answer] of Object.entries(answers.value)) {
      if (answer !== '' && answer.length > 0) {
        await exerciseApi.submitAnswer(recordId.value, {
          question_id: parseInt(questionId),
          answer: Array.isArray(answer) ? answer.join(',') : answer
        })
      }
    }
    
    // 完成练习
    const completeResponse = await exerciseApi.completeExercise(recordId.value)
    if (completeResponse.data && completeResponse.data.code === 200) {
      const score = completeResponse.data.data.total_score
      alert(`练习完成！你的得分是：${score}分`)
      router.push('/dashboard/exercises')
    }
  } catch (err) {
    console.error('提交答案失败:', err)
    alert(err.response?.data?.message || '提交答案失败，请重试')
  } finally {
    submitting.value = false
  }
}

const parseOptions = (optionsStr) => {
  try {
    return JSON.parse(optionsStr)
  } catch {
    return []
  }
}

onMounted(() => {
  loadExercise()
})
</script>

<style scoped>
.exercise-detail-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  color: #f56c6c;
}

.exercise-content {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.exercise-info {
  display: flex;
  gap: 20px;
  margin: 20px 0;
  flex-wrap: wrap;
}

.start-section, .submit-section {
  text-align: center;
  margin: 30px 0;
}

.start-btn, .submit-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 30px;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
  transition: opacity 0.3s;
}

.start-btn:hover:not(:disabled), .submit-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.start-btn:disabled, .submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.exercise-questions {
  margin-top: 30px;
}

.question-item {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
}

.question-content {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  margin: 10px 0;
}

.options {
  margin: 15px 0;
}

.option {
  display: flex;
  align-items: center;
  margin: 10px 0;
  gap: 10px;
}

.answer-input textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  resize: vertical;
}
</style> 