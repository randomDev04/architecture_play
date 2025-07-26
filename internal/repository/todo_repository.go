package repository

import "go_todo_backend/internal/model"

type TodoRepository interface {
	GetTodos(userId int) ([]model.Todo, error)
	GetTodoByID(userID, todoID int) (*model.Todo, error)
	CreateTodo(todo *model.Todo) error
	UpdateTodo(todo *model.Todo) error
	DeleteTodo(userID, todoID int) error
}
