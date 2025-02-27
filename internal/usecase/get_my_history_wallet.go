package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
)

type (
	GetMyHistoryWalletUseCase interface {
		Execute(ctx context.Context, input GetMyHistoryWalletInput) (GetMyHistoryWalletOutput, error)
	}

	GetMyHistoryWalletInput struct {
		UserID int
	}

	GetMyHistoryWalletOutput struct {
		Transactions []entities.UserWalletHistory `json:"transactions"`
	}

	getMyHistoryWalletInteractor struct {
		repo entities.UserWalletHistoryRepo
	}
)

func NewGetMyHistoryWalletInteractor(
	repo entities.UserWalletHistoryRepo,
) GetMyHistoryWalletUseCase {
	return getMyHistoryWalletInteractor{
		repo: repo,
	}
}

func (uc getMyHistoryWalletInteractor) Execute(ctx context.Context, input GetMyHistoryWalletInput) (GetMyHistoryWalletOutput, error) {
	transactions, err := uc.repo.GetByUserID(input.UserID)
	if err != nil {
		return GetMyHistoryWalletOutput{}, err
	}

	return GetMyHistoryWalletOutput{
		Transactions: transactions,
	}, nil
}
