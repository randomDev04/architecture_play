package db

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

func Init() {
	connStr := os.Getenv("DATABASE_URL")
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
}
