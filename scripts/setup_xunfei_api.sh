#!/bin/bash

# 讯飞星火 API 配置脚本

set -e

echo "🚀 讯飞星火 API 配置向导"
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

# 使用提供的凭据
APP_ID="04107cca"
API_SECRET="NmYyYjc2OTk1Yjc4ZGMwZDhkYWM5YjBj"
API_KEY="1a6c4989345073f44140f47aa57e5ae8"

echo ""
echo "📋 使用提供的讯飞星火凭据："
echo "  - AppID: $APP_ID"
echo "  - APISecret: ${API_SECRET:0:10}..."
echo "  - APIKey: ${API_KEY:0:10}..."

# 设置环境变量
export XUNFEI_APP_ID="$APP_ID"
export XUNFEI_API_SECRET="$API_SECRET"
export XUNFEI_API_KEY="$API_KEY"

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
  # 主要AI提供商
  provider: "xunfei"  # xunfei, deepseek, openai
  
  # DeepSeek API配置
  deepseek_api_key: "your-deepseek-api-key"
  deepseek_base_url: "https://api.deepseek.com/v1"
  deepseek_model: "deepseek-coder"
  
  # OpenAI配置（备用）
  openai_api_key: "your-openai-api-key"
  openai_base_url: "https://api.openai.com/v1"
  openai_model: "gpt-3.5-turbo"
  
  # 通用配置
  max_tokens: 2000
  temperature: 0.7
  timeout: 60

# 讯飞星火配置
xunfei:
  app_id: "$APP_ID"
  api_secret: "$API_SECRET"
  api_key: "$API_KEY"
  base_url: "https://spark-api.xf-yun.com/v3.1/chat"
  model: "spark-v3.1"
  max_tokens: 2000
  timeout: 60
EOF

echo "✅ 配置文件已更新"

# 创建 .env 文件
cat > .env << EOF
# 讯飞星火 API
XUNFEI_APP_ID=$APP_ID
XUNFEI_API_SECRET=$API_SECRET
XUNFEI_API_KEY=$API_KEY
XUNFEI_BASE_URL=https://spark-api.xf-yun.com/v3.1/chat
XUNFEI_MODEL=spark-v3.1

# DeepSeek API (备用)
DEEPSEEK_API_KEY=your-deepseek-api-key
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
echo "🧪 测试讯飞星火 API..."

if python3 scripts/test_xunfei_api.py "$APP_ID" "$API_SECRET" "$API_KEY"; then
    echo ""
    echo "🎉 讯飞星火 API 配置成功！"
    echo ""
    echo "📋 配置信息："
    echo "  - AppID: $APP_ID"
    echo "  - APISecret: ${API_SECRET:0:10}..."
    echo "  - APIKey: ${API_KEY:0:10}..."
    echo "  - Base URL: https://spark-api.xf-yun.com/v3.1/chat"
    echo "  - 模型: spark-v3.1"
    echo ""
    echo "📖 下一步："
    echo "1. 启动后端服务: cd backend && go run main.go"
    echo "2. 启动前端服务: cd frontend && npm run dev"
    echo "3. 访问应用: http://localhost:3000"
    echo ""
    echo "📚 更多信息请查看: docs/XUNFEI_API_SETUP.md"
else
    echo ""
    echo "❌ API 测试失败，请检查："
    echo "1. 网络连接是否正常"
    echo "2. 账户余额是否充足"
    echo "3. API 凭据是否正确"
    echo ""
    echo "🔧 故障排除："
    echo "- 查看详细文档: docs/XUNFEI_API_SETUP.md"
    echo "- 检查 API 状态: https://www.xfyun.cn/service/spark"
    exit 1
fi 