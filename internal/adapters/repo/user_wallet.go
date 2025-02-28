package repo

import (
	"b8boost/backend/internal/entities"
	"fmt"

	"gorm.io/gorm"
)

type userWallet struct {
	db *gorm.DB
}

func NewUserWallet(db *gorm.DB) entities.UserWalletRepo {
	return userWallet{db: db}
}

func (r userWallet) Create(wallet entities.UserWallet) error {
	return r.db.Create(&wallet).Error
}

func (r userWallet) DownBalance(user_ids []int, price float64) error {
	if price < 0 {
		return fmt.Errorf("cannot reduce balance by negative amount")
	}

	result := r.db.Model(&entities.UserWallet{}).
		Where("user_id IN ?", user_ids).
		Update("price", gorm.Expr("price - ?", price))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user wallet not found for user_id: %d", user_ids)
	}

	return nil
}

func (r userWallet) GetTopBalance() ([]struct {
	entities.UserWallet
	entities.User
}, error) {
	var results []struct {
		entities.UserWallet
		entities.User
	}

	err := r.db.Table("user_wallet").
		Select("user_wallet.*, product.*").
		Joins("JOIN users ON user_wallet.user_id = users.user_id").
		Order("user_wallet.price DESC").Limit(10).
		Scan(&results).Error

	fmt.Println(results)

	return results, err
}

func (r userWallet) GetWallet(userId uint) (entities.UserWallet, error) {
	var userWallet entities.UserWallet
	err := r.db.Where("user_id = ?", userId).First(&userWallet).Error
	if err != nil {
		return entities.UserWallet{}, err
	}
	return userWallet, nil
}

func (r userWallet) UpBalance(user_ids []int, price float64) error {
	if price < 0 {
		return fmt.Errorf("cannot reduce balance by negative amount")
	}

	result := r.db.Model(&entities.UserWallet{}).
		Where("user_id IN ?", user_ids).
		Update("price", gorm.Expr("price + ?", price))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user wallet not found for user_id: %d", user_ids)
	}

	return nil
}
