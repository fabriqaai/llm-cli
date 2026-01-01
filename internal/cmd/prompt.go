package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fabriqaai/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt [prompt-text]",
	Short: "Send a prompt to the LLM",
	Long: `Send a prompt to the specified LLM model.
The prompt can be provided as an argument, via stdin, or both.`,
	Args: cobra.ArbitraryArgs,
	PreRunE: preRun,
	RunE: runPrompt,
}

var (
	systemFlag  string
	continueFlag bool
)

func init() {
	rootCmd.AddCommand(promptCmd)

	promptCmd.Flags().StringVarP(&systemFlag, "system", "s", "", "System prompt to set context")
	promptCmd.Flags().BoolVarP(&continueFlag, "continue", "c", false, "Continue the previous conversation")
}

func runPrompt(cmd *cobra.Command, args []string) error {
	// Get model configuration
	modelConfig := config.GetModelConfig(modelFlag)
	if modelConfig.CLI == "" {
		return fmt.Errorf("unknown model: %s", modelFlag)
	}

	// Build prompt from args and/or stdin
	prompt, err := buildPrompt(args)
	if err != nil {
		return err
	}

	// Execute the appropriate CLI
	return executeLLM(modelConfig, prompt)
}

// buildPrompt constructs the prompt from args and stdin
func buildPrompt(args []string) (string, error) {
	var parts []string

	// Add argument prompt
	if len(args) > 0 {
		parts = append(parts, strings.Join(args, " "))
	}

	// Add stdin if available (piped input)
	if isPiped() {
		stdinContent, err := readStdin()
		if err != nil {
			return "", err
		}
		if stdinContent != "" {
			parts = append(parts, stdinContent)
		}
	}

	// If no prompt provided and not piped, return empty (will wait for input?)
	if len(parts) == 0 {
		return "", nil
	}

	return strings.Join(parts, "\n\n"), nil
}

// isPiped checks if input is being piped in
func isPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// readStdin reads all content from stdin
func readStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var buffer bytes.Buffer

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		buffer.WriteString(line)
	}

	return buffer.String(), nil
}

// executeLLM runs the appropriate CLI with the prompt
func executeLLM(modelConfig config.ModelConfig, prompt string) error {
	if prompt == "" {
		return fmt.Errorf("no prompt provided")
	}

	// Build the command
	// claude CLI: claude -m MODEL "prompt"
	// gemini CLI: gemini -m MODEL "prompt"
	command := exec.Command(modelConfig.CLI, "-m", modelConfig.ModelArg, prompt)

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
		return fmt.Errorf("failed to start %s: %w (is %s installed?)", modelConfig.CLI, err, modelConfig.CLI)
	}

	// Stream stdout to our stdout in real-time
	go streamToStdout(stdout, os.Stdout)
	// Stream stderr to our stderr
	go streamToStdout(stderr, os.Stderr)

	// Wait for command to finish
	if err := command.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// streamToStdout copies from reader to writer in real-time
func streamToStdout(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Fprintln(writer, scanner.Text())
	}
}
