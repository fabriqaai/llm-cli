package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// OptionsFile returns the options config file path
func OptionsFile() string {
	return filepath.Join(ConfigDir(), "options.json")
}

// Options represents the application options
type Options struct {
	// RunOnCurrentDirectory determines whether to run CLI in current directory
	// or in the sessions directory. Defaults to false (sessions).
	RunOnCurrentDirectory bool `json:"run_on_current_directory"`
}

// defaultOptions returns the default options
func defaultOptions() *Options {
	return &Options{
		RunOnCurrentDirectory: false,
	}
}

// LoadOptions reads the options from file, creating default if it doesn't exist
func LoadOptions() (*Options, error) {
	optionsPath := OptionsFile()

	// Check if options file exists
	if _, err := os.Stat(optionsPath); os.IsNotExist(err) {
		// Create default options
		return createDefaultOptions()
	}

	// Read existing options
	data, err := os.ReadFile(optionsPath)
	if err != nil {
		return nil, err
	}

	var options Options
	if err := json.Unmarshal(data, &options); err != nil {
		return nil, err
	}

	return &options, nil
}

// createDefaultOptions creates the default options file
func createDefaultOptions() (*Options, error) {
	configDir := ConfigDir()
	optionsPath := OptionsFile()

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	options := defaultOptions()

	// Write options file
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(optionsPath, data, 0644); err != nil {
		return nil, err
	}

	return options, nil
}

// SaveOptions writes the options to file
func SaveOptions(options *Options) error {
	optionsPath := OptionsFile()

	// Ensure directory exists
	if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(optionsPath, data, 0644)
}

// SessionsDir returns the sessions directory path
func SessionsDir() string {
	return filepath.Join(ConfigDir(), "sessions")
}

// EnsureSessionsDir creates the sessions directory if it doesn't exist
func EnsureSessionsDir() error {
	sessionsDir := SessionsDir()
	return os.MkdirAll(sessionsDir, 0755)
}

// GetWorkingDirectory returns the working directory for the CLI session
// runOnTempDir: when true, force use of temp/sessions directory (overrides config)
func GetWorkingDirectory(runOnTempDir bool) (string, error) {
	// If temp dir flag is set, use sessions directory
	if runOnTempDir {
		sessionsDir := SessionsDir()
		if err := EnsureSessionsDir(); err != nil {
			return "", err
		}
		return sessionsDir, nil
	}

	// Load options
	options, err := LoadOptions()
	if err != nil {
		// If options can't be loaded, default to sessions directory
		sessionsDir := SessionsDir()
		if err := EnsureSessionsDir(); err != nil {
			return "", err
		}
		return sessionsDir, nil
	}

	// If run_on_current_directory is true, use current directory
	if options.RunOnCurrentDirectory {
		return os.Getwd()
	}

	// Otherwise (default), use sessions directory
	sessionsDir := SessionsDir()
	if err := EnsureSessionsDir(); err != nil {
		return "", err
	}

	return sessionsDir, nil
}
