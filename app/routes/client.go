package routes

import (
	"github.com/gin-gonic/gin"
	"url_shortener/app/controller"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", shortener.WelcomeFunc)
	return router
}