package routes

import (
    "github.com/gofiber/fiber/v2"
    "url_shortener/app/controller"
)

func SetupRoutes(app *fiber.App) {
    app.Post("/", shortener.ShortenURL)
    //app.Get("/", shortener.ListURLs)
	app.Get("/", shortener.WelcomeFunc)
    app.Get("/:short_url", shortener.RedirectURL)
    app.Get("/:short_url/stats", shortener.URLStats)
    app.Delete("/:short_url", shortener.DeleteURL)
}