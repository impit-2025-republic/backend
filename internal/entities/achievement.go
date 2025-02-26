package entities

import "time"

type (
	AchievementRepo interface {
		GetByNotAchievementTypeIDsAndAchievementTypeIDAndEndDs(achievementTypeIDs []int, achievementTypeID int) ([]Achievement, error)
	}

	Achievement struct {
		AchievementID     uint       `gorm:"primaryKey;column:achievement_id;autoIncrement"`
		CompanyID         *int       `gorm:"column:company_id"`
		Logo              *string    `gorm:"column:logo;type:varchar(255)"`
		Name              *string    `gorm:"column:name;type:varchar(255)"`
		Description       *string    `gorm:"column:description;type:text"`
		AchievementTypeID *int       `gorm:"column:achievement_type_id"`
		TreshholdValue    *int       `gorm:"column:treshhold_value"`
		Coin              *float64   `gorm:"column:coin;type:numeric(10,2)"`
		EndDs             *time.Time `gorm:"column:end_ds;type:date"`
	}
)

func (Achievement) TableName() string {
	return "achievement"
}
