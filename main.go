package main

import (
	"fmt"
	"log"

	"github.com/RSUD-Daha-Husada/polda-be/config"
	"github.com/RSUD-Daha-Husada/polda-be/routes"
	"github.com/RSUD-Daha-Husada/polda-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
	}

	db := config.ConnectDB()

	// Generate hashed password dan print (hanya untuk testing)
	password := "password123"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hashed password:", hashedPassword)

	app := fiber.New()

	// middleware cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // atau "*" untuk semua origin
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Static("/uploads", "./public/uploads")

	routes.RegisterRoutes(app, db)

	log.Fatal(app.Listen(":3000"))
}
