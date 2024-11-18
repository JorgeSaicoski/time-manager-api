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
	UserID    int64 `json:"userId" binding:"required"`
	CompanyID int64 `json:"companyId"`
}

func (h *TotalTimeHandler) CreateTotalTime(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var companyID *int64
	if req.CompanyID != 0 {
		companyID = &req.CompanyID
	}

	totalTime := models.TotalTime{
		UserID:    req.UserID,
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
	fmt.Println(c)
}

/*
type TotalTime struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	UserID     int64  `gorm:"not null"`
	CompanyID  *int64 // Optional company association
	StartTime  time.Time
	FinishTime time.Time
	WorkTimes  []WorkTime `gorm:"foreignKey:TotalTimeID"`
	BreakTime  *BreakTime `gorm:"foreignKey:TotalTimeID;constraint:OnDelete:CASCADE"`
	Brb        *Brb       `gorm:"foreignKey:TotalTimeID;constraint:OnDelete:CASCADE"`
	Closed     bool
}
*/
