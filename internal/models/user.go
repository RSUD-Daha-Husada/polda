package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	Username  string     `gorm:"type:varchar(100);not null;unique" json:"username"`
	Password  *string    `gorm:"type:text" json:"password"`
	App       *string    `gorm:"type:text" json:"app"`
	Gender    string     `gorm:"type:varchar(10)" json:"gender"`
	Email     string     `gorm:"type:varchar(100);unique" json:"email"`
	Telephone string     `gorm:"type:varchar(20)" json:"telephone"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	RoleID    uuid.UUID  `gorm:"type:uuid" json:"role_id"`
	Photo     string     `gorm:"type:text" json:"photo"`
	LastLogin *time.Time `gorm:"type:timestamp" json:"last_login"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	UpdatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy uuid.UUID  `gorm:"type:uuid" json:"updated_by"`

	// Relasi ke tabel pivot
	UserApps []UserApplication `gorm:"foreignKey:UserID" json:"-"`

	// Ini hanya untuk response agar langsung dapat list aplikasi (tanpa tag GORM)
	AccessApps []Application `gorm:"-" json:"access_apps"`
}
