package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

//db config
func Connect() error {
    connection := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_PORT", "5432"),
        getEnv("DB_USER", "test"),
        getEnv("DB_PASSWORD", "test"),
        getEnv("DB_NAME", "url_shortener"),
    )
    db, err := sql.Open("postgres", connection)
    if err != nil {
        return err
    }
    if err := db.Ping(); err != nil {
        return err
    }
    DB = db
    return nil
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}