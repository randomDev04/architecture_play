package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Installed sql and goose
func Init() {
	var err error

	DB, err = sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	log.Println("Connected to PostgreSQL")
}
