package helpers

import (
	"fmt"
	"time"

	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SaveLoginLog(db *gorm.DB, userID *uuid.UUID, access, userAgent, ipAddress, status, message string) error {
	log := model.Log{
		LogID:     uuid.New(),
		Access:    access,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		Status:    status,
		Message:   message,
		CreatedAt: time.Now(),
		UserID:    userID, // ini sudah *uuid.UUID, jadi nil akan disimpan NULL
	}

	err := db.Create(&log).Error
	if err != nil {
		fmt.Println("Error saving log:", err)
	}
	return err
}
