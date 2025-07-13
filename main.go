package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to run server", err)
	}
}
