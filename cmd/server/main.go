package main

import (
	"url_shortener/app/routes"
)

func main(){
	router := routes.SetupRouter()

	router.Run() //on 0.0.0.0:8080 
}