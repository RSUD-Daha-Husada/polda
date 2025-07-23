package routes

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/handlers"
	"github.com/RSUD-Daha-Husada/polda-be/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoleRoutes(router fiber.Router, db *gorm.DB) {
	role := router.Group("/roles", middleware.JWTProtected(db))

	role.Get("/", func(c *fiber.Ctx) error {
		return handler.GetAllRoles(c, db)
	})
}