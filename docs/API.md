# 智能教学实训平台 - API文档

## 📋 概述

本文档详细说明智能教学实训平台的后端API接口，包括认证、用户管理、课程管理、练习管理、AI对话等功能。

## 🔐 认证机制

### JWT Token认证
所有需要认证的接口都需要在请求头中携带JWT Token：

```
Authorization: Bearer <your_jwt_token>
```

### Token获取
通过登录接口获取Token，Token有效期为24小时。

## 📊 响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 具体数据
  }
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

## 认证相关

### 用户登录
```
POST /auth/login
```

请求体:
```json
{
  "username": "string",
  "password": "string"
}
```

响应:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "string",
    "user": {
      "id": 1,
      "username": "string",
      "email": "string",
      "real_name": "string",
      "role": "string",
      "avatar": "string"
    }
  }
}
```

### 用户注册
```
POST /auth/register
```

请求体:
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "real_name": "string",
  "role": "string"
}
```

## 用户相关

### 获取当前用户信息
```
GET /user/profile
```

### 修改密码
```
PUT /user/password
```

请求体:
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

## 课程相关

### 获取课程列表
```
GET /courses
```

### 获取课程详情
```
GET /courses/{id}
```

### 创建课程 (教师)
```
POST /courses
```

请求体:
```json
{
  "name": "string",
  "description": "string",
  "subject": "string",
  "grade": "string",
  "cover_image": "string"
}
```

### 更新课程 (教师)
```
PUT /courses/{id}
```

### 删除课程 (教师)
```
DELETE /courses/{id}
```

### 生成备课内容 (教师)
```
POST /courses/{courseId}/chapters/{chapterId}/lesson-plan
```

### 获取课程统计
```
GET /courses/{id}/stats
```

## 练习相关

### 获取练习列表
```
GET /exercises?course_id={courseId}&chapter_id={chapterId}
```

### 获取练习详情
```
GET /exercises/{id}
```

### 创建练习 (教师)
```
POST /exercises
```

请求体:
```json
{
  "title": "string",
  "description": "string",
  "course_id": 1,
  "chapter_id": 1,
  "type": "string",
  "duration": 60,
  "total_score": 100
}
```

### 生成练习题 (教师)
```
POST /exercises/{courseId}/chapters/{chapterId}/generate?type={type}&count={count}
```

### 开始练习
```
POST /exercises/{id}/start
```

### 提交答案
```
POST /exercises/{recordId}/answers
```

请求体:
```json
{
  "question_id": 1,
  "answer": "string"
}
```

### 完成练习
```
POST /exercises/{recordId}/complete
```

### 获取练习统计
```
GET /exercises/stats
```

## 聊天相关

### 获取聊天会话列表
```
GET /chat/sessions
```

### 创建聊天会话
```
POST /chat/sessions
```

请求体:
```json
{
  "title": "string",
  "type": "string",
  "course_id": 1,
  "chapter_id": 1
}
```

### 获取聊天会话详情
```
GET /chat/sessions/{id}
```

### 发送消息
```
POST /chat/sessions/{sessionId}/messages
```

请求体:
```json
{
  "content": "string"
}
```

### 删除聊天会话
```
DELETE /chat/sessions/{id}
```

### 获取学习建议
```
GET /chat/advice
```

## 错误响应

```json
{
  "code": 400,
  "message": "错误信息"
}
```

常见状态码:
- 200: 成功
- 400: 请求参数错误
- 401: 未认证
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误 