package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rendomDev/task-manager-api/config"
	"github.com/rendomDev/task-manager-api/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("6BsfFdczBVP9jHCG4bEIHIMrPwqsXl9f8lNIWoNitC0=")

// Register handles POST /api/v1/auth/register
func Register(c *gin.Context) {
	var req models.RegisterRequest

	// Step 1. Bind JSON request to struct + validate
	// binding:"required" and binding:"email" are cheked here
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Step 2. Check if email already exists
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email =$1)"
	err := config.DB.QueryRow(checkQuery, req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})

		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already registered",
		})
		return
	}

	// Step 3. Hash password with bcrypt
	// Cose 10 = good balance between security and speed
	// Higer = more secute but slower
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		10, // cost
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash passowrd",
		})
		return
	}

	// Step 4. Insert user into database
	userID := uuid.New().String()
	insertQuery :=
		`
	INSERT INTO users (id, name, email, password_hash, token_version)
	VALUES ($1, $2, $3, $4, 0)
	RETURNING created_at, updated_at
	`

	var user models.User
	err = config.DB.QueryRow(
		insertQuery,
		userID,
		req.Name,
		req.Email,
		string(hashedPassword),
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"details": err.Error(),
		})

		return
	}

	// Set user fields for response
	user.ID = userID
	user.Name = req.Name
	user.Email = req.Email
	user.TokenVersion = 0

	// Step 5. Generate JWT token
	token, err := generateJWT(userID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})

		return
	}

	// Step 6: Return success response
	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// generateJWT creates a JWT token with user_id and token_version
func generateJWT(userID string, tokenVersion int) (string, error) {
	// JWT claims = data stored inside the token
	claims := jwt.MapClaims{
		"user_id":       userID,
		"token_version": tokenVersion,
		"exp":           time.Now().Add(24 * time.Hour).Unix(), // expires in 24 hours
		"iat":           time.Now().Unix(),                     // issued at
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
