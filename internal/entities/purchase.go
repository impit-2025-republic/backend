package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	Purchase struct {
		ID           uint      `gorm:"primaryKey"`
		UserID       uint      `gorm:"index;not null"`
		ProductID    uint      `gorm:"index;not null"`
		Quantity     int       `gorm:"not null"`
		TotalPrice   float64   `gorm:"type:decimal(10,2);not null"`
		PurchaseDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		Status       string    `gorm:"size:20;default:pending"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    gorm.DeletedAt `gorm:"index"`

		// Relationships
		User    User    `gorm:"foreignKey:UserID"`
		Product Product `gorm:"foreignKey:ProductID"`
	}
)
