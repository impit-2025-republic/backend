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
		Period *string `json:"period" validate:"oneof=today tomorrow week month"`
	}

	UpcomingEventList struct {
		Events []entities.Event `json:"events"`
		Total  int              `json:"total" `
	}

	upcomingEventsInteractor struct {
		eventsRepo entities.EventRepo
	}
)

func NewUpcomingEventsInteractor(eventsRepo entities.EventRepo) UpcomingEventsUseCase {
	return upcomingEventsInteractor{
		eventsRepo: eventsRepo,
	}
}

func (uc upcomingEventsInteractor) Execute(ctx context.Context, input UpcomingEventInput) (UpcomingEventList, error) {
	events, err := uc.eventsRepo.GetUpcomingEvents(input.Period)
	if err != nil {
		fmt.Println(err)
		return UpcomingEventList{}, err
	}

	return UpcomingEventList{
		Events: events,
		Total:  len(events),
	}, nil
}
