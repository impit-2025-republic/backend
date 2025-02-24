package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
	"time"
)

type (
	UpcomingEventsUseCase interface {
		Execute(ctx context.Context) ([]UpcomingEventList, error)
	}

	UpcomingEventOutput struct {
		ID          uint      `json:"id"`
		Points      *int      `json:"points"`
		Title       string    `json:"title"`
		StartDate   time.Time `json:"startDate"`
		EventStatus string    `json:"eventStatus"`
	}

	UpcomingEventList struct {
		Events []UpcomingEventOutput `json:"events"`
		Total  int                   `json:"total"`
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

func (uc upcomingEventsInteractor) Execute(ctx context.Context) ([]UpcomingEventList, error) {
	events, err := uc.eventsRepo.GetUpcomingEvents()
	if err != nil {
		return nil, err
	}

	var upcomingEvents []UpcomingEventOutput
	for _, event := range events {
		upcomingEvents = append(upcomingEvents, UpcomingEventOutput{
			ID:          event.ID,
			Points:      event.Points,
			Title:       event.Title,
			StartDate:   event.StartDate,
			EventStatus: event.EventStatus,
		})
	}

	return []UpcomingEventList{
		{
			Events: upcomingEvents,
			Total:  len(upcomingEvents),
		},
	}, nil
}
