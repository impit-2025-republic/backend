package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
)

type (
	ClosedEventsUseCase interface {
		Execute(ctx context.Context) (ClosedEventsOutput, error)
	}

	ClosedEventsOutput struct {
		Events []entities.Event `json:"events"`
	}

	closedEventsInteractor struct {
		eventRepo entities.EventRepo
	}
)

func NewClosedEventsInteractor(
	eventRepo entities.EventRepo,
) ClosedEventsUseCase {
	return closedEventsInteractor{
		eventRepo: eventRepo,
	}
}

func (uc closedEventsInteractor) Execute(ctx context.Context) (ClosedEventsOutput, error) {
	events, err := uc.eventRepo.GetClosedEvents()
	if err != nil {
		return ClosedEventsOutput{}, err
	}

	return ClosedEventsOutput{
		Events: events,
	}, nil
}
