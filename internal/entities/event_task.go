package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	EventTask struct {
		ID          uint      `gorm:"primaryKey"`
		EventID     uint      `gorm:"index;not null"`
		Title       string    `gorm:"size:200;not null"`
		Description string    `gorm:"type:text"`
		TaskType    string    `gorm:"size:50;not null"`
		Points      int       `gorm:"default:0"`
		CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt   time.Time
		DeletedAt   gorm.DeletedAt `gorm:"index"`

		// Relationships
		Event           Event            `gorm:"foreignKey:EventID"`
		TaskCompletions []TaskCompletion `gorm:"foreignKey:TaskID"`
	}
)
