package repository

import (
	"time"

	"github.com/JorgeSaicoski/time-manager-api/internal/models"
)

func (r *Repository) GetRunningBreak(userID int64) (*models.BreakTime, error) {
	var breakTime models.BreakTime
	if err := GetCurrentItem(userID, &breakTime); err != nil {
		return nil, err
	}
	if breakTime.ID == nil {
		return nil, nil
	}

	return &breakTime, nil
}

func (r *Repository) StopCurrentBreak(userID int64) error {
	breakTime, err := r.GetRunningBreak(userID)
	if err != nil {
		return fmt.Errorf("error finding break time: %w", err)
	}

	if breakTime == nil {
		return nil
	}

	breakTime.Closed = true
	now := time.Now()
	diff := now.Sub(breakTime.StartTime)
	breakTime.Duration.Add(diff)

	return r.db.Save(breakTime).Error
}

func (r *Repository) CreateBreakTime(bt *models.TotalTime) error {
	return r.db.Create(bt).Error
}
