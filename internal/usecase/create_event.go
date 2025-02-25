package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
	"time"
)

type (
	CreateEventUseCase interface {
		Execute(ctx context.Context, input CreateEventInput) (CreateEventOutput, error)
	}

	CreateEventInput struct {
		Title           string    `json:"title"`
		Description     string    `json:"description"`
		EventType       string    `json:"eventType"`
		StartDate       time.Time `json:"startDate"`
		EndDate         time.Time `json:"endDate"`
		Location        string    `json:"location"`
		Points          *int      `json:"points"`
		MaxParticipants int       `json:"maxParticipants"`
		UserID          uint      `json:"-"`
	}

	CreateEventOutput struct {
		ID uint `json:"id"`
	}

	CreateEventInteractor struct {
		eventRepo entities.EventRepo
		userRepo  entities.UserRepo
	}
)

func NewCreateEventInteractor(eventRepo entities.EventRepo, userRepo entities.UserRepo) CreateEventUseCase {
	return &CreateEventInteractor{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

func (uc CreateEventInteractor) Execute(ctx context.Context, input CreateEventInput) (CreateEventOutput, error) {
	var event entities.Event
	event.Title = input.Title
	event.Description = input.Description
	event.EventType = input.EventType
	event.StartDate = input.StartDate
	event.EndDate = input.EndDate
	event.Location = input.Location
	event.Points = input.Points
	event.MaxParticipants = input.MaxParticipants

	user, err := uc.userRepo.GetByID(input.UserID)
	if err != nil {
		return CreateEventOutput{}, err
	}
	event.Creator = &user

	err = uc.eventRepo.CreateEvent(event)
	if err != nil {
		return CreateEventOutput{}, err
	}

	return CreateEventOutput{ID: event.ID}, nil
}
