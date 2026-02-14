package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rendomDev/task-manager-api/config"
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

	//Health Check
	r.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"message": "server is running",
		})
	})

	// Test 404 handling
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error": "route not found",
		})
	})

	// Start server
	log.Fatal(r.Run(":8080"))
}
