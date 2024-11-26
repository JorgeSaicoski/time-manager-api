package repository

import (
	"fmt"
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) GetCurrentTotalTime(userID int64) (*models.TotalTime, error) {
	var totalTime models.TotalTime
	if err := GetCurrentItem(userID, &totalTime); err  ! = n il {
		return nil, err
	}
	if totalTime.ID == nil {
		return nil, nil
	}
	return &totalTime, nil
}

func (r *Repository) StopCurrentTotalTime(userID int64) error {
	totalTime, err := r.GetCurrentTotalTime(userID)
	if err != nil {
		return fmt.Errorf("error finding total time: %w", err)
	}

	if totalTime == nil {
		return nil
	}

	totalTime.Closed = true
	totalTime.FinishTime = time.Now()
	return r.db.Save(totalTime).Error
}

func (r *Repository) CreateTotalTime(tt *models.TotalTime) error {
	return r.db.Create(tt).Error
}
