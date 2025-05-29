package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

    "github.com/joho/godotenv"
    
    _ "url_shortener/internal/opentelemetry"
)

var DB *gorm.DB

func init() {
    if err := Connect(); err != nil {
        log.Fatalf("DB connection failed: %v", err)
        os.Exit(1)
    }
    log.Println("PostgreSQL OK")
}


//db config
//localhost ve docker
func Connect(models ...interface{}) error {
    if err := godotenv.Load(); err != nil {
        log.Println(err)
    }
    connection := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        getEnv("DB_HOST", "null"),
        getEnv("DB_PORT", "null"),
        getEnv("DB_USER", "null"),
        getEnv("DB_PASSWORD", "null"),
        getEnv("DB_NAME", "null"),
    )

    db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
    if err != nil {
        return err
    }
    DB = db

    if err := AutoMigrate(models...); err != nil {
        return err
    }

    return nil
}
//kubernetes
/* func Connect() error {
    if err := godotenv.Load(); err != nil {
        log.Println(err)
    }
    connection := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        getEnv("POSTGRES_HOST", "null"),
        getEnv("DB_PORT", "null"),
        getEnv("DB_USER", "null"),
        getEnv("DB_PASSWORD", "null"),
        getEnv("DB_NAME", "null"),
    )

    db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
    if err != nil {
        return err
    }
    DB = db

    if err := AutoMigrate(models...); err != nil {
        return err
    }
    
    return nil
} */

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
func AutoMigrate(models ...interface{}) error {
    if DB == nil {
        return fmt.Errorf("database connection is not initialized")
    }
    return DB.AutoMigrate(models...)
}