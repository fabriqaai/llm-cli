package cmd

import (
	"fmt"
	"os"

	"github.com/fabriqaai/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	modelFlag string
	version   string
	commit    string
	date      string
)

var rootCmd = &cobra.Command{
	Use:   "llm-cli",
	Short: "A simple CLI wrapper for Claude and Gemini LLMs",
	Long: `llm-cli is a simple wrapper around the claude and gemini CLIs.
It provides a unified interface to prompt different models without worrying
about which underlying CLI to use.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "dev" {
			fmt.Printf("llm-cli version %s\n", version)
		} else {
			fmt.Printf("llm-cli version %s (commit: %s, built: %s)\n", version, commit, date)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&modelFlag, "model", "m", "", "Model to use (e.g., claude-sonnet-4, gemini-pro)")
	rootCmd.AddCommand(versionCmd)
}

// PreRun runs before any command to set defaults
func preRun(cmd *cobra.Command, args []string) error {
	// Set default model from config if not provided
	if modelFlag == "" {
		modelFlag = config.GetDefaultModel()
	}
	return nil
}

func Execute(v, c, d string) {
	version = v
	commit = c
	date = d
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
