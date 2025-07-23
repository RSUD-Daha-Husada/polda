package repository

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
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

func (r *UserRepository) FindPaginated(limit, offset int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.DB.Model(&model.User{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR username ILIKE ?", searchPattern, searchPattern)
	}

	// Hitung total dulu
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data user + relasi ke Application lewat UserApplication
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("UserApps.Application"). // Preload relasi aplikasi
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// Assign access_apps
	for i := range users {
		for _, ua := range users[i].UserApps {
			users[i].AccessApps = append(users[i].AccessApps, ua.Application)
		}
	}

	return users, total, nil
}

func (r *UserRepository) CreateUser(db *gorm.DB, user *model.User) error {
    return db.Create(user).Error
}

func (r *UserRepository) CreateUserApplication(tx *gorm.DB, userApp *model.UserApplication) error {
	return tx.Create(userApp).Error
}
