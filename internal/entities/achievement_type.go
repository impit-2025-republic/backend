package entities

type (
	AchievementType struct {
		AchievementTypeID int    `gorm:"column:achievement_type_id;primaryKey;autoIncrement"`
		Name              string `gorm:"column:name;type:character varying(255);not null"`

		Events []Event `gorm:"foreignKey:AchievementTypeID"`
	}
)

func (AchievementType) TableName() string {
	return "achievement_type"
}
