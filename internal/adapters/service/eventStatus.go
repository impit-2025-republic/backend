package service

import (
	"b8boost/backend/internal/entities"
	"b8boost/backend/internal/infra/tgbot"
	"fmt"
	"time"
)

type EventStatusService struct {
	eventRepo             entities.EventRepo
	eventUserVisitRepo    entities.EventUserVisitRepo
	userWalletRepo        entities.UserWalletRepo
	userWalletHistoryRepo entities.UserWalletHistoryRepo
	tgbot                 tgbot.TgBot
}

const ()

func NewEventStatusService(
	eventRepo entities.EventRepo,
	eventUserVisitRepo entities.EventUserVisitRepo,
	userWalletRepo entities.UserWalletRepo,
	userWalletHistoryRepo entities.UserWalletHistoryRepo,
	tgbot tgbot.TgBot,
) EventStatusService {
	return EventStatusService{
		eventRepo:             eventRepo,
		eventUserVisitRepo:    eventUserVisitRepo,
		userWalletHistoryRepo: userWalletHistoryRepo,
		userWalletRepo:        userWalletRepo,
		tgbot:                 tgbot,
	}
}

func (s EventStatusService) Start() {
	events, err := s.eventRepo.GetAllEventsOpenAndRunning()
	if err != nil {
		return
	}

	now := time.Now()

	for _, event := range events {
		if event.StartDs != nil && event.Status != nil {
			eventStartDs := *event.StartDs
			eventEndDs := *event.EndDs
			fmt.Println("Running")
			if *event.Status == entities.EventStatusOpen && now.Compare(eventStartDs) <= 0 {
				fmt.Println("Status Running")
				newStatus := entities.EventStatusRunning
				event.Status = &newStatus
			}

			if *event.Status == entities.EventStatusRunning && now.Compare(eventEndDs) <= 0 {
				fmt.Println("Status closed")
				newStatus := entities.EventStatusClosed
				event.Status = &newStatus

				users, err := s.eventUserVisitRepo.GetByEventIDAndVisit(int(event.EventID))
				if err != nil {
					continue
				}

				var userIds []int
				for _, user := range users {
					s.tgbot.SendMessage(int64(user.UserID), fmt.Sprintf("За посещение мероприятия %s. Вам начислили %.2f", event.Title, event.Coin))
					s.userWalletHistoryRepo.Create(
						entities.UserWalletHistory{
							UserID:      int(user.UserID),
							Coin:        event.Coin,
							RefillType:  "plus",
							Description: fmt.Sprintf("За посещение мероприятия %s. Вам начислили %.2f", event.Title, event.Coin),
						},
					)
					userIds = append(userIds, user.UserID)
				}

				s.userWalletRepo.UpBalance(userIds, event.Coin)
			}
		}
	}

	err = s.eventRepo.UpdateMany(events)
	if err != nil {
		fmt.Println(err)
	}
}
