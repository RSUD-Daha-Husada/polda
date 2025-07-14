package repository

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"github.com/google/uuid"
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

func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, "user_id = ?", id).Error
	return &user, err
}

func (r *UserRepository) Update(id uuid.UUID, data map[string]interface{}) error {
	return r.DB.Model(&model.User{}).Where("user_id = ?", id).Updates(data).Error
}