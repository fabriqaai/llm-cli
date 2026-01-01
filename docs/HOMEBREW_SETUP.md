# Homebrew Tap Setup

This document describes how to set up automatic Homebrew formula publishing for llm-cli.

## GitHub Actions Workflow

The `.github/workflows/homebrew-release.yml` workflow runs on releases and publishes the Homebrew formula to the `fabriqaai/homebrew-tap` repository.

## Required Setup

### 1. Create Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Create a new token with `repo` scope (write access to repositories)
3. Name it `Homebrew Tap Token`

### 2. Add Secret to llm-cli Repository

1. Go to `github.com/fabriqaai/llm-cli/settings/secrets/actions`
2. Add a new secret named `HOMEBREW_TAP_TOKEN`
3. Paste the personal access token value

### 3. Update homebrew-tap README

The homebrew-tap README needs to include llm-cli. Here's the updated content:

```markdown
# Homebrew Tap for FabriqaAI Tools

This repository contains Homebrew formulae for FabriqaAI tools.

## Installation

```bash
# Add the tap
brew tap fabriqaai/tap

# Install a tool
brew install llm-cli
brew install claude-code-logs
```

## Available Formulae

### llm-cli

A simple CLI wrapper for Claude and Gemini CLIs. Provides a unified interface to prompt different models without worrying about which underlying CLI to use.

```bash
brew install llm-cli
```

**Usage:**

```bash
# Simple prompt with default model (haiku)
llm-cli "what is the capital of france?"

# Use a specific model alias
llm-cli opus "explain go interfaces"
llm-cli sonnet "write a python function"
llm-cli gemini "what is 2+2?"
llm-cli flash "translate hello to spanish"

# Using flags
llm-cli -m opus -s "You are a Go expert" "how do I use interfaces?"

# List all available models
llm-cli models
```

**Repository:** [github.com/fabriqaai/llm-cli](https://github.com/fabriqaai/llm-cli)

### claude-code-logs

Browse and search Claude Code chat logs.

```bash
brew install claude-code-logs
```

**Usage:**

```bash
# Generate HTML from chat logs
claude-logs generate

# Start the local server
claude-logs serve

# Watch for changes and auto-regenerate
claude-logs watch
```

**Repository:** [github.com/fabriqaai/claude-code-logs](https://github.com/fabriqaai/claude-code-logs)

## Updating

```bash
brew update
brew upgrade llm-cli
brew upgrade claude-code-logs
```

## Troubleshooting

### Formula not found

```bash
brew tap fabriqaai/tap
brew update
```

### SHA256 mismatch

```bash
brew update
brew reinstall llm-cli
brew reinstall claude-code-logs
```

## License

MIT
```

To apply this update:

```bash
cd /path/to/homebrew-tap
# Update README.md with the content above
git add README.md
git commit -m "Add llm-cli to README"
git push origin main
```

## How It Works

1. When a new release is created in `llm-cli`, the `homebrew-release.yml` workflow triggers
2. GoReleaser builds the binaries and creates a Homebrew formula
3. The formula is pushed to `fabriqaai/homebrew-tap` repository using the `HOMEBREW_TAP_TOKEN`
4. Users can then `brew upgrade llm-cli` to get the latest version

## Testing

To test the Homebrew formula locally before release:

```bash
# Build with GoReleaser (skip publish)
goreleaser build --snapshot --clean

# Or test install from local build
brew install ./Formula/llm-cli.rb --HEAD
```
