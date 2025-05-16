package main

import (
	"log"
	"url_shortener/internal/database"
	"url_shortener/app/routes"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("DB connection failed")
	}

	log.Println("DB connection successful")
	router := routes.SetupRouter()
	router.Run() // on 0.0.0.0:8080
}