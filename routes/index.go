package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	// Modular route register
	RegisterAuthRoutes(api, db)
	RegisterUserRoutes(api, db)
	// Tambah modul lain di sini
}
