package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
	"fmt"
)

type (
	UpcomingEventsUseCase interface {
		Execute(ctx context.Context, input UpcomingEventInput) (UpcomingEventList, error)
	}

	UpcomingEventInput struct {
		Period *string `form:"period" validate:"oneof=today tomorrow week month"`
		UserID *int
	}

	EventWithRegistration struct {
		entities.Event
		IsRegistered bool `json:"is_registered"`
	}

	UpcomingEventList struct {
		Events []EventWithRegistration `json:"events"`
		Total  int                     `json:"total"`
	}

	upcomingEventsInteractor struct {
		eventsRepo      entities.EventRepo
		eventUserVisits entities.EventUserVisitRepo
	}
)

func NewUpcomingEventsInteractor(
	eventsRepo entities.EventRepo,
	eventUserVisits entities.EventUserVisitRepo,
) UpcomingEventsUseCase {
	return upcomingEventsInteractor{
		eventsRepo:      eventsRepo,
		eventUserVisits: eventUserVisits,
	}
}

func (uc upcomingEventsInteractor) Execute(ctx context.Context, input UpcomingEventInput) (UpcomingEventList, error) {
	events, err := uc.eventsRepo.GetUpcomingEvents(input.Period)
	if err != nil {
		fmt.Println(err)
		return UpcomingEventList{}, err
	}

	eventsWithRegistration := make([]EventWithRegistration, 0, len(events))

	if input.UserID != nil {

		userVisits, err := uc.eventUserVisits.GetByUserID(uint(*input.UserID))
		if err != nil {
			return UpcomingEventList{}, err
		}

		registrationMap := make(map[uint]bool)
		for _, visit := range userVisits {
			registrationMap[uint(visit.EventID)] = true
		}

		for _, event := range events {
			isRegistered := registrationMap[uint(event.EventID)]
			eventsWithRegistration = append(eventsWithRegistration, EventWithRegistration{
				Event:        event,
				IsRegistered: isRegistered,
			})
		}
	} else {

		for _, event := range events {
			eventsWithRegistration = append(eventsWithRegistration, EventWithRegistration{
				Event:        event,
				IsRegistered: false,
			})
		}
	}

	return UpcomingEventList{
		Events: eventsWithRegistration,
		Total:  len(eventsWithRegistration),
	}, nil
}
