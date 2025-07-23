package model

import (
	"github.com/google/uuid"
	"time"
)

type Log struct {
	LogID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"log_id"`
	UserID    *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Access    string    `gorm:"type:varchar(50)" json:"access"` // contoh: "login"
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	IPAddress string    `gorm:"type:varchar(100)" json:"ip_address"`
	Status    string    `gorm:"type:varchar(20)" json:"status"` // misal "success", "failed"
	Message   string    `gorm:"type:text" json:"message"`    
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
