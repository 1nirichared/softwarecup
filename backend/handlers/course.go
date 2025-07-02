package handlers

import (
	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateCourseRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Subject     string `json:"subject" binding:"required"`
	Grade       string `json:"grade"`
	CoverImage  string `json:"cover_image"`
}

// 创建课程
func CreateCourse(c *gin.Context) {
	teacherID := middleware.GetCurrentUserID(c)

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	course := models.Course{
		Name:        req.Name,
		Description: req.Description,
		Subject:     req.Subject,
		Grade:       req.Grade,
		CoverImage:  req.CoverImage,
		TeacherID:   teacherID,
		Status:      1,
	}

	if err := database.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "课程创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "课程创建成功",
		"data":    course,
	})
}

// 获取课程列表
func GetCourses(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	userRole := middleware.GetCurrentUserRole(c)

	var courses []models.Course
	query := database.DB.Preload("Teacher")

	// 根据角色过滤课程
	if userRole == string(models.RoleTeacher) {
		query = query.Where("teacher_id = ?", userID)
	}

	if err := query.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取课程列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    courses,
	})
}

// 获取课程详情
func GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	var course models.Course
	if err := database.DB.Preload("Teacher").Preload("Chapters").First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    course,
	})
}

// 更新课程
func UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")
	teacherID := middleware.GetCurrentUserID(c)

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 检查课程是否存在且属于当前教师
	var course models.Course
	if err := database.DB.Where("id = ? AND teacher_id = ?", courseID, teacherID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在或无权限",
		})
		return
	}

	// 更新课程信息
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"subject":     req.Subject,
		"grade":       req.Grade,
		"cover_image": req.CoverImage,
	}

	if err := database.DB.Model(&course).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "课程更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "课程更新成功",
		"data":    course,
	})
}

// 删除课程
func DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")
	teacherID := middleware.GetCurrentUserID(c)

	// 检查课程是否存在且属于当前教师
	var course models.Course
	if err := database.DB.Where("id = ? AND teacher_id = ?", courseID, teacherID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在或无权限",
		})
		return
	}

	if err := database.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "课程删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "课程删除成功",
	})
}

// 生成备课内容
func GenerateLessonPlan(c *gin.Context) {
	courseID := c.Param("courseId")
	chapterID := c.Param("chapterId")

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

	// 获取知识点
	var knowledge []models.Knowledge
	if err := database.DB.Where("chapter_id = ?", chapterID).Find(&knowledge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取知识点失败",
		})
		return
	}

	// 调用AI服务生成备课内容
	aiService := services.NewAIService()
	lessonPlan, err := aiService.GenerateLessonPlan(&course, &chapter, knowledge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成备课内容失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成成功",
		"data": gin.H{
			"lesson_plan": lessonPlan,
		},
	})
}

// 获取课程统计信息
func GetCourseStats(c *gin.Context) {
	courseID := c.Param("id")

	// 获取课程基本信息
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在",
		})
		return
	}

	// 统计章节数量
	var chapterCount int64
	database.DB.Model(&models.Chapter{}).Where("course_id = ?", courseID).Count(&chapterCount)

	// 统计练习数量
	var exerciseCount int64
	database.DB.Model(&models.Exercise{}).Where("course_id = ?", courseID).Count(&exerciseCount)

	// 统计学习进度
	var progressCount int64
	database.DB.Model(&models.LearningProgress{}).Where("course_id = ?", courseID).Count(&progressCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"course":         course,
			"chapter_count":  chapterCount,
			"exercise_count": exerciseCount,
			"progress_count": progressCount,
		},
	})
}

// 上传教学资料（教师）
func UploadCourseMaterial(c *gin.Context) {
	teacherID := middleware.GetCurrentUserID(c)
	courseID := c.Param("courseId")

	// 角色校验（仅教师）
	role, _ := c.Get("role")
	if role != "teacher" {
		c.JSON(403, gin.H{"message": "只有教师可以上传资料"})
		return
	}

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": "未选择文件"})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		title = file.Filename
	}

	// 保存文件到 uploads 目录
	uploadDir := "uploads/materials"
	os.MkdirAll(uploadDir, os.ModePerm)
	filename := time.Now().Format("20060102150405") + "_" + file.Filename
	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(500, gin.H{"message": "文件保存失败"})
		return
	}

	// 写入数据库
	material := models.CourseMaterial{
		CourseID:   parseUint(courseID),
		TeacherID:  teacherID,
		Title:      title,
		FileURL:    "/" + savePath,
		FileType:   filepath.Ext(file.Filename),
		UploadedAt: time.Now(),
	}
	if err := database.DB.Create(&material).Error; err != nil {
		c.JSON(500, gin.H{"message": "数据库写入失败"})
		return
	}

	c.JSON(200, gin.H{"message": "上传成功", "data": material})
}

// 获取课程资料列表
func GetCourseMaterials(c *gin.Context) {
	courseID := c.Param("courseId")
	var materials []models.CourseMaterial
	if err := database.DB.Where("course_id = ?", courseID).Order("uploaded_at desc").Find(&materials).Error; err != nil {
		c.JSON(500, gin.H{"message": "获取资料失败"})
		return
	}
	c.JSON(200, gin.H{"data": materials})
}

// 删除资料（仅教师）
func DeleteCourseMaterial(c *gin.Context) {
	materialID := c.Param("materialId")
	teacherID := middleware.GetCurrentUserID(c)
	role, _ := c.Get("role")
	if role != "teacher" {
		c.JSON(403, gin.H{"message": "只有教师可以删除资料"})
		return
	}
	var material models.CourseMaterial
	if err := database.DB.First(&material, materialID).Error; err != nil {
		c.JSON(404, gin.H{"message": "资料不存在"})
		return
	}
	if material.TeacherID != teacherID {
		c.JSON(403, gin.H{"message": "只能删除自己上传的资料"})
		return
	}
	// 删除文件
	os.Remove(material.FileURL[1:]) // 去掉前面的/
	if err := database.DB.Delete(&material).Error; err != nil {
		c.JSON(500, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(200, gin.H{"message": "删除成功"})
}

// 工具函数
func parseUint(s string) uint {
	u, _ := strconv.ParseUint(s, 10, 64)
	return uint(u)
}
