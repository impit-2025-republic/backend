package entities

type (
	Product struct {
		ProductID   int     `gorm:"column:product_id;primaryKey;autoIncrement"`
		CompanyID   *int    `gorm:"column:company_id"`
		Name        string  `gorm:"column:name;type:character varying(255);not null"`
		Price       float64 `gorm:"column:price;type:numeric(10,2)"`
		Description string  `gorm:"column:description;type:text"`
		Image       string  `gorm:"column:image;type:character varying(255)"`
		Avalibility *int    `gorm:"column:avalibility"`

		Company *Company `gorm:"foreignKey:CompanyID"`
	}
)

func (Product) TableName() string {
	return "product"
}
