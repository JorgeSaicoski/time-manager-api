package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/JorgeSaicoski/time-manager-api/internal/models"
)

type TotalTimeHandler struct {
	db *gorm.DB
}

func NewTotalTimerHandler(db *gorm.DB) *TotalTimeHandler {
	return &TotalTimeHandler{db: db}
}

type CreateRequest struct {
	CompanyID int64 `json:"companyId"`
}

func (h *TotalTimeHandler) CreateTotalTime(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var companyID *int64
	if req.CompanyID != 0 {
		companyID = &req.CompanyID
	}

	fmt.Printf("Authenticated User ID: %d\n", userID)

	totalTime := models.TotalTime{
		UserID:    userID,
		CompanyID: companyID,
		StartTime: time.Now(),
		Closed:    false,
	}

	if result := h.db.Create(&totalTime); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create total time"})
		return
	}

	c.JSON(http.StatusOK, totalTime)
}

func (h *TotalTimeHandler) CloseTotalTime(c *gin.Context) {
	// waiting
	fmt.Println(c)
}
