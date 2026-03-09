package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rendomDev/task-manager-api/config"
)

// Must match the secret in handlers/auth.go
var jwtSecret = []byte("6BsfFdczBVP9jHCG4bEIHIMrPwqsXl9f8lNIWoNitC0=")

// AuthMiddleware verifies JWT token on every request
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Get Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			c.Abort() // Stop the request chain, don't call handler
			return
		}

		// Step 2: Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Step 3: Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Step 4: Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// Get user_id and token_version from JWT
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id in token",
			})
			c.Abort()
			return
		}

		tokenVersion, ok := claims["token_version"].(float64) // JSON numbers are float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token_version in token",
			})
			c.Abort()
			return
		}

		// Step 5: Verify token_version against database
		var dbTokenVersion int
		query := "SELECT token_version FROM users WHERE id = $1"
		err = config.DB.QueryRow(query, userID).Scan(&dbTokenVersion)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			c.Abort()
			return
		}

		// Check if token has been invalidated (logout)
		if int(tokenVersion) != dbTokenVersion {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token has been revoked",
			})
			c.Abort()
			return
		}

		// Step 6: Store user_id in context for handlers to use
		c.Set("user_id", userID)

		// Continue to the handler
		c.Next()
	}
}
