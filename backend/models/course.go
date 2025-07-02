package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	Subject     string         `json:"subject" gorm:"size:50"` // 学科
	Grade       string         `json:"grade" gorm:"size:20"`   // 适用年级
	CoverImage  string         `json:"cover_image" gorm:"size:255"`
	TeacherID   uint           `json:"teacher_id"`
	Status      int            `json:"status" gorm:"default:1"` // 1: 正常, 0: 禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Teacher     User           `json:"teacher" gorm:"foreignKey:TeacherID"`
	Chapters    []Chapter      `json:"chapters" gorm:"foreignKey:CourseID"`
}

type Chapter struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CourseID    uint           `json:"course_id"`
	Title       string         `json:"title" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	Order       int            `json:"order"`
	Content     string         `json:"content" gorm:"type:text"` // 章节内容
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Course      Course         `json:"course" gorm:"foreignKey:CourseID"`
	Knowledge   []Knowledge    `json:"knowledge" gorm:"foreignKey:ChapterID"`
}

type Knowledge struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ChapterID  uint           `json:"chapter_id"`
	Title      string         `json:"title" gorm:"not null;size:100"`
	Content    string         `json:"content" gorm:"type:text"`
	Keywords   string         `json:"keywords" gorm:"type:text"` // 关键词，逗号分隔
	Difficulty int            `json:"difficulty"`                // 难度等级 1-5
	Order      int            `json:"order"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Chapter    Chapter        `json:"chapter" gorm:"foreignKey:ChapterID"`
}

type CourseMaterial struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CourseID   uint      `json:"course_id"`
	TeacherID  uint      `json:"teacher_id"`
	Title      string    `json:"title" gorm:"size:200"`
	FileURL    string    `json:"file_url" gorm:"size:255"`
	FileType   string    `json:"file_type" gorm:"size:50"`
	UploadedAt time.Time `json:"uploaded_at"`
	Course     Course    `json:"course" gorm:"foreignKey:CourseID"`
	Teacher    User      `json:"teacher" gorm:"foreignKey:TeacherID"`
}
