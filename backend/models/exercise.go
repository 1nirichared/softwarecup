package models

import (
	"time"

	"gorm.io/gorm"
)

type QuestionType string

const (
	QuestionTypeSingle   QuestionType = "single"   // 单选题
	QuestionTypeMultiple QuestionType = "multiple" // 多选题
	QuestionTypeFill     QuestionType = "fill"     // 填空题
	QuestionTypeEssay    QuestionType = "essay"    // 简答题
	QuestionTypeCode     QuestionType = "code"     // 编程题
)

type Exercise struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	CourseID    uint           `json:"course_id"`
	ChapterID   uint           `json:"chapter_id"`
	Type        string         `json:"type" gorm:"size:20"` // practice, exam
	Duration    int            `json:"duration"`            // 时长(分钟)
	TotalScore  int            `json:"total_score"`
	Status      int            `json:"status" gorm:"default:1"` // 1: 正常, 0: 禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Course      Course         `json:"course" gorm:"foreignKey:CourseID"`
	Chapter     Chapter        `json:"chapter" gorm:"foreignKey:ChapterID"`
	Questions   []Question     `json:"questions" gorm:"foreignKey:ExerciseID"`
}

type Question struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ExerciseID uint           `json:"exercise_id"`
	Type       QuestionType   `json:"type" gorm:"not null;size:20"`
	Title      string         `json:"title" gorm:"not null;size:200"`
	Content    string         `json:"content" gorm:"type:text"`
	Options    string         `json:"options" gorm:"type:text"`  // JSON格式存储选项
	Answer     string         `json:"answer" gorm:"type:text"`   // 正确答案
	Analysis   string         `json:"analysis" gorm:"type:text"` // 解析
	Score      int            `json:"score"`                     // 分值
	Difficulty int            `json:"difficulty"`                // 难度等级 1-5
	Order      int            `json:"order"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Exercise   Exercise       `json:"exercise" gorm:"foreignKey:ExerciseID"`
}

type StudentAnswer struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	ExerciseID uint           `json:"exercise_id"`
	QuestionID uint           `json:"question_id"`
	Answer     string         `json:"answer" gorm:"type:text"`
	Score      int            `json:"score"`                     // 得分
	IsCorrect  bool           `json:"is_correct"`                // 是否正确
	Feedback   string         `json:"feedback" gorm:"type:text"` // AI反馈
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	Exercise   Exercise       `json:"exercise" gorm:"foreignKey:ExerciseID"`
	Question   Question       `json:"question" gorm:"foreignKey:QuestionID"`
}

type ExerciseRecord struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	ExerciseID uint           `json:"exercise_id"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    *time.Time     `json:"end_time"`
	Score      int            `json:"score"`
	Status     string         `json:"status" gorm:"size:20"` // ongoing, completed, timeout
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	Exercise   Exercise       `json:"exercise" gorm:"foreignKey:ExerciseID"`
}
