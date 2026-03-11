package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rendomDev/task-manager-api/config"
	"github.com/rendomDev/task-manager-api/models"
)

// GetTasks handles GET /api/v1/tasks
// Returns all tasks for the authenticated user
func GetTasks(c *gin.Context) {
	// Get user_id from context (set by middleware)
	userID, _ := c.Get("user_id")

	// Pagination parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	offset := (pageNum - 1) * limitNum

	// Query tasks for this user only
	query := `
		SELECT id, user_id, title, details, is_completed, due_date, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := config.DB.Query(query, userID, limitNum, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch tasks",
		})
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Details,
			&task.IsCompleted,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	// Get total count for pagination
	var total int
	countQuery := "SELECT COUNT(*) FROM tasks WHERE user_id = $1"
	config.DB.QueryRow(countQuery, userID).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

// CreateTask handles POST /api/v1/tasks
func CreateTask(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Insert task
	query := `
		INSERT INTO tasks (user_id, title, details, due_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	var task models.Task
	err := config.DB.QueryRow(
		query,
		userID,
		req.Title,
		req.Details,
		req.DueDate,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create task",
		})
		return
	}

	// Set remaining fields
	task.UserID = userID.(string)
	task.Title = req.Title
	task.Details = req.Details
	task.DueDate = req.DueDate
	task.IsCompleted = false

	c.JSON(http.StatusCreated, gin.H{
		"task": task,
	})
}

// GetTask handles GET /api/v1/tasks/:id
func GetTask(c *gin.Context) {
	userID, _ := c.Get("user_id")
	taskID := c.Param("id")

	query := `
		SELECT id, user_id, title, details, is_completed, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`

	var task models.Task
	err := config.DB.QueryRow(query, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Details,
		&task.IsCompleted,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

// UpdateTask handles PATCH /api/v1/tasks/:id
func UpdateTask(c *gin.Context) {
	userID, _ := c.Get("user_id")
	taskID := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Verify task exists and belongs to user
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1 AND user_id = $2)"
	config.DB.QueryRow(checkQuery, taskID, userID).Scan(&exists)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}

	// Build dynamic UPDATE query based on what fields are sent
	query := "UPDATE tasks SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argCount := 1

	if req.Title != nil {
		query += ", title = $" + strconv.Itoa(argCount)
		args = append(args, *req.Title)
		argCount++
	}

	if req.Details != nil {
		query += ", details = $" + strconv.Itoa(argCount)
		args = append(args, *req.Details)
		argCount++
	}

	if req.IsCompleted != nil {
		query += ", is_completed = $" + strconv.Itoa(argCount)
		args = append(args, *req.IsCompleted)
		argCount++
	}

	if req.DueDate != nil {
		query += ", due_date = $" + strconv.Itoa(argCount)
		args = append(args, *req.DueDate)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, taskID)
	argCount++

	query += " AND user_id = $" + strconv.Itoa(argCount)
	args = append(args, userID)

	_, err := config.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update task",
		})
		return
	}

	// Fetch updated task
	var task models.Task
	selectQuery := `
		SELECT id, user_id, title, details, is_completed, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	config.DB.QueryRow(selectQuery, taskID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Details,
		&task.IsCompleted,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

// DeleteTask handles DELETE /api/v1/tasks/:id
func DeleteTask(c *gin.Context) {
	userID, _ := c.Get("user_id")
	taskID := c.Param("id")

	// Delete only if task belongs to user
	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"
	result, err := config.DB.Exec(query, taskID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete task",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
