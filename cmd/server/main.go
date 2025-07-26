package main

import (
	"go_todo_backend/internal/handler"
	"go_todo_backend/internal/repository"
	"go_todo_backend/pkg/db"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	// leading env file
	err := godotenv.Load("../../.env") // pass the file to connect to the env
	if err != nil {
		log.Println("⚠️  No .env file found. Using system environment.")
	}

	db.Init()
	defer db.DB.Close()

	todoRepo := repository.NewPostgresTodoRepo(db.DB)

	todoHandler := handler.NewTodoHandler(todoRepo)

	r := chi.NewRouter()

	r.Route("/api/todos", func(r chi.Router) {
		r.Get("/", todoHandler.GetTodos)
		r.Post("/", todoHandler.CreateTodo)
	})

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to run server", err)
	}
}
