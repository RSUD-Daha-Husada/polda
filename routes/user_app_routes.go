package routes

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/handlers"
	"github.com/RSUD-Daha-Husada/polda-be/internal/services"
	"github.com/RSUD-Daha-Husada/polda-be/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserAppRoutes(router fiber.Router, db *gorm.DB) {
	userApps := router.Group("/user-apps", middleware.JWTProtected(db))

	appSvc := service.NewUserAppService(db)
	appHandler := handler.NewUserAppHandler(appSvc)

	userApps.Get("/me", appHandler.GetUserAppsByUser)
	userApps.Post("/create", appHandler.CreateUserApp)
}
