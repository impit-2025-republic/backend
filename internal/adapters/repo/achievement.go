package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type achievementRepo struct {
	db *gorm.DB
}

func NewAchievementRepo(db *gorm.DB) entities.AchievementRepo {
	return achievementRepo{
		db: db,
	}
}

func (r achievementRepo) GetByNotAchievementTypeIDsAndAchievementTypeIDAndEndDs(achievementTypeIDs []int, achievementTypeID int) ([]entities.Achievement, error) {
	var achievements []entities.Achievement
	err := r.db.Where("achievement_id NOT IN (?) AND achievement_type_id = ? AND (end_ds > CURRENT_DATE OR end_ds IS NULL)").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}
