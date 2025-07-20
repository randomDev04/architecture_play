package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	// add database here

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to run server", err)
	}
}
