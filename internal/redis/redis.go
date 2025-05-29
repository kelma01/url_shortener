package redis

import (
	"log"
	"os"
	
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisNil = redis.Nil

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}
	RedisAddr := getEnv("REDIS_ADDR", "null")
	RedisPort := getEnv("REDIS_PORT", "null")

	RedisClient = redis.NewClient(&redis.Options{
		Addr: RedisAddr + ":" + RedisPort,
	})

	log.Println("Redis OK")
}
