package shortener

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WelcomeFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H {
		"message": "Welcome to URL Shortener",
	})
}