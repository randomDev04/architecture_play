package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rendomDev/task-manager-api/config"
	"github.com/rendomDev/task-manager-api/handlers"
	"github.com/rendomDev/task-manager-api/middleware"
)

func main() {
	// Connect to database FIRST
	// Of DB fails, server should not start at all
	config.ConnectDB()

	r := gin.Default()
	r.SetTrustedProxies(nil) // trust no proxies for now

	// Tell Gin to return 405 instead of 404
	// when method doesn't match
	r.HandleMethodNotAllowed = true

	// Custom 405 handler
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error": "method not allowed",
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "route not found"})
	})

	//Health Check
	r.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"message": "server is running",
		})
	})

	// Auth routes
	r.POST("/api/v1/auth/register", handlers.Register)
	r.POST("/api/v1/auth/login", handlers.Login)

	protected := r.Group("./api/v1")
	protected.Use(middleware.AuthMiddleware()) // Apply middleware to all routes in this group
	{
		// Test endpoint to verify middlerware works
		protected.GET("/me", func(ctx *gin.Context) {
			// Get user_id that middleware stored in context
			userID, _ := ctx.Get("user_id")

			ctx.JSON(200, gin.H{
				"message": "you are authenticated",
				"user_id": userID,
			})
		})
	}

	// Task endpoints
	protected.GET("/tasks", handlers.GetTasks)
	protected.POST("/tasks", handlers.CreateTask)
	protected.GET("/tasks/:id", handlers.GetTask)
	protected.PATCH("/tasks/:id", handlers.UpdateTask)
	protected.DELETE("/tasks/:id", handlers.DeleteTask)

	// Start server
	log.Fatal(r.Run(":8080"))
}
