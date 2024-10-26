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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

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

	if newClass.Classnum == "" || newClass.Passwd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Classnum and password are required"})
		return
	}

	var existingClass models.Class
	result := h.DB.Where("classnum = ?", newClass.Classnum).First(&existingClass)

	if result.Error == nil {
		if checkPassword(newClass.Passwd, existingClass.Passwd) {
			token, err := generateToken(existingClass.ID, existingClass.Classnum)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":  "Login successful",
				"id":       existingClass.ID,
				"classnum": existingClass.Classnum,
				"token":    token,
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPassword, err := hashPassword(newClass.Passwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	newClass.Passwd = hashedPassword

	if err := h.DB.Create(&newClass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create class"})
		return
	}

	token, err := generateToken(newClass.ID, newClass.Classnum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Class created successfully",
		"id":       newClass.ID,
		"classnum": newClass.Classnum,
		"token":    token,
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
// @Failure 404,403 {object} map[string]string
func (h *ClassHandler) GetClass(c *gin.Context) {
	id := c.Param("id")
	var class models.Class

	// Check authentication
	classID, exists := c.Get("class_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
		return
	}

	if err := h.DB.Preload("Problems").First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// Verify ownership
	if classID.(uint) != class.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to access this class"})
		return
	}

	class.Passwd = "" // Remove password from response
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
// @Failure 404,403 {object} map[string]string
// func (h *ClassHandler) GetClassByClassnum(c *gin.Context) {
// 	classnum := c.Param("classnum")
// 	var class models.Class

// 	// Check authentication
// 	authenticatedClassnum, exists := c.Get("classnum")
// 	if !exists {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
// 		return
// 	}

// 	if err := h.DB.Preload("Problems").Where("classnum = ?", classnum).First(&class).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
// 		return
// 	}

// 	// Verify ownership
// 	if authenticatedClassnum.(string) != class.Classnum {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to access this class"})
// 		return
// 	}

//		class.Passwd = "" // Remove password from response
//		c.JSON(http.StatusOK, class)
//	}
func (h *ClassHandler) GetClassByClassnum(c *gin.Context) {
	classnum := c.Param("classnum")
	var class models.Class

	if err := h.DB.Preload("Problems").Where("classnum = ?", classnum).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch class"})
		return
	}

	class.Passwd = "" // Remove password from response
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
// @Failure 400,404,403,500 {object} map[string]string
func (h *ClassHandler) UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var existingClass models.Class

	// Check authentication
	classID, exists := c.Get("class_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
		return
	}

	if err := h.DB.First(&existingClass, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// Verify ownership
	if classID.(uint) != existingClass.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this class"})
		return
	}

	var updateData models.Class
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle password update
	if updateData.Passwd != "" {
		hashedPassword, err := hashPassword(updateData.Passwd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		updateData.Passwd = hashedPassword
	} else {
		updateData.Passwd = existingClass.Passwd
	}

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
// @Failure 400,404,403,500 {object} map[string]string
func (h *ClassHandler) DeleteClass(c *gin.Context) {
	id := c.Param("id")

	// Check authentication
	classID, exists := c.Get("class_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
		return
	}

	var class models.Class
	if err := h.DB.First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// Verify ownership
	if classID.(uint) != class.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this class"})
		return
	}

	tx := h.DB.Begin()

	if err := tx.Where("class_id = ?", id).Delete(&models.Problem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related problems"})
		return
	}

	if err := tx.Delete(&class).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete class"})
		return
	}

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

	// Remove passwords from response
	for i := range classes {
		classes[i].Passwd = ""
	}

	c.JSON(http.StatusOK, classes)
}
