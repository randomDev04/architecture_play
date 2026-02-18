package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" = never sent in JSON response
	TokenVersion int       `json:"token_version"`
	CreatedAt    time.Time `json:"create_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Registration Request = when client sends to Auth/registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login Request
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse = what server sends back after login/register
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
