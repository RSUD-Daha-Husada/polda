package model

import (
	"time"

	"github.com/google/uuid"
)

type UserApplication struct {
	UserApplicationID uuid.UUID    `gorm:"type:uuid;primaryKey" json:"user_application_id"`
	UserID            uuid.UUID    `gorm:"type:uuid;not null" json:"user_id"`
	ApplicationID     uuid.UUID    `gorm:"type:uuid;not null" json:"application_id"`
	CreatedAt         time.Time    `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time    `gorm:"type:timestamp ;default:CURRENT_TIMESTAMP" json:"updated_at"`

	Application       Application `gorm:"foreignKey:ApplicationID;references:ApplicationID"`
}

func (UserApplication) TableName() string {
	return "user_applications"
}
