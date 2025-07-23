package service

import (
	"context"

	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserApplicationService struct {
	DB *gorm.DB
}

func NewUserAppService(db *gorm.DB) *UserApplicationService {
	return &UserApplicationService{DB: db}
}

func (s *UserApplicationService) CreateUserApp(ctx context.Context, userApp *model.UserApplication) error {
	userApp.UserApplicationID = uuid.New()
	return s.DB.WithContext(ctx).Create(userApp).Error
}

func (s *UserApplicationService) GetAllUserApps(ctx context.Context) ([]model.UserApplication, error) {
	var userApps []model.UserApplication
	err := s.DB.WithContext(ctx).Preload("App").Find(&userApps).Error
	return userApps, err
}

func (s *UserApplicationService) GetUserAppByID(ctx context.Context, id uuid.UUID) (*model.UserApplication, error) {
	var userApp model.UserApplication
	err := s.DB.WithContext(ctx).First(&userApp, "user_application_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &userApp, nil
}

func (s *UserApplicationService) GetUserAppsByUserID(ctx context.Context, userID uuid.UUID) ([]model.UserApplication, error) {
	var userApps []model.UserApplication
	err := s.DB.WithContext(ctx).
		Preload("Application"). // ini preload agar datanya langsung ada
		Where("user_id = ?", userID).
		Find(&userApps).Error

	return userApps, err
}

func (s *UserApplicationService) UpdateUserApp(ctx context.Context, id uuid.UUID, updated *model.UserApplication) error {
	return s.DB.WithContext(ctx).Model(&model.UserApplication{}).Where("user_application_id = ?", id).Updates(updated).Error
}

func (s *UserApplicationService) DeleteUserApp(ctx context.Context, id uuid.UUID) error {
	return s.DB.WithContext(ctx).Delete(&model.UserApplication{}, "user_application_id = ?", id).Error
}
