package handler

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllRoles(c *fiber.Ctx, db *gorm.DB) error {
	roles, err := service.GetAllRoles(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch roles",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Roles fetched successfully",
		"data":    roles,
	})
}