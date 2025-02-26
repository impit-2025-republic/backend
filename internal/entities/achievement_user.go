package entities

type (
	AchievementUserRepo interface {
		Create(user AchievementUser) error
		GetAll(userID int) ([]AchievementUser, error)
	}

	AchievementUser struct {
		UserID        int `gorm:"column:user_id"`
		AchievementID int `gorm:"column:achievement_id"`
	}
)

func (AchievementUser) TableName() string {
	return "user_achievement"
}
