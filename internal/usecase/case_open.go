package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type (
	CaseOpenUseCase interface {
		Execute(ctx context.Context, input CaseOpenInput) (CaseOpenOutput, error)
	}

	CaseOpenInput struct {
		ProductID uint `json:"productId"`
		UserID    int  `json:"-"`
	}

	CaseOpenOutput struct {
		ProductID uint `jsom:"productId"`
	}

	caseOpenInteractor struct {
		caseProductProbabilityRepo entities.CaseProductProbabilityRepo
		productRepo                entities.ProductRepo
		userWalletRepo             entities.UserWalletRepo
		userHistoryWalletRepo      entities.UserWalletHistoryRepo
		userWinningRepo            entities.UserWinningRepo
	}
)

func NewCaseOpenInteractor(
	caseProductProbabilityRepo entities.CaseProductProbabilityRepo,
	productRepo entities.ProductRepo,
	userWalletRepo entities.UserWalletRepo,
	userHistoryWalletRepo entities.UserWalletHistoryRepo,
	userWinningRepo entities.UserWinningRepo,
) CaseOpenUseCase {
	return caseOpenInteractor{
		caseProductProbabilityRepo: caseProductProbabilityRepo,
		productRepo:                productRepo,
		userWalletRepo:             userWalletRepo,
		userHistoryWalletRepo:      userHistoryWalletRepo,
		userWinningRepo:            userWinningRepo,
	}
}

func (uc caseOpenInteractor) Execute(ctx context.Context, input CaseOpenInput) (CaseOpenOutput, error) {
	product, err := uc.productRepo.GetByID(input.ProductID)
	if err != nil {
		return CaseOpenOutput{}, err
	}

	if product.CaseTypeID == nil {
		return CaseOpenOutput{}, errors.New("is_not_case")
	}

	wallet, err := uc.userWalletRepo.GetWallet(uint(input.UserID))
	if err != nil {
		return CaseOpenOutput{}, err
	}

	if wallet.Price < product.Price {
		return CaseOpenOutput{}, errors.New("not_enoungh_coin")
	}

	products, err := uc.caseProductProbabilityRepo.GetAll(uint(*product.CaseTypeID))
	if err != nil {
		return CaseOpenOutput{}, err
	}

	productId, err := GetRandomProductByProbability(products)

	if err != nil {
		return CaseOpenOutput{}, err
	}

	err = uc.userWalletRepo.DownBalance([]int{input.UserID}, product.Price)
	if err != nil {
		return CaseOpenOutput{}, err
	}

	err = uc.userHistoryWalletRepo.Create(
		entities.UserWalletHistory{
			UserID:      input.UserID,
			Coin:        product.Price,
			RefillType:  "minus",
			Description: fmt.Sprintf("Открытие кейса %s. Отняли %.2f", product.Name, product.Price),
		},
	)

	if err != nil {
		return CaseOpenOutput{}, err
	}

	uc.userWinningRepo.Create(entities.UserWinning{
		UserID:    wallet.UserID,
		ProductID: int(productId),
		WonAt:     time.Now(),
		WinType:   "case",
	})

	return CaseOpenOutput{
		ProductID: productId,
	}, nil
}

func GetRandomProductByProbability(probabilities []entities.CaseProductProbability) (uint, error) {
	var filteredProbs []entities.CaseProductProbability
	filteredProbs = append(filteredProbs, probabilities...)

	if len(filteredProbs) == 0 {
		return 0, fmt.Errorf("no products found for case type ID")
	}

	var totalProb float64
	for _, prob := range filteredProbs {
		totalProb += prob.DropProbability
	}

	if totalProb < 99.5 || totalProb > 100.5 {
		return 0, fmt.Errorf("sum of probabilities should be 100, but got: %.2f", totalProb)
	}

	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Float64() * 100

	var cumulativeProb float64
	for _, prob := range filteredProbs {
		cumulativeProb += prob.DropProbability
		if randomValue <= cumulativeProb {
			return prob.ProductID, nil
		}
	}

	return filteredProbs[len(filteredProbs)-1].ProductID, nil
}
