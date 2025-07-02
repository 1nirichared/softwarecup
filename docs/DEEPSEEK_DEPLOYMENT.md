# DeepSeek AI 本地部署指南

## 概述

DeepSeek 是一个优秀的开源大语言模型，特别擅长代码生成和编程任务。本指南将帮助您在本地部署 DeepSeek AI 模型，并集成到教学培训平台中。

## 系统要求

### 硬件要求

#### 最低配置
- **GPU**: NVIDIA RTX 3090 (24GB VRAM) 或更高
- **内存**: 32GB RAM
- **存储**: 100GB 可用空间
- **CPU**: 8核心以上

#### 推荐配置
- **GPU**: NVIDIA RTX 4090 (24GB VRAM) 或 A100 (40GB VRAM)
- **内存**: 64GB RAM
- **存储**: 200GB SSD
- **CPU**: 16核心以上

### 软件要求

- **操作系统**: Windows 10/11, Linux (Ubuntu 20.04+), macOS
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **NVIDIA 驱动**: 470+
- **CUDA**: 11.8 或 12.1 (可选，用于GPU加速)

## 部署方法

### 方法一：使用 Docker (推荐)

#### 1. 快速部署

```bash
# Linux/macOS
chmod +x scripts/deploy_deepseek.sh
./scripts/deploy_deepseek.sh

# Windows
scripts\deploy_deepseek.bat
```

#### 2. 手动部署

```bash
# 创建部署目录
mkdir deepseek_deploy
cd deepseek_deploy

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

# 启动服务
docker-compose up -d

# 下载模型
docker exec ollama ollama pull deepseek-coder:6.7b
```

### 方法二：直接安装 Ollama

#### Linux/macOS

```bash
# 安装 Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# 启动服务
ollama serve

# 下载模型
ollama pull deepseek-coder:6.7b
```

#### Windows

```powershell
# 下载并安装 Ollama
# 访问 https://ollama.ai/download 下载 Windows 版本

# 启动服务
ollama serve

# 下载模型
ollama pull deepseek-coder:6.7b
```

### 方法三：使用 vLLM

```bash
# 安装 vLLM
pip install vllm

# 启动服务
python -m vllm.entrypoints.openai.api_server \
    --model deepseek-ai/deepseek-coder-6.7b-base \
    --host 0.0.0.0 \
    --port 8000
```

## 配置教学平台

### 1. 更新配置文件

编辑 `backend/config/config.yaml`：

```yaml
# 本地AI模型配置
local_ai:
  enabled: true
  type: "ollama"  # ollama, vllm, transformers
  base_url: "http://localhost:11434"
  model: "deepseek-coder:6.7b"
  max_tokens: 2000
  temperature: 0.7
  timeout: 60
```

### 2. 测试连接

```bash
# 测试 API 连接
curl http://localhost:11434/api/tags

# 测试文本生成
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder:6.7b",
    "prompt": "Write a Python function to calculate fibonacci numbers",
    "stream": false
  }'
```

## 可用模型

### DeepSeek 模型列表

| 模型名称 | 大小 | 内存需求 | 特点 |
|---------|------|----------|------|
| `deepseek-coder:6.7b` | 6.7B | 16GB | 代码生成，平衡性能 |
| `deepseek-coder:33b` | 33B | 64GB | 更高质量，需要更多资源 |
| `deepseek-coder:1.3b` | 1.3B | 4GB | 轻量级，快速响应 |
| `deepseek-llm:7b` | 7B | 16GB | 通用对话模型 |
| `deepseek-llm:67b` | 67B | 128GB | 最高质量，需要大量资源 |

### 下载模型

```bash
# 使用 Ollama
ollama pull deepseek-coder:6.7b
ollama pull deepseek-coder:33b

# 查看已安装模型
ollama list
```

## API 使用

### 文本生成

```bash
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder:6.7b",
    "prompt": "Write a Python function to calculate fibonacci numbers",
    "stream": false,
    "options": {
      "num_predict": 2000,
      "temperature": 0.7,
      "top_p": 0.9,
      "top_k": 40
    }
  }'
```

### 聊天对话

```bash
curl -X POST http://localhost:11434/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder:6.7b",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"},
      {"role": "assistant", "content": "I am doing well, thank you!"},
      {"role": "user", "content": "Can you help me with Python programming?"}
    ],
    "stream": false,
    "options": {
      "num_predict": 2000,
      "temperature": 0.7
    }
  }'
```

### 流式输出

```bash
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-coder:6.7b",
    "prompt": "Write a Python function to calculate fibonacci numbers",
    "stream": true
  }'
```

## 性能优化

### GPU 优化

1. **启用 GPU 加速**
   ```bash
   # 检查 GPU 状态
   nvidia-smi
   
   # 设置环境变量
   export CUDA_VISIBLE_DEVICES=0
   ```

