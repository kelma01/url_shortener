package server

import (
	"log"
	"url_shortener/internal/database"
	"url_shortener/app/routes"
	"github.com/gofiber/fiber/v2"
)

func RunServer(){
	//db kontrolu yapiliyor
	if err := database.Connect(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	log.Println("DB connection successful")

	app := fiber.New()
	//routes pathinde belirtilen routelarin tanimlanmasi
	routes.SetupRoutes(app)
	app.Listen(":8080")
}