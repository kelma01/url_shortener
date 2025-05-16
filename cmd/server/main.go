package main

import (
	"log"
	"url_shortener/internal/database"
	"url_shortener/app/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("DB connection failed: %v", err)
	}

	log.Println("DB connection successful")
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":8080")
}