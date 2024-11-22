package repository

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &Repository{db: db}
}
