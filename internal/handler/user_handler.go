package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/RSUD-Daha-Husada/polda-be/helpers"
	"github.com/RSUD-Daha-Husada/polda-be/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{
		UserSvc: userSvc,
	}
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, nil, "get_profile", c.Get("User-Agent"), c.IP(), "failed", "unauthorized access - user_id missing")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, nil, "get_profile", c.Get("User-Agent"), c.IP(), "failed", "invalid user_id format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	user, err := h.UserSvc.GetProfile(userID)
	if err != nil {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, &userID, "get_profile", c.Get("User-Agent"), c.IP(), "failed", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// Tidak perlu logging kalau sukses
	return c.JSON(fiber.Map{
		"id":        user.UserID,
		"username":  user.Username,
		"email":     user.Email,
		"telephone": user.Telephone,
		"gender":    user.Gender,
		"roleId":    user.RoleID,
		"poto":      user.Poto,
		"isActive":  user.IsActive,
	})
}

func (h *UserHandler) EditProfile(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)

	userAgent := c.Get("User-Agent")
	ip := c.IP()

	if !ok {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, nil, "edit_profile", userAgent, ip, "failed", "unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, nil, "edit_profile", userAgent, ip, "failed", "invalid user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	newPassword := c.FormValue("password")
	oldPassword := c.FormValue("old_password")

	// ✅ Validasi password (jika diisi)
	if newPassword != "" {
		var (
			hasMinLen = false
			hasUpper  = false
			hasLower  = false
			hasNumber = false
		)

		if len(newPassword) >= 8 {
			hasMinLen = true
		}

		for _, ch := range newPassword {
			switch {
			case unicode.IsUpper(ch):
				hasUpper = true
			case unicode.IsLower(ch):
				hasLower = true
			case unicode.IsDigit(ch):
				hasNumber = true
			}
		}

		if !(hasMinLen && hasUpper && hasLower && hasNumber) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password harus minimal 8 karakter, mengandung huruf kapital, huruf kecil, dan angka",
			})
		}

		if oldPassword == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password lama wajib diisi untuk mengganti password",
			})
		}
	}

	// 📷 Ambil file avatar
	fileHeader, err := c.FormFile("avatar")
	var potoURL string
	if err == nil && fileHeader != nil {
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format file harus JPG, JPEG, atau PNG"})
		}

		if fileHeader.Size > 2*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ukuran file maksimal 2MB"})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuka file"})
		}
		defer file.Close()

		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca file"})
		}

		contentType := http.DetectContentType(buffer)
		if contentType != "image/jpeg" && contentType != "image/png" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tipe file tidak valid"})
		}

		filename := fmt.Sprintf("%s_%d%s", userID.String(), time.Now().Unix(), ext)
		savePath := filepath.Join("public", "uploads", "avatars", filename)

		if err := c.SaveFile(fileHeader, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan file"})
		}

		potoURL = filename
	}

	// Kirim ke service
	input := service.UpdateProfileInput{
		UserID:      userID,
		OldPassword: oldPassword,
		Password:    newPassword,
		Poto:        potoURL,
	}

	updatedUser, err := h.UserSvc.UpdateProfile(input)
	if err != nil {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, &userID, "edit_profile", userAgent, ip, "failed", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_ = helpers.SaveLoginLog(h.UserSvc.DB, &userID, "edit_profile", userAgent, ip, "success", "profile updated")
	return c.JSON(fiber.Map{
		"message":  "Profile updated successfully",
		"username": updatedUser.Username,
		"poto":     updatedUser.Poto,
	})
}
