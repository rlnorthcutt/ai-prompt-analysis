package llm

import (
	"errors"
)

// Common errors
var (
	ErrAPIKeyNotSet     = errors.New("API key not set")
	ErrInvalidResponse  = errors.New("invalid response from LLM API")
	ErrRequestFailed    = errors.New("request to LLM API failed")
	ErrResponseParsing  = errors.New("failed to parse LLM API response")
	ErrInvalidPrompt    = errors.New("invalid or empty prompt")
)

// PromptAnalysis represents the structured analysis of a prompt
type PromptAnalysis struct {
	TokenCount   int    `json:"tokenCount"`
	PromptType   string `json:"promptType"`
	ContainsPII  bool   `json:"containsPII"`
	IsSuspicious bool   `json:"isSuspicious"`
	RiskScore    int    `json:"riskScore"`
}

// LLM defines the interface for language model providers
type LLM interface {
	// Name returns the name of the LLM provider
	Name() string
	
	// AnalyzePrompt analyzes a prompt and returns a structured analysis
	AnalyzePrompt(prompt string) (*PromptAnalysis, error)
	
	// IsAvailable checks if the LLM provider is available (API key set, etc.)
	IsAvailable() bool
}