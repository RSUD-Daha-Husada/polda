package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/RSUD-Daha-Husada/polda-be/helpers"
	"github.com/RSUD-Daha-Husada/polda-be/internal/dto"
	"github.com/RSUD-Daha-Husada/polda-be/internal/services"
	"github.com/RSUD-Daha-Husada/polda-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB      *gorm.DB
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

	return c.JSON(fiber.Map{
		"id":        user.UserID,
		"username":  user.Username,
		"app":       user.App,
		"email":     user.Email,
		"telephone": user.Telephone,
		"gender":    user.Gender,
		"role_id":   user.RoleID,
		"photo":     user.Photo,
		"is_active": user.IsActive,
		"lastLogin": user.LastLogin,
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
	telephone := c.FormValue("telephone")

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

	fileHeader, err := c.FormFile("avatar")
	var photoURL string
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
			_ = helpers.SaveLoginLog(h.UserSvc.DB, &userID, "edit_profile", userAgent, ip, "failed", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuka file"})
		}
		defer file.Close()

		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			_ = helpers.SaveLoginLog(h.UserSvc.DB, &userID, "edit_profile", userAgent, ip, "failed", err.Error())
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

		photoURL = filename
	}

	input := service.UpdateProfileInput{
		UserID:      userID,
		OldPassword: oldPassword,
		Password:    newPassword,
		Photo:       photoURL,
		Telephone:   telephone,
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
		"photo":    updatedUser.Photo,
	})
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	if limit < 1 {
		limit = 8
	}
	offset := (page - 1) * limit
	search := c.Query("search", "")

	users, total, err := h.UserSvc.GetPaginatedUsers(limit, offset, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	input := new(dto.CreateUserInput)

	roleIDStr := c.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role_id",
		})
	}
	input.RoleID = roleID

	input.Name = c.FormValue("name")
	input.Email = c.FormValue("email")
	input.Username = c.FormValue("username")
	input.Password = c.FormValue("password")
	input.Telephone = c.FormValue("telephone")
	input.Gender = c.FormValue("gender")

	file, err := c.FormFile("photo")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
		if !allowedExt[ext] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Only .jpg, .jpeg, .png files are allowed for photo",
			})
		}

		if file.Size > 2*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Photo size must be less than 2MB",
			})
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := fmt.Sprintf("./public/uploads/avatars/%s", filename)

		if err := c.SaveFile(file, savePath); err != nil {
			_ = helpers.SaveLoginLog(h.UserSvc.DB, nil, "create_user", c.Get("User-Agent"), c.IP(), "failed", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save photo file",
			})
		}

		appURL := os.Getenv("APP_URL")
		photoURL := fmt.Sprintf("%s/uploads/avatars/%s", appURL, filename)
		input.Photo = photoURL
	}

	createdBy := c.Locals("user_id")
	if createdByStr, ok := createdBy.(string); ok {
		parsedCreatedBy, err := uuid.Parse(createdByStr)
		if err == nil {
			input.CreatedBy = parsedCreatedBy
		}
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid multipart form data",
		})
	}
	applicationIDs := form.Value["application_ids"]
	for _, idStr := range applicationIDs {
		id, err := uuid.Parse(idStr)
		if err == nil {
			input.ApplicationIDs = append(input.ApplicationIDs, id)
		}
	}

	if err := utils.ValidateStruct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.UserSvc.CreateUser(*input, input.CreatedBy)
	if err != nil {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, nil, "create_user", c.Get("User-Agent"), c.IP(), "failed", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_ = helpers.SaveLoginLog(h.UserSvc.DB, &input.CreatedBy, "create_user", c.Get("User-Agent"), c.IP(), "success", fmt.Sprintf("username %s created", user.Username))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

func (h *UserHandler) ToggleUserActive(c *fiber.Ctx) error {
	userID := c.Params("id")

	var updatedBy *uuid.UUID
	if updatedByVal := c.Locals("user_id"); updatedByVal != nil {
		if updatedByStr, ok := updatedByVal.(string); ok {
			if parsed, err := uuid.Parse(updatedByStr); err == nil {
				updatedBy = &parsed
			}
		}
	}

	if err := h.UserSvc.ToggleUserActiveStatus(userID); err != nil {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, updatedBy, "toggle_user_active", c.Get("User-Agent"), c.IP(), "failed", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to toggle active status",
			"error":   err.Error(),
		})
	}

	_ = helpers.SaveLoginLog(
		h.UserSvc.DB,
		updatedBy,
		"toggle_user_active",
		c.Get("User-Agent"),
		c.IP(),
		"success",
		fmt.Sprintf("user_id %s toggled active status", userID),
	)

	return c.JSON(fiber.Map{
		"message": "User status updated successfully",
	})
}

func (h *UserHandler) EditUser(c *fiber.Ctx) error {
	input := new(dto.EditUserInput)

	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id",
		})
	}

	input.UserID = userID

	roleIDStr := c.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role_id",
		})
	}
	input.RoleID = roleID

	input.Name = c.FormValue("name")
	input.Email = c.FormValue("email")
	input.Username = c.FormValue("username")
	password := c.FormValue("password")
	if password != "" {
		input.Password = &password
	}
	input.Telephone = c.FormValue("telephone")
	input.Gender = c.FormValue("gender")

	file, err := c.FormFile("photo")
	if err == nil {
		savePath := "./public/uploads/avatars/"
		_ = utils.CreateDirIfNotExist(savePath)

		filename := uuid.New().String() + "_" + file.Filename
		fullPath := savePath + filename

		if err := c.SaveFile(file, fullPath); err == nil {
			input.Photo = filename
		}
	}

	updatedBy := c.Locals("user_id")
	if updatedByStr, ok := updatedBy.(string); ok {
		parsedUpdatedBy, err := uuid.Parse(updatedByStr)
		if err == nil {
			input.UpdatedBy = parsedUpdatedBy
		}
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid multipart form data",
		})
	}

	applicationIDs := form.Value["application_ids"]
	for _, idStr := range applicationIDs {
		id, err := uuid.Parse(idStr)
		if err == nil {
			input.ApplicationIDs = append(input.ApplicationIDs, id)
		}
	}

	if err := utils.ValidateStruct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	updatedUser, err := h.UserSvc.UpdateUser(*input, input.UpdatedBy)
	if err != nil {
		_ = helpers.SaveLoginLog(h.UserSvc.DB, &input.UserID, "edit_user", c.Get("User-Agent"), c.IP(), "failed", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_ = helpers.SaveLoginLog(h.UserSvc.DB, &input.UpdatedBy, "edit_user", c.Get("User-Agent"), c.IP(), "success", fmt.Sprintf("username %s updated", input.Username))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    updatedUser,
	})
}
