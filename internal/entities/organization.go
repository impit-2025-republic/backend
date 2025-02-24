package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	Organization struct {
		ID          uint      `gorm:"primaryKey"`
		Name        string    `gorm:"size:200;not null"`
		Description string    `gorm:"type:text"`
		LogoURL     string    `gorm:"size:255"`
		Website     string    `gorm:"size:255"`
		CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt   time.Time
		DeletedAt   gorm.DeletedAt `gorm:"index"`

		// Relationships
		Users                    []User                    `gorm:"foreignKey:OrganizationID"`
		OrganizationAchievements []OrganizationAchievement `gorm:"foreignKey:OrganizationID"`
	}
)
