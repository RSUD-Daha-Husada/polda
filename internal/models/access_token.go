package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Token      string     `gorm:"type:text;not null" json:"token"`
	IPAddress  string     `gorm:"type:text" json:"ip_address"`
	UserAgent  string     `gorm:"type:text" json:"user_agent"`
	IsRevoked  bool       `gorm:"default:false" json:"is_revoked"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiredAt  time.Time  `json:"expired_at"`
}