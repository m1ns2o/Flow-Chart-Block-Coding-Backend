package models

import (
	"time"
)

type Class struct {
	ID       uint   `gorm:"primaryKey"`
	Classnum string `gorm:"unique"`
	Passwd   string
	Problems []Problem
	// Teachers []Teacher
}

type Problem struct {
	ID             uint `gorm:"primaryKey"`
	Title          string
	Content        string
	TestcaseInput  string
	TestcaseOutput string
	ClassID        uint
	// Class          Class `gorm:"foreignKey:ClassID"`
	// SolvedProblems []Solved
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique"`
	Classnum string `gorm:"type:string;foreignKey:Classnum;references:Classnum"`
	// Solved   []Solved
}

type Solved struct {
    ID        uint      `gorm:"primaryKey"`
    ProblemID uint      `gorm:"foreignKey:ProblemID;references:ID"`
    UserID    uint      `gorm:"foreignKey:UserID;references:ID"`
    SolvedAt  time.Time
}