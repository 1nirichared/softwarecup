package handlers

import (
	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/services"
	"backend/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateChatSessionRequest struct {
	Title     string `json:"title"`
	Type      string `json:"type" binding:"required"`
	CourseID  *uint  `json:"course_id"`
	ChapterID *uint  `json:"chapter_id"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// 创建聊天会话
func CreateChatSession(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req CreateChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	session := models.ChatSession{
		UserID:    userID,
		Title:     req.Title,
		Type:      req.Type,
		CourseID:  req.CourseID,
		ChapterID: req.ChapterID,
		Status:    1,
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建会话失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "会话创建成功",
		"data":    session,
	})
}

// 获取聊天会话列表
func GetChatSessions(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var sessions []models.ChatSession
	if err := database.DB.Preload("Course").Preload("Chapter").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取会话列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    sessions,
	})
}

// 获取聊天会话详情
func GetChatSession(c *gin.Context) {
	sessionID := c.Param("id")
	userID := middleware.GetCurrentUserID(c)

	var session models.ChatSession
	if err := database.DB.Preload("Course").Preload("Chapter").Preload("Messages").
		Where("id = ? AND user_id = ?", sessionID, userID).
		First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "会话不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    session,
	})
}

// 发送消息
func SendMessage(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID := middleware.GetCurrentUserID(c)

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 验证会话是否存在且属于当前用户
	var session models.ChatSession
	if err := database.DB.Preload("Course").Preload("Chapter").
		Where("id = ? AND user_id = ?", sessionID, userID).
		First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "会话不存在",
		})
		return
	}

	// 保存用户消息
	userMessage := models.ChatMessage{
		SessionID:   session.ID,
		Role:        "user",
		Content:     req.Content,
		MessageType: "text",
	}

	if err := database.DB.Create(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存消息失败",
		})
		return
	}

	// 获取相关知识库
	var knowledgeBase []models.KnowledgeBase
	if session.CourseID != nil {
		database.DB.Where("type = ? OR type = ?", "course", "general").Find(&knowledgeBase)
	} else {
		database.DB.Where("type = ?", "general").Find(&knowledgeBase)
	}

	// 调用AI服务生成回复
	aiService := services.NewAIService()
	aiResponse, err := aiService.ChatWithAI(&session, req.Content, knowledgeBase)
	if err != nil {
		log.Printf("AI服务调用失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "AI回复生成失败: " + err.Error(),
		})
		return
	}

	// 保存AI回复
	aiMessage := models.ChatMessage{
		SessionID:   session.ID,
		Role:        "assistant",
		Content:     aiResponse,
		MessageType: "text",
	}

	if err := database.DB.Create(&aiMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存AI回复失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "消息发送成功",
		"data": gin.H{
			"user_message": userMessage,
			"ai_message":   aiMessage,
		},
	})
}

// 删除聊天会话
func DeleteChatSession(c *gin.Context) {
	sessionID := c.Param("id")
	userID := middleware.GetCurrentUserID(c)

	// 验证会话是否存在且属于当前用户
	var session models.ChatSession
	if err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "会话不存在",
		})
		return
	}

	// 删除会话及其消息
	if err := database.DB.Delete(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除会话失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "会话删除成功",
	})
}

// 获取学习建议
func GetLearningAdvice(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	// 获取用户学习进度
	var progress []models.LearningProgress
	if err := database.DB.Preload("Course").Preload("Chapter").
		Where("user_id = ?", userID).
		Find(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学习进度失败",
		})
		return
	}

	// 获取用户答题情况
	var answers []models.StudentAnswer
	if err := database.DB.Preload("Question").
		Where("user_id = ?", userID).
		Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取答题情况失败",
		})
		return
	}

	// 调用AI服务生成学习建议
	aiService := services.NewAIService()
	advice, err := aiService.GenerateLearningAdvice(userID, progress, answers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成学习建议失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"advice": advice,
		},
	})
}

// SSE流式AI回复
func StreamAIChat(c *gin.Context) {
	// 设置SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	// 1. 先从query参数获取token
	token := c.Query("token")
	if token == "" {
		c.SSEvent("error", "未登录")
		c.Writer.Flush()
		return
	}
	// 2. 校验token，获取用户ID
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.SSEvent("error", "无效token")
		c.Writer.Flush()
		return
	}
	userID := claims.UserID
	_ = userID // 如后续需用userID可去掉此行
	// 3. 继续原有逻辑
	var req struct {
		Message string `json:"message"`
	}
	// 兼容query传参
	msg := c.Query("message")
	if msg != "" {
		req.Message = msg
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.SSEvent("error", "参数错误")
		c.Writer.Flush()
		return
	}

	ai := services.NewAIService()
	if ai.Xunfei == nil {
		c.SSEvent("error", "AI服务未配置")
		return
	}
	// 使用XunfeiX1Service的ChatStream方法
	err = ai.Xunfei.ChatStream(req.Message, func(chunk string) error {
		c.SSEvent("message", chunk)
		c.Writer.Flush()
		return nil
	})
	if err != nil {
		c.SSEvent("error", err.Error())
	}
	// 结束流
	time.Sleep(100 * time.Millisecond)
	c.SSEvent("end", "done")
}
