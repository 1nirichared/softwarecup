package handlers

import (
	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateExerciseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CourseID    uint   `json:"course_id" binding:"required"`
	ChapterID   uint   `json:"chapter_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Duration    int    `json:"duration"`
	TotalScore  int    `json:"total_score"`
}

type SubmitAnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

// 创建练习
func CreateExercise(c *gin.Context) {
	teacherID := middleware.GetCurrentUserID(c)

	var req CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证课程是否属于当前教师
	var course models.Course
	if err := database.DB.Where("id = ? AND teacher_id = ?", req.CourseID, teacherID).First(&course).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程不存在或无权限",
		})
		return
	}

	exercise := models.Exercise{
		Title:       req.Title,
		Description: req.Description,
		CourseID:    req.CourseID,
		ChapterID:   req.ChapterID,
		Type:        req.Type,
		Duration:    req.Duration,
		TotalScore:  req.TotalScore,
		Status:      1,
	}

	if err := database.DB.Create(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "练习创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "练习创建成功",
		"data":    exercise,
	})
}

// 获取练习列表
func GetExercises(c *gin.Context) {
	courseID := c.Query("course_id")
	chapterID := c.Query("chapter_id")

	var exercises []models.Exercise
	query := database.DB.Preload("Course").Preload("Chapter")

	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	if chapterID != "" {
		query = query.Where("chapter_id = ?", chapterID)
	}

	if err := query.Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取练习列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    exercises,
	})
}

// 获取练习详情
func GetExercise(c *gin.Context) {
	exerciseID := c.Param("id")

	var exercise models.Exercise
	if err := database.DB.Preload("Course").Preload("Chapter").Preload("Questions").First(&exercise, exerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "练习不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    exercise,
	})
}

// 生成练习题
func GenerateExercises(c *gin.Context) {
	courseID := c.Param("courseId")
	chapterID := c.Param("chapterId")
	questionType := c.Query("type")
	countStr := c.Query("count")

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		count = 5 // 默认生成5道题
	}

	// 获取课程信息
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在",
		})
		return
	}

	// 获取章节信息
	var chapter models.Chapter
	if err := database.DB.First(&chapter, chapterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "章节不存在",
		})
		return
	}

	// 调用AI服务生成练习题
	aiService := services.NewAIService()
	questions, err := aiService.GenerateExercises(&course, &chapter, models.QuestionType(questionType), count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成练习题失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成成功",
		"data":    questions,
	})
}

// 开始练习
func StartExercise(c *gin.Context) {
	exerciseID := c.Param("id")
	userID := middleware.GetCurrentUserID(c)

	// 检查练习是否存在
	var exercise models.Exercise
	if err := database.DB.First(&exercise, exerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "练习不存在",
		})
		return
	}

	// 创建练习记录
	record := models.ExerciseRecord{
		UserID:     userID,
		ExerciseID: exercise.ID,
		Status:     "ongoing",
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "开始练习失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "开始练习成功",
		"data": gin.H{
			"record_id": record.ID,
			"exercise":  exercise,
		},
	})
}

// 提交答案
func SubmitAnswer(c *gin.Context) {
	recordID := c.Param("recordId")
	userID := middleware.GetCurrentUserID(c)

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 验证练习记录
	var record models.ExerciseRecord
	if err := database.DB.Where("id = ? AND user_id = ?", recordID, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "练习记录不存在",
		})
		return
	}

	// 获取题目信息
	var question models.Question
	if err := database.DB.First(&question, req.QuestionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "题目不存在",
		})
		return
	}

	// 调用AI服务评估答案
	aiService := services.NewAIService()
	score, isCorrect, feedback, err := aiService.EvaluateAnswer(&question, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "答案评估失败",
		})
		return
	}

	// 保存学生答案
	studentAnswer := models.StudentAnswer{
		UserID:     userID,
		ExerciseID: record.ExerciseID,
		QuestionID: req.QuestionID,
		Answer:     req.Answer,
		Score:      score,
		IsCorrect:  isCorrect,
		Feedback:   feedback,
	}

	if err := database.DB.Create(&studentAnswer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存答案失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "答案提交成功",
		"data": gin.H{
			"score":      score,
			"is_correct": isCorrect,
			"feedback":   feedback,
		},
	})
}

// 完成练习
func CompleteExercise(c *gin.Context) {
	recordID := c.Param("recordId")
	userID := middleware.GetCurrentUserID(c)

	// 验证练习记录
	var record models.ExerciseRecord
	if err := database.DB.Where("id = ? AND user_id = ?", recordID, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "练习记录不存在",
		})
		return
	}

	// 计算总分
	var totalScore int
	database.DB.Model(&models.StudentAnswer{}).
		Where("user_id = ? AND exercise_id = ?", userID, record.ExerciseID).
		Select("SUM(score)").
		Scan(&totalScore)

	// 更新练习记录
	updates := map[string]interface{}{
		"score":  totalScore,
		"status": "completed",
	}

	if err := database.DB.Model(&record).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "完成练习失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "练习完成",
		"data": gin.H{
			"total_score": totalScore,
			"record":      record,
		},
	})
}

// 获取练习统计
func GetExerciseStats(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	// 统计完成的练习数量
	var completedCount int64
	database.DB.Model(&models.ExerciseRecord{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Count(&completedCount)

	// 统计平均分
	var avgScore float64
	database.DB.Model(&models.ExerciseRecord{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Select("AVG(score)").
		Scan(&avgScore)

	// 统计正确率
	var totalAnswers int64
	var correctAnswers int64
	database.DB.Model(&models.StudentAnswer{}).
		Where("user_id = ?", userID).
		Count(&totalAnswers)
	database.DB.Model(&models.StudentAnswer{}).
		Where("user_id = ? AND is_correct = ?", userID, true).
		Count(&correctAnswers)

	var accuracy float64
	if totalAnswers > 0 {
		accuracy = float64(correctAnswers) / float64(totalAnswers) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"completed_count": completedCount,
			"avg_score":       avgScore,
			"accuracy":        accuracy,
			"total_answers":   totalAnswers,
			"correct_answers": correctAnswers,
		},
	})
}
