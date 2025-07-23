package middleware

import (
	"os"
	"strings"

	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func JWTProtected(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1 — header check
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		// 2 — verify JWT signature & expiry
		token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// 3 — Cek token di DB (revoked / tidak)
		var at model.AccessToken
		if err := db.Where("token = ?", rawToken).First(&at).Error; err != nil {
			// Tidak ada di DB  → tolak
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token not recognized",
			})
		}
		if at.IsRevoked {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token has been revoked",
			})
		}

		// 4 — set user_id ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if uid, ok := claims["user_id"].(string); ok {
				c.Locals("user_id", uid)
			}
		}

		return c.Next()
	}
}
