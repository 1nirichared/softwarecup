package database

import (
	"backend/config"
	"backend/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	cfg := config.GlobalConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	// 自动迁移数据库表
	err = AutoMigrate()
	if err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Course{},
		&models.Chapter{},
		&models.Knowledge{},
		&models.CourseMaterial{},
		&models.Exercise{},
		&models.Question{},
		&models.StudentAnswer{},
		&models.ExerciseRecord{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.KnowledgeBase{},
		&models.LearningProgress{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
