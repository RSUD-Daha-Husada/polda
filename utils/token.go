package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID) (string, time.Time, error) {
	expiration := time.Now().Add(3 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     expiration.Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiration, nil
}
