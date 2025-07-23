package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	RoleID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"role_id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
