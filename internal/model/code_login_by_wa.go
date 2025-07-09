package model

import (
	"time"

	"github.com/google/uuid"
)

type CodeLoginByWA struct {
	CodeLoginByWAID uuid.UUID `gorm:"type:uuid;primaryKey" json:"code_login_by_wa_id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Code   string    `gorm:"type:varchar(6);not null" json:"code"`

	ValidUntil time.Time `gorm:"not null" json:"valid_until"` 
	Used       bool      `gorm:"default:false" json:"used"`   
	Status     bool      `gorm:"default:true" json:"status"` 

	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (CodeLoginByWA) TableName() string {
	return "code_login_by_wa"
}