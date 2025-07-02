# 智能教学实训平台

一个基于 Go + Vue.js 的现代化智能教学实训平台，集成了多种AI服务，支持教师教学管理和学生学习训练。

## 🚀 主要功能

### 教师端功能
- **课程管理**：创建、编辑、管理课程内容和章节
- **练习管理**：创建练习题，支持AI自动生成练习
- **学生管理**：查看学生学习进度和成绩分析
- **资料管理**：上传和管理教学资料
- **AI助手**：智能备课助手和教学辅助工具
- **数据分析**：学生学习数据可视化分析

### 学生端功能
- **课程学习**：浏览和学习课程内容
- **练习训练**：完成各种类型的练习题
- **AI对话**：与AI助手进行学习交流
- **进度跟踪**：查看个人学习进度和成绩
- **资料下载**：下载课程相关学习资料

### AI服务集成
- **讯飞星火X1**：支持流式对话，实时AI回复
- **DeepSeek**：高质量文本生成和问答
- **本地AI**：离线AI服务支持
- **OpenAI**：备用AI服务

## 🏗️ 技术架构

### 后端技术栈
- **语言**：Go 1.21+
- **框架**：Gin (Web框架)
- **ORM**：GORM (数据库操作)
- **数据库**：MySQL 8.0+
- **缓存**：Redis 6.0+
- **认证**：JWT (JSON Web Token)
- **AI集成**：WebSocket + SSE (流式通信)

### 前端技术栈
- **框架**：Vue 3 + Composition API
- **构建工具**：Vite
- **UI组件库**：Element Plus
- **状态管理**：Pinia
- **路由**：Vue Router 4
- **HTTP客户端**：Axios
- **实时通信**：EventSource (SSE)

## 📁 项目结构

```
softwarecup/
├── backend/                 # 后端Go项目
│   ├── config/             # 配置文件
│   ├── database/           # 数据库相关
│   ├── handlers/           # HTTP处理器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由配置
│   ├── services/           # 业务服务层
│   ├── utils/              # 工具函数
│   └── uploads/            # 文件上传目录
├── frontend/               # 前端Vue项目
│   ├── src/
│   │   ├── api/            # API接口
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # 状态管理
│   │   ├── views/          # 页面组件
│   │   └── components/     # 通用组件
│   └── public/             # 静态资源
├── docs/                   # 项目文档
├── scripts/                # 部署脚本
└── README.md              # 项目说明
```

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 后端启动
```bash
cd backend
go mod tidy
go run main.go
```

### 前端启动
```bash
cd frontend
npm install
npm run dev
```

### 数据库初始化
```bash
# 执行数据库初始化脚本
mysql -u root -p < scripts/init.sql
```

## 🔧 配置说明

### 后端配置
复制 `backend/config/config.example.yaml` 为 `backend/config/config.yaml`，并修改相关配置：

```yaml
server:
  port: 3002
  mode: debug

database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: softwarecup

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

ai:
  provider: "xunfei"  # xunfei, deepseek, local, openai
  openai_api_key: "your_openai_key"
  max_tokens: 2048
  temperature: 0.7

xunfei:
  app_id: "your_app_id"
  api_key: "your_api_key"
  api_secret: "your_api_secret"
```

### 前端配置
在 `frontend/src/api/index.js` 中配置API基础URL：

```javascript
const api = axios.create({
  baseURL: 'http://localhost:3002/api/v1',
  timeout: 10000
})
```

## 🔐 AI服务配置

### 讯飞星火X1配置
1. 注册讯飞开放平台账号
2. 创建应用获取 AppID、APIKey、APISecret
3. 在配置文件中填入相关信息
4. 参考 `docs/XUNFEI_API_SETUP.md` 详细配置

### DeepSeek配置
1. 注册DeepSeek平台账号
2. 获取API密钥
3. 参考 `docs/DEEPSEEK_API_SETUP.md` 配置

## 📚 功能特性

### 流式AI对话
- 支持SSE (Server-Sent Events) 流式通信
- 实时AI回复，类似ChatGPT体验
- 支持多种AI服务切换

### 角色权限管理
- 教师端：完整的教学管理功能
- 学生端：专注学习体验
- 管理员：系统管理权限

### 文件管理
- 支持课程资料上传下载
- 文件类型和大小限制
- 安全的文件存储

### 数据可视化
- 学生学习进度图表
- 成绩统计分析
- 实时数据更新

## 🚀 部署

### 开发环境
```bash
# 使用提供的脚本
./scripts/start.sh  # Linux/Mac
./scripts/start.bat # Windows
```

### 生产环境
参考 `DEPLOYMENT.md` 进行生产环境部署配置。

## 📖 文档

- [API文档](docs/API.md) - 后端API接口说明
- [讯飞配置](docs/XUNFEI_API_SETUP.md) - 讯飞AI服务配置
- [DeepSeek配置](docs/DEEPSEEK_API_SETUP.md) - DeepSeek服务配置
- [部署指南](DEPLOYMENT.md) - 生产环境部署
- [项目总结](PROJECT_SUMMARY.md) - 项目功能总结

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进项目。

## 📄 许可证

本项目采用 MIT 许可证。 