<template>
  <div class="course-detail-container">
    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="course" class="course-content">
      <div class="course-header">
        <div class="course-image">
          <img :src="course.coverImage || '/default-course.jpg'" :alt="course.name" />
        </div>
        <div class="course-info">
          <h1>{{ course.name }}</h1>
          <p class="description">{{ course.description }}</p>
          <div class="meta">
            <span class="subject">{{ course.subject }}</span>
            <span class="grade">{{ course.grade }}</span>
            <span class="teacher">教师: {{ course.teacher?.real_name || '未知' }}</span>
          </div>
        </div>
      </div>
      
      <div class="course-sections">
        <!-- 课程材料部分 -->
        <div class="materials-section">
          <div class="section-header">
            <h2>课程材料</h2>
            <button v-if="isTeacher" @click="showUploadDialog = true" class="upload-btn">
              📁 上传材料
            </button>
          </div>
          
          <div v-if="materials.length > 0" class="materials-list">
            <div v-for="material in materials" :key="material.id" class="material-item">
              <div class="material-info">
                <div class="material-icon">📄</div>
                <div class="material-details">
                  <h4>{{ material.title }}</h4>
                  <p class="material-meta">
                    上传时间: {{ formatDate(material.uploaded_at) }} | 
                    文件类型: {{ material.file_type }}
                  </p>
                </div>
              </div>
              <div class="material-actions">
                <a :href="material.file_url" target="_blank" class="download-btn">
                  下载
                </a>
                <button v-if="isTeacher" @click="deleteMaterial(material.id)" class="delete-btn">
                  删除
                </button>
              </div>
            </div>
          </div>
          
          <div v-else class="empty-materials">
            <p>暂无课程材料</p>
            <button v-if="isTeacher" @click="showUploadDialog = true" class="upload-btn">
              上传第一个材料
            </button>
          </div>
        </div>
        
        <div class="chapters-section">
          <h2>课程章节</h2>
          <div class="chapters-list">
            <div v-for="chapter in chapters" :key="chapter.id" class="chapter-item">
              <div class="chapter-info">
                <h3>{{ chapter.title }}</h3>
                <p>{{ chapter.description }}</p>
              </div>
              <div class="chapter-actions">
                <button class="view-btn">查看内容</button>
              </div>
            </div>
          </div>
        </div>
        
        <div class="exercises-section">
          <h2>相关练习</h2>
          <div class="exercises-list">
            <div v-for="exercise in exercises" :key="exercise.id" class="exercise-item">
              <div class="exercise-info">
                <h3>{{ exercise.title }}</h3>
                <p>{{ exercise.description }}</p>
                <div class="exercise-meta">
                  <span>时长: {{ exercise.duration }}分钟</span>
                  <span>总分: {{ exercise.totalScore }}分</span>
                </div>
              </div>
              <router-link :to="`/dashboard/exercises/${exercise.id}`" class="start-btn">
                开始练习
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="error">课程不存在</div>

    <!-- 上传材料对话框 -->
    <div v-if="showUploadDialog" class="upload-dialog-overlay" @click="showUploadDialog = false">
      <div class="upload-dialog" @click.stop>
        <h3>上传课程材料</h3>
        <form @submit.prevent="uploadMaterial" class="upload-form">
          <div class="form-group">
            <label for="title">材料标题</label>
            <input
              id="title"
              v-model="uploadForm.title"
              type="text"
              placeholder="请输入材料标题"
              required
            />
          </div>
          
          <div class="form-group">
            <label for="file">选择文件</label>
            <input
              id="file"
              ref="fileInput"
              type="file"
              @change="handleFileSelect"
              required
            />
            <p class="file-info" v-if="selectedFile">
              已选择: {{ selectedFile.name }}
            </p>
          </div>
          
          <div class="dialog-actions">
            <button type="button" @click="showUploadDialog = false" class="cancel-btn">
              取消
            </button>
            <button type="submit" :disabled="uploading" class="submit-btn">
              {{ uploading ? '上传中...' : '上传' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { courseApi } from '@/api/course'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const authStore = useAuthStore()
const loading = ref(true)
const course = ref(null)
const chapters = ref([])
const exercises = ref([])
const materials = ref([])
const showUploadDialog = ref(false)
const uploading = ref(false)
const selectedFile = ref(null)
const fileInput = ref(null)

const uploadForm = ref({
  title: ''
})

const isTeacher = computed(() => {
  return authStore.user?.role === 'teacher'
})

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const handleFileSelect = (event) => {
  selectedFile.value = event.target.files[0]
}

const uploadMaterial = async () => {
  if (!selectedFile.value) {
    ElMessage.error('请选择文件')
    return
  }

  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    formData.append('title', uploadForm.title)

    await courseApi.uploadCourseMaterial(route.params.id, formData)
    
    ElMessage.success('材料上传成功')
    showUploadDialog.value = false
    uploadForm.value.title = ''
    selectedFile.value = null
    if (fileInput.value) {
      fileInput.value.value = ''
    }
    
    // 重新加载材料列表
    await loadMaterials()
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败: ' + (error.response?.data?.message || error.message))
  } finally {
    uploading.value = false
  }
}

const deleteMaterial = async (materialId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个材料吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await courseApi.deleteCourseMaterial(materialId)
    ElMessage.success('删除成功')
    await loadMaterials()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败: ' + (error.response?.data?.message || error.message))
    }
  }
}

