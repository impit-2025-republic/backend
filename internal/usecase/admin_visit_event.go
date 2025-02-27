package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
	"fmt"
	"time"
)

type (
	AdminVisitEventUseCase interface {
		Execute(ctx context.Context, input AdminVisitEventInput) error
	}

	AdminVisitEventInput struct {
		EventID           int `json:"eventID"`
		UserID            int `json:"userID"`
		AchievementTypeID int `json:"achievement_type_id"`
	}

	adminVisitEventInteractor struct {
		eventRepo             entities.EventRepo
		eventUserVisit        entities.EventUserVisitRepo
		achivmentUserRepo     entities.AchievementUserRepo
		achievementRepo       entities.AchievementRepo
		userWalletRepo        entities.UserWalletRepo
		userHistoryWalletRepo entities.UserWalletHistoryRepo
	}
)

func NewAdminVisitEventInteractor(
	eventRepo entities.EventRepo,
	eventUserVisit entities.EventUserVisitRepo,
	achivmentUserRepo entities.AchievementUserRepo,
	achievementRepo entities.AchievementRepo,
	userWalletRepo entities.UserWalletRepo,
	userHistoryWalletRepo entities.UserWalletHistoryRepo,
) AdminVisitEventUseCase {
	return adminVisitEventInteractor{
		eventRepo:             eventRepo,
		eventUserVisit:        eventUserVisit,
		achivmentUserRepo:     achivmentUserRepo,
		achievementRepo:       achievementRepo,
		userHistoryWalletRepo: userHistoryWalletRepo,
		userWalletRepo:        userWalletRepo,
	}
}

func (uc adminVisitEventInteractor) Execute(ctx context.Context, input AdminVisitEventInput) error {
	event, err := uc.eventRepo.GetByID(input.EventID)
	if err != nil {
		return err
	}

	if event.Coin != 0 {
		wallet, err := uc.userWalletRepo.GetWallet(uint(input.UserID))
		if err != nil {
			return err
		}
		uc.userWalletRepo.UpBalance([]int{wallet.UserID}, event.Coin)
		uc.userHistoryWalletRepo.Create(entities.UserWalletHistory{
			UserID:      input.UserID,
			Coin:        event.Coin,
			RefillType:  "plus",
			Description: fmt.Sprintf("За прохождение мероприятия %s. Вам дали %.2f", event.Title, event.Coin),
			CreatedAt:   time.Now(),
		})
	}

	achievementsUser, err := uc.achivmentUserRepo.GetAll(input.UserID)
	if err != nil {
		return err
	}

	mapAchievementIds := make(map[int]struct{})
	achievementIds := make([]int, 0)

	for _, achievementUser := range achievementsUser {
		_, exists := mapAchievementIds[achievementUser.AchievementID]
		if !exists {
			mapAchievementIds[achievementUser.AchievementID] = struct{}{}
		}
	}

	for key, _ := range mapAchievementIds {
		achievementIds = append(achievementIds, key)
	}

	achievements, err := uc.achievementRepo.GetByNotAchievementTypeIDsAndAchievementTypeIDAndEndDs(achievementIds, input.AchievementTypeID)
	if err != nil {
		return err
	}

	visits, err := uc.eventUserVisit.GetByAchievemenTypeIDAndUserIDAndVisited(input.AchievementTypeID, input.UserID)
	if err != nil {
		return err
	}

	countVisit := len(visits)

	for _, achievement := range achievements {
		if achievement.TreshholdValue != nil && countVisit >= *achievement.TreshholdValue {
			uc.achivmentUserRepo.Create(entities.AchievementUser{
				UserID:        input.UserID,
				AchievementID: int(achievement.AchievementID),
			})

			wallet, err := uc.userWalletRepo.GetWallet(uint(input.UserID))
			if err != nil {
				return err
			}
			uc.userWalletRepo.UpBalance([]int{wallet.UserID}, event.Coin)
			uc.userHistoryWalletRepo.Create(entities.UserWalletHistory{
				UserID:      input.UserID,
				Coin:        event.Coin,
				RefillType:  "plus",
				Description: fmt.Sprintf("Вы получили достижение %s. Вам дали %.2f", achievement.Name, event.Coin),
				CreatedAt:   time.Now(),
			})

		}
	}

	return nil
}
