package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&entities.User{},
		&entities.Event{},
		&entities.EventParticipant{},
		&entities.TaskCompletion{},
		&entities.UserAchievement{},
		&entities.EventCategory{},
		&entities.Product{},
		&entities.Purchase{},
	)
}
