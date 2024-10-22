package handlers

import (
	"Flow-Chart-Block-Coding-Backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ClassHandler struct {
	DB *gorm.DB
}

func NewClassHandler(db *gorm.DB) *ClassHandler {
	return &ClassHandler{DB: db}
}

// 비밀번호 해싱을 위한 유틸리티 함수
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 비밀번호 검증을 위한 유틸리티 함수
func checkPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CreateClass godoc
// @Summary Create a new class or verify existing class
// @Description Create a new class if it doesn't exist, or verify password if it exists
// @Tags classes
// @Accept  json
// @Produce  json
// @Param class body models.Class true "Create class"
// @Success 201 {object} models.Class
// @Success 200 {object} models.Class
// @Failure 400,401,500 {object} map[string]string
func (h *ClassHandler) CreateClass(c *gin.Context) {
	var newClass models.Class
	if err := c.ShouldBindJSON(&newClass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if newClass.Classnum == "" || newClass.Passwd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Classnum and password are required"})
		return
	}

	// 기존 클래스 확인
	var existingClass models.Class
	result := h.DB.Where("classnum = ?", newClass.Classnum).First(&existingClass)

	if result.Error == nil {
		// 클래스가 이미 존재하는 경우
		if checkPassword(newClass.Passwd, existingClass.Passwd) {
			// 비밀번호가 일치하는 경우
			c.JSON(http.StatusOK, gin.H{
				"message":  "Class already exists",
				"id":       existingClass.ID,
				"classnum": existingClass.Classnum,
			})
			return
		} else {
			// 비밀번호가 일치하지 않는 경우
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password for existing class"})
			return
		}
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 데이터베이스 오류
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 비밀번호 해싱
	hashedPassword, err := hashPassword(newClass.Passwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	newClass.Passwd = hashedPassword

	// 새로운 클래스 생성
	if err := h.DB.Create(&newClass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create class"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Class created successfully",
		"id":       newClass.ID,
		"classnum": newClass.Classnum,
	})
}

// GetClass godoc
// @Summary Get a class
// @Description Get a class by its ID
// @Tags classes
// @Accept  json
// @Produce  json
// @Param id path int true "Class ID"
// @Success 200 {object} models.Class
// @Failure 404 {object} map[string]string
func (h *ClassHandler) GetClass(c *gin.Context) {
	id := c.Param("id")
	var class models.Class

	if err := h.DB.Preload("Problems").First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// 비밀번호는 응답에서 제외
	class.Passwd = ""
	c.JSON(http.StatusOK, class)
}

// GetClassByClassnum godoc
// @Summary Get a class by classnum
// @Description Get a class by its classnum
// @Tags classes
// @Accept  json
// @Produce  json
// @Param classnum path string true "Class Number"
// @Success 200 {object} models.Class
// @Failure 404 {object} map[string]string
func (h *ClassHandler) GetClassByClassnum(c *gin.Context) {
	classnum := c.Param("classnum")
	var class models.Class

	if err := h.DB.Preload("Problems").Where("classnum = ?", classnum).First(&class).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// 비밀번호는 응답에서 제외
	class.Passwd = ""
	c.JSON(http.StatusOK, class)
}

// UpdateClass godoc
// @Summary Update a class
// @Description Update a class's information
// @Tags classes
// @Accept  json
// @Produce  json
// @Param id path int true "Class ID"
// @Param class body models.Class true "Update class"
// @Success 200 {object} models.Class
// @Failure 400,404,500 {object} map[string]string
func (h *ClassHandler) UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var existingClass models.Class

	// 기존 클래스 확인
	if err := h.DB.First(&existingClass, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	var updateData models.Class
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 비밀번호 업데이트가 있는 경우
	if updateData.Passwd != "" {
		hashedPassword, err := hashPassword(updateData.Passwd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		updateData.Passwd = hashedPassword
	} else {
		// 비밀번호가 제공되지 않은 경우 기존 비밀번호 유지
		updateData.Passwd = existingClass.Passwd
	}

	// 업데이트 수행
	if err := h.DB.Model(&existingClass).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update class"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Class updated successfully",
		"id":       existingClass.ID,
		"classnum": updateData.Classnum,
	})
}

// DeleteClass godoc
// @Summary Delete a class
// @Description Delete a class by its ID
// @Tags classes
// @Accept  json
// @Produce  json
// @Param id path int true "Class ID"
// @Success 200 {object} map[string]string
// @Failure 400,500 {object} map[string]string
func (h *ClassHandler) DeleteClass(c *gin.Context) {
	id := c.Param("id")

	// 클래스 존재 여부 확인
	var class models.Class
	if err := h.DB.First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// 트랜잭션 시작
	tx := h.DB.Begin()

	// 관련된 문제들 삭제
	if err := tx.Where("class_id = ?", id).Delete(&models.Problem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related problems"})
		return
	}

	// 클래스 삭제
	if err := tx.Delete(&class).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete class"})
		return
	}

	// 트랜잭션 커밋
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Class and related problems deleted successfully"})
}

// ListClasses godoc
// @Summary List all classes
// @Description Get a list of all classes
// @Tags classes
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Class
// @Failure 500 {object} map[string]string
func (h *ClassHandler) ListClasses(c *gin.Context) {
	var classes []models.Class
	if err := h.DB.Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list classes"})
		return
	}

	// 모든 클래스의 비밀번호 정보 제외
	for i := range classes {
		classes[i].Passwd = ""
	}

	c.JSON(http.StatusOK, classes)
}

// VerifyClassPassword godoc
// @Summary Verify class password
// @Description Verify the password for a specific class
// @Tags classes
// @Accept  json
// @Produce  json
// @Param id path int true "Class ID"
// @Param password body map[string]string true "Password"
// @Success 200 {object} map[string]string
// @Failure 401,404 {object} map[string]string
// func (h *ClassHandler) VerifyClassPassword(c *gin.Context) {
// 	id := c.Param("id")
// 	var class models.Class
// 	if err := h.DB.First(&class, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
// 		return
// 	}

// 	var input struct {
// 		Password string `json:"password"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if !checkPassword(input.Password, class.Passwd) {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":  "Password verified successfully",
// 		"id":       class.ID,
// 		"classnum": class.Classnum,
// 	})
// }
