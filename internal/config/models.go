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

// ModelConfig maps a model alias to CLI and model ID
type ModelConfig struct {
	// CLI is the command to run (e.g., "claude", "gemini", "llm")
	CLI string `json:"cli"`
	// ModelID is the model ID to pass to the CLI
	ModelID string `json:"model_id"`
}

// AppConfig represents the application configuration
type AppConfig struct {
	DefaultModel string                 `json:"default_model"`
	Models       map[string]ModelConfig `json:"models"`
}

// defaultConfig returns the default configuration
// Uses shortcuts for claude and gemini CLIs
func defaultConfig() *AppConfig {
	return &AppConfig{
		DefaultModel: "haiku",
		Models: map[string]ModelConfig{
			// Shortcut aliases
			"haiku":   {CLI: "claude", ModelID: "claude-haiku-4-5-20251001"},
			"haiku45": {CLI: "claude", ModelID: "claude-haiku-4-5-20251001"},
			"haiku35": {CLI: "claude", ModelID: "claude-3-5-haiku-20241022"},

			// Claude models
			"opus":    {CLI: "claude", ModelID: "claude-opus-4-5-20251101"},
			"opus45":  {CLI: "claude", ModelID: "claude-opus-4-5-20251101"},
			"opus41":  {CLI: "claude", ModelID: "claude-opus-4-20250101"},

			"sonnet":   {CLI: "claude", ModelID: "claude-sonnet-4-5-20251001"},
			"sonnet45": {CLI: "claude", ModelID: "claude-sonnet-4-5-20251001"},
			"sonnet35": {CLI: "claude", ModelID: "claude-3-5-sonnet-20241022"},
			"sonnet37": {CLI: "claude", ModelID: "claude-3-7-sonnet-20250219"},
			"claude":   {CLI: "claude", ModelID: "claude-sonnet-4-5-20251001"},

			// Gemini models
			"gemini":        {CLI: "gemini", ModelID: "gemini-3-pro-preview"},
			"gemini30":      {CLI: "gemini", ModelID: "gemini-3-pro-preview"},
			"gemini30pro":   {CLI: "gemini", ModelID: "gemini-3-pro-preview"},
			"gemini30flash": {CLI: "gemini", ModelID: "gemini-3-flash-preview"},
			"flash":         {CLI: "gemini", ModelID: "gemini-3-flash-preview"},
			"flash30":       {CLI: "gemini", ModelID: "gemini-3-flash-preview"},
			"pro25":         {CLI: "gemini", ModelID: "gemini-2.5-pro-exp"},
			"flash25":         {CLI: "gemini", ModelID: "gemini-2.5-flash"},

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

// GetModelConfig returns the CLI and model ID for a given alias
func GetModelConfig(alias string) ModelConfig {
	config, err := Load()
	if err != nil {
		// Return the alias as-is if config fails
		return ModelConfig{CLI: "llm", ModelID: alias}
	}

	if cfg, ok := config.Models[alias]; ok {
		return cfg
	}
	// Pass through unknown aliases to llm CLI
	return ModelConfig{CLI: "llm", ModelID: alias}
}

// GetDefaultModel returns the default model alias
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
