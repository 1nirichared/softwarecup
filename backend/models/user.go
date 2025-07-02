package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100"`
	RealName  string         `json:"real_name" gorm:"size:50"`
	Role      UserRole       `json:"role" gorm:"not null;default:'student';size:20"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Status    int            `json:"status" gorm:"default:1"` // 1: 正常, 0: 禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserProfile struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	Department string    `json:"department" gorm:"size:100"` // 院系
	Major      string    `json:"major" gorm:"size:100"`      // 专业
	Grade      string    `json:"grade" gorm:"size:20"`       // 年级
	Class      string    `json:"class" gorm:"size:50"`       // 班级
	StudentID  string    `json:"student_id" gorm:"size:50"`  // 学号
	TeacherID  string    `json:"teacher_id" gorm:"size:50"`  // 教师工号
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
}
