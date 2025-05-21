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
		"short_url":    "http://localhost:8080/" + short,
		"expires_at":   expiresAt,
		"usage_count":  0,
	})
}

func generateShortURL(size int) string {
	for {
		rand.Seed(time.Now().UnixNano())
		b := make([]byte, size)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		short := string(b)

		//shorted url'lerin unique olmasi gerekli, bu yuzden her uretilen icin oncelikle db'de kontrolu yapiliyor
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM url_table WHERE short_url=$1)", short).Scan(&exists)
		if err != nil {
			return short
		}
		if !exists {
			return short
		}
	}
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

	//rediste ttl kontorlü
    val, err := redisClient.Get(context.Background(), shortURL).Result()
    if err == redis.Nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL expired or not found"})
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

	//usage_count arttirici
    _, _ = database.DB.Exec("UPDATE url_table SET usage_count = usage_count + 1 WHERE short_url = $1", shortURL)

    return c.Redirect(val, fiber.StatusTemporaryRedirect)
}

//yuz kez denedim biseler olmadi copilotla bitirdim
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
func StatsURL(c *fiber.Ctx) error {
	shortURL := c.Params("short_url")
	if shortURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "status bad request"})
	}

	type Result struct {
		OriginalURL string     `json:"original_url"`
		ShortURL    string     `json:"short_url"`
		UsageCount  int        `json:"usage_count"`
		CreatedAt   time.Time  `json:"created_at"`
		DeletedAt   *time.Time `json:"deleted_at"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	var res Result

	err := database.DB.QueryRow("SELECT original_url, short_url, usage_count, created_at, deleted_at, expires_at FROM url_table WHERE short_url=$1", shortURL,).Scan(&res.OriginalURL, &res.ShortURL, &res.UsageCount, &res.CreatedAt, &res.DeletedAt, &res.ExpiresAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	}
	return c.JSON(res)
}