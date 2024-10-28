package models

import (
	"time"
)

type Class struct {
	ID       uint   `gorm:"primaryKey"`
	Classnum string `gorm:"unique;varchar(100)"`
	Passwd   string `gorm:"type:varchar(100)"`
	Problems []Problem
	// Teachers []Teacher
}

type Problem struct {
	ID             uint   `gorm:"primaryKey"`
	Title          string `gorm:"type:varchar(200)"`
	Content        string `gorm:"type:varchar(500)"`
	TestcaseInput  string `gorm:"type:varchar(100)"`
	TestcaseOutput string `gorm:"type:varchar(100)"`
	ClassID        uint
	// Class          Class `gorm:"foreignKey:ClassID"`
	// SolvedProblems []Solved
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(50)"`
	Classnum string `gorm:"type:string;foreignKey:Classnum;references:Classnum"`
	// Solved   []Solved
}

type Solved struct {
	ID        uint      `gorm:"primaryKey"`
	ProblemID uint      `gorm:"foreignKey:ProblemID;references:ID"`
	UserID    uint      `gorm:"foreignKey:UserID;references:ID"`
	UserName  string    `gorm:"foreignKey:UserName;references:UserName;type:varchar(50)"`
	SolvedAt  time.Time `gorm:"autoCreateTime"`
}
