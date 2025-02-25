package usecase

import (
	"b8boost/backend/internal/infra/jwt"
	"context"
)

type (
	FindJwtsUseCase interface {
		Execute(ctx context.Context) FindJwtsUseCaseOutput
	}

	FindJwtsUseCaseOutput map[string]interface{}

	FindJwtsInteractor struct {
		jwt jwt.JWKSHandler
	}
)

func NewFindJwtsInteractor(jwt jwt.JWKSHandler) FindJwtsUseCase {
	return &FindJwtsInteractor{
		jwt: jwt,
	}
}

func (uc FindJwtsInteractor) Execute(ctx context.Context) FindJwtsUseCaseOutput {
	jwts := uc.jwt.Validate()
	return FindJwtsUseCaseOutput(jwts)
}
