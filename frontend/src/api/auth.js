import api from './index'

export const authApi = {
  // 用户登录
  login: (credentials) => {
    return api.post('/auth/login', credentials)
  },

  // 用户注册
  register: (userData) => {
    return api.post('/auth/register', userData)
  },

  // 获取当前用户信息
  getProfile: () => {
    return api.get('/user/profile')
  },

  // 更新用户信息
  updateProfile: (profileData) => {
    return api.put('/user/profile', profileData)
  },

  // 修改密码
  changePassword: (passwordData) => {
    return api.put('/user/password', passwordData)
  }
} 