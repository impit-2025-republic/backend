package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type userWinningRepo struct {
	db *gorm.DB
}

func NewUserWinningRepo(db *gorm.DB) entities.UserWinningRepo {
	return userWinningRepo{db: db}
}

func (r userWinningRepo) Create(history entities.UserWinning) error {
	return r.db.Create(&history).Error
}
