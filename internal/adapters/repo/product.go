package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) entities.ProductRepo {
	return productRepo{db: db}
}

func (r productRepo) GetAll() ([]entities.Product, error) {
	var products []entities.Product
	err := r.db.Find(&products).Error
	if err != nil {
		return []entities.Product{}, err
	}
	return products, nil
}
