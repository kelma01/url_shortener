package migrations

import (
	"database/sql"
	"log"
)

//main ile db eslenmesi icin kullanilacak db querysi
func Migrate(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS url_table(
		id SERIAL PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_url VARCHAR(16) UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP,
		expires_at TIMESTAMP,
		usage_count INTEGER DEFAULT 0);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}
	return err
}
