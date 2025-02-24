package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	GlobalAchievement struct {
		ID              uint      `gorm:"primaryKey"`
		Name            string    `gorm:"size:100;not null"`
		Description     string    `gorm:"type:text"`
		AchievementType string    `gorm:"size:50;not null"`
		PointsValue     int       `gorm:"default:0"`
		BadgeIcon       string    `gorm:"size:255"`
		Requirements    string    `gorm:"type:text"`
		CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt       time.Time
		DeletedAt       gorm.DeletedAt `gorm:"index"`

		// Relationships
		Users         []User         `gorm:"many2many:user_global_achievements"`
		Organizations []Organization `gorm:"many2many:organization_achievements"`
	}
)
