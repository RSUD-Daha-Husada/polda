package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	RegisterAuthRoutes(api, db)
	RegisterUserRoutes(api, db)
	RegisterAppRoutes(api, db)
	RegisterUserAppRoutes(api, db)
	RegisterRoleRoutes(api, db)
}
