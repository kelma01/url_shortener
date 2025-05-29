package server

import (
	"log"
	"os"
	"url_shortener/app/routes"
	"url_shortener/internal/database"
	otel "url_shortener/internal/opentelemetry"
	"github.com/gofiber/fiber/v2"
)

func RunServer(){
	//opentelemetry initialize etme kismi
	shutdown := otel.InitTracer()
	defer shutdown()

	//db kontrolu yapiliyor
	if err := database.Connect(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
		os.Exit(1)
	}
	log.Println("DB connection successful")


	app := fiber.New()
	routes.SetupRoutes(app)
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Fiber failed to start: %v", err)
	}
}
