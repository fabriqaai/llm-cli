package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConfigDir returns the llm-cli config directory path
func ConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".llm-cli"
	}
	return filepath.Join(homeDir, ".llm-cli")
}

// ConfigFile returns the models config file path
func ConfigFile() string {
	return filepath.Join(ConfigDir(), "models.config")
}

// ModelConfig maps a model ID to its CLI and model name
type ModelConfig struct {
	CLI      string `json:"cli"`
	ModelArg string `json:"model_arg"`
}

// AppConfig represents the application configuration
type AppConfig struct {
	DefaultModel string                 `json:"default_model"`
	Models       map[string]ModelConfig `json:"models"`
}

// defaultConfig returns the default configuration
func defaultConfig() *AppConfig {
	return &AppConfig{
		DefaultModel: "claude-sonnet-4",
		Models: map[string]ModelConfig{
			// Claude models
			"claude-opus-4-5":   {CLI: "claude", ModelArg: "claude-opus-4-5-20250101"},
			"claude-opus-4":     {CLI: "claude", ModelArg: "claude-opus-4-20250101"},
			"claude-sonnet-4-5": {CLI: "claude", ModelArg: "claude-3-7-sonnet-20250219"},
			"claude-sonnet-4":   {CLI: "claude", ModelArg: "claude-3-5-sonnet-20241022"},
			"claude-sonnet-3-5": {CLI: "claude", ModelArg: "claude-3-5-sonnet-20241022"},
			"claude-sonnet-3":   {CLI: "claude", ModelArg: "claude-3-sonnet-20240229"},
			"claude-haiku-4":    {CLI: "claude", ModelArg: "claude-3-5-haiku-20241022"},
			"claude-haiku-3-5":  {CLI: "claude", ModelArg: "claude-3-5-haiku-20241022"},
			"claude-haiku-3":    {CLI: "claude", ModelArg: "claude-3-haiku-20240307"},
			"claude-opus":       {CLI: "claude", ModelArg: "claude-3-opus-20240229"},
			"claude-sonnet":     {CLI: "claude", ModelArg: "claude-3-5-sonnet-20241022"},
			"claude-haiku":      {CLI: "claude", ModelArg: "claude-3-5-haiku-20241022"},

			// Gemini models
			"gemini-pro":         {CLI: "gemini", ModelArg: "gemini-2.0-flash-exp"},
			"gemini-flash":       {CLI: "gemini", ModelArg: "gemini-2.0-flash-exp"},
			"gemini-2-flash":     {CLI: "gemini", ModelArg: "gemini-2.0-flash-exp"},
			"gemini-2-pro":       {CLI: "gemini", ModelArg: "gemini-1.5-pro"},
			"gemini-1-5-pro":     {CLI: "gemini", ModelArg: "gemini-1.5-pro"},
			"gemini-1-5-flash":   {CLI: "gemini", ModelArg: "gemini-1.5-flash"},
			"gemini-1-5-flash-8b": {CLI: "gemini", ModelArg: "gemini-1.5-flash-8b"},
		},
	}
}

// Load reads the config from file, creating default if it doesn't exist
func Load() (*AppConfig, error) {
	configPath := ConfigFile()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		return createDefaultConfig()
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure models map is initialized
	if config.Models == nil {
		config.Models = make(map[string]ModelConfig)
	}

	return &config, nil
}

// createDefaultConfig creates the config directory and default config file
func createDefaultConfig() (*AppConfig, error) {
	configDir := ConfigDir()
	configPath := ConfigFile()

	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	config := defaultConfig()

	// Write config file
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	return config, nil
}

// GetModelConfig returns the CLI and model arg for a given model ID
// Returns empty config if model not found
func GetModelConfig(modelID string) ModelConfig {
	config, err := Load()
	if err != nil {
		// Fallback to default if load fails
		defaultCfg := defaultConfig()
		if cfg, ok := defaultCfg.Models[modelID]; ok {
			return cfg
		}
		return ModelConfig{}
	}

	if cfg, ok := config.Models[modelID]; ok {
		return cfg
	}
	return ModelConfig{}
}

// GetDefaultModel returns the default model ID
func GetDefaultModel() string {
	config, err := Load()
	if err != nil {
		return defaultConfig().DefaultModel
	}
	if config.DefaultModel == "" {
		return defaultConfig().DefaultModel
	}
	return config.DefaultModel
}

// Save writes the current config to file
func Save(config *AppConfig) error {
	configPath := ConfigFile()

	// Ensure directory exists
	if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
