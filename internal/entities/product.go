package entities

type (
	ProductRepo interface {
		GetAll() ([]Product, error)
		GetByID(productId uint) (Product, error)
	}
	Product struct {
		ProductID       uint    `gorm:"column:product_id;primaryKey;autoIncrement"`
		CompanyID       *int    `gorm:"column:company_id"`
		Name            string  `gorm:"column:name;type:varchar(255);not null"`
		Price           float64 `gorm:"column:price;type:numeric(10,2)"`
		Description     string  `gorm:"column:description;type:text"`
		Image           string  `gorm:"column:image;type:varchar(255)"`
		Availability    int     `gorm:"column:avalibility"`
		ProductCategory string  `gorm:"column:product_category;type:varchar(50);default:merch"`
		CaseTypeID      *int    `gorm:"column:case_type_id"`
	}
)

func (Product) TableName() string {
	return "product"
}
