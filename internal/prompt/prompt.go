package prompt

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// Request represents the incoming prompt analysis request
type Request struct {
	Prompt string `json:"prompt"`
}

// Validate validates a prompt request
func (r *Request) Validate() error {
	if strings.TrimSpace(r.Prompt) == "" {
		return errors.New("prompt cannot be empty")
	}
	return nil
}

// ExtractJSON attempts to extract valid JSON from text that might contain additional content
func ExtractJSON(text string) string {
	// Look for JSON between curly braces
	re := regexp.MustCompile(`\{.*\}`)
	match := re.FindString(text)
	if match != "" {
		return match
	}
	return text
}

// ParseJSON parses a JSON string into a structured type
func ParseJSON(jsonStr string, target interface{}) error {
	cleanJSON := ExtractJSON(jsonStr)
	return json.Unmarshal([]byte(cleanJSON), target)
}