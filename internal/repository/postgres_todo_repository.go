package repository

import (
	"database/sql"
	"go_todo_backend/internal/model"
)

type postgresTodoRepo struct {
	DB *sql.DB
}

func NewPostgresTodoRepo(db *sql.DB) TodoRepository {
	return &postgresTodoRepo{DB: db}
}
func (r *postgresTodoRepo) GetTodos(userID int) ([]model.Todo, error) {
	rows, err := r.DB.Query("SELECT id, user_id, title, done FROM todos WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (r *postgresTodoRepo) GetTodoByID(userID, todoID int) (*model.Todo, error) {
	row := r.DB.QueryRow("SELECT id, user_id, title, done FROM todos WHERE id=$1 AND user_id=$2", todoID, userID)
	var t model.Todo
	if err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Done); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *postgresTodoRepo) CreateTodo(todo *model.Todo) error {
	return r.DB.QueryRow(
		"INSERT INTO todos (user_id, title, done) VALUES ($1, $2, $3) RETURNING id",
		todo.UserID, todo.Title, todo.Done,
	).Scan(&todo.ID)
}

func (r *postgresTodoRepo) UpdateTodo(todo *model.Todo) error {
	_, err := r.DB.Exec(
		"UPDATE todos SET title=$1, done=$2 WHERE id=$3 AND user_id=$4",
		todo.Title, todo.Done, todo.ID, todo.UserID,
	)
	return err
}

func (r *postgresTodoRepo) DeleteTodo(userID, todoID int) error {
	_, err := r.DB.Exec("DELETE FROM todos WHERE id=$1 AND user_id=$2", todoID, userID)
	return err
}
