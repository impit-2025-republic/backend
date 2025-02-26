package entities

type (
	UserWalletRepo interface {
		Create(wallet UserWallet) error
		UpBalance(user_ids []int, price float64) error
		DownBalance(user_ids []int, price float64) error
		GetWallet(userID uint) (UserWallet, error)
	}

	UserWallet struct {
		UserID int     `gorm:"column:user_id"`
		Price  float64 `gorm:"column:price"`
	}
)

func (UserWallet) TableName() string {
	return "user_wallet"
}
