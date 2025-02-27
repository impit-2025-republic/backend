package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type caseProductProbabilityRepo struct {
	db *gorm.DB
}

func NewCaseProductProbabilityRepo(db *gorm.DB) entities.CaseProductProbabilityRepo {
	return caseProductProbabilityRepo{db: db}
}

func (r caseProductProbabilityRepo) GetAll(caseTypeID uint) ([]entities.CaseProductProbability, error) {
	var products []entities.CaseProductProbability
	err := r.db.Where("case_type_id = ?", caseTypeID).Find(&products).Error
	if err != nil {
		return []entities.CaseProductProbability{}, err
	}
	return products, nil
}
