package cmd

import (
	"os"

	"github.com/fabriqaai/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "llm-cli",
	Short: "A simple CLI wrapper for Claude and Gemini LLMs",
	Long: `llm-cli is a simple wrapper around the claude and gemini CLIs.
It provides a unified interface to prompt different models without worrying
about which underlying CLI to use.`,
}

var (
	modelFlag string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&modelFlag, "model", "m", "", "Model to use (e.g., claude-sonnet-4, gemini-pro)")
}

// PreRun runs before any command to set defaults
func preRun(cmd *cobra.Command, args []string) error {
	// Set default model from config if not provided
	if modelFlag == "" {
		modelFlag = config.GetDefaultModel()
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
