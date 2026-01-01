package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/fabriqaai/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	modelFlag  string
	promptFlag string
	systemFlag string
	version    string
	commit     string
	date       string
)

var rootCmd = &cobra.Command{
	Use:   "llm-cli [model-alias] [prompt]",
	Short: "A simple CLI wrapper for claude and gemini CLIs",
	Long: `llm-cli is a simple wrapper around the claude and gemini CLIs.
It provides shortcuts for common models and aliases.

Examples:
  llm-cli haiku "what is the capital of france?"
  llm-cli opus "explain go interfaces"
  llm-cli gemini "what is 2+2?"
  llm-cli models`,
	Args:              cobra.ArbitraryArgs,
	DisableFlagParsing: false,
	RunE:              runRoot,
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
	rootCmd.Flags().StringVarP(&modelFlag, "model", "m", "", "Model to use (e.g., haiku, opus, sonnet, gemini)")
	rootCmd.Flags().StringVarP(&promptFlag, "prompt", "p", "", "Prompt text")
	rootCmd.Flags().StringVarP(&systemFlag, "system", "s", "", "System prompt for context")
	rootCmd.AddCommand(versionCmd)
}

// runRoot handles the root command with positional args
func runRoot(cmd *cobra.Command, args []string) error {
	modelAlias := modelFlag
	promptText := promptFlag

	if len(args) == 0 && promptText == "" {
		// No args at all, show help
		return cmd.Help()
	}

	// Parse positional args based on count
	if modelAlias == "" && promptText == "" && len(args) == 1 {
		// Just one arg: it's the prompt, use default model
		promptText = args[0]
		modelAlias = config.GetDefaultModel()
	} else if modelAlias == "" && len(args) >= 1 {
		// Multiple args without model flag
		// Check if first arg is a known model alias
		cfg, _ := config.Load()
		if _, isModel := cfg.Models[args[0]]; isModel && len(args) >= 2 {
			// First arg is a model, second is prompt
			modelAlias = args[0]
			promptText = args[1]
		} else {
			// First arg is the prompt, use default model
			promptText = args[0]
			modelAlias = config.GetDefaultModel()
		}
	} else if promptText == "" && len(args) > 0 {
		// Model provided via flag, first arg is prompt
		promptText = args[0]
	}

	// If no model specified, use default
	if modelAlias == "" {
		modelAlias = config.GetDefaultModel()
	}

	// If no prompt, show help
	if promptText == "" {
		return cmd.Help()
	}

	// Get model configuration
	modelConfig := config.GetModelConfig(modelAlias)

	// Execute the appropriate CLI
	return executeLLM(modelConfig.CLI, modelConfig.ModelID, promptText, systemFlag)
}

// executeLLM runs the appropriate CLI with the prompt
func executeLLM(cli, modelID, prompt, systemPrompt string) error {
	if prompt == "" {
		return fmt.Errorf("no prompt provided")
	}

	// Build the command based on CLI type
	var args []string

	switch cli {
	case "claude":
		// claude CLI: claude --model MODEL -p "prompt"
		args = []string{"--model", modelID, "-p", prompt}
		if systemPrompt != "" {
			args = []string{"--model", modelID, "--system-prompt", systemPrompt, "-p", prompt}
		}
	case "gemini":
		// gemini CLI: gemini -m MODEL "prompt"
		args = []string{"-m", modelID, prompt}
		if systemPrompt != "" {
			args = []string{"-m", modelID, "--system", systemPrompt, prompt}
		}
	default:
		// llm CLI: llm prompt -m MODEL "prompt"
		args = []string{"prompt", "-m", modelID}
		if systemPrompt != "" {
			args = append(args, "--system", systemPrompt)
		}
		args = append(args, prompt)
	}

	command := exec.Command(cli, args...)

	// Set up pipes for stdout and stderr
	stdout, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := command.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w (is %s installed?)", cli, err, cli)
	}

	// Start a spinner in the background
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		spinner(done, cli)
	}()

	// Stream stdout to our stdout in real-time
	outputDone := &sync.WaitGroup{}
	outputDone.Add(2)
	go func() {
		defer outputDone.Done()
		streamToStdout(stdout, os.Stdout, done)
	}()
	// Stream stderr to our stderr
	go func() {
		defer outputDone.Done()
		streamToStdout(stderr, os.Stderr, done)
	}()

	// Wait for output to finish
	outputDone.Wait()

	// Stop spinner and wait for it to finish
	select {
	case <-done:
	default:
		close(done)
	}
	wg.Wait()

	// Wait for command to finish
	if err := command.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// spinner shows a running indicator until done channel is closed
func spinner(done chan struct{}, cli string) {
	base := "running " + cli
	frames := []string{base + ".", base + "..", base + "..."}
	i := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// Clear the spinner line when done
			fmt.Print("\r\033[K")
			return
		case <-ticker.C:
			fmt.Printf("\r\033[K%s", frames[i])
			i = (i + 1) % len(frames)
		}
	}
}

// streamToStdout copies from reader to writer in real-time
func streamToStdout(reader io.Reader, writer io.Writer, done chan struct{}) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// Stop spinner on first output
		select {
		case <-done:
		default:
			close(done)
			fmt.Print("\r\033[K") // Clear spinner line
		}
		fmt.Fprintln(writer, scanner.Text())
	}
}

func Execute(v, c, d string) {
	version = v
	commit = c
	date = d
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
