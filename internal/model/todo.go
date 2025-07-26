package model

type Todo struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"` // Linked to user
	Title  string `json:"title"`
	Done   bool   `json:"done"`
}
