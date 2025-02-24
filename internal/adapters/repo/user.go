package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetByUID(uid string) (entities.User, error) {
	var user entities.User
	if err := r.db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
