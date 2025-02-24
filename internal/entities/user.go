package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	UserRepo interface {
		GetByUID(uid string) (User, error)
	}
	User struct {
		ID         uint   `gorm:"primaryKey"`
		Uid        string `gorm:"uid"`
		TelegramID *int   `gorm:"unique;not null"`

		Role             string    `gorm:"size:20;not null"`
		RegistrationDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		LastLogin        *time.Time
		CreatedAt        time.Time
		OrganizationID   *uint `gorm:"index"`
		UpdatedAt        time.Time
		DeletedAt        gorm.DeletedAt `gorm:"index"`

		// Relationships
		Organization      *Organization      `gorm:"foreignKey:OrganizationID"`
		CreatedEvents     []Event            `gorm:"foreignKey:CreatorID"`
		EventParticipants []EventParticipant `gorm:"foreignKey:UserID"`
		TaskCompletions   []TaskCompletion   `gorm:"foreignKey:UserID"`
		Achievements      []UserAchievement  `gorm:"foreignKey:UserID"`
		Purchases         []Purchase         `gorm:"foreignKey:UserID"`
	}
)
