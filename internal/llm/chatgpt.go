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

// ChatGPT implements the LLM interface for OpenAI's ChatGPT API
type ChatGPT struct {
	config *config.Config
}

// ChatGPTRequest represents the request structure for ChatGPT API
type ChatGPTRequest struct {
	Model       string           `json:"model"`
	Messages    []ChatGPTMessage `json:"messages"`
	MaxTokens   int              `json:"max_tokens"`
	Temperature float64          `json:"temperature"`
}

// ChatGPTMessage represents a message in ChatGPT API request
type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTResponse represents the response structure from ChatGPT API
type ChatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewChatGPT creates a new ChatGPT instance
func NewChatGPT(config *config.Config) *ChatGPT {
	return &ChatGPT{
		config: config,
	}
}

// Name returns the name of the LLM provider
func (c *ChatGPT) Name() string {
	return "ChatGPT"
}

// IsAvailable checks if the ChatGPT API is available
func (c *ChatGPT) IsAvailable() bool {
	apiKey := os.Getenv("OPENAI_API_KEY")
	return apiKey != ""
}

// AnalyzePrompt analyzes a prompt using ChatGPT API
func (c *ChatGPT) AnalyzePrompt(promptText string) (*PromptAnalysis, error) {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyNotSet
	}

	// Create ChatGPT API request payload
	chatGPTReq := ChatGPTRequest{
		Model: c.config.ChatGPT.ModelID,
		Messages: []ChatGPTMessage{
			{
				Role:    "system",
				Content: c.config.Analysis.SystemPrompt,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Analyze this prompt: %s", promptText),
			},
		},
		MaxTokens:   c.config.ChatGPT.MaxTokens,
		Temperature: c.config.ChatGPT.Temperature,
	}

	// Convert request to JSON
	reqBody, err := json.Marshal(chatGPTReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.config.ChatGPT.APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

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

	// Parse ChatGPT's response
	var chatGPTResp ChatGPTResponse
	if err := json.Unmarshal(body, &chatGPTResp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParsing, err)
	}

	// Extract and parse the JSON response from ChatGPT
	if len(chatGPTResp.Choices) == 0 {
		return nil, ErrInvalidResponse
	}

	// Parse the analysis
	jsonText := chatGPTResp.Choices[0].Message.Content
	var analysis PromptAnalysis
	if err := prompt.ParseJSON(jsonText, &analysis); err != nil {
		return nil, fmt.Errorf("%w: %v\nRaw response: %s", ErrResponseParsing, err, jsonText)
	}

	return &analysis, nil
}