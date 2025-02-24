package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	TaskCompletion struct {
		ID             uint      `gorm:"primaryKey"`
		TaskID         uint      `gorm:"uniqueIndex:idx_task_user_event;not null"`
		UserID         uint      `gorm:"uniqueIndex:idx_task_user_event;not null"`
		EventID        uint      `gorm:"uniqueIndex:idx_task_user_event;not null"`
		CompletionTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		PointsEarned   int       `gorm:"default:0"`
		Status         string    `gorm:"size:20;default:completed"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		DeletedAt      gorm.DeletedAt `gorm:"index"`

		// Relationships
		Task  EventTask `gorm:"foreignKey:TaskID"`
		User  User      `gorm:"foreignKey:UserID"`
		Event Event     `gorm:"foreignKey:EventID"`
	}
)
