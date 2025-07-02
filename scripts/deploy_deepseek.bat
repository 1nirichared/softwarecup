@echo off
setlocal enabledelayedexpansion

echo 🚀 开始部署 DeepSeek AI 本地服务...

REM 检查系统要求
echo 📋 检查系统要求...

REM 检查Docker
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker 未安装，请先安装 Docker Desktop
    pause
    exit /b 1
)
echo ✅ Docker 已安装

REM 检查Docker Compose
docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker Compose 未安装，请先安装 Docker Compose
    pause
    exit /b 1
)
echo ✅ Docker Compose 已安装

REM 创建部署目录
set DEPLOY_DIR=deepseek_deploy
if not exist %DEPLOY_DIR% mkdir %DEPLOY_DIR%
cd %DEPLOY_DIR%

echo 📁 创建部署目录: %DEPLOY_DIR%

REM 创建 docker-compose.yml
echo version: '3.8' > docker-compose.yml
echo. >> docker-compose.yml
echo services: >> docker-compose.yml
echo   ollama: >> docker-compose.yml
echo     image: ollama/ollama:latest >> docker-compose.yml
echo     container_name: ollama >> docker-compose.yml
echo     ports: >> docker-compose.yml
echo       - "11434:11434" >> docker-compose.yml
echo     volumes: >> docker-compose.yml
echo       - ollama_data:/root/.ollama >> docker-compose.yml
echo     environment: >> docker-compose.yml
echo       - OLLAMA_HOST=0.0.0.0 >> docker-compose.yml
echo     restart: unless-stopped >> docker-compose.yml
echo. >> docker-compose.yml
echo   ollama-webui: >> docker-compose.yml
echo     image: ghcr.io/ollama-webui/ollama-webui:main >> docker-compose.yml
echo     container_name: ollama-webui >> docker-compose.yml
echo     ports: >> docker-compose.yml
echo       - "3000:8080" >> docker-compose.yml
echo     environment: >> docker-compose.yml
echo       - OLLAMA_API_BASE_URL=http://ollama:11434/api >> docker-compose.yml
echo     depends_on: >> docker-compose.yml
echo       - ollama >> docker-compose.yml
echo     restart: unless-stopped >> docker-compose.yml
echo. >> docker-compose.yml
echo volumes: >> docker-compose.yml
echo   ollama_data: >> docker-compose.yml

echo 📄 创建 docker-compose.yml

REM 创建启动脚本
echo @echo off > start.bat
echo echo 🚀 启动 DeepSeek AI 服务... >> start.bat
echo. >> start.bat
echo REM 启动服务 >> start.bat
echo docker-compose up -d >> start.bat
echo. >> start.bat
echo echo ⏳ 等待服务启动... >> start.bat
echo timeout /t 10 /nobreak ^>nul >> start.bat
echo. >> start.bat
echo REM 检查服务状态 >> start.bat
echo curl -s http://localhost:11434/api/tags ^>nul 2^>^&1 >> start.bat
echo if %%errorlevel%% equ 0 ( >> start.bat
echo   echo ✅ Ollama 服务启动成功 >> start.bat
echo ) else ( >> start.bat
echo   echo ❌ Ollama 服务启动失败 >> start.bat
echo   pause >> start.bat
echo   exit /b 1 >> start.bat
echo ) >> start.bat
echo. >> start.bat
echo REM 下载 DeepSeek 模型 >> start.bat
echo echo 📥 下载 DeepSeek 模型... >> start.bat
echo docker exec ollama ollama pull deepseek-coder:6.7b >> start.bat
echo. >> start.bat
echo echo 🎉 DeepSeek AI 部署完成！ >> start.bat
echo echo. >> start.bat
echo echo 📊 服务信息： >> start.bat
echo echo   - Ollama API: http://localhost:11434 >> start.bat
echo echo   - Web UI: http://localhost:3000 >> start.bat
echo echo   - 模型: deepseek-coder:6.7b >> start.bat
echo echo. >> start.bat
echo echo 🔧 常用命令： >> start.bat
echo echo   - 查看日志: docker-compose logs -f >> start.bat
echo echo   - 停止服务: docker-compose down >> start.bat
echo echo   - 重启服务: docker-compose restart >> start.bat
echo echo   - 更新模型: docker exec ollama ollama pull deepseek-coder:6.7b >> start.bat
echo pause >> start.bat

echo 📄 创建 start.bat

REM 创建停止脚本
echo @echo off > stop.bat
echo echo 🛑 停止 DeepSeek AI 服务... >> stop.bat
echo docker-compose down >> stop.bat
echo echo ✅ 服务已停止 >> stop.bat
echo pause >> stop.bat

echo 📄 创建 stop.bat

REM 创建测试脚本
echo @echo off > test.bat
echo echo 🧪 测试 DeepSeek AI 服务... >> test.bat
echo. >> test.bat
echo REM 测试 API 连接 >> test.bat
echo curl -s http://localhost:11434/api/tags ^>nul 2^>^&1 >> test.bat
echo if %%errorlevel%% equ 0 ( >> test.bat
echo   echo ✅ API 连接正常 >> test.bat
echo ) else ( >> test.bat
echo   echo ❌ API 连接失败 >> test.bat
echo   pause >> test.bat
echo   exit /b 1 >> test.bat
echo ) >> test.bat
echo. >> test.bat
echo echo 📝 测试文本生成... >> test.bat
echo curl -X POST http://localhost:11434/api/generate -H "Content-Type: application/json" -d "{\"model\": \"deepseek-coder:6.7b\", \"prompt\": \"Hello, how are you?\", \"stream\": false}" >> test.bat
echo. >> test.bat
echo echo ✅ 测试完成 >> test.bat
echo pause >> test.bat

echo 📄 创建 test.bat

echo.
echo 🎉 部署文件创建完成！
echo.
echo 📋 下一步操作：
echo 1. 进入部署目录: cd %DEPLOY_DIR%
echo 2. 启动服务: start.bat
echo 3. 测试服务: test.bat
echo.
echo 📖 详细说明请查看 README.md
pause 