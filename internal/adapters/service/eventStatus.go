package service

import (
	"b8boost/backend/internal/entities"
	"time"
)

type EventStatusService struct {
	eventRepo          entities.EventRepo
	eventUserVisitRepo entities.EventUserVisitRepo
	userWalletRepo     entities.UserWalletRepo
}

const ()

func NewEventStatusService(
	eventRepo entities.EventRepo,
	eventUserVisitRepo entities.EventUserVisitRepo,
	userWalletRepo entities.UserWalletRepo,
) EventStatusService {
	return EventStatusService{
		eventRepo:          eventRepo,
		eventUserVisitRepo: eventUserVisitRepo,
		userWalletRepo:     userWalletRepo,
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

				users, err := s.eventUserVisitRepo.GetByEventIDAndVisit(int(event.EventID))
				if err != nil {
					continue
				}

				var userIds []int
				for _, user := range users {
					userIds = append(userIds, user.UserID)
				}

				s.userWalletRepo.UpBalance(userIds, event.Coin)
			}
		}
	}

	s.eventRepo.UpdateMany(events)
}
