package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type userWalletHistoryRepo struct {
	db *gorm.DB
}

func NewUserWalletHistoryRepo(db *gorm.DB) entities.UserWalletHistoryRepo {
	return userWalletHistoryRepo{db: db}
}

func (r userWalletHistoryRepo) Create(history entities.UserWalletHistory) error {
	return r.db.Create(history).Error
}

func (r userWalletHistoryRepo) GetByUserID(userID int) ([]entities.UserWalletHistory, error) {
	var transactions []entities.UserWalletHistory
	err := r.db.Where("user_id = ?", userID).Find(&transactions).Error
	if err != nil {
		return []entities.UserWalletHistory{}, err
	}
	return transactions, nil
}
