package routes

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/handler"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	auth := api.Group("/auth") // semua login lewat /api/auth/...

	authHandler := handler.NewAuthHandler(db)

	// Login biasa (username + password)
	auth.Post("/login", authHandler.Login)

	// Kirim kode ke WA / Email
	auth.Post("/request-code", authHandler.RequestLoginCode)

	// Login pakai kode
	auth.Post("/login-code", authHandler.LoginWithCode)
}

