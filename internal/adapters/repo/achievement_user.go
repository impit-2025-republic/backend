package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type achievementUserRepo struct {
	db *gorm.DB
}

func NewachievementUserRepo(db *gorm.DB) entities.AchievementUserRepo {
	return achievementUserRepo{
		db: db,
	}
}

func (r achievementUserRepo) Create(user entities.AchievementUser) error {
	return r.db.Create(&user).Error
}

func (r achievementUserRepo) GetAll(userID int) ([]entities.AchievementUser, error) {
	var achievements []entities.AchievementUser
	err := r.db.Where("user_id = ?", userID).Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}
