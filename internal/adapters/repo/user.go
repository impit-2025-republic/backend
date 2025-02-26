package repo

import (
	"b8boost/backend/internal/entities"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) entities.UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetByUID(uid string) (entities.User, error) {
	var user entities.User
	if err := r.db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepo) GetByID(id uint) (entities.User, error) {
	var user entities.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepo) GetAll() ([]entities.User, error) {
	var users []entities.User
	if err := r.db.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (r *userRepo) Create(user entities.User) error {
	return r.db.Create(&user).Error
}

// Update implements entities.UserRepo.
func (r *userRepo) Update(user entities.User) error {
	return r.db.Save(&user).Error
}
