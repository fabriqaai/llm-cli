package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fabriqaai/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available models",
	PreRunE: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		listModels()
	},
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}

func listModels() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	fmt.Println("Available models:")
	fmt.Println()

	// Group models by CLI
	claudeModels := []string{}
	geminiModels := []string{}

	for modelID := range cfg.Models {
		if strings.HasPrefix(modelID, "claude-") {
			claudeModels = append(claudeModels, modelID)
		} else if strings.HasPrefix(modelID, "gemini-") {
			geminiModels = append(geminiModels, modelID)
		}
	}

	sort.Strings(claudeModels)
	sort.Strings(geminiModels)

	fmt.Println("Claude models:")
	for _, m := range claudeModels {
		modelCfg := cfg.Models[m]
		defaultMark := ""
		if m == cfg.DefaultModel {
			defaultMark = " (default)"
		}
		fmt.Printf("  %s%s -> %s -m %s\n", m, defaultMark, modelCfg.CLI, modelCfg.ModelArg)
	}

	fmt.Println()
	fmt.Println("Gemini models:")
	for _, m := range geminiModels {
		modelCfg := cfg.Models[m]
		defaultMark := ""
		if m == cfg.DefaultModel {
			defaultMark = " (default)"
		}
		fmt.Printf("  %s%s -> %s -m %s\n", m, defaultMark, modelCfg.CLI, modelCfg.ModelArg)
	}

	fmt.Println()
	fmt.Printf("Config file: %s\n", config.ConfigFile())
}
