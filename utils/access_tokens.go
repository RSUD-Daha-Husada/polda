package utils

import (
	"fmt"
	"time"

	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SaveAccessToken(db *gorm.DB, userID uuid.UUID, token, userAgent, ipAddress string, expiredAt time.Time) error {
	accessToken := model.AccessToken{
		ID:         uuid.New(),
		UserID:     userID,
		Token:      token,
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
		IsRevoked:  false,
		CreatedAt:  time.Now(),
		ExpiredAt:  expiredAt,
	}

	err := db.Create(&accessToken).Error
	if err != nil {
		fmt.Println("Error saving access token:", err)
	}
	return err
}
