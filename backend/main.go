package main

import (
	"backend/config"
	"backend/database"
	"backend/redis"
	"backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化Redis
	if err := redis.InitRedis(); err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// 设置Gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	log.Printf("Server starting on port %s", config.GlobalConfig.Server.Port)
	if err := r.Run(config.GlobalConfig.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
