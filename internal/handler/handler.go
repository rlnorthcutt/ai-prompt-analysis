package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rlnorthcutt/ai-prompt-analysis/internal/config"
	"github.com/rlnorthcutt/ai-prompt-analysis/internal/llm"
	"github.com/rlnorthcutt/ai-prompt-analysis/internal/prompt"
)

// Handler provides HTTP handlers for the API
type Handler struct {
	claudeAPI  llm.LLM
	chatGPTAPI llm.LLM
	config     *config.Config
	templates  *template.Template
	routes     Routes
}

// Routes defines the API endpoints
type Routes struct {
	Claude     string
	ChatGPT    string
	Demo       string
	DemoSubmit string
}

// AnalysisResponse extends the prompt analysis with latency information
type AnalysisResponse struct {
	llm.PromptAnalysis
	Latency int64 `json:"latency"` // Response latency in milliseconds
}

// NewHandler creates a new Handler instance with initialized LLM providers
func NewHandler(cfg *config.Config) *Handler {
	// Initialize LLM providers
	claudeAPI := llm.NewClaude(cfg)
	chatGPTAPI := llm.NewChatGPT(cfg)

	// Load templates
	templates := template.Must(template.ParseGlob("internal/handler/templates/*.html"))

	// Define routes
	routes := Routes{
		Claude:     "/analyze/claude",
		ChatGPT:    "/analyze/chatgpt",
		Demo:       "/analyze",
		DemoSubmit: "/analyze/submit",
	}

	return &Handler{
		claudeAPI:  claudeAPI,
		chatGPTAPI: chatGPTAPI,
		config:     cfg,
		templates:  templates,
		routes:     routes,
	}
}

// GetLLMProviders returns the initialized LLM providers
func (h *Handler) GetLLMProviders() (claude llm.LLM, chatgpt llm.LLM) {
	return h.claudeAPI, h.chatGPTAPI
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes() {
	// Register API endpoints
	http.HandleFunc(h.routes.Claude, h.ClaudeHandler())
	http.HandleFunc(h.routes.ChatGPT, h.ChatGPTHandler())
	
	// Demo UI (only if enabled in config)
	if h.config.Server.DemoUI {
		http.HandleFunc(h.routes.Demo, h.HandleDemoUI())
		http.HandleFunc(h.routes.DemoSubmit, h.HandleFormSubmit())
	}
}

// StartServer starts the HTTP server
func (h *Handler) StartServer() error {
	// Get server address
	serverAddr := fmt.Sprintf(":%s", h.config.Server.Port)
	baseURL := fmt.Sprintf("http://localhost%s", serverAddr)
	
	// Log server information
	log.Printf("Starting server on %s...\n", serverAddr)
	log.Printf("Claude API available: %v", h.claudeAPI.IsAvailable())
	log.Printf("ChatGPT API available: %v", h.chatGPTAPI.IsAvailable())
	log.Printf("Demo UI enabled: %v", h.config.Server.DemoUI)
	log.Printf("Endpoints:")
	log.Printf("  - Claude:  %s%s", baseURL, h.routes.Claude)
	log.Printf("  - ChatGPT: %s%s", baseURL, h.routes.ChatGPT)
	if h.config.Server.DemoUI {
		log.Printf("  - Demo UI: %s%s", baseURL, h.routes.Demo)
	}
	
	// Start the server
	return http.ListenAndServe(serverAddr, nil)
}

// HandleAnalyze handles the generic prompt analysis endpoint
func (h *Handler) HandleAnalyze(provider llm.LLM) http.HandlerFunc {
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

		// Start timing the response
		startTime := time.Now()

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

		// Calculate latency in milliseconds
		latency := time.Since(startTime).Milliseconds()

		// Create extended response with latency
		response := AnalysisResponse{
			PromptAnalysis: *analysis,
			Latency:        latency,
		}

		// Return the analysis as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// ClaudeHandler returns the handler for Claude analysis
func (h *Handler) ClaudeHandler() http.HandlerFunc {
	return h.HandleAnalyze(h.claudeAPI)
}

// ChatGPTHandler returns the handler for ChatGPT analysis
func (h *Handler) ChatGPTHandler() http.HandlerFunc {
	return h.HandleAnalyze(h.chatGPTAPI)
}