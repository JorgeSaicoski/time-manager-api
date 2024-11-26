package repository

import (
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
)

func (r *Repository) GetRunningBreak(userID int64) (*models.BreakTime, error){
  var break models.BreakTime
  if err := GetCurrentItem(userID, &break); err!=nil{
    return nil, err
  }
  if break.ID == nil {
    return nil, nil
  }

	return &break, nil  
}
