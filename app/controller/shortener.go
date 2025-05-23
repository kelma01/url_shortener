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
//rediste tutulma structure'inda key degerleri short url kismi value'si ise original_url kisim oluyor
// {"zRWzOuQT": "https://www.google.com"} gibi gibi.
var redisClient = redis.NewClient(&redis.Options{
	Addr: "url-shortener-redis:6379",
})

//convert algosu icin gerekli
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ShortenURL(c *fiber.Ctx) error {
	ip := c.IP()

    count, err := redisClient.Incr(context.Background(), ip).Result()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
    }
    if count == 1 {
        redisClient.Expire(context.Background(), ip, time.Minute)	//1 dakika icinde rate limit sifirlanir
    }
    if count > 5 {	//belirilen sure icerisinde(1dk) maks 5 post req yollanabilir
        return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Rate limit exceeded. Please try again later."})
    }

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


	_ , err = database.DB.Exec("INSERT INTO url_table (original_url, short_url, created_at, expires_at) VALUES ($1, $2, $3, $4)", body.OriginalURL, short, now, expiresAt)
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

func DeleteURL(c *fiber.Ctx) error {
	shortURL := c.Params("short_url")
	if shortURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no url detected for redirecting"})
	}
	
	//redisi sifirlamak yeterli cunku zaten redirect ederken redis'teki ttl'i de kontrol ederek redirect ediyor, redis.del islemi yeterli
	_, _ = redisClient.Del(context.Background(), shortURL).Result()

	_, _ = database.DB.Exec("UPDATE url_table SET deleted_at=$1 WHERE short_url = $2", time.Now(), shortURL)
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "URL deleted successfully"})
}