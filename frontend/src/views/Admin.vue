<template>
  <div class="admin-container">
    <h1>管理员面板</h1>
    <div class="admin-stats">
      <div class="stat-card">
        <h3>用户总数</h3>
        <p>{{ stats.userCount }}</p>
      </div>
      <div class="stat-card">
        <h3>课程总数</h3>
        <p>{{ stats.courseCount }}</p>
      </div>
      <div class="stat-card">
        <h3>练习总数</h3>
        <p>{{ stats.exerciseCount }}</p>
      </div>
    </div>
    
    <div class="admin-actions">
      <h2>管理操作</h2>
      <div class="action-buttons">
        <button @click="manageUsers">用户管理</button>
        <button @click="manageCourses">课程管理</button>
        <button @click="manageExercises">练习管理</button>
        <button @click="systemSettings">系统设置</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const stats = ref({
  userCount: 156,
  courseCount: 23,
  exerciseCount: 89
})

const manageUsers = () => alert('用户管理功能待实现')
const manageCourses = () => alert('课程管理功能待实现')
const manageExercises = () => alert('练习管理功能待实现')
const systemSettings = () => alert('系统设置功能待实现')
</script>

<style scoped>
.admin-container {
  padding: 20px;
}

.admin-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  text-align: center;
}

.stat-card h3 {
  margin: 0 0 10px 0;
  color: #666;
}

.stat-card p {
  margin: 0;
  font-size: 2rem;
  color: #667eea;
  font-weight: bold;
}

.admin-actions {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.action-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 15px;
  margin-top: 15px;
}

.action-buttons button {
  background: #667eea;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.action-buttons button:hover {
  opacity: 0.9;
}
</style> 