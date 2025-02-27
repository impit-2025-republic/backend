package usecase

import (
	"b8boost/backend/internal/entities"
	"b8boost/backend/internal/infra/jwt"
	"b8boost/backend/internal/infra/ldap"
	"context"
	"errors"
	"fmt"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type (
	LoginUsecase interface {
		Execute(ctx context.Context, input LoginInput) (LoginOutput, error)
	}

	LoginInput struct {
		InitData string
	}

	LoginOutput struct {
		Token string `json:"token"`
	}

	loginInteractor struct {
		botToken string
		jwt      jwt.JWKSHandler
		ldap     ldap.LDAP
		userRepo entities.UserRepo
	}
)

func NewLoginInteractor(
	botToken string,
	jwt jwt.JWKSHandler,
	ldap ldap.LDAP,
	userRepo entities.UserRepo,
) LoginUsecase {
	return loginInteractor{
		botToken: botToken,
		jwt:      jwt,
		ldap:     ldap,
		userRepo: userRepo,
	}
}

func (uc loginInteractor) Execute(ctx context.Context, input LoginInput) (LoginOutput, error) {
	expIn := 24 * time.Hour
	if err := initdata.Validate(input.InitData, uc.botToken, expIn); err != nil {
		return LoginOutput{}, errors.New("invalid_data")
	}

	data, err := initdata.Parse(input.InitData)
	if err != nil {
		return LoginOutput{}, errors.New("invalid_data")
	}

	userData := uc.ldap.FindTelegramID(data.User.ID)

	var entryUUID string = "0"
	if userData != nil {
		entryUUID := ldap.GetFirstValueOrDefault(userData, "entryUUID", "")

		if entryUUID == "" {
			return LoginOutput{}, errors.New("user_not_found1")
		}
	}

	user, err := uc.userRepo.GetByLdapID(entryUUID)
	if err != nil {
		fmt.Println(err)
		return LoginOutput{}, errors.New("user_not_found2")
	}

	token, err := uc.jwt.Generate(fmt.Sprintf("%d", user.UserID))
	if err != nil {
		return LoginOutput{}, errors.New("token_generation_error")
	}

	return LoginOutput{Token: token}, nil
}
