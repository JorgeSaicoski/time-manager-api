package handlers

import (
	"fmt"
	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
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

	userID, err := middleware.GetUserRequesting(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := middleware.StopCurrentTotalTime(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop current total time"})
		return
	}

	var companyID *int64
	if req.CompanyID != 0 {
		companyID = &req.CompanyID
	}

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
	userID, err := middleware.GetUserRequesting(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := middleware.StopCurrentTotalTime(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop current total time"})
		return
	}

	totalTime, err := middleware.GetCurrentTotalTime(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current total time"})
		return
	}

	c.JSON(http.StatusOK, totalTime)
}

func (h *TotalTimeHandler) GetTotalTime(c *gin.Context) {
	userID, err := middleware.GetUserRequesting(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalTime, err := middleware.GetCurrentTotalTime(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current total time"})
		return
	}

	c.JSON(http.StatusOK, totalTime)
}
