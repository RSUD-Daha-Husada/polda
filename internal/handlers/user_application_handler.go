package handler

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/RSUD-Daha-Husada/polda-be/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserApplicationHandler struct {
	Service *service.UserApplicationService
}

func NewUserAppHandler(service *service.UserApplicationService) *UserApplicationHandler {
	return &UserApplicationHandler{Service: service}
}

func (h *UserApplicationHandler) CreateUserApp(c *fiber.Ctx) error {
	var input model.UserApplication
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.Service.CreateUserApp(c.Context(), &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(input)
}

func (h *UserApplicationHandler) GetAllUserApps(c *fiber.Ctx) error {
	userApps, err := h.Service.GetAllUserApps(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(userApps)
}

func (h *UserApplicationHandler) GetUserAppByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	userApp, err := h.Service.GetUserAppByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "UserApp not found"})
	}

	return c.JSON(userApp)
}

func (h *UserApplicationHandler) GetUserAppsByUser(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	userApps, err := h.Service.GetUserAppsByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(userApps)
}

func (h *UserApplicationHandler) UpdateUserApp(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	var input model.UserApplication
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.Service.UpdateUserApp(c.Context(), id, &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Updated successfully"})
}

func (h *UserApplicationHandler) DeleteUserApp(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	if err := h.Service.DeleteUserApp(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Deleted successfully"})
}
