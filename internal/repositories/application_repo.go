package repository

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppRepository struct {
	DB *gorm.DB
}

func NewAppRepository(db *gorm.DB) *AppRepository {
	return &AppRepository{DB: db}
}

func (r *AppRepository) FindPaginated(limit, offset int, search string) ([]model.Application, int64, error) {
	var apps []model.Application
	var total int64

	query := r.DB.Model(&model.Application{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ?", searchPattern)
	}

	// Hitung total data dulu
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data paginated
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

func (r *AppRepository) FindByID(id uuid.UUID) (*model.Application, error) {
	var app model.Application
	if err := r.DB.First(&app, "application_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *AppRepository) Update(app *model.Application) error {
	return r.DB.Save(app).Error
}

func (r *AppRepository) ToggleActiveStatus(appID string) error {
	var app model.Application
	if err := r.DB.First(&app, "application_id = ?", appID).Error; err != nil {
		return err
	}

	app.IsActive = !app.IsActive
	return r.DB.Save(&app).Error
}

func (r *AppRepository) FindAll() ([]*model.Application, error) {
	var apps []*model.Application
	if err := r.DB.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}
