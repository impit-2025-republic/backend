package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
)

type (
	TopBalanceUseCase interface {
		Execute(ctx context.Context) (TopBalanceOutput, error)
	}

	TopBalance struct {
		UserWallet entities.UserWallet `json:"wallet"`
		User       entities.User       `json:"user"`
	}

	TopBalanceOutput struct {
		Wallets []TopBalance `json:"wallets"`
	}

	topBalanceInteractor struct {
		userWalletRepo entities.UserWalletRepo
	}
)

func NewTopBalanceInteractor(
	userWalletRepo entities.UserWalletRepo,
) TopBalanceUseCase {
	return topBalanceInteractor{
		userWalletRepo: userWalletRepo,
	}
}

func (uc topBalanceInteractor) Execute(ctx context.Context) (TopBalanceOutput, error) {

	userWinns, err := uc.userWalletRepo.GetTopBalance()
	if err != nil {
		return TopBalanceOutput{}, err
	}

	var userWinnings []TopBalance
	for _, uw := range userWinns {
		userWinnings = append(userWinnings, TopBalance{
			UserWallet: uw.UserWallet,
			User:       uw.User,
		})
	}

	return TopBalanceOutput{
		Wallets: userWinnings,
	}, nil
}
