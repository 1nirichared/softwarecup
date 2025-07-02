import api from './index'

export const courseApi = {
  // 获取课程列表
  getCourses: () => {
    return api.get('/courses')
  },

  // 获取课程详情
  getCourse: (id) => {
    return api.get(`/courses/${id}`)
  },

  // 获取课程统计
  getCourseStats: (id) => {
    return api.get(`/courses/${id}/stats`)
  },

  // 创建课程（教师）
  createCourse: (data) => {
    return api.post('/courses', data)
  },

  // 更新课程（教师）
  updateCourse: (id, data) => {
    return api.put(`/courses/${id}`, data)
  },

  // 删除课程（教师）
  deleteCourse: (id) => {
    return api.delete(`/courses/${id}`)
  },

  // 生成备课内容（教师）
  generateLessonPlan: (courseId, chapterId) => {
    return api.post(`/courses/${courseId}/chapters/${chapterId}/lesson-plan`)
  },

  // 课程材料相关接口
  // 获取课程材料列表
  getCourseMaterials: (courseId) => {
    return api.get(`/course-materials/${courseId}`)
  },

  // 上传课程材料（教师）
  uploadCourseMaterial: (courseId, formData) => {
    return api.post(`/course-materials/${courseId}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 删除课程材料（教师）
  deleteCourseMaterial: (materialId) => {
    return api.delete(`/course-materials/${materialId}`)
  }
} 