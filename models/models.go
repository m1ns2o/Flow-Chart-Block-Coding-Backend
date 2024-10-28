package models

import (
	"time"
)

type Class struct {
	ID       uint      `gorm:"primaryKey"`
	Classnum string    `gorm:"unique;type:varchar(100)"` // 명확한 길이 지정
	Passwd   string    `gorm:"type:varchar(100)"`
	Problems []Problem `gorm:"constraint:OnDelete:CASCADE"`
}

type Problem struct {
	ID             uint   `gorm:"primaryKey"`
	Title          string `gorm:"type:varchar(200)"`
	Content        string `gorm:"type:varchar(500)"`
	TestcaseInput  string `gorm:"type:varchar(100)"`
	TestcaseOutput string `gorm:"type:varchar(100)"`
	ClassID        uint   `gorm:"foreignKey:ClassID;references:ID"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(50)"`
	Classnum string `gorm:"type:varchar(100)"` // 명확한 길이 지정
}

type Solved struct {
	ID        uint      `gorm:"primaryKey"`
	ProblemID uint      `gorm:"foreignKey:ProblemID;references:ID"`
	UserID    uint      `gorm:"foreignKey:UserID;references:ID"`
	UserName  string    `gorm:"type:varchar(50)"`
	SolvedAt  time.Time `gorm:"autoCreateTime"`
}
