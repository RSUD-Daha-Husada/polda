package routes

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/handlers"
	"github.com/RSUD-Daha-Husada/polda-be/internal/services"
	"github.com/RSUD-Daha-Husada/polda-be/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAppRoutes(router fiber.Router, db *gorm.DB) {
    app := router.Group("/apps", middleware.JWTProtected(db))

    appSvc := service.NewAppService(db)
    appHandler := handler.NewAppHandler(appSvc)

    app.Get("/", appHandler.GetApplications)
    app.Get("/get-all", appHandler.GetAllApplicationsForTable)
    app.Post("/create", appHandler.CreateApplication)
    app.Put("/edit/:id", appHandler.EditApplication)
    app.Patch("/toggle-active/:id", appHandler.ToggleAppActive)
}
