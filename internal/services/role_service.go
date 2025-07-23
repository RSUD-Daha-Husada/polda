package service

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/RSUD-Daha-Husada/polda-be/internal/repositories"
	"gorm.io/gorm"
)

func GetAllRoles(db *gorm.DB) ([]model.Role, error) {
	return repository.GetAllRoles(db)
}