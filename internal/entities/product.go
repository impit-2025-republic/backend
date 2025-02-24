package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	Product struct {
		ID                uint      `gorm:"primaryKey"`
		Name              string    `gorm:"size:200;not null"`
		Description       string    `gorm:"type:text"`
		Price             float64   `gorm:"type:decimal(10,2);not null"`
		AvailableQuantity int       `gorm:"default:0"`
		CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt         time.Time
		DeletedAt         gorm.DeletedAt `gorm:"index"`

		// Relationships
		Purchases []Purchase `gorm:"foreignKey:ProductID"`
	}
)
