package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Role(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role") // role harus diset di JWT middleware

		for _, r := range allowedRoles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access forbidden: insufficient role",
		})
	}
}
