package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/RSUD-Daha-Husada/polda-be/helpers"
	"github.com/RSUD-Daha-Husada/polda-be/internal/models"
	"github.com/RSUD-Daha-Husada/polda-be/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ApplicationHandler struct {
	AppSvc *service.ApplicationService
}

func NewAppHandler(appSvc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{AppSvc: appSvc}
}

func (h *ApplicationHandler) GetApplications(c *fiber.Ctx) error {
	apps, err := h.AppSvc.ListApplications()
	if err != nil {
		helpers.SaveLoginLog(h.AppSvc.DB, nil, "get_applications", c.Get("User-Agent"), c.IP(), "failed", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var resp []model.Application
	for _, a := range apps {
		resp = append(resp, model.Application{
			ApplicationID: a.ApplicationID,
			Name:          a.Name,
			Icon:          a.Icon,
			RedirectURI:   a.RedirectURI,
			LogoutURL:     a.LogoutURL,
			CreatedAt:     a.CreatedAt,
		})
	}

	return c.JSON(resp)
}

func (h *ApplicationHandler) GetAllApplicationsForTable(c *fiber.Ctx) error {
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

	apps, total, err := h.AppSvc.GetPaginatedApplications(limit, offset, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch applications",
		})
	}

	return c.JSON(fiber.Map{
		"data":  apps,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	name := c.FormValue("name")
	redirectURI := c.FormValue("redirect_uri")
	logoutURL := c.FormValue("logout_url")

	if name == "" || redirectURI == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and Redirect URI are required",
		})
	}

	file, err := c.FormFile("icon")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Icon file is required",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExt[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only .jpg, .jpeg, .png files are allowed",
		})
	}

	if file.Size > 2*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size must be less than 2MB",
		})
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	publicPath := fmt.Sprintf("./public/uploads/icons/%s", filename)

	if err := c.SaveFile(file, publicPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save icon file",
		})
	}

	appURL := os.Getenv("APP_URL")
	iconURL := fmt.Sprintf("%s/uploads/icons/%s", appURL, filename)

	createdBy := c.Locals("user_id")
	var parsedCreatedBy uuid.UUID
	if createdByStr, ok := createdBy.(string); ok {
        parsedCreatedBy, _ = uuid.Parse(createdByStr) 
	}
    fmt.Println(parsedCreatedBy)
    
	app := &model.Application{
		ApplicationID: uuid.New(),
		Name:          name,
		RedirectURI:   redirectURI,
		LogoutURL:     logoutURL,
		Icon:          iconURL,
		CreatedAt:     time.Now(),
		CreatedBy:     parsedCreatedBy,
	}

	if err := h.AppSvc.CreateApplication(app); err != nil {
		_ = helpers.SaveLoginLog(h.AppSvc.DB, &app.CreatedBy, "create_application", c.Get("User-Agent"), c.IP(), "failed", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create application",
		})
	}

	_ = helpers.SaveLoginLog(h.AppSvc.DB, &app.CreatedBy, "create_application", c.Get("User-Agent"), c.IP(), "success", fmt.Sprintf("application_id %s created", app.ApplicationID))
	return c.Status(fiber.StatusCreated).JSON(app)
}

func (h *ApplicationHandler) ToggleAppActive(c *fiber.Ctx) error {
	appID := c.Params("id")

    var updatedBy *uuid.UUID
    if updatedByVal := c.Locals("user_id"); updatedByVal != nil {
        if updatedByStr, ok := updatedByVal.(string); ok {
            if parsed, err := uuid.Parse(updatedByStr); err == nil {
                updatedBy = &parsed
            }
        }
    }

	if err := h.AppSvc.ToggleAppActiveStatus(appID); err != nil {
        _ = helpers.SaveLoginLog(h.AppSvc.DB, updatedBy, "toggle_app_active", c.Get("User-Agent"), c.IP(), "failed", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to toggle active status",
			"error":   err.Error(),
		})
	}

    _ = helpers.SaveLoginLog(h.AppSvc.DB, updatedBy, "toggle_app_active", c.Get("User-Agent"), c.IP(), "success", fmt.Sprintf("app_id %s toggled active status", appID))
	return c.JSON(fiber.Map{
		"message": "Application status updated successfully",
	})
}

func (h *ApplicationHandler) EditApplication(c *fiber.Ctx) error {
	// Ambil ID aplikasi dari parameter URL
	appIDParam := c.Params("id")
	appID, err := uuid.Parse(appIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Ambil data lama dari DB
	existingApp, err := h.AppSvc.GetApplicationByID(appID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Application not found",
		})
	}

	// Ambil data form
	name := c.FormValue("name")
	redirectURI := c.FormValue("redirect_uri")
	logoutURL := c.FormValue("logout_url")

	if name == "" || redirectURI == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and Redirect URI are required",
		})
	}

	iconURL := existingApp.Icon

	// Handle file icon jika ada upload baru
	file, err := c.FormFile("icon")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
		if !allowedExt[ext] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Only .jpg, .jpeg, .png files are allowed",
			})
		}

		if file.Size > 2*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "File size must be less than 2MB",
			})
		}

		// Simpan file baru
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		publicPath := fmt.Sprintf("./public/uploads/icons/%s", filename)

		if err := c.SaveFile(file, publicPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save icon file",
			})
		}

		appURL := os.Getenv("APP_URL")
		iconURL = fmt.Sprintf("%s/uploads/icons/%s", appURL, filename)
	}

	// Update data aplikasi
	existingApp.Name = name
	existingApp.RedirectURI = redirectURI
	existingApp.LogoutURL = logoutURL
	existingApp.Icon = iconURL

	if err := h.AppSvc.UpdateApplication(existingApp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update application",
		})
	}

	return c.Status(fiber.StatusOK).JSON(existingApp)
}
