package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectPostgres() {

	dsn := os.Getenv("DATABASE_URL")

	var err error

	for i := 1; i <= 10; i++ {

		DB, err = pgxpool.New(context.Background(), dsn)
		if err == nil {

			err = DB.Ping(context.Background())
			if err == nil {
				fmt.Println("Connected to PostgreSQL")

				initSchema() // initialize schema design
				return
			}
		}

		log.Printf("Database not ready... retrying (%d/10)\n", i)
		time.Sleep(3 * time.Second)
	}

	log.Fatal("Could not connect to PostgreSQL after retries")
}

func initSchema() {
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	payload JSONB NOT NULL,
	status TEXT NOT NULL,
	retries INT DEFAULT 0,
	max_retries INT DEFAULT 3,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
	);
	`

	_, err := DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Failed to initialize database schema:", err)
	}

	fmt.Println("Jobs table ready")
}
