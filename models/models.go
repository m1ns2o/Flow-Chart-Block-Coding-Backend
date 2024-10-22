package models

import (
	"time"
)

type Class struct {
	ID       uint `gorm:"primaryKey"`
	Classnum string  `gorm:"unique"`
	Passwd  string
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
	Class          Class `gorm:"foreignKey:ClassID"`
	SolvedProblems []Solved
}

type User struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"unique"`
	ClassID uint
	Class   Class `gorm:"foreignKey:ClassID"`
	Solved  []Solved
}

type Solved struct {
	ID        uint `gorm:"primaryKey"`
	ProblemID uint
	Problem   Problem `gorm:"foreignKey:ProblemID"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	SolvedAt  time.Time
}
