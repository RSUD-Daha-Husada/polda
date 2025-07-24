package service

import (
	"errors"

	"github.com/RSUD-Daha-Husada/polda-be/internal/dto"
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/RSUD-Daha-Husada/polda-be/internal/repositories"
	"github.com/RSUD-Daha-Husada/polda-be/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB       *gorm.DB
	UserRepo *repository.UserRepository
}

type CreateUserInput struct {
	Name      string    `json:"name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	Telephone string    `json:"telephone"`
	RoleID    uuid.UUID `json:"role_id" validate:"required"`
	IsActive  bool      `json:"is_active"`
	CreatedBy uuid.UUID `json:"created_by"`

	ApplicationIDs []uuid.UUID `json:"application_ids"`
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

func (s *UserService) UpdateProfile(input dto.UpdateProfileInput) (*model.User, error) {
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

	if input.Photo != "" {
		updates["photo"] = input.Photo
	}

	if input.Telephone != "" {
		updates["telephone"] = input.Telephone
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

func (s *UserService) GetPaginatedUsers(limit, offset int, search string) ([]model.User, int64, error) {
	return s.UserRepo.FindPaginated(limit, offset, search)
}

func (s *UserService) CreateUser(input dto.CreateUserInput, createdBy uuid.UUID) (*model.User, error) {
	// 1. Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// 2. Buat user
	user := &model.User{
		UserID:    uuid.New(),
		Name:      input.Name,
		Username:  input.Username,
		Password:  &hashedPassword,
		Gender:    input.Gender,
		Email:     input.Email,
		Telephone: input.Telephone,
		Photo:     input.Photo,
		RoleID:    input.RoleID,
		CreatedBy: createdBy,
	}

	// 3. Simpan user ke DB
	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// 4. Tambahkan relasi ke aplikasi (jika ada)
	if len(input.ApplicationIDs) > 0 {
		var userApps []model.UserApplication

		for _, appID := range input.ApplicationIDs {
			userApps = append(userApps, model.UserApplication{
				UserApplicationID: uuid.New(),
				UserID:            user.UserID,
				ApplicationID:     appID,
			})
		}

		if err := s.DB.Create(&userApps).Error; err != nil {
			return nil, err
		}
	}

	// 5. Return user
	return user, nil
}

func (s *UserService) ToggleUserActiveStatus(userID string) error {
	var user model.User

	if err := s.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return err
	}

	user.IsActive = !user.IsActive

	if err := s.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(input dto.EditUserInput, updatedBy uuid.UUID) (*model.User, error) {
	// 1. Ambil user yang mau diupdate
	var user model.User
	if err := s.DB.First(&user, "user_id = ?", input.UserID).Error; err != nil {
		return nil, err
	}

	// 2. Update field yang dikirim
	user.Name = input.Name
	user.Username = input.Username
	user.Gender = input.Gender
	user.Email = input.Email
	user.Telephone = input.Telephone
	user.Photo = input.Photo
	user.RoleID = input.RoleID
	user.UpdatedBy = updatedBy

	// 3. Update password jika dikirim
	if input.Password != nil && *input.Password != "" {
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		user.Password = &hashedPassword
	}

	// 4. Simpan perubahan user
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	// 5. Update relasi aplikasi (hapus lama, simpan baru jika ada)
	if len(input.ApplicationIDs) > 0 {
		// Hapus relasi sebelumnya
		if err := s.DB.Where("user_id = ?", user.UserID).Delete(&model.UserApplication{}).Error; err != nil {
			return nil, err
		}

		// Tambahkan relasi baru
		var userApps []model.UserApplication
		for _, appID := range input.ApplicationIDs {
			userApps = append(userApps, model.UserApplication{
				UserApplicationID: uuid.New(),
				UserID:            user.UserID,
				ApplicationID:     appID,
			})
		}

		if err := s.DB.Create(&userApps).Error; err != nil {
			return nil, err
		}
	}

	// 6. Return user yang sudah diperbarui
	return &user, nil
}
