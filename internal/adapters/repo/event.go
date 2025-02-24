package repo

import (
	"b8boost/backend/internal/entities"
	"time"

	"gorm.io/gorm"
)

type eventRepo struct {
	db *gorm.DB
}

func NewEventRepo(db *gorm.DB) entities.EventRepo {
	return eventRepo{
		db: db,
	}
}

func (r eventRepo) GetUpcomingEvents() ([]entities.Event, error) {
	var events []entities.Event
	err := r.db.Where("start_date BETWEEN ? AND ?", time.Now(), time.Now().AddDate(0, 0, 5)).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}
