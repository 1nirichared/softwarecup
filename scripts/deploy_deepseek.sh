#!/bin/bash

# DeepSeek AI 本地部署脚本
# 适用于 Linux/macOS

set -e

echo "🚀 开始部署 DeepSeek AI 本地服务..."

# 检查系统要求
echo "📋 检查系统要求..."

# 检查GPU
if command -v nvidia-smi &> /dev/null; then
    echo "✅ 检测到 NVIDIA GPU"
    nvidia-smi --query-gpu=name,memory.total --format=csv,noheader,nounits
else
    echo "⚠️  未检测到 NVIDIA GPU，将使用CPU模式（性能较慢）"
fi

# 检查Docker
if command -v docker &> /dev/null; then
    echo "✅ Docker 已安装"
else
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查Docker Compose
if command -v docker-compose &> /dev/null; then
    echo "✅ Docker Compose 已安装"
else
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 创建部署目录
DEPLOY_DIR="./deepseek_deploy"
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR

echo "📁 创建部署目录: $DEPLOY_DIR"

# 创建 docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  ollama:
    image: ollama/ollama:latest
    container_name: ollama
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama
    environment:
      - OLLAMA_HOST=0.0.0.0
    restart: unless-stopped
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]

  # 可选：添加 Web UI
  ollama-webui:
    image: ghcr.io/ollama-webui/ollama-webui:main
    container_name: ollama-webui
    ports:
      - "3000:8080"
    environment:
      - OLLAMA_API_BASE_URL=http://ollama:11434/api
    depends_on:
      - ollama
    restart: unless-stopped

volumes:
  ollama_data:
EOF

echo "📄 创建 docker-compose.yml"

# 创建启动脚本
cat > start.sh << 'EOF'
#!/bin/bash

echo "🚀 启动 DeepSeek AI 服务..."

# 启动服务
docker-compose up -d

echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
if curl -s http://localhost:11434/api/tags > /dev/null; then
    echo "✅ Ollama 服务启动成功"
else
    echo "❌ Ollama 服务启动失败"
    exit 1
fi

# 下载 DeepSeek 模型
echo "📥 下载 DeepSeek 模型..."
docker exec ollama ollama pull deepseek-coder:6.7b

echo "🎉 DeepSeek AI 部署完成！"
echo ""
echo "📊 服务信息："
echo "  - Ollama API: http://localhost:11434"
echo "  - Web UI: http://localhost:3000"
echo "  - 模型: deepseek-coder:6.7b"
echo ""
echo "🔧 常用命令："
echo "  - 查看日志: docker-compose logs -f"
echo "  - 停止服务: docker-compose down"
echo "  - 重启服务: docker-compose restart"
echo "  - 更新模型: docker exec ollama ollama pull deepseek-coder:6.7b"
EOF

chmod +x start.sh

# 创建停止脚本
cat > stop.sh << 'EOF'
#!/bin/bash

echo "🛑 停止 DeepSeek AI 服务..."
docker-compose down
echo "✅ 服务已停止"
EOF

chmod +x stop.sh

# 创建测试脚本
cat > test.sh << 'EOF'
#!/bin/bash

echo "🧪 测试 DeepSeek AI 服务..."

# 测试 API 连接
if curl -s http://localhost:11434/api/tags > /dev/null; then
    echo "✅ API 连接正常"
else
    echo "❌ API 连接失败"
    exit 1
fi

# 测试模型生成
echo "📝 测试文本生成..."
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder:6.7b",
    "prompt": "Hello, how are you?",
    "stream": false
  }' | jq -r '.response'

echo ""
echo "✅ 测试完成"
EOF

chmod +x test.sh

# 创建 README
cat > README.md << 'EOF'
# DeepSeek AI 本地部署

## 快速开始

1. 启动服务：
   ```bash
   ./start.sh
   ```

2. 测试服务：
   ```bash
   ./test.sh
   ```

3. 停止服务：
   ```bash
   ./stop.sh
   ```

## 服务地址

- **Ollama API**: http://localhost:11434
- **Web UI**: http://localhost:3000

## 可用模型

- `deepseek-coder:6.7b` - 代码生成模型
- `deepseek-coder:33b` - 更大模型（需要更多显存）

## API 使用示例

### 文本生成
```bash
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder:6.7b",
    "prompt": "Write a Python function to calculate fibonacci numbers",
    "stream": false
  }'
```

### 聊天对话
```bash
curl -X POST http://localhost:11434/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder:6.7b",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ],
    "stream": false
  }'
```

## 故障排除

1. **GPU 内存不足**：
   - 使用更小的模型：`deepseek-coder:6.7b`
   - 增加系统内存
   - 使用 CPU 模式（较慢）

2. **服务启动失败**：
   - 检查 Docker 是否运行
   - 检查端口是否被占用
   - 查看日志：`docker-compose logs -f`

3. **模型下载失败**：
   - 检查网络连接
   - 手动下载：`docker exec ollama ollama pull deepseek-coder:6.7b`

## 性能优化

1. **GPU 优化**：
   - 使用 NVIDIA GPU
   - 安装最新驱动
   - 配置 CUDA 环境

2. **内存优化**：
   - 关闭不必要的服务
   - 增加系统内存
   - 使用更小的模型

3. **网络优化**：
   - 使用本地网络
   - 配置代理（如需要）
EOF

echo "📚 创建 README.md"

echo ""
echo "🎉 部署文件创建完成！"
echo ""
echo "📋 下一步操作："
echo "1. 进入部署目录: cd $DEPLOY_DIR"
echo "2. 启动服务: ./start.sh"
echo "3. 测试服务: ./test.sh"
echo ""
echo "�� 详细说明请查看 README.md" 