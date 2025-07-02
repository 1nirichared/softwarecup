# 教学实训智能体软件

基于开源大模型的教学实训智能体软件，实现教师智能备课、学生在线学习、智能考核等功能。

## 项目架构

- **后端**: Go + Gin + GORM + MySQL + Redis
- **前端**: Vue.js + Element UI
- **AI模型**: 开源大模型 (支持本地部署)
- **知识库**: 本地向量数据库

## 核心功能

### 教师侧
- 智能备课设计
- 考核内容生成
- 学情数据分析

### 学生侧
- 在线学习助手
- 实时练习评测

### 管理侧
- 用户管理
- 课件资源管理
- 大屏数据概览

## 快速开始

### 环境要求
- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- Node.js 16+
- Python 3.7+ (用于API测试)

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd softwarecup
   ```

2. **配置 AI 服务**
   ```bash
   # 讯飞星火 API（推荐）
   # Linux/macOS
   chmod +x scripts/setup_xunfei_api.sh
   ./scripts/setup_xunfei_api.sh
   
   # Windows
   scripts\setup_xunfei_api.bat
   
   # 或配置 DeepSeek API（备用）
   # Linux/macOS
   chmod +x scripts/setup_deepseek_api.sh
   ./scripts/setup_deepseek_api.sh
   
   # Windows
   scripts\setup_deepseek_api.bat
   ```

3. **配置数据库**
   ```bash
   # 创建数据库
   mysql -u root -p < scripts/init.sql
   
   # 或手动创建
   mysql -u root -p
   CREATE DATABASE teaching_platform;
   ```

4. **启动后端**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

5. **启动前端**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

6. **访问应用**
   - 前端: http://localhost:3000
   - 后端API: http://localhost:8080

### AI 服务配置

本项目支持多种 AI 服务：

#### 讯飞星火 API (推荐)
- 优秀的中文理解和生成能力
- 专精教学场景和知识问答
- 成本合理，性能稳定
- 支持多种教学应用场景

#### DeepSeek API (备用)
- 专精代码生成和编程教学
- 成本较低，性能优秀
- 支持中文对话

#### OpenAI API (备用)
- 通用对话能力强
- 稳定性高
- 支持多种模型

详细配置请查看：
- [讯飞星火 API 配置指南](docs/XUNFEI_API_SETUP.md)
- [DeepSeek API 配置指南](docs/DEEPSEEK_API_SETUP.md)
- [本地部署指南](docs/DEEPSEEK_DEPLOYMENT.md)

## 技术栈

- **后端框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **AI模型**: 开源大模型API
- **前端框架**: Vue.js
- **UI组件**: Element UI

## 项目结构

```
softwarecup/
├── backend/          # 后端代码
├── frontend/         # 前端代码
├── docs/            # 文档
├── scripts/         # 脚本文件
└── README.md
``` 