package cmd

import (
	"fmt"
	"sort"

	"github.com/fabriqaai/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available model aliases",
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

	fmt.Println("Available model aliases:")
	fmt.Println()

	// Group aliases by CLI
	shortcutModels := []string{}
	claudeModels := []string{}
	geminiModels := []string{}
	otherModels := []string{}

	for alias, modelCfg := range cfg.Models {
		if alias == "haiku" {
			shortcutModels = append(shortcutModels, alias)
		} else if modelCfg.CLI == "claude" {
			claudeModels = append(claudeModels, alias)
		} else if modelCfg.CLI == "gemini" {
			geminiModels = append(geminiModels, alias)
		} else {
			otherModels = append(otherModels, alias)
		}
	}

	sort.Strings(shortcutModels)
	sort.Strings(claudeModels)
	sort.Strings(geminiModels)
	sort.Strings(otherModels)

	printModelGroup("Shortcuts", shortcutModels, cfg)
	printModelGroup("Claude", claudeModels, cfg)
	printModelGroup("Gemini", geminiModels, cfg)
	if len(otherModels) > 0 {
		printModelGroup("Other", otherModels, cfg)
	}

	fmt.Printf("\nConfig file: %s\n", config.ConfigFile())
	fmt.Println("\nAdd custom aliases by editing the config file.")
}

func printModelGroup(category string, aliases []string, cfg *config.AppConfig) {
	if len(aliases) == 0 {
		return
	}
	fmt.Printf("%s:\n", category)
	for _, alias := range aliases {
		modelCfg := cfg.Models[alias]
		defaultMark := ""
		if alias == cfg.DefaultModel {
			defaultMark = " (default)"
		}
		fmt.Printf("  %s%s -> %s -m %s\n", alias, defaultMark, modelCfg.CLI, modelCfg.ModelID)
	}
	fmt.Println()
}
