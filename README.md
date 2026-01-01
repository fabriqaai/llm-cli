# llm-cli

A simple CLI wrapper for [Claude](https://github.com/anthropics/claude-cli) and [Gemini](https://github.com/google/generative-ai-python) CLIs. Provides a unified interface to prompt different models without worrying about which underlying CLI to use.

## Features

- **Unified Interface** - Single CLI for both Claude and Gemini models
- **Model Aliases** - Easy-to-remember model names that map to underlying model IDs
- **Config File** - Add custom models and change defaults via `~/.llm-cli/models.config`
- **Streaming Output** - Real-time response streaming
- **Piped Input** - Accept input via stdin for scripting

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap fabriqaai/tap
brew install llm-cli
```

### From Source

```bash
go install github.com/fabriqaai/llm-cli@latest
```

Or build manually:

```bash
git clone https://github.com/fabriqaai/llm-cli.git
cd llm-cli
go build -o llm-cli .
mv llm-cli /usr/local/bin/
```

## First Run

On first run, llm-cli creates a default config file at `~/.llm-cli/models.config` with common Claude and Gemini models.

## Usage

```bash
# Check version
llm-cli version

# Send a prompt (uses default model: haiku)
llm-cli "what is the capital of france?"

# Use a specific model alias
llm-cli opus "explain go interfaces"
llm-cli gemini "what is 2+2?"

# List all available models
llm-cli models

# Use flags
llm-cli -m haiku "hello"
llm-cli -m sonnet -s "You are a Go expert" "how do I use interfaces?"
```

## Options

| Flag | Short | Description |
|------|-------|-------------|
| `--model` | `-m` | Model to use (e.g., haiku, opus, sonnet, gemini) |
| `--prompt` | `-p` | Prompt text |
| `--system` | `-s` | System prompt for context |

## Configuration

Edit `~/.llm-cli/models.config` to add custom models or change the default:

```json
{
  "default_model": "haiku",
  "models": {
    "my-custom-model": {
      "cli": "claude",
      "model_id": "some-custom-model-id"
    }
  }
}
```

## Available Models

Run `llm-cli models` to see all available models.

### Shortcuts

| Alias | Model |
|-------|-------|
| `haiku` | claude-haiku-4-5-20251001 |
| `opus` | claude-opus-4-5-20251101 |
| `sonnet` | claude-sonnet-4-5-20251001 |
| `claude` | claude-sonnet-4-5-20251001 |
| `gemini` | gemini-3-pro-preview |
| `flash` | gemini-3-flash-preview |

### Claude 4.5 Models

| Alias | Model |
|-------|-------|
| `haiku45` | claude-haiku-4-5-20251001 |
| `opus45` | claude-opus-4-5-20251101 |
| `sonnet45` | claude-sonnet-4-5-20251001 |

### Claude 3.5 / 3.7 Models

| Alias | Model |
|-------|-------|
| `haiku35` | claude-3-5-haiku-20241022 |
| `opus41` | claude-opus-4-20250101 |
| `sonnet35` | claude-3-5-sonnet-20241022 |
| `sonnet37` | claude-3-7-sonnet-20250219 |

### Gemini 3.0 Models

| Alias | Model |
|-------|-------|
| `gemini30` | gemini-3-pro-preview |
| `gemini30pro` | gemini-3-pro-preview |
| `gemini30flash` | gemini-3-flash-preview |
| `flash30` | gemini-3-flash-preview |

### Gemini 2.5 Models

| Alias | Model |
|-------|-------|
| `pro25` | gemini-2.5-pro-exp |
| `flash25` | gemini-2.5-flash |

## Requirements

- **Claude CLI** - Install via `npm install -g @anthropic-ai/claude-cli` (for Claude models)
- **Gemini CLI** - Install via `npm install -g @google/generative-ai-cli` (for Gemini models)

## License

MIT
