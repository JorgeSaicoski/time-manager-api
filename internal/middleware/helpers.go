package middleware

import (
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
	"gorm.io/gorm"

	"strings"
	"time"
)

var DB *gorm.DB

func StopCurrentlyTotalTime(userID int64) error {
	var totalTime models.TotalTime
	result := DB.Where("user_id = ? AND closed = ?", userID, false).First(&totalTime)
	if result.Error != nil {
		return fmt.Errorf("error finding total time: %w", result.Error)
	}

	totalTime.Closed = true
	totalTime.FinishTime = time.Now()

	return DB.Save(&totalTime).Error
}
