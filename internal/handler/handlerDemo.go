package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/rlnorthcutt/ai-prompt-analysis/internal/llm"
)

// TemplateData holds data for UI templates
type TemplateData struct {
	ClaudeAvailable  bool
	ChatGPTAvailable bool
	Error            string
	TokenCount       int
	PromptType       string
	ContainsPII      bool
	IsSuspicious     bool
	RiskScore        int
	Latency          int64
	RawJSON          string
}

// HandleDemoUI handles the demo UI page
func (h *Handler) HandleDemoUI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only serve if demo UI is enabled
		if !h.config.Server.DemoUI {
			http.NotFound(w, r)
			return
		}

		// Prepare template data
		data := TemplateData{
			ClaudeAvailable:  h.claudeAPI.IsAvailable(),
			ChatGPTAvailable: h.chatGPTAPI.IsAvailable(),
		}

		// Render template
		if err := h.templates.ExecuteTemplate(w, "analyze.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// HandleFormSubmit handles form submission from the demo UI
func (h *Handler) HandleFormSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only serve if demo UI is enabled
		if !h.config.Server.DemoUI {
			http.NotFound(w, r)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			renderErrorResult(w, h.templates, "Failed to parse form: "+err.Error())
			return
		}

		// Get form values
		providerName := r.FormValue("provider")
		promptText := r.FormValue("prompt")

		// Validate input
		if promptText == "" {
			renderErrorResult(w, h.templates, "Prompt cannot be empty")
			return
		}

		// Choose provider
		var selectedProvider llm.LLM
		switch providerName {
		case "claude":
			selectedProvider = h.claudeAPI
		case "chatgpt":
			selectedProvider = h.chatGPTAPI
		default:
			renderErrorResult(w, h.templates, "Invalid provider selected")
			return
		}

		// Check if provider is available
		if !selectedProvider.IsAvailable() {
			renderErrorResult(w, h.templates, selectedProvider.Name()+" API key not set")
			return
		}

		// Start timing the response
		startTime := time.Now()

		// Analyze the prompt
		analysis, err := selectedProvider.AnalyzePrompt(promptText)
		if err != nil {
			renderErrorResult(w, h.templates, "Error analyzing prompt: "+err.Error())
			return
		}

		// Calculate latency in milliseconds
		latency := time.Since(startTime).Milliseconds()

		// Create extended response with latency
		response := AnalysisResponse{
			PromptAnalysis: *analysis,
			Latency:        latency,
		}

		// Convert to JSON for raw display
		rawJSON, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			renderErrorResult(w, h.templates, "Error formatting JSON: "+err.Error())
			return
		}

		// Prepare template data
		data := TemplateData{
			TokenCount:   analysis.TokenCount,
			PromptType:   analysis.PromptType,
			ContainsPII:  analysis.ContainsPII,
			IsSuspicious: analysis.IsSuspicious,
			RiskScore:    analysis.RiskScore,
			Latency:      latency,
			RawJSON:      string(rawJSON),
		}

		// Render template
		if err := h.templates.ExecuteTemplate(w, "result.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// renderErrorResult renders an error response template
func renderErrorResult(w http.ResponseWriter, tmpl *template.Template, errMsg string) {
	data := TemplateData{
		Error: errMsg,
	}
	if err := tmpl.ExecuteTemplate(w, "result.html", data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}