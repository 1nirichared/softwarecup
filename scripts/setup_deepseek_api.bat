@echo off
setlocal enabledelayedexpansion

echo 🚀 DeepSeek API 配置向导
echo ==================================================

REM 检查是否已安装 Python
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Python 未安装，请先安装 Python
    pause
    exit /b 1
)

REM 检查是否已安装 requests
python -c "import requests" >nul 2>&1
if %errorlevel% neq 0 (
    echo 📦 安装 requests 库...
    pip install requests
)

REM 获取 API Key
echo.
echo 📋 配置 DeepSeek API
echo 请访问 https://platform.deepseek.com/api_keys 获取 API Key
echo.

set /p api_key="请输入您的 DeepSeek API Key: "

if "%api_key%"=="" (
    echo ❌ API Key 不能为空
    pause
    exit /b 1
)

REM 设置环境变量
set DEEPSEEK_API_KEY=%api_key%

echo.
echo 🔧 更新配置文件...

REM 备份原配置文件
if exist "backend\config\config.yaml" (
    copy "backend\config\config.yaml" "backend\config\config.yaml.backup"
    echo ✅ 已备份原配置文件
)

REM 更新配置文件
echo server: > backend\config\config.yaml
echo   port: ":8080" >> backend\config\config.yaml
echo   mode: "debug" >> backend\config\config.yaml
echo. >> backend\config\config.yaml
echo database: >> backend\config\config.yaml
echo   host: "localhost" >> backend\config\config.yaml
echo   port: "3306" >> backend\config\config.yaml
echo   username: "root" >> backend\config\config.yaml
echo   password: "password" >> backend\config\config.yaml
echo   dbname: "teaching_platform" >> backend\config\config.yaml
echo. >> backend\config\config.yaml
echo redis: >> backend\config\config.yaml
echo   host: "localhost" >> backend\config\config.yaml
echo   port: "6379" >> backend\config\config.yaml
echo   password: "" >> backend\config\config.yaml
echo   db: 0 >> backend\config\config.yaml
echo. >> backend\config\config.yaml
echo jwt: >> backend\config\config.yaml
echo   secret: "your-secret-key-here" >> backend\config\config.yaml
echo   expire: 24 >> backend\config\config.yaml
echo. >> backend\config\config.yaml
echo # AI服务配置 >> backend\config\config.yaml
echo ai: >> backend\config\config.yaml
echo   # DeepSeek API配置 >> backend\config\config.yaml
echo   provider: "deepseek" >> backend\config\config.yaml
echo   api_key: "%api_key%" >> backend\config\config.yaml
echo   base_url: "https://api.deepseek.com/v1" >> backend\config\config.yaml
echo   model: "deepseek-coder" >> backend\config\config.yaml
echo   max_tokens: 2000 >> backend\config\config.yaml
echo   temperature: 0.7 >> backend\config\config.yaml
echo   timeout: 60 >> backend\config\config.yaml
echo. >> backend\config\config.yaml
echo   # OpenAI配置（备用） >> backend\config\config.yaml
echo   openai_api_key: "your-openai-api-key" >> backend\config\config.yaml
echo   openai_base_url: "https://api.openai.com/v1" >> backend\config\config.yaml
echo   openai_model: "gpt-3.5-turbo" >> backend\config\config.yaml

echo ✅ 配置文件已更新

REM 创建 .env 文件
echo # DeepSeek API > .env
echo DEEPSEEK_API_KEY=%api_key% >> .env
echo DEEPSEEK_BASE_URL=https://api.deepseek.com/v1 >> .env
echo DEEPSEEK_MODEL=deepseek-coder >> .env
echo. >> .env
echo # OpenAI API (备用) >> .env
echo OPENAI_API_KEY=your-openai-api-key >> .env
echo OPENAI_BASE_URL=https://api.openai.com/v1 >> .env
echo OPENAI_MODEL=gpt-3.5-turbo >> .env

echo ✅ 环境变量文件已创建

REM 测试 API
echo.
echo 🧪 测试 DeepSeek API...

python scripts\test_deepseek_api.py "%api_key%"
if %errorlevel% equ 0 (
    echo.
    echo 🎉 DeepSeek API 配置成功！
    echo.
    echo 📋 配置信息：
    echo   - API Key: %api_key:~0,10%...
    echo   - Base URL: https://api.deepseek.com/v1
    echo   - 模型: deepseek-coder
    echo.
    echo 📖 下一步：
    echo 1. 启动后端服务: cd backend ^&^& go run main.go
    echo 2. 启动前端服务: cd frontend ^&^& npm run dev
    echo 3. 访问应用: http://localhost:3000
    echo.
    echo 📚 更多信息请查看: docs\DEEPSEEK_API_SETUP.md
) else (
    echo.
    echo ❌ API 测试失败，请检查：
    echo 1. API Key 是否正确
    echo 2. 网络连接是否正常
    echo 3. 账户余额是否充足
    echo.
    echo 🔧 故障排除：
    echo - 查看详细文档: docs\DEEPSEEK_API_SETUP.md
    echo - 检查 API 状态: https://status.deepseek.com
)

pause 