# llm-cli

A simple CLI wrapper for [Claude](https://github.com/anthropics/claude-cli) and [Gemini](https://github.com/google/generative-ai-python) CLIs. Provides a unified interface to prompt different models without worrying about which underlying CLI to use.

## Installation

```bash
# Clone the repo
git clone https://github.com/fabriqaai/llm-cli.git
cd llm-cli

# Build
go build -o llm-cli .

# Or use the build script
chmod +x build.sh
./build.sh

# Install globally (optional)
mv llm-cli /usr/local/bin/
```

## First Run

On first run, llm-cli creates a default config file at `~/.llm-cli/models.config` with common Claude and Gemini models. You can edit this file to add your own models or change the default.

## Usage

```bash
# Send a prompt (uses default model: claude-sonnet-4)
llm-cli prompt "explain go interfaces"

# Use a specific model
llm-cli prompt -m claude-opus-4-5 "explain go interfaces"
llm-cli prompt -m gemini-pro "explain go interfaces"

# List all available models
llm-cli models

# Pipe input
echo "what is 2+2?" | llm-cli prompt
cat file.txt | llm-cli prompt -m gemini-flash

# Set a system prompt (for context)
llm-cli prompt -s "You are a Go expert" "how do I use interfaces?"
```

## Configuration

Edit `~/.llm-cli/models.config` to add custom models:

```json
{
  "default_model": "claude-sonnet-4",
  "models": {
    "claude-opus-4-5": {
      "cli": "claude",
      "model_arg": "claude-opus-4-5-20250101"
    },
    "my-custom-model": {
      "cli": "claude",
      "model_arg": "some-custom-model-id"
    }
  }
}
```

## Available Models

Run `llm-cli models` to see all available models. By default includes:

### Claude Models
- `claude-opus-4-5` - Latest Opus
- `claude-sonnet-4-5` - Latest Sonnet
- `claude-haiku-4` - Latest Haiku
- And more...

### Gemini Models
- `gemini-pro` - Gemini 2.0 Flash
- `gemini-2-pro` - Gemini 1.5 Pro
- And more...

## Requirements

- Go 1.23+
- [claude CLI](https://github.com/anthropics/claude-cli) installed (for Claude models)
- [gemini CLI](https://github.com/google/generative-ai-python) installed (for Gemini models)
