package server

import (
	"log"
	"os"
	"url_shortener/app/routes"
	"url_shortener/internal/database"
	"github.com/gofiber/fiber/v2"
)

func RunServer(){
	//db kontrolu yapiliyor
	if err := database.Connect(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
		os.Exit(1)
	}
	log.Println("DB connection successful")


	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":8080")
}