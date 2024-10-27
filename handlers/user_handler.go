package handlers

import (
	"Flow-Chart-Block-Coding-Backend/models" // 프로젝트에 맞는 경로로 수정 필요
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

// CreateUserRequest는 사용자 생성 시 필요한 요청 구조체입니다
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Classnum string `json:"classnum" binding:"required"`
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetAllUsers 모든 사용자 조회
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser 특정 사용자 조회
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.Preload("Solved.Problem").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser 새로운 사용자 생성
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Classnum 존재 여부 확인
	var class models.Class
	if err := h.db.Where("classnum = ?", req.Classnum).First(&class).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classnum"})
		return
	}

	// 동일한 name과 classnum을 가진 사용자가 있는지 확인
	var existingUser models.User
	result := h.db.Where("name = ? AND classnum = ?", req.Name, req.Classnum).First(&existingUser)

	if result.Error == nil {
		// 이미 존재하는 사용자인 경우
		c.JSON(http.StatusOK, gin.H{
			"id":      existingUser.ID,
			"message": "User already exists", //응답코드 200
		})
		return
	}

	// 존재하지 않는 경우 새로운 User 생성
	newUser := models.User{
		Name:     req.Name,
		Classnum: req.Classnum,
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{ // 응답코드 201
		"id":      newUser.ID,
		"message": "User created successfully",
	})
}

// DeleteUser 사용자 삭제
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// 사용자 존재 여부 확인
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Solved 레코드와 함께 삭제
	if err := h.db.Select("Solved").Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetUsersByClass(c *gin.Context) {
	classnum := c.Param("classnum")

	// 해당 classnum이 존재하는지 먼저 확인
	var class models.Class
	if err := h.db.Where("classnum = ?", classnum).First(&class).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// 해당 classnum을 가진 모든 사용자 조회
	var users []models.User
	if err := h.db.Where("classnum = ?", classnum).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"classnum": classnum,
		"users":    users,
	})
}