2. **内存优化**
   ```bash
   # 使用量化模型
   ollama pull deepseek-coder:6.7b-q4_0
   
   # 设置内存限制
   docker run --gpus all --memory=32g ollama/ollama
   ```

3. **批处理优化**
   ```bash
   # 设置批处理大小
   export OLLAMA_NUM_PARALLEL=4
   ```

### 网络优化

1. **本地网络**
   ```bash
   # 使用本地地址
   export OLLAMA_HOST=127.0.0.1
   ```

2. **代理设置**
   ```bash
   # 设置代理
   export HTTP_PROXY=http://proxy:port
   export HTTPS_PROXY=http://proxy:port
   ```

## 监控和维护

### 健康检查

```bash
# 检查服务状态
curl http://localhost:11434/api/tags

# 检查模型状态
ollama list

# 查看日志
docker-compose logs -f ollama
```

### 性能监控

```bash
# 监控 GPU 使用
watch -n 1 nvidia-smi

# 监控内存使用
htop

# 监控网络
iftop
```

### 备份和恢复

```bash
# 备份模型
docker exec ollama ollama cp deepseek-coder:6.7b backup-model

# 恢复模型
docker exec ollama ollama cp backup-model deepseek-coder:6.7b
```

## 故障排除

### 常见问题

1. **GPU 内存不足**
   ```bash
   # 解决方案：使用更小的模型
   ollama pull deepseek-coder:1.3b
   
   # 或使用量化模型
   ollama pull deepseek-coder:6.7b-q4_0
   ```

2. **服务启动失败**
   ```bash
   # 检查端口占用
   netstat -tulpn | grep 11434
   
   # 重启服务
   docker-compose restart
   ```

3. **模型下载失败**
   ```bash
   # 检查网络连接
   ping ollama.ai
   
   # 使用镜像源
   export OLLAMA_ORIGINS=https://mirror.ghproxy.com/https://github.com/ollama/ollama/releases
   ```

4. **响应速度慢**
   ```bash
   # 优化参数
   export OLLAMA_NUM_PARALLEL=1
   export OLLAMA_GPU_LAYERS=35
   ```

### 日志分析

```bash
# 查看详细日志
docker-compose logs -f ollama

# 查看错误日志
docker-compose logs ollama | grep ERROR

# 查看性能日志
docker-compose logs ollama | grep "tokens per second"
```

## 安全考虑

### 网络安全

1. **防火墙设置**
   ```bash
   # 只允许本地访问
   ufw allow from 127.0.0.1 to any port 11434
   ```

2. **API 认证**
   ```bash
   # 设置 API 密钥
   export OLLAMA_API_KEY=your-secret-key
   ```

3. **HTTPS 配置**
   ```bash
   # 使用反向代理
   nginx -s reload
   ```

### 数据安全

1. **模型隔离**
   ```bash
   # 使用独立容器
   docker run --name ollama-isolated ollama/ollama
   ```

2. **数据加密**
   ```bash
   # 加密存储
   docker run -v encrypted-data:/root/.ollama ollama/ollama
   ```

## 扩展功能

### 多模型支持

```bash
# 安装多个模型
ollama pull deepseek-coder:6.7b
ollama pull deepseek-llm:7b
ollama pull llama2:7b

# 切换模型
curl -X POST http://localhost:11434/api/generate \
  -d '{"model": "deepseek-llm:7b", "prompt": "Hello"}'
```

### 模型微调

```bash
# 创建自定义模型
ollama create my-model -f Modelfile

# 微调模型
ollama run my-model "Your training data"
```

### 集成其他服务

```bash
# 集成 LangChain
pip install langchain-ollama

# 集成 FastAPI
pip install fastapi uvicorn
```

## 更新和维护

### 模型更新

```bash
# 更新模型
ollama pull deepseek-coder:6.7b

# 删除旧模型
ollama rm deepseek-coder:6.7b
```

### 服务更新

```bash
# 更新 Ollama
docker-compose pull
docker-compose up -d

# 更新 Web UI
docker-compose pull ollama-webui
docker-compose up -d ollama-webui
```

### 系统维护

```bash
# 清理缓存
docker system prune

# 更新系统
sudo apt update && sudo apt upgrade

# 重启服务
docker-compose restart
```

## 总结

通过本指南，您已经成功在本地部署了 DeepSeek AI 模型，并将其集成到教学培训平台中。本地部署的优势包括：

- **数据安全**: 数据不离开本地环境
- **成本控制**: 无需支付 API 调用费用
- **性能优化**: 可以根据硬件进行优化
- **离线使用**: 不依赖网络连接
- **自定义**: 可以微调和定制模型

建议根据实际需求和硬件配置选择合适的模型大小，并定期进行性能监控和维护。 