package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rlnorthcutt/ai-prompt-analysis/internal/config"
	"github.com/rlnorthcutt/ai-prompt-analysis/internal/prompt"
)

// Claude implements the LLM interface for Claude API
type Claude struct {
	config *config.Config
}

// ClaudeRequest represents the request structure for Claude API
type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []ClaudeMessage `json:"messages"`
	System      string          `json:"system"`
	Temperature float64         `json:"temperature"`
}

// ClaudeMessage represents a message in Claude API request
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents the response structure from Claude API
type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// NewClaude creates a new Claude instance
func NewClaude(config *config.Config) *Claude {
	return &Claude{
		config: config,
	}
}

// Name returns the name of the LLM provider
func (c *Claude) Name() string {
	return "Claude"
}

// IsAvailable checks if the Claude API is available
func (c *Claude) IsAvailable() bool {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	return apiKey != ""
}

// AnalyzePrompt analyzes a prompt using Claude API
func (c *Claude) AnalyzePrompt(promptText string) (*PromptAnalysis, error) {
	// Get API key from environment
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyNotSet
	}

	// Create Claude API request payload
	claudeReq := ClaudeRequest{
		Model:     c.config.Claude.ModelID,
		MaxTokens: c.config.Claude.MaxTokens,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("Analyze this prompt: %s", promptText),
			},
		},
		System:      c.config.Analysis.SystemPrompt,
		Temperature: c.config.Claude.Temperature,
	}

	// Convert request to JSON
	reqBody, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.config.Claude.APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", c.config.Claude.Version)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w with status %d: %s", ErrRequestFailed, resp.StatusCode, string(body))
	}

	// Parse Claude's response
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParsing, err)
	}

	// Extract and parse the JSON response from Claude
	if len(claudeResp.Content) == 0 {
		return nil, ErrInvalidResponse
	}

	// Parse the analysis
	jsonText := claudeResp.Content[0].Text
	var analysis PromptAnalysis
	if err := prompt.ParseJSON(jsonText, &analysis); err != nil {
		return nil, fmt.Errorf("%w: %v\nRaw response: %s", ErrResponseParsing, err, jsonText)
	}

	return &analysis, nil
}
