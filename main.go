package main

import (
	"Flow-Chart-Block-Coding-Backend/db"
	"Flow-Chart-Block-Coding-Backend/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 데이터베이스 연결
	dsn := "flowchart_user:minseo0128@tcp(localhost:3306)/flowchart_db?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := db.InitDB(dsn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Gin 라우터 생성
	router := gin.Default()

	router.Use(cors.Default())

	// 핸들러 초기화
	problemHandler := handlers.NewProblemHandler(database)
	classHandler := handlers.NewClassHandler(database)
	// teacherHandler := handlers.NewTeacherHandler(database)
	// userHandler := handlers.NewUserHandler(database)
	// solvedHandler := handlers.NewSolvedHandler(database)

	// 라우트 설정
	// Problem 라우트
	router.POST("/problems", problemHandler.CreateProblem)
	router.GET("/problems/:id", problemHandler.GetProblem)
	router.PUT("/problems/:id", problemHandler.UpdateProblem)
	router.DELETE("/problems/:id", problemHandler.DeleteProblem)
	router.GET("/problems", problemHandler.ListProblems)

	// Class 라우트
	router.POST("/classes", classHandler.CreateClass)
	router.GET("/classes/:id", classHandler.GetClass)
	router.GET("/classes/number/:classnum", classHandler.GetClassByClassnum)
	router.PUT("/classes/:id", classHandler.UpdateClass)
	router.DELETE("/classes/:id", classHandler.DeleteClass)
	router.GET("/classes", classHandler.ListClasses)
	// router.POST("/classes/:id/verify-password", classHandler.VerifyClassPassword)

	// // Teacher 라우트
	// router.POST("/teachers", teacherHandler.CreateTeacher)
	// router.GET("/teachers/:id", teacherHandler.GetTeacher)
	// router.PUT("/teachers/:id", teacherHandler.UpdateTeacher)
	// router.DELETE("/teachers/:id", teacherHandler.DeleteTeacher)
	// router.GET("/teachers", teacherHandler.ListTeachers)

	// // User 라우트
	// router.POST("/users", userHandler.CreateUser)
	// router.GET("/users/:id", userHandler.GetUser)
	// router.PUT("/users/:id", userHandler.UpdateUser)
	// router.DELETE("/users/:id", userHandler.DeleteUser)
	// router.GET("/users", userHandler.ListUsers)

	// // Solved 라우트
	// router.POST("/solved", solvedHandler.CreateSolved)
	// router.GET("/solved/:id", solvedHandler.GetSolved)
	// router.PUT("/solved/:id", solvedHandler.UpdateSolved)
	// router.DELETE("/solved/:id", solvedHandler.DeleteSolved)
	// router.GET("/solved", solvedHandler.ListSolved)

	// 서버 시작
	log.Fatal(router.Run(":8080"))
}
