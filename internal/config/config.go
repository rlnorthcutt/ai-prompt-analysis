package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Server struct {
		Port   string `mapstructure:"port"`
		DemoUI bool   `mapstructure:"demoui"`
	} `mapstructure:"server"`

	Claude struct {
		APIURL      string  `mapstructure:"api_url"`
		ModelID     string  `mapstructure:"model_id"`
		MaxTokens   int     `mapstructure:"max_tokens"`
		Temperature float64 `mapstructure:"temperature"`
		Version     string  `mapstructure:"version"`
	} `mapstructure:"claude"`

	ChatGPT struct {
		APIURL      string  `mapstructure:"api_url"`
		ModelID     string  `mapstructure:"model_id"`
		MaxTokens   int     `mapstructure:"max_tokens"`
		Temperature float64 `mapstructure:"temperature"`
	} `mapstructure:"chatgpt"`

	Analysis struct {
		SystemPrompt string `mapstructure:"system_prompt"`
	} `mapstructure:"analysis"`
}

// Load loads configuration from config.yaml
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	// Load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		return viper.ReadInConfig()
	}
	return nil
}

// GetEnv gets an environment variable
func GetEnv(key string) string {
	return os.Getenv(key)
}