package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type eventUserVisits struct {
	db *gorm.DB
}

func NewEventUserVisits(db *gorm.DB) entities.EventUserVisitRepo {
	return eventUserVisits{
		db: db,
	}
}

func (r eventUserVisits) Create(event entities.EventUserVisit) error {
	return r.db.Create(&event).Error
}

func (r eventUserVisits) GetByEventIDAndUserID(eventID uint, userID uint) (entities.EventUserVisit, error) {
	var event entities.EventUserVisit
	err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&event).Error
	if err != nil {
		return entities.EventUserVisit{}, err
	}
	return event, nil
}

func (r eventUserVisits) GetByAchievemenTypeIDAndUserIDAndVisited(achievementTypeID int, userID int) ([]entities.EventUserVisit, error) {
	var events []entities.EventUserVisit
	err := r.db.Where("achievement_type_id = ? AND user_id = ?", achievementTypeID, userID).Find(&events).Error
	if err != nil {
		return []entities.EventUserVisit{}, err
	}
	return events, nil
}

func (r eventUserVisits) GetByEventIDAndVisit(eventID int) ([]entities.EventUserVisit, error) {
	var events []entities.EventUserVisit
	err := r.db.Where("event_id = ? AND visit = 'attended'", eventID).Find(&events).Error
	if err != nil {
		return []entities.EventUserVisit{}, err
	}
	return events, nil
}

func (r eventUserVisits) GetByUserID(userID uint) ([]entities.EventUserVisit, error) {
	var events []entities.EventUserVisit
	err := r.db.Where("user_id = ?", userID).Find(&events).Error
	if err != nil {
		return []entities.EventUserVisit{}, err
	}
	return events, nil
}
