package repository

import "gorm.io/gorm"

func (r *Repository) GetCurrentItem(userID int64, model interface{}) error {
	if model == nil {
		return errors.New("model cannot be nil")
	}

	result := r.db.Where("user_id = ? AND closed = ?", userID, false).First(model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}

	return nil
}
