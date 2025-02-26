package service

import (
	"b8boost/backend/internal/entities"
	"time"
)

type EventStatusService struct {
	eventRepo entities.EventRepo
}

const ()

func NewEventStatusService(
	eventRepo entities.EventRepo,
) EventStatusService {
	return EventStatusService{
		eventRepo: eventRepo,
	}
}

func (s EventStatusService) Start() {
	events, err := s.eventRepo.GetClosedEvents()
	if err != nil {
		return
	}

	now := time.Now()

	for _, event := range events {
		if event.StartDs != nil && event.Status != nil {
			eventStartDs := *event.StartDs
			eventEndDs := *event.EndDs
			if *event.Status == entities.EventStatusOpen && !now.Before(eventStartDs) {
				newStatus := entities.EventStatusRunning
				event.Status = &newStatus
			}

			if *event.Status == entities.EventStatusRunning && !now.Before(eventEndDs) {
				newStatus := entities.EventStatusClosed
				event.Status = &newStatus
			}
		}
	}

	s.eventRepo.UpdateMany(events)
}
