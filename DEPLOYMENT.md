# 部署说明

## 环境要求

### 系统要求
- 操作系统: Windows 10+, macOS 10.15+, Ubuntu 18.04+
- 内存: 最少 4GB RAM
- 存储: 最少 10GB 可用空间

### 软件要求
- Go 1.19+
- Node.js 16+
- MySQL 8.0+
- Redis 6.0+

## 安装步骤

### 1. 克隆项目
```bash
git clone <repository-url>
cd softwarecup
```

### 2. 安装依赖

#### 后端依赖
```bash
cd backend
go mod tidy
```

#### 前端依赖
```bash
cd frontend
npm install
```

### 3. 配置数据库

#### 创建数据库
```sql
CREATE DATABASE teaching_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 运行初始化脚本
```bash
mysql -u root -p teaching_platform < scripts/init.sql
```

### 4. 配置Redis
确保Redis服务正在运行：
```bash
# Ubuntu/Debian
sudo systemctl start redis-server

# macOS
brew services start redis

# Windows
redis-server
```

### 5. 配置文件

#### 后端配置
编辑 `backend/config/config.yaml`：
```yaml
server:
  port: ":8080"
  mode: "debug"

database:
  host: "localhost"
  port: "3306"
  username: "your_username"
  password: "your_password"
  dbname: "teaching_platform"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-here"
  expire: 24

ai:
  api_key: "your-openai-api-key"
  base_url: "https://api.openai.com/v1"
  model: "gpt-3.5-turbo"
  max_tokens: 2000
  temperature: 0.7
```

### 6. 启动服务

#### 使用脚本启动（推荐）
```bash
# Linux/macOS
chmod +x scripts/start.sh
./scripts/start.sh

# Windows
scripts/start.bat
```

#### 手动启动
```bash
# 启动后端
cd backend
go run main.go

# 启动前端（新终端）
cd frontend
npm run dev
```

## 访问应用

- 前端地址: http://localhost:3000
- 后端API: http://localhost:8080/api/v1

## 默认用户

- 管理员: `admin` / `password`
- 教师: `teacher` / `password`
- 学生: `student` / `password`

## 生产环境部署

### 1. 构建前端
```bash
cd frontend
npm run build
```

### 2. 配置Nginx
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端API代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. 使用PM2管理后端进程
```bash
cd backend
go build -o main
pm2 start main --name "teaching-platform"
```

### 4. 配置SSL证书（可选）
使用Let's Encrypt或其他SSL证书提供商。

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查MySQL服务是否运行
   - 验证数据库配置信息
   - 确认用户权限

2. **Redis连接失败**
   - 检查Redis服务是否运行
   - 验证Redis配置信息

3. **前端无法访问后端API**
   - 检查后端服务是否启动
   - 验证代理配置
   - 检查CORS设置

4. **AI功能无法使用**
   - 检查OpenAI API密钥配置
   - 验证网络连接
   - 确认API配额

### 日志查看

#### 后端日志
```bash
# 如果使用PM2
pm2 logs teaching-platform

# 直接运行
go run main.go
```

#### 前端日志
浏览器开发者工具控制台

### 性能优化

1. **数据库优化**
   - 添加适当的索引
   - 优化查询语句
   - 配置连接池

2. **缓存优化**
   - 使用Redis缓存热点数据
   - 配置合理的缓存策略

3. **前端优化**
   - 启用Gzip压缩
   - 使用CDN加速静态资源
   - 优化打包配置

## 监控和维护

### 系统监控
- 使用Prometheus + Grafana监控系统性能
- 配置日志收集和分析
- 设置告警机制

### 数据备份
```bash
# 数据库备份
mysqldump -u root -p teaching_platform > backup.sql

# 文件备份
tar -czf backup.tar.gz uploads/
```

### 更新部署
```bash
# 拉取最新代码
git pull

# 重新构建
cd backend && go build
cd frontend && npm run build

# 重启服务
pm2 restart teaching-platform
``` 