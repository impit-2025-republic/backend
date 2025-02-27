package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
	"errors"
	"fmt"
	"time"
)

type (
	BuyProductUseCase interface {
		Execute(ctx context.Context, input BuyProductInput) error
	}

	BuyProductInput struct {
		ProductID uint `json:"productId"`
		UserID    uint
	}

	buyProductInteractor struct {
		productRepo           entities.ProductRepo
		userWinningRepo       entities.UserWinningRepo
		userWalletRepo        entities.UserWalletRepo
		userWalletHistoryRepo entities.UserWalletHistoryRepo
	}
)

func NewBuyProductInteractor(
	productRepo entities.ProductRepo,
	userWinningRepo entities.UserWinningRepo,
	userWalletRepo entities.UserWalletRepo,
	userWalletHistoryRepo entities.UserWalletHistoryRepo,
) BuyProductUseCase {
	return buyProductInteractor{
		productRepo:           productRepo,
		userWinningRepo:       userWinningRepo,
		userWalletRepo:        userWalletRepo,
		userWalletHistoryRepo: userWalletHistoryRepo,
	}
}

func (uc buyProductInteractor) Execute(ctx context.Context, input BuyProductInput) error {
	product, err := uc.productRepo.GetByID(input.ProductID)
	if err != nil {
		return err
	}

	if product.CaseTypeID != nil {
		return errors.New("is_case")
	}

	wallet, err := uc.userWalletRepo.GetWallet(uint(input.UserID))
	if err != nil {
		return err
	}

	if wallet.Price < product.Price {
		return errors.New("not_enoungh_coin")
	}

	err = uc.userWalletRepo.DownBalance([]int{int(input.UserID)}, product.Price)
	if err != nil {
		return err
	}

	err = uc.userWalletHistoryRepo.Create(
		entities.UserWalletHistory{
			UserID:      int(input.UserID),
			Coin:        product.Price,
			RefillType:  "minus",
			Description: fmt.Sprintf("Покупка товара %s. Отняли %.2f", product.Name, product.Price),
		},
	)

	if err != nil {
		return err
	}

	uc.userWinningRepo.Create(entities.UserWinning{
		UserID:    wallet.UserID,
		ProductID: int(product.ProductID),
		WonAt:     time.Now(),
		WinType:   "case",
	})

	product.Availability = product.Availability - 1

	uc.productRepo.Update(product)

	return nil
}
