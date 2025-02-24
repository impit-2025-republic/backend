package entities

type (
	EventCategoryMapping struct {
		EventID    uint `gorm:"primaryKey"`
		CategoryID uint `gorm:"primaryKey"`
	}
)
