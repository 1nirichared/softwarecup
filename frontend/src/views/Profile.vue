<template>
  <div class="profile-container">
    <h1>个人资料</h1>
    <div class="profile-card">
      <div class="avatar-section">
        <div class="avatar">{{ user?.real_name?.charAt(0) || user?.username?.charAt(0) }}</div>
        <h2>{{ user?.real_name || user?.username }}</h2>
        <p>{{ user?.email }}</p>
      </div>
      
      <div class="profile-form">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="form.username" disabled />
        </div>
        
        <div class="form-group">
          <label>真实姓名</label>
          <input v-model="form.real_name" />
        </div>
        
        <div class="form-group">
          <label>邮箱</label>
          <input v-model="form.email" type="email" />
        </div>
        
        <div class="form-group">
          <label>角色</label>
          <input v-model="form.role" disabled />
        </div>
        
        <button @click="saveProfile" class="save-btn">保存修改</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

const form = ref({
  username: user.value?.username || '',
  real_name: user.value?.real_name || '',
  email: user.value?.email || '',
  role: user.value?.role || ''
})

const saveProfile = () => {
  alert('保存功能待实现')
}
</script>

<style scoped>
.profile-container {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.profile-card {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.avatar-section {
  text-align: center;
  margin-bottom: 30px;
}

.avatar {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin: 0 auto 15px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.form-group input:disabled {
  background: #f5f5f5;
  color: #666;
}

.save-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  width: 100%;
}
</style> 