const loadMaterials = async () => {
  try {
    const response = await courseApi.getCourseMaterials(route.params.id)
    materials.value = response.data.data || []
  } catch (error) {
    console.error('获取材料失败:', error)
  }
}

onMounted(async () => {
  try {
    // 模拟数据
    course.value = {
      id: route.params.id,
      name: 'JavaScript基础',
      description: '学习JavaScript编程语言的基础知识和核心概念，包括变量、函数、对象、数组等核心内容。',
      subject: '计算机科学',
      grade: '大学',
      coverImage: '/js-course.jpg',
      teacher: { real_name: '张老师' }
    }
    
    chapters.value = [
      {
        id: 1,
        title: '第一章：JavaScript简介',
        description: '了解JavaScript的历史、特点和应用场景'
      },
      {
        id: 2,
        title: '第二章：变量和数据类型',
        description: '学习JavaScript中的变量声明和基本数据类型'
      },
      {
        id: 3,
        title: '第三章：函数',
        description: '掌握函数的定义、调用和参数传递'
      }
    ]
    
    exercises.value = [
      {
        id: 1,
        title: 'JavaScript基础练习',
        description: '测试对JavaScript基础知识的掌握程度',
        duration: 30,
        totalScore: 100
      },
      {
        id: 2,
        title: '函数编程练习',
        description: '练习函数的定义和使用',
        duration: 45,
        totalScore: 100
      }
    ]

    // 加载课程材料
    await loadMaterials()
  } catch (error) {
    console.error('获取课程详情失败:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.course-detail-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  color: #666;
}

.course-header {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 30px;
  margin-bottom: 40px;
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.course-image img {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 8px;
}

.course-info h1 {
  margin: 0 0 15px 0;
  color: #333;
}

.description {
  color: #666;
  line-height: 1.6;
  margin-bottom: 20px;
}

.meta {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.subject, .grade, .teacher {
  background: #f0f0f0;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 0.9rem;
  color: #666;
}

.course-sections {
  display: grid;
  gap: 30px;
}

.materials-section, .chapters-section, .exercises-section {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0;
  color: #333;
}

.upload-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: opacity 0.3s;
}

.upload-btn:hover {
  opacity: 0.9;
}

.materials-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.material-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  transition: border-color 0.3s;
}

.material-item:hover {
  border-color: #667eea;
}

.material-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.material-icon {
  font-size: 24px;
}

.material-details h4 {
  margin: 0 0 5px 0;
  color: #333;
}

.material-meta {
  margin: 0;
  font-size: 0.9rem;
  color: #666;
}

.material-actions {
  display: flex;
  gap: 10px;
}

.download-btn, .delete-btn {
  padding: 6px 12px;
  border-radius: 4px;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  transition: opacity 0.3s;
}

.download-btn {
  background: #67c23a;
  color: white;
}

.delete-btn {
  background: #f56c6c;
  color: white;
  border: none;
}

.download-btn:hover, .delete-btn:hover {
  opacity: 0.8;
}

.empty-materials {
  text-align: center;
  padding: 40px;
  color: #666;
}

.chapters-list, .exercises-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.chapter-item, .exercise-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  transition: border-color 0.3s;
}

.chapter-item:hover, .exercise-item:hover {
  border-color: #667eea;
}

.chapter-info h3, .exercise-info h3 {
  margin: 0 0 8px 0;
  color: #333;
}

.chapter-info p, .exercise-info p {
  margin: 0 0 10px 0;
  color: #666;
}

.exercise-meta {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: #666;
}

.view-btn, .start-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  transition: opacity 0.3s;
}

.view-btn:hover, .start-btn:hover {
  opacity: 0.9;
}

/* 上传对话框样式 */
.upload-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.upload-dialog {
  background: white;
  padding: 30px;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
}

.upload-dialog h3 {
  margin: 0 0 20px 0;
  color: #333;
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-weight: 500;
  color: #333;
}

.form-group input {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.file-info {
  margin: 5px 0 0 0;
  font-size: 0.9rem;
  color: #666;
}

.dialog-actions {
  display: flex;
  gap: 15px;
  justify-content: flex-end;
  margin-top: 20px;
}

.cancel-btn, .submit-btn {
  padding: 10px 20px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  transition: opacity 0.3s;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.submit-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.cancel-btn:hover, .submit-btn:hover {
  opacity: 0.8;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .course-header {
    grid-template-columns: 1fr;
  }
  
  .material-item, .chapter-item, .exercise-item {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
  }
  
  .material-actions {
    justify-content: center;
  }
  
  .section-header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
}
</style> 