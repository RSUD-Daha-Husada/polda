package routes

import (
	"github.com/RSUD-Daha-Husada/polda-be/internal/handler"
	"github.com/RSUD-Daha-Husada/polda-be/internal/service"
	"github.com/RSUD-Daha-Husada/polda-be/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router fiber.Router, db *gorm.DB) {
	user := router.Group("/user", middleware.JWTProtected(db))

	userSvc := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userSvc)

	user.Get("/me", userHandler.Me)
	user.Put("/edit-profile", userHandler.EditProfile)
}
