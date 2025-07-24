package handler

import (
	"strings"
	"time"

	"github.com/RSUD-Daha-Husada/polda-be/helpers"
	model "github.com/RSUD-Daha-Husada/polda-be/internal/models"
	service "github.com/RSUD-Daha-Husada/polda-be/internal/services"
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

func (h *AuthHandler) RequestLoginCode(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&input); err != nil || input.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	var user model.User
	if err := h.Service.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	if user.Telephone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user does not have a registered phone number"})
	}

	err := h.Service.GenerateLoginCode(user.Telephone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "login code sent"})
}

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

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No token provided"})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	ip := c.IP()
	userAgent := c.Get("User-Agent")

	err := h.Service.InvalidateToken(token, userAgent, ip)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "logout successful"})
}

func (h *AuthHandler) CheckToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	ip := c.IP()
	userAgent := c.Get("User-Agent")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		_ = helpers.SaveLoginLog(h.Service.DB, nil, "check_token", userAgent, ip, "failed", "no authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	rawToken := strings.TrimPrefix(authHeader, "Bearer ")

	var at model.AccessToken
	if err := h.Service.DB.Where("token = ?", rawToken).First(&at).Error; err != nil {
		_ = helpers.SaveLoginLog(h.Service.DB, nil, "check_token", userAgent, ip, "failed", "token not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is revoked or invalid",
		})
	}

	if at.IsRevoked {
		_ = helpers.SaveLoginLog(h.Service.DB, &at.UserID, "check_token", userAgent, ip, "failed", "token is revoked")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token has been revoked",
		})
	}

	if time.Now().After(at.ExpiredAt) {
		_ = h.Service.DB.Model(&at).Update("is_revoked", true).Error
		_ = helpers.SaveLoginLog(h.Service.DB, &at.UserID, "check_token", userAgent, ip, "failed", "token expired")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token has expired",
		})
	}

	return c.JSON(fiber.Map{"valid": true})
}
