package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB is the global database connection pool
// Accessible from anywhere in the application

var DB *sql.DB

func ConnectDB() {
	// Connection string - tells GO where PostgreSQL is
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"admin",
		"password",
		"task_manager",
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}

	// sql.Open doesn't actually connect yet
	// Ping forces and actual connection to the database
	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	// Connection pool settings -this is architectural
	DB.SetMaxOpenConns(25)     // MAX 25 simualtaneous connections to DB
	DB.SetMaxIdleConns(5)      // KEEP 5 connections ready even when idle
	DB.SetConnMaxLifetime(300) // Recycle connections after 5 mins

	fmt.Println("Database connected successfully")
}
