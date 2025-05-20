package shortener

import (
    "github.com/gofiber/fiber/v2"
	"math/rand"
	"time"
	"url_shortener/internal/database"
	"github.com/redis/go-redis/v9"
	"context"
)

//redis client setupu
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

//convert algosu icin gerekli
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

/*
postmanda girilecek payload:
{
	"original_url": "https://www.google.com"
}
*/
//TODO: short url create ederken unique değilse lineer probing ile yenisi aranacak. 
func ShortenURL(c *fiber.Ctx) error {
	now := time.Now()
	expiresAt := now.Add(5 * time.Minute)

	type reqBody struct {
		OriginalURL string     `json:"original_url"`
	}
	var body reqBody

	if err := c.BodyParser(&body); err != nil || body.OriginalURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	short := generateShortURL(8)

	redisClient.Set(context.Background(), short, body.OriginalURL, 5*time.Minute)


	_, err := database.DB.Exec("INSERT INTO url_table (original_url, short_url, created_at, expires_at) VALUES ($1, $2, $3, $4)", body.OriginalURL, short, now, expiresAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"created_at":   now,
		"deleted_at":   nil,
		"original_url": body.OriginalURL,
		"short_url":    "http://localhost:3000/" + short,
		"expires_at":   expiresAt,
		"usage_count":  0,
	})
}

func generateShortURL(size int) string {
	rand.Seed(time.Now().UnixNano())
	b:=   make([]byte, size)
	for i := range size {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func ListURLs(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT original_url, short_url, usage_count, created_at, deleted_at, expires_at FROM url_table")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type URL struct {
		OriginalURL string     `json:"original_url"`
		ShortURL    string     `json:"short_url"`
		UsageCount	int		   `json:"usage_count"`
		CreatedAt   time.Time  `json:"created_at"`
		DeletedAt	*time.Time `json:"deleted_at"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	//gosterilecek olan url tuplelarini tutacak
	var urls []URL

	for rows.Next() {
		var u URL
		if err := rows.Scan(&u.OriginalURL, &u.ShortURL, &u.UsageCount, &u.CreatedAt, &u.DeletedAt, &u.ExpiresAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		urls = append(urls, u)
	}

	return c.JSON(urls)
}

func RedirectURL(c *fiber.Ctx) error {
	shortURL := c.Params("short_url")
	if shortURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no url detected for redirecting"})
	}
	var originalURL string
	err := database.DB.QueryRow("SELECT original_url FROM url_table WHERE short_url=$1", shortURL).Scan(&originalURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Redirect(originalURL, fiber.StatusTemporaryRedirect)
}

//yuz kez denedim biseler olmadi copilot bitirdi
func DeleteURL(c *fiber.Ctx) error {
	shortURL := c.Params("short_url")
	if shortURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no url detected for redirecting"})
	}
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM url_table WHERE short_url=$1)", shortURL).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	}

	result, err := database.DB.Exec("DELETE FROM url_table WHERE short_url=$1", shortURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "URL deleted successfully"})
}
