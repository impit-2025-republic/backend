package entities

type (
	CaseProductProbabilityRepo interface {
		GetAll(caseTypeID uint) ([]CaseProductProbability, error)
	}

	CaseProductProbability struct {
		ID              uint    `gorm:"column:id;primaryKey;autoIncrement"`
		CaseTypeID      uint    `gorm:"column:case_type_id;not null"`
		ProductID       uint    `gorm:"column:product_id;not null"`
		DropProbability float64 `gorm:"column:drop_probability;type:numeric(5,2);not null"`

		// Define association (optional)
		CaseType CaseType `gorm:"foreignKey:CaseTypeID"`
	}
)

func (CaseProductProbability) TableName() string {
	return "case_product_probability"
}
