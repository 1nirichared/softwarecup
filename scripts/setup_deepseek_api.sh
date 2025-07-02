#!/bin/bash

# DeepSeek API 快速配置脚本

set -e

echo "🚀 DeepSeek API 配置向导"
echo "=" * 50

# 检查是否已安装 Python
if ! command -v python3 &> /dev/null; then
    echo "❌ Python3 未安装，请先安装 Python3"
    exit 1
fi

# 检查是否已安装 requests
if ! python3 -c "import requests" &> /dev/null; then
    echo "📦 安装 requests 库..."
    pip3 install requests
fi

# 获取 API Key
echo ""
echo "📋 配置 DeepSeek API"
echo "请访问 https://platform.deepseek.com/api_keys 获取 API Key"
echo ""

read -p "请输入您的 DeepSeek API Key: " api_key

if [ -z "$api_key" ]; then
    echo "❌ API Key 不能为空"
    exit 1
fi

# 设置环境变量
export DEEPSEEK_API_KEY="$api_key"

echo ""
echo "🔧 更新配置文件..."

# 备份原配置文件
if [ -f "backend/config/config.yaml" ]; then
    cp backend/config/config.yaml backend/config/config.yaml.backup
    echo "✅ 已备份原配置文件"
fi

# 更新配置文件
cat > backend/config/config.yaml << EOF
server:
  port: ":8080"
  mode: "debug"

database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "password"
  dbname: "teaching_platform"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-here"
  expire: 24

# AI服务配置
ai:
  # DeepSeek API配置
  provider: "deepseek"
  api_key: "$api_key"
  base_url: "https://api.deepseek.com/v1"
  model: "deepseek-coder"
  max_tokens: 2000
  temperature: 0.7
  timeout: 60
  
  # OpenAI配置（备用）
  openai_api_key: "your-openai-api-key"
  openai_base_url: "https://api.openai.com/v1"
  openai_model: "gpt-3.5-turbo"
EOF

echo "✅ 配置文件已更新"

# 创建 .env 文件
cat > .env << EOF
# DeepSeek API
DEEPSEEK_API_KEY=$api_key
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-coder

# OpenAI API (备用)
OPENAI_API_KEY=your-openai-api-key
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-3.5-turbo
EOF

echo "✅ 环境变量文件已创建"

# 测试 API
echo ""
echo "🧪 测试 DeepSeek API..."

if python3 scripts/test_deepseek_api.py "$api_key"; then
    echo ""
    echo "🎉 DeepSeek API 配置成功！"
    echo ""
    echo "📋 配置信息："
    echo "  - API Key: ${api_key:0:10}..."
    echo "  - Base URL: https://api.deepseek.com/v1"
    echo "  - 模型: deepseek-coder"
    echo ""
    echo "📖 下一步："
    echo "1. 启动后端服务: cd backend && go run main.go"
    echo "2. 启动前端服务: cd frontend && npm run dev"
    echo "3. 访问应用: http://localhost:3000"
    echo ""
    echo "📚 更多信息请查看: docs/DEEPSEEK_API_SETUP.md"
else
    echo ""
    echo "❌ API 测试失败，请检查："
    echo "1. API Key 是否正确"
    echo "2. 网络连接是否正常"
    echo "3. 账户余额是否充足"
    echo ""
    echo "🔧 故障排除："
    echo "- 查看详细文档: docs/DEEPSEEK_API_SETUP.md"
    echo "- 检查 API 状态: https://status.deepseek.com"
    exit 1
fi 