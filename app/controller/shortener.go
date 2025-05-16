package shortener

import (
    "github.com/gofiber/fiber/v2"
)

func WelcomeFunc(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Welcome to URL Shortener"})
}
func ShortenURL(c *fiber.Ctx) error {
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "ShortenURL"})
}

func ListURLs(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "ListURLs"})
}

func RedirectURL(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "RedirectURL"})
}

func URLStats(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "URLStats"})
}

func DeleteURL(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "DeleteURL"})
}