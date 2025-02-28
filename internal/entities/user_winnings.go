package entities

import "time"

type (
	UserWinningRepo interface {
		Create(winner UserWinning) error
		GetMyWinnings(userID uint) ([]struct {
			UserWinning
			Product
		}, error)
	}

	UserWinning struct {
		UserWinningID uint       `gorm:"column:user_winning_id;primaryKey;autoIncrement"`
		UserID        int        `gorm:"column:user_id;not null"`
		ProductID     int        `gorm:"column:product_id;not null"`
		WonAt         time.Time  `gorm:"column:won_at;default:CURRENT_TIMESTAMP"`
		Delivered     bool       `gorm:"column:delivered;default:false"`
		DeliveredAt   *time.Time `gorm:"column:delivered_at"`
		DeliveredBy   *int       `gorm:"column:delivered_by"`
		WinType       string     `gorm:"column:win_type"`
	}
)

func (UserWinning) TableName() string {
	return "user_winnings"
}
