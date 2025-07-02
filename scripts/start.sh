#!/bin/bash

echo "启动教学实训智能体软件..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go，请先安装Go 1.19+"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "错误: 未找到Node.js，请先安装Node.js 16+"
    exit 1
fi

# 检查MySQL是否运行
if ! mysqladmin ping -h localhost -u root -p &> /dev/null; then
    echo "警告: MySQL可能未运行，请确保MySQL服务已启动"
fi

# 检查Redis是否运行
if ! redis-cli ping &> /dev/null; then
    echo "警告: Redis可能未运行，请确保Redis服务已启动"
fi

# 启动后端
echo "启动后端服务..."
cd backend
go mod tidy
go run main.go &
BACKEND_PID=$!

# 等待后端启动
sleep 3

# 启动前端
echo "启动前端服务..."
cd ../frontend
npm install
npm run dev &
FRONTEND_PID=$!

echo "服务启动完成！"
echo "前端地址: http://localhost:3000"
echo "后端地址: http://localhost:8080"
echo ""
echo "默认用户:"
echo "管理员: admin / password"
echo "教师: teacher / password"
echo "学生: student / password"
echo ""
echo "按 Ctrl+C 停止服务"

# 等待用户中断
trap "echo '正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID; exit" INT
wait 