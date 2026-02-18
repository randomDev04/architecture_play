package models

import "time"

type Task struct {
	ID          int        `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Details     string     `json:"details"`
	IsCompleted bool       `json:"is_completed"`
	DueDate     *time.Time `json:"due_date"` // pointer = nullable
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTaskRequest = what client sends to POST /tasks
type CreateTaskRequest struct {
	Title   string     `json:"title"   binding:"required"`
	Details string     `json:"details"`
	DueDate *time.Time `json:"due_date"`
}

// UpdateTaskRequest = what client sends to PATCH /tasks/:id
type UpdateTaskRequest struct {
	Title       *string    `json:"title"`        // pointer = optional
	Details     *string    `json:"details"`      // pointer = optional
	IsCompleted *bool      `json:"is_completed"` // pointer = optional
	DueDate     *time.Time `json:"due_date"`     // pointer = optional
}
