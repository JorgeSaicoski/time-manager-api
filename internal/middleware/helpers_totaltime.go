package middleware

import (
	"fmt"
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func GetCurrentTotalTime(userID int64) (*models.TotalTime, error) {
	var totalTime models.TotalTime
	result := DB.Where("user_id = ? AND closed = ?", userID, false).First(&totalTime)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &totalTime, nil
}

func StopCurrentTotalTime(userID int64) error {
	totalTime, err := GetCurrentTotalTime(userID)
	if err != nil {
		return fmt.Errorf("error finding total time: %w", err)
	}

	totalTime.Closed = true
	totalTime.FinishTime = time.Now()
	return DB.Save(&totalTime).Error
}
