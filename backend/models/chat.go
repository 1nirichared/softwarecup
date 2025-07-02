package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatSession struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	Title     string         `json:"title" gorm:"size:100"`
	Type      string         `json:"type" gorm:"size:20"` // learning, practice, general
	CourseID  *uint          `json:"course_id"`
	ChapterID *uint          `json:"chapter_id"`
	Status    int            `json:"status" gorm:"default:1"` // 1: 活跃, 0: 结束
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Course    *Course        `json:"course" gorm:"foreignKey:CourseID"`
	Chapter   *Chapter       `json:"chapter" gorm:"foreignKey:ChapterID"`
	Messages  []ChatMessage  `json:"messages" gorm:"foreignKey:SessionID"`
}

type ChatMessage struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	SessionID   uint           `json:"session_id"`
	Role        string         `json:"role" gorm:"size:20"` // user, assistant, system
	Content     string         `json:"content" gorm:"type:text"`
	MessageType string         `json:"message_type" gorm:"size:20"` // text, image, file
	Metadata    string         `json:"metadata" gorm:"type:text"`   // JSON格式存储额外信息
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Session     ChatSession    `json:"session" gorm:"foreignKey:SessionID"`
}

type KnowledgeBase struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	Type        string         `json:"type" gorm:"size:20"`        // course, general, custom
	Content     string         `json:"content" gorm:"type:text"`   // 知识库内容
	Keywords    string         `json:"keywords" gorm:"type:text"`  // 关键词
	Embedding   string         `json:"embedding" gorm:"type:text"` // 向量化数据
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type LearningProgress struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	CourseID    uint           `json:"course_id"`
	ChapterID   uint           `json:"chapter_id"`
	Progress    float64        `json:"progress"`   // 学习进度 0-100
	TimeSpent   int            `json:"time_spent"` // 学习时长(分钟)
	LastStudyAt time.Time      `json:"last_study_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Course      Course         `json:"course" gorm:"foreignKey:CourseID"`
	Chapter     Chapter        `json:"chapter" gorm:"foreignKey:ChapterID"`
}
