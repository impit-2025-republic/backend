package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	Event struct {
		ID              uint      `gorm:"primaryKey"`
		Title           string    `gorm:"size:200;not null"`
		Description     string    `gorm:"type:text"`
		EventType       string    `gorm:"size:50;not null"`
		StartDate       time.Time `gorm:"not null"`
		EndDate         time.Time `gorm:"not null"`
		Location        string    `gorm:"size:255"`
		CreatorID       *uint     `gorm:"index"`
		Points          *int
		MaxParticipants int
		Status          string    `gorm:"size:20;default:active"`
		CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt       time.Time
		DeletedAt       gorm.DeletedAt `gorm:"index"`

		Creator           *User              `gorm:"foreignKey:CreatorID"`
		EventParticipants []EventParticipant `gorm:"foreignKey:EventID"`
		TaskCompletions   []TaskCompletion   `gorm:"foreignKey:EventID"`
		Categories        []EventCategory    `gorm:"many2many:event_category_mappings"`
	}
)
