@echo off
echo 启动教学实训智能体软件...

REM 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到Go，请先安装Go 1.19+
    pause
    exit /b 1
)

REM 检查Node.js是否安装
node --version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到Node.js，请先安装Node.js 16+
    pause
    exit /b 1
)

REM 启动后端
echo 启动后端服务...
cd backend
go mod tidy
start "Backend" go run main.go

REM 等待后端启动
timeout /t 3 /nobreak >nul

REM 启动前端
echo 启动前端服务...
cd ..\frontend
npm install
start "Frontend" npm run dev

echo 服务启动完成！
echo 前端地址: http://localhost:3000
echo 后端地址: http://localhost:8080
echo.
echo 默认用户:
echo 管理员: admin / password
echo 教师: teacher / password
echo 学生: student / password
echo.
echo 按任意键退出...
pause >nul 