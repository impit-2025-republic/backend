package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
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
		eventUserVisit    entities.EventUserVisitRepo
		achivmentUserRepo entities.AchievementUserRepo
		achievementRepo   entities.AchievementRepo
	}
)

func NewAdminVisitEventInteractor(
	eventUserVisit entities.EventUserVisitRepo,
	achivmentUserRepo entities.AchievementUserRepo,
	achievementRepo entities.AchievementRepo,
) AdminVisitEventUseCase {
	return adminVisitEventInteractor{
		eventUserVisit:    eventUserVisit,
		achivmentUserRepo: achivmentUserRepo,
		achievementRepo:   achievementRepo,
	}
}

func (uc adminVisitEventInteractor) Execute(ctx context.Context, input AdminVisitEventInput) error {
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

			// TODO give coin
		}
	}

	return nil
}
