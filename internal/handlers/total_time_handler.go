package handlers

import (
	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type TotalTimeHandler struct {
	repo *repository.Repository
}

func NewTotalTimerHandler(repo *repository.Repository) *TotalTimeHandler {
	return &TotalTimeHandler{repo: repo}
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

	if err := h.repo.StopCurrentTotalTime(userID); err != nil {
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

	if result := h.repo.CreateTotalTime(&totalTime); result.Error != nil {
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

	if err := h.repo.StopCurrentTotalTime(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop current total time"})
		return
	}

	totalTime, err := h.repo.GetCurrentTotalTime(userID)
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

	totalTime, err := h.repo.GetCurrentTotalTime(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current total time"})
		return
	}

	if totalTime == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No current total time found"})
		return
	}

	c.JSON(http.StatusOK, totalTime)
}
