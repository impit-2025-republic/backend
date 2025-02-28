package repo

import (
	"b8boost/backend/internal/entities"
	"fmt"

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

func (r userWinningRepo) GetMyWinnings(userID uint) ([]struct {
	entities.UserWinning
	Product entities.Product
}, error) {
	var results []struct {
		entities.UserWinning
		Product entities.Product
	}

	err := r.db.Table("user_winnings").
		Select("user_winnings.*, product.*").
		Joins("JOIN product ON user_winnings.product_id = product.product_id").
		Where("user_winnings.user_id = ?", userID).
		Scan(&results).Error

	fmt.Println(results)

	return results, err
}
