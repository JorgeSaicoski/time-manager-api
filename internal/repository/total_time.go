package repository

import (
	"fmt"
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) GetCurrentTotalTime(userID int64) (*models.TotalTime, error) {
	var totalTime models.TotalTime
	result := r.db.Where("user_id = ? AND closed = ?", userID, false).First(&totalTime)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &totalTime, nil
}

func (r *Repository) StopCurrentTotalTime(userID int64) error {
	totalTime, err := r.GetCurrentTotalTime(userID)
	if err != nil {
		return fmt.Errorf("error finding total time: %w", err)
	}

	totalTime.Closed = true
	totalTime.FinishTime = time.Now()
	return h.db.Save(&totalTime).Error
}

func (r *Repository) CreateTotalTime(tt *models.TotalTime) error {
	return r.db.Create(tt).Error
}
