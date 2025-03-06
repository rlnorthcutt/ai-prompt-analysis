package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rlnorthcutt/ai-prompt-analysis/internal/config"
	"github.com/rlnorthcutt/ai-prompt-analysis/internal/handler"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create handler with LLM providers
	h := handler.NewHandler(cfg)

	// Register routes
	h.RegisterRoutes()

	// Start the server
	if err := h.StartServer(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}