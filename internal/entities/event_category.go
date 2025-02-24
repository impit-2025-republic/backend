package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	EventCategory struct {
		ID          uint   `gorm:"primaryKey"`
		Name        string `gorm:"size:100;not null"`
		Description string `gorm:"type:text"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   gorm.DeletedAt `gorm:"index"`

		// Relationships
		Events []Event `gorm:"many2many:event_category_mappings"`
	}
)
