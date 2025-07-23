package repository

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"gorm.io/gorm"
)

type UserAppRepository interface {
	FindAllWithRelations(db *gorm.DB) ([]model.UserApplication, error)
}

type userAppRepository struct{}

func NewUserAppRepository() UserAppRepository {
	return &userAppRepository{}
}

func (r *userAppRepository) FindAllWithRelations(db *gorm.DB) ([]model.UserApplication, error) {
	var userApps []model.UserApplication
	err := db.Preload("User").Preload("Application").Find(&userApps).Error
	return userApps, err
}
