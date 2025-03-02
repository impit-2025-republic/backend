package usecase

import (
	"b8boost/backend/internal/entities"
	"b8boost/backend/internal/infra/tgbot"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type (
	VisitEventUseCase interface {
		Execute(ctx context.Context, input VisitEventInput) error
	}

	VisitEventInput struct {
		EventID int `json:"eventID"`
		UserID  int `json:"-"`
	}

	visitEventInteractor struct {
		eventRepo           entities.EventRepo
		eventUserVisitsRepo entities.EventUserVisitRepo
		userRepo            entities.UserRepo
		tgbot               tgbot.TgBot
	}
)

func NewVisitEventInteractor(
	eventRepo entities.EventRepo,
	eventUserVisitsRepo entities.EventUserVisitRepo,
	userRepo entities.UserRepo,
	tgbot tgbot.TgBot,
) VisitEventUseCase {
	return visitEventInteractor{
		eventRepo:           eventRepo,
		eventUserVisitsRepo: eventUserVisitsRepo,
		userRepo:            userRepo,
		tgbot:               tgbot,
	}
}

func (uc visitEventInteractor) Execute(ctx context.Context, input VisitEventInput) error {
	event, err := uc.eventRepo.GetByID(input.EventID)
	if err != nil {
		return err
	}

	if event.Status != nil && *event.Status == entities.EventStatusClosed {
		return errors.New("event_is_closed")
	}

	eventUserVisit, err := uc.eventUserVisitsRepo.GetByEventIDAndUserID(uint(input.EventID), uint(input.UserID))
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if eventUserVisit.EventID != 0 {
		return errors.New("event_is_visit")
	}

	uc.tgbot.SendMessage(int64(input.UserID), fmt.Sprintf("Вы записались на %s", event.Title))

	eventUserVisit.EventID = int(event.EventID)
	eventUserVisit.UserID = input.UserID
	eventUserVisit.Visit = "signed"
	eventUserVisit.AchievementTypeID = 1

	err = uc.eventUserVisitsRepo.Create(eventUserVisit)
	if err != nil {
		return err
	}
	return nil
}
