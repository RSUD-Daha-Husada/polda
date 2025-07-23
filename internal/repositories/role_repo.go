package repository

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"gorm.io/gorm"
)
func GetAllRoles(db *gorm.DB) ([]model.Role, error) {
	var roles []model.Role
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}