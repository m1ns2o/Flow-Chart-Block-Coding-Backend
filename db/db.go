package db

import (
	"Flow-Chart-Block-Coding-Backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 모델 마이그레이션
	err = db.AutoMigrate(&models.Class{}, &models.Problem{}, &models.User{}, &models.Solved{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
