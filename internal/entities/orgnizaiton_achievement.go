package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	OrganizationAchievement struct {
		ID             uint      `gorm:"primaryKey"`
		OrganizationID uint      `gorm:"uniqueIndex:idx_org_achievement;not null"`
		AchievementID  uint      `gorm:"uniqueIndex:idx_org_achievement;not null"`
		EarnedDate     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		DeletedAt      gorm.DeletedAt `gorm:"index"`

		// Relationships
		Organization Organization      `gorm:"foreignKey:OrganizationID"`
		Achievement  GlobalAchievement `gorm:"foreignKey:AchievementID"`
	}
)
