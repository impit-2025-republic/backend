package entities

type (
	Company struct {
		CompanyID   int    `gorm:"column:company_id;primaryKey;autoIncrement"`
		Company     string `gorm:"column:company;type:character varying(255);not null"`
		Description string `gorm:"column:description;type:text"`
		Title       string `gorm:"column:title;type:text"`
		Logo        string `gorm:"column:logo;type:character varying(255)"`

		Events   []Event   `gorm:"foreignKey:CompanyID"`
		Products []Product `gorm:"foreignKey:CompanyID"`
	}
)

func (Company) TableName() string {
	return "company"
}
