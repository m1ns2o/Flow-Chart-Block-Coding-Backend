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
	// gin.SetMode(gin.ReleaseMode)
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
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API 그룹 생성
	api := router.Group("/api")
	{
		// 핸들러 초기화
		problemHandler := handlers.NewProblemHandler(database)
		classHandler := handlers.NewClassHandler(database)
		handler := handlers.NewUserHandler(database)
		solvedHandler := handlers.SolvedHandler(database)

		// Solve 그룹
		solve := api.Group("/solve")
		{
			solve.POST("", solvedHandler)
			solve.GET("/user/:username", handlers.GetUserSolvedProblems(database))
			solve.GET("/problem/:problem_id", handlers.GetProblemSolvedUsers(database))
		}

		// Problems 그룹
		problems := api.Group("/problems")
		{
			problems.GET("/:id", problemHandler.GetProblem)

			// 보호된 라우트
			protected := problems.Group("")
			protected.Use(handlers.AuthMiddleware())
			{
				protected.POST("", problemHandler.CreateProblem)
				protected.PUT("/:id", problemHandler.UpdateProblem)
				protected.DELETE("/:id", problemHandler.DeleteProblem)
				protected.GET("", problemHandler.ListProblems)
			}
		}

		// Classes 그룹
		classes := api.Group("/classes")
		{
			classes.POST("", classHandler.CreateClass) // 로그인/회원가입용
			classes.GET("/number/:classnum", classHandler.GetClassByClassnum)

			// 보호된 라우트
			protected := classes.Group("")
			protected.Use(handlers.AuthMiddleware())
			{
				protected.GET("/:id", classHandler.GetClass)
				protected.PUT("/:id", classHandler.UpdateClass)
				protected.DELETE("/:id", classHandler.DeleteClass)
				protected.GET("", classHandler.ListClasses)
			}
		}

		// Users 그룹
		users := api.Group("/users")
		{
			users.POST("", handler.CreateUser)
			users.GET("", handler.GetAllUsers)
			users.GET(":id", handler.GetUser)
			users.GET("class/:classnum", handler.GetUsersByClass)
			users.DELETE(":id", handler.DeleteUser)
		}
	}

	// 서버 시작
	log.Printf("Server starting on :8080")
	log.Fatal(router.Run(":8080"))
}
