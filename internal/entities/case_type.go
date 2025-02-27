package entities

type (
	CaseType struct {
		CaseTypeID  uint   `gorm:"column:case_type_id;primaryKey;autoIncrement"`
		Name        string `gorm:"column:name;type:varchar(50);not null"`
		Description string `gorm:"column:description;type:text"`
	}
)

func (CaseType) TableName() string {
	return "case_type"
}
