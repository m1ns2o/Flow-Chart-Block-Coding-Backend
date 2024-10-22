package handlers

import (
	"net/http"
	"strconv"

	"Flow-Chart-Block-Coding-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProblemHandler handles operations on Problem model
type ProblemHandler struct {
	DB *gorm.DB
}

func NewProblemHandler(db *gorm.DB) *ProblemHandler {
	return &ProblemHandler{DB: db}
}

func (h *ProblemHandler) CreateProblem(c *gin.Context) {
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create problem"})
		return
	}
	c.JSON(http.StatusCreated, problem)
}

func (h *ProblemHandler) GetProblem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var problem models.Problem
	if err := h.DB.First(&problem, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}
	c.JSON(http.StatusOK, problem)
}

func (h *ProblemHandler) UpdateProblem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var problem models.Problem
	if err := h.DB.First(&problem, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.DB.Save(&problem)
	c.JSON(http.StatusOK, problem)
}

func (h *ProblemHandler) DeleteProblem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.DB.Delete(&models.Problem{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete problem"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Problem deleted successfully"})
}

func (h *ProblemHandler) ListProblems(c *gin.Context) {
	var problems []models.Problem
	if err := h.DB.Find(&problems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list problems"})
		return
	}
	c.JSON(http.StatusOK, problems)
}

