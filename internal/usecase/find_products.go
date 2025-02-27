package usecase

import (
	"b8boost/backend/internal/entities"
	"context"
)

type (
	FindProductUseCase interface {
		Execute(ctx context.Context) (FindProductOutput, error)
	}

	FindProductOutput struct {
		Products []entities.Product `json:"products"`
	}

	findProductInteractor struct {
		productRepo entities.ProductRepo
	}
)

func NewFindProductInteractor(
	productRepo entities.ProductRepo,
) FindProductUseCase {
	return findProductInteractor{
		productRepo: productRepo,
	}
}

func (uc findProductInteractor) Execute(ctx context.Context) (FindProductOutput, error) {
	products, err := uc.productRepo.GetAll()
	if err != nil {
		return FindProductOutput{}, err
	}

	return FindProductOutput{
		Products: products,
	}, nil
}
