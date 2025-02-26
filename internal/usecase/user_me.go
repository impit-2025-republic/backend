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

	UserMeOutput struct {
		UserID      int              `json:"user_id"`
		Surname     *string          `json:"surname"`
		Name        *string          `json:"name"`
		LastSurname *string          `json:"l_surname"`
		BirthDate   *time.Time       `json:"birth_date"`
		Email       *string          `json:"email"`
		Phone       *string          `json:"phone"`
		Events      []entities.Event `json:"events"`
		Coin        int              `json:"coin"`
	}

	userMeInteractor struct {
		userRepo entities.UserRepo
	}
)

func NewUserMeInteractor(
	userRepo entities.UserRepo,
) UserMeUseCase {
	return userMeInteractor{
		userRepo: userRepo,
	}
}

func (uc userMeInteractor) Execute(ctx context.Context, input UserMeInput) (UserMeOutput, error) {
	user, err := uc.userRepo.GetByID(uint(input.UserID))
	if err != nil {
		return UserMeOutput{}, err
	}

	return UserMeOutput{
		UserID:      int(user.UserID),
		Surname:     user.Surname,
		Name:        user.Name,
		LastSurname: user.LastSurname,
		BirthDate:   user.BirthDate,
		Email:       user.Email,
		Phone:       user.Phone,
		Coin:        1,
	}, nil
}
