package handler

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/model"
	"github.com/RSUD-Daha-Husada/polda-be/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		Service: service.NewAuthService(db),
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	ip := c.IP()
	userAgent := string(c.Request().Header.UserAgent())

	token, err := h.Service.Login(input.Username, input.Password, ip, userAgent)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "login successful",
		"token":   token,
	})
}

// RequestLoginCode handles code request via WA/email
func (h *AuthHandler) RequestLoginCode(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
	}

	// Parse input dari body
	if err := c.BodyParser(&input); err != nil || input.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Cari user berdasarkan username
	var user model.User
	if err := h.Service.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Cek apakah user punya nomor telepon
	if user.Telephone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user does not have a registered phone number"})
	}

	// Kirim kode login via WA menggunakan nomor telepon user
	err := h.Service.GenerateLoginCode(user.Telephone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "login code sent"})
}

// LoginWithCode handles actual login via code
func (h *AuthHandler) LoginWithCode(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}

	if err := c.BodyParser(&input); err != nil || input.Username == "" || input.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	token, err := h.Service.LoginWithCode(input.Username, input.Code, userAgent, ipAddress)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "login via code successful",
		"token":   token,
	})
}

