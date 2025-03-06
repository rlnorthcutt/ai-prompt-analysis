package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
	"github.com/yourusername/prompt-analyzer/config"
	"github.com/yourusername/prompt-analyzer/llm"
	"github.com/yourusername/prompt-analyzer/prompt"
)

// Handler provides HTTP handlers for the API
type Handler struct {
	claudeAPI  llm.LLM
	chatGPTAPI llm.LLM
}

// NewHandler creates a new Handler instance
func NewHandler(claudeAPI, chatGPTAPI llm.LLM) *Handler {
	return &Handler{
		claudeAPI:  claudeAPI,
		chatGPTAPI: chatGPTAPI,
	}
}

// handleAnalyze handles the generic prompt analysis endpoint
func (h *Handler) handleAnalyze(provider llm.LLM) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if provider is available
		if !provider.IsAvailable() {
			http.Error(w, fmt.Sprintf("%s API key not set", provider.Name()), http.StatusServiceUnavailable)
			return
		}

		// Parse request body
		var req prompt.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Validate input
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Analyze the prompt
		analysis, err := provider.AnalyzePrompt(req.Prompt)
		if err != nil {
			// Handle specific errors
			switch {
			case strings.Contains(err.Error(), "API key not set"):
				http.Error(w, fmt.Sprintf("%s API key not set", provider.Name()), http.StatusServiceUnavailable)
			default:
				http.Error(w, fmt.Sprintf("Error analyzing prompt: %v", err), http.StatusInternalServerError)
			}
			return
		}

		// Return the analysis as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(analysis); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

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

	// Initialize LLM providers
	claudeAPI := llm.NewClaude(cfg)
	chatGPTAPI := llm.NewChatGPT(cfg)

	// Create handler
	handler := NewHandler(claudeAPI, chatGPTAPI)

	// Register endpoints
	http.HandleFunc("/analyze/claude", handler.handleAnalyze(claudeAPI))
	http.HandleFunc("/analyze/chatgpt", handler.handleAnalyze(chatGPTAPI))

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s...\n", serverAddr)
	log.Printf("Claude API available: %v", claudeAPI.IsAvailable())
	log.Printf("ChatGPT API available: %v", chatGPTAPI.IsAvailable())
	
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}