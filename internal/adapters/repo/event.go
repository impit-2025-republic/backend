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

func (r eventRepo) GetUpcomingEvents(period *string) ([]entities.Event, error) {
	var events []entities.Event
	if period != nil {
		per := *period
		switch per {
		case "today":
			today := time.Now().Truncate(24 * time.Hour)
			tomorrow := today.AddDate(0, 0, 1)

			err := r.db.Where("start_ds >= ? AND start_ds < ?", today, tomorrow).Find(&events).Error
			if err != nil {
				return nil, err
			}
		case "tomorrow":
			tomorrow := time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1)
			dayAfterTomorrow := tomorrow.AddDate(0, 0, 1)

			err := r.db.Where("start_ds >= ? AND start_ds < ?", tomorrow, dayAfterTomorrow).Find(&events).Error
			if err != nil {
				return nil, err
			}
		case "week":
			now := time.Now()
			today := now.Truncate(24 * time.Hour)

			daysUntilMonday := int((8 - int(now.Weekday())) % 7)
			if daysUntilMonday == 0 {
				daysUntilMonday = 7
			}

			nextWeekStart := today.AddDate(0, 0, daysUntilMonday)

			err := r.db.Where("start_ds >= ? AND start_ds < ?", today, nextWeekStart).Find(&events).Error
			if err != nil {
				return nil, err
			}
		case "month":
			now := time.Now()

			today := now.Truncate(24 * time.Hour)

			nextMonthStart := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

			err := r.db.Where("start_ds >= ? AND start_ds < ?", today, nextMonthStart).Find(&events).Error
			if err != nil {
				return nil, err
			}
		}
	} else {
		err := r.db.Where("start_ds BETWEEN ? AND ?", time.Now(), time.Now().AddDate(0, 0, 5)).Find(&events).Error
		if err != nil {
			return nil, err
		}
	}
	return events, nil
}

func (r eventRepo) GetClosedEvents() ([]entities.Event, error) {
	var events []entities.Event
	err := r.db.Where("status = 'closed'").Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r eventRepo) GetAllEventsOpenAndRunning() ([]entities.Event, error) {
	var events []entities.Event
	err := r.db.Where("status = 'runnig' OR status = 'open'").Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r eventRepo) UpdateMany(events []entities.Event) error {
	db := r.db

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, event := range events {
		result := tx.Model(&entities.Event{}).
			Where("event_id = ?", event.EventID).
			Updates(map[string]interface{}{
				"event_name":          event.EventName,
				"description":         event.Description,
				"title":               event.Title,
				"start_ds":            event.StartDs,
				"end_ds":              event.EndDs,
				"status":              event.Status,
				"event_type":          event.EventType,
				"max_users":           event.MaxUsers,
				"coin":                event.Coin,
				"achievement_type_id": event.AchievementTypeID,
				"company_id":          event.CompanyID,
			})

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r eventRepo) GetByID(id int) (entities.Event, error) {
	var event entities.Event
	err := r.db.Where("event_id = ?", id).Find(&event).Error
	if err != nil {
		return entities.Event{}, err
	}
	return event, nil
}
