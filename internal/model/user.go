package main

// using capital letter to make it exportable
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // hashed
}
