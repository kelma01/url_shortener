package server

import (
	"log"

	"url_shortener/app/routes"
	
	"github.com/gofiber/fiber/v2"
	
	"github.com/gofiber/contrib/otelfiber"
	otel "url_shortener/internal/opentelemetry"

)

func RunServer(){
	//opentelemetry initialize etme kismi
	shutdown := otel.InitTracer()
	defer shutdown()

	app := fiber.New()

	app.Use(otelfiber.Middleware())
	
	routes.SetupRoutes(app)
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Fiber failed to start: %v", err)
	}
}
