package routes

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/handler"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(router fiber.Router, db *gorm.DB) {
	auth := router.Group("/auth")
	authHandler := handler.NewAuthHandler(db)

	auth.Post("/login", authHandler.Login)
	auth.Post("/request-code", authHandler.RequestLoginCode)
	auth.Post("/login-code", authHandler.LoginWithCode)
	auth.Post("/logout", authHandler.Logout)
	auth.Post("/check-token", authHandler.CheckToken)
}
