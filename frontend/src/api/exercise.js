import api from './index'

export const exerciseApi = {
  // 获取练习列表
  getExercises: (params = {}) => {
    return api.get('/exercises', { params })
  },

  // 获取练习详情
  getExercise: (id) => {
    return api.get(`/exercises/${id}`)
  },

  // 创建练习（教师）
  createExercise: (data) => {
    return api.post('/exercises', data)
  },

  // 生成练习题（教师）
  generateExercises: (courseId, chapterId, params = {}) => {
    return api.post(`/exercises/${courseId}/chapters/${chapterId}/generate`, null, { params })
  },

  // 开始练习
  startExercise: (exerciseId) => {
    return api.post(`/exercise-records/start/${exerciseId}`)
  },

  // 提交答案
  submitAnswer: (recordId, data) => {
    return api.post(`/exercise-records/${recordId}/answers`, data)
  },

  // 完成练习
  completeExercise: (recordId) => {
    return api.post(`/exercise-records/${recordId}/complete`)
  },

  // 获取练习统计
  getExerciseStats: () => {
    return api.get('/exercises/stats')
  }
} 