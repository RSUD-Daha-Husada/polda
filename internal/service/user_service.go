package service

import (
	"errors"

	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"github.com/RSUD-Daha-Husada/polda-be/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB       *gorm.DB
	UserRepo *repository.UserRepository
}

type UpdateProfileInput struct {
	UserID      uuid.UUID
	OldPassword string
	Password    string
	Poto        string
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB:       db,
		UserRepo: repository.NewUserRepository(db),
	}
}

func (s *UserService) GetProfile(userID uuid.UUID) (*model.User, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(input UpdateProfileInput) (*model.User, error) {
	user, err := s.UserRepo.FindByID(input.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	updates := make(map[string]interface{})

	if input.Password != "" {
		if user.Password == nil {
			return nil, errors.New("password not set")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.OldPassword)); err != nil {
			return nil, errors.New("old password is incorrect")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		updates["password"] = string(hashedPassword)
	}

	if input.Poto != "" {
		updates["poto"] = input.Poto
	}

	if len(updates) == 0 {
		return user, nil
	}

	if err := s.UserRepo.Update(input.UserID, updates); err != nil {
		return nil, errors.New("failed to update profile")
	}

	updatedUser, err := s.UserRepo.FindByID(input.UserID)
	if err != nil {
		return nil, errors.New("failed to fetch updated user")
	}
	return updatedUser, nil
}
