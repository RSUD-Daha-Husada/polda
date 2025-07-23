package routes

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/handlers"
	"github.com/RSUD-Daha-Husada/polda-be/internal/services"
	"github.com/RSUD-Daha-Husada/polda-be/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router fiber.Router, db *gorm.DB) {
	user := router.Group("/user", middleware.JWTProtected(db))

	userSvc := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userSvc)

	user.Get("/me", userHandler.Me)
	user.Get("/get-all", userHandler.GetAllUsers)
	user.Get("/search", userHandler.GetAllUsers)
	user.Put("/edit-profile", userHandler.EditProfile)
	user.Put("/edit/:id", userHandler.EditUser)
	user.Post("/create", userHandler.CreateUser)
	user.Patch("/toggle-active/:id", userHandler.ToggleUserActive)
}
