package entities

type (
	UserWalletHistoryRepo interface {
		Create(history UserWalletHistory) error
		GetByUserID(userID int) ([]UserWalletHistory, error)
	}

	UserWalletHistory struct {
		ID          uint    `gorm:"column:id;primaryKey;autoIncrement"`
		UserID      int     `gorm:"column:user_id"`
		Coin        float64 `gorm:"column:coin;type:numeric(10,2)"`
		RefillType  string  `gorm:"column:refill_type;type:varchar(255)"`
		Description string  `gorm:"column:description;type:text"`
	}
)

func (UserWalletHistory) TableName() string {
	return "user_wallet_history"
}
