package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	UserAchievement struct {
		ID             uint      `gorm:"primaryKey"`
		UserID         uint      `gorm:"index;not null"`
		TotalPoints    int       `gorm:"default:0"`
		EventsAttended int       `gorm:"default:0"`
		TasksCompleted int       `gorm:"default:0"`
		LastUpdated    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		DeletedAt      gorm.DeletedAt `gorm:"index"`

		// Relationships
		User User `gorm:"foreignKey:UserID"`
	}
)
