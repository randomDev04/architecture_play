package handler

import (
	"encoding/json"
	"go_todo_backend/internal/repository"
	"net/http"
)

type TodoHandler struct {
	Repo repository.TodoRepository
}

func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{Repo: repo}
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	// hardcoded user ID for now
	todos, err := h.Repo.GetTodos(1)
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// TODO: decode JSON + insert
	w.WriteHeader(http.StatusNotImplemented)
}
