package model

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ApplicationID uuid.UUID `gorm:"type:uuid;primaryKey" json:"application_id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Icon          string    `gorm:"type:varchar(255)" json:"icon"`
	RedirectURI   string    `gorm:"type:varchar(255)" json:"redirect_uri"`
	LogoutURL     string    `gorm:"type:varchar(255)" json:"logout_url"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy     uuid.UUID `gorm:"type:uuid" json:"created_by"`
}
