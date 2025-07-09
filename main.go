package main

import (
	"fmt"
	"log"

	"github.com/RSUD-Daha-Husada/polda-be/config"
	"github.com/RSUD-Daha-Husada/polda-be/routes"
	"github.com/RSUD-Daha-Husada/polda-be/utils"
	"github.com/gofiber/fiber/v2"
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
	routes.RegisterRoutes(app, db)

	log.Fatal(app.Listen(":3000"))
}

