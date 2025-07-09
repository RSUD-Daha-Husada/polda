package repository

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByTelephone(telephone string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("telephone = ?", telephone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}