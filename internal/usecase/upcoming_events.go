package usecase

import (
	"context"
	"time"
)

type (
	UpcomingEventsUseCase interface {
		Execute(ctx context.Context) ([]UpcomingEvent, error)
	}

	UpcomingEvent struct {
		ID        uint      `json:"id"`
		Points    int32     `json:"points"`
		Title     string    `json:"title"`
		StartDate time.Time `json:"startDate"`
		Desc
	}
)
