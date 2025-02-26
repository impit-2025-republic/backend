package entities

import (
	"time"
)

const (
	EventStatusOpen    = "open"
	EventStatusClosed  = "closed"
	EventStatusRunning = "running"
)

type (
	EventRepo interface {
		GetUpcomingEvents(period *string) ([]Event, error)
		GetClosedEvents() ([]Event, error)
		UpdateMany(events []Event) error
		GetByEventsIds(eventIds []int) ([]Event, error)
		GetAllEventsOpenAndRunning() ([]Event, error)
		GetByID(id int) (Event, error)
		FindByErpID(erpId int) (Event, error)
		Update(event Event) error
		Create(event Event) error
	}

	Event struct {
		EventID           int64      `gorm:"column:event_id;primaryKey;autoIncrement"`
		EventName         string     `gorm:"column:event_name;type:text"`
		CreatedAt         time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
		Description       string     `gorm:"column:description;type:text"`
		Title             string     `gorm:"column:title;type:text"`
		StartDs           *time.Time `gorm:"column:start_ds"`
		EndDs             *time.Time `gorm:"column:end_ds"`
		Status            *string    `gorm:"column:status"`
		EventType         *string    `gorm:"column:event_type"`
		MaxUsers          *int       `gorm:"column:max_users"`
		Coin              float64    `gorm:"column:coin;type:numeric(10,2)"`
		AchievementTypeID *int       `gorm:"column:achievement_type_id"`
		CompanyID         *int       `gorm:"column:company_id"`
		ErpID             *int       `gorm:"column:erp_id"`

		AchievementType *AchievementType `gorm:"foreignKey:AchievementTypeID"`
		Company         *Company         `gorm:"foreignKey:CompanyID"`
	}
)

func (Event) TableName() string {
	return "events"
}
