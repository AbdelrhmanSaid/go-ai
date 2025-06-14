package main

import (
	"fmt"
	"os"

	"github.com/AbdelrhmanSaid/go-ai/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize gin
	r := gin.Default()

	// Register routes
	r.POST("/chat/completions", handlers.ChatCompletions)

	// Define port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server on port
	r.Run(fmt.Sprintf(":%s", port))
}
