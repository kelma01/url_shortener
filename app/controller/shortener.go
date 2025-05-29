package shortener

import (
	"context"
	"math/rand"
	"time"

	"url_shortener/app/entities"
	"url_shortener/internal/database"
	"url_shortener/internal/redis"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

//redis client setupu
//rediste tutulma structure'inda key degerleri short url kismi value'si ise original_url kisim oluyor
// {"zRWzOuQT": "https://www.google.com"} gibi gibi.
var redisClient = redis.RedisClient



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

	url := entities.URL{
		OriginalURL: body.OriginalURL,
		ShortURL: short,
		CreatedAt: now,
		ExpiresAt: &expiresAt,
	}
	if err := database.DB.Create(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
		err := database.DB.Model(&entities.URL{}).Select("count(*) > 0").Where("short_url = ?", short).Find(&exists).Error
		if err != nil {
			continue
		} else {
			return short
		}
	}
}
func ListURLs(c *fiber.Ctx) error {
    var urls []entities.URL

    if err := database.DB.Find(&urls).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
    if err == redis.RedisNil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL expired or not found"})
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

	//usage_count arttirici
	if err := database.DB.Model(&entities.URL{}).Where("short_url = ?", shortURL).Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

    return c.Redirect(val, fiber.StatusTemporaryRedirect)
}

func StatsURL(c *fiber.Ctx) error {
	shortURL := c.Params("short_url")
	if shortURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "status bad request"})
	}

	var res entities.URL

	if err := database.DB.Where("short_url = ?", shortURL).First(&res).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	if err := database.DB.Model(&entities.URL{}).Where("short_url = ?", shortURL).Update("deleted_at", time.Now()).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "URL deleted succeessfully"})
}