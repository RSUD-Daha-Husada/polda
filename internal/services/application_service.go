package service

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/RSUD-Daha-Husada/polda-be/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationService struct {
	DB      *gorm.DB
	AppRepo *repository.AppRepository
}

func NewAppService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{
		DB:      db,
		AppRepo: repository.NewAppRepository(db),
	}
}

func (s *ApplicationService) ListApplications() ([]*model.Application, error) {
	return s.AppRepo.FindAll()
}

func (s *ApplicationService) GetPaginatedApplications(limit, offset int, search string) ([]model.Application, int64, error) {
	return s.AppRepo.FindPaginated(limit, offset, search)
}

func (s *ApplicationService) CreateApplication(app *model.Application) error {
	return s.DB.Create(app).Error
}

func (s *ApplicationService) ToggleAppActiveStatus(appID string) error {
	return s.AppRepo.ToggleActiveStatus(appID)
}

func (s *ApplicationService) GetApplicationByID(id uuid.UUID) (*model.Application, error) {
	return s.AppRepo.FindByID(id)
}

func (s *ApplicationService) UpdateApplication(app *model.Application) error {
	return s.AppRepo.Update(app)
}
