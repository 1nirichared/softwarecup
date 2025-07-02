package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// CORS配置
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 静态文件服务 - 提供上传的文件访问
	r.Static("/uploads", "./uploads")

	// API路由组
	api := r.Group("/api/v1")

	// 公开路由
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
		}

		// SSE流式AI回复 - 需要从query参数获取token，所以放在公开路由
		api.GET("/chat/stream", handlers.StreamAIChat)
	}

	// 需要认证的路由
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		user := authenticated.Group("/user")
		{
			user.GET("/profile", handlers.GetCurrentUser)
			user.PUT("/password", handlers.ChangePassword)
		}

		// 课程相关
		courses := authenticated.Group("/courses")
		{
			courses.GET("", handlers.GetCourses)
			courses.GET("/:id", handlers.GetCourses)
			courses.GET("/:id/stats", handlers.GetCourseStats)

			// 教师专用
			courses.POST("", middleware.RoleMiddleware("teacher"), handlers.CreateCourse)
			courses.PUT("/:id", middleware.RoleMiddleware("teacher"), handlers.UpdateCourse)
			courses.DELETE("/:id", middleware.RoleMiddleware("teacher"), handlers.DeleteCourse)
			courses.POST("/:courseId/chapters/:chapterId/lesson-plan", middleware.RoleMiddleware("teacher"), handlers.GenerateLessonPlan)
		}

		// 课程材料相关 - 使用不同的路径结构避免冲突
		authenticated.POST("/course-materials/:courseId", middleware.RoleMiddleware("teacher"), handlers.UploadCourseMaterial)
		authenticated.GET("/course-materials/:courseId", handlers.GetCourseMaterials)
		authenticated.DELETE("/course-materials/:materialId", middleware.RoleMiddleware("teacher"), handlers.DeleteCourseMaterial)

		// 练习相关
		exercises := authenticated.Group("/exercises")
		{
			exercises.GET("", handlers.GetExercises)
			exercises.GET("/:id", handlers.GetExercises)
			exercises.GET("/stats", handlers.GetExerciseStats)

			// 教师专用
			exercises.POST("", middleware.RoleMiddleware("teacher"), handlers.CreateExercise)
			exercises.POST("/:courseId/chapters/:chapterId/generate", middleware.RoleMiddleware("teacher"), handlers.GenerateExercises)
		}

		// 练习记录相关（学生答题）
		exerciseRecords := authenticated.Group("/exercise-records")
		{
			exerciseRecords.POST("/start/:exerciseId", handlers.StartExercise)
			exerciseRecords.POST("/:recordId/answers", handlers.SubmitAnswer)
			exerciseRecords.POST("/:recordId/complete", handlers.CompleteExercise)
		}

		// 聊天相关
		chat := authenticated.Group("/chat")
		{
			chat.GET("/sessions", handlers.GetChatSessions)
			chat.POST("/sessions", handlers.CreateChatSession)
			chat.GET("/sessions/:id", handlers.GetChatSession)
			chat.POST("/sessions/:sessionId/messages", handlers.SendMessage)
			chat.DELETE("/sessions/:id", handlers.DeleteChatSession)
			chat.GET("/advice", handlers.GetLearningAdvice)
		}

		// 管理相关
		admin := authenticated.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// 这里可以添加管理员专用接口
		}
	}

	return r
}
