package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm/logger"

	"go.opentelemetry.io/otel/attribute"

    "github.com/joho/godotenv"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
    
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
        os.Getenv("DB_HOST"), //k8s ise POSTGRES_HOST olmali
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
    if err != nil {
        return err
    }
    if err := db.Use(otelgorm.NewPlugin(otelgorm.WithAttributes(attribute.String("DB_HOST", os.Getenv("DB_HOST"))))); err != nil {
		panic(err)
	}
    DB = db

    if err := AutoMigrate(models...); err != nil {
        return err
    }

    return nil
}

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