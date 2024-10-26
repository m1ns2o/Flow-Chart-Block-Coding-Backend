package main

import (
	"Flow-Chart-Block-Coding-Backend/config"
	"Flow-Chart-Block-Coding-Backend/db"
	"Flow-Chart-Block-Coding-Backend/handlers"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 설정 로드
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// JWT 시크릿 키 설정
	handlers.SetJWTSecret(cfg.GetJWTSecret())

	// 데이터베이스 연결
	database, err := db.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Gin 라우터 생성
	router := gin.Default()

	// CORS 설정
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With"}, // 헤더 추가
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},                       // 노출할 헤더 추가
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // preflight 캐시 시간 설정
	}))

	// 핸들러 초기화
	problemHandler := handlers.NewProblemHandler(database)
	classHandler := handlers.NewClassHandler(database)
	handler := handlers.NewUserHandler(database)

	// 공개 라우트 (인증 불필요)
	router.POST("/classes", classHandler.CreateClass) // 로그인/회원가입용

	// 보호된 라우트 그룹 (인증 필요)
	protected := router.Group("")
	protected.Use(handlers.AuthMiddleware())
	{
		// Problem 라우트
		protected.POST("/problems", problemHandler.CreateProblem)
		// protected.GET("/problems/:id", problemHandler.GetProblem)
		protected.PUT("/problems/:id", problemHandler.UpdateProblem)
		protected.DELETE("/problems/:id", problemHandler.DeleteProblem)
		protected.GET("/problems", problemHandler.ListProblems)

		// Class 보호된 라우트
		protected.GET("/classes/:id", classHandler.GetClass)

		protected.PUT("/classes/:id", classHandler.UpdateClass)
		protected.DELETE("/classes/:id", classHandler.DeleteClass)
		protected.GET("/classes", classHandler.ListClasses)
	}

	router.GET("/problems/:id", problemHandler.GetProblem)
	router.GET("/classes/number/:classnum", classHandler.GetClassByClassnum)

	users := router.Group("/users")
	{
		users.POST("", handler.CreateUser)
		users.GET("", handler.GetAllUsers)
		users.GET(":id", handler.GetUser)
		// users.PUT("/:id", handler.UpdateUser)
		users.DELETE(":id", handler.DeleteUser)
		// users.GET("/:id/solved", handler.GetUserSolved)
	}

	// 서버 시작
	log.Printf("Server starting on :8080")
	log.Fatal(router.Run(":8080"))
}
