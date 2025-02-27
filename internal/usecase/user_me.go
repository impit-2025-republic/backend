package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
	"time"
)

type (
	UserMeUseCase interface {
		Execute(ctx context.Context, input UserMeInput) (UserMeOutput, error)
	}

	UserMeInput struct {
		UserID int
	}

	UserWinnings struct {
		UserWinns entities.UserWinning `json:"buy"`
		Product   entities.Product     `json:"product"`
	}

	UserMeOutput struct {
		UserID       int              `json:"user_id"`
		Surname      *string          `json:"surname"`
		Name         *string          `json:"name"`
		LastSurname  *string          `json:"l_surname"`
		BirthDate    *time.Time       `json:"birth_date"`
		Email        *string          `json:"email"`
		Phone        *string          `json:"phone"`
		Events       []entities.Event `json:"events"`
		Coin         float64          `json:"coin"`
		UserWinnings []UserWinnings   `json:"buys"`
	}

	userMeInteractor struct {
		userRepo           entities.UserRepo
		userWalletRepo     entities.UserWalletRepo
		eventUserVisitRepo entities.EventUserVisitRepo
		eventRepo          entities.EventRepo
		userWinningRepo    entities.UserWinningRepo
	}
)

func NewUserMeInteractor(
	userRepo entities.UserRepo,
	userWalletRepo entities.UserWalletRepo,
	eventUserVisitRepo entities.EventUserVisitRepo,
	eventRepo entities.EventRepo,
	userWinningRepo entities.UserWinningRepo,
) UserMeUseCase {
	return userMeInteractor{
		userRepo:           userRepo,
		userWalletRepo:     userWalletRepo,
		eventUserVisitRepo: eventUserVisitRepo,
		eventRepo:          eventRepo,
		userWinningRepo:    userWinningRepo,
	}
}

func (uc userMeInteractor) Execute(ctx context.Context, input UserMeInput) (UserMeOutput, error) {
	user, err := uc.userRepo.GetByID(uint(input.UserID))
	if err != nil {
		return UserMeOutput{}, err
	}

	wallet, err := uc.userWalletRepo.GetWallet(user.UserID)
	if err != nil {
		return UserMeOutput{}, err
	}
	eventsUserVisit, err := uc.eventUserVisitRepo.GetByUserID(user.UserID)
	if err != nil {
		return UserMeOutput{}, err
	}

	var eventIds []int
	for _, eventsUser := range eventsUserVisit {
		eventIds = append(eventIds, eventsUser.EventID)
	}

	events, err := uc.eventRepo.GetByEventsIds(eventIds)
	if err != nil {
		return UserMeOutput{}, err
	}

	userWinns, err := uc.userWinningRepo.GetMyWinnings(uint(input.UserID))
	if err != nil {
		return UserMeOutput{}, err
	}

	var userWinnings []UserWinnings
	for _, uw := range userWinns {
		userWinnings = append(userWinnings, UserWinnings{
			UserWinns: uw.UserWinning,
			Product:   uw.Product,
		})
	}

	return UserMeOutput{
		UserID:       int(user.UserID),
		Surname:      user.Surname,
		Name:         user.Name,
		LastSurname:  user.LastSurname,
		BirthDate:    user.BirthDate,
		Email:        user.Email,
		Phone:        user.Phone,
		Events:       events,
		Coin:         wallet.Price,
		UserWinnings: userWinnings,
	}, nil
}
