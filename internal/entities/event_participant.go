package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	EventParticipant struct {
		ID               uint      `gorm:"primaryKey"`
		EventID          uint      `gorm:"uniqueIndex:idx_event_user;not null"`
		UserID           uint      `gorm:"uniqueIndex:idx_event_user;not null"`
		RegistrationDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		Status           string    `gorm:"size:20;default:registered"`
		AttendanceTime   *time.Time
		CreatedAt        time.Time
		UpdatedAt        time.Time
		DeletedAt        gorm.DeletedAt `gorm:"index"`

		// Relationships
		Event Event `gorm:"foreignKey:EventID"`
		User  User  `gorm:"foreignKey:UserID"`
	}
)
