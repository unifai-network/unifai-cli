---
name: ucli
description: Unifai CLI for searching and invoking blockchain services across multiple chains including Solana, Base, and Ethereum.
version: 1.0.0
metadata:
  openclaw:
    requires:
      env:
        - UNIFAI_AGENT_API_KEY
      bins:
        - ucli
    primaryEnv: UNIFAI_AGENT_API_KEY
    emoji: "🦄"
    homepage: https://github.com/unifai-network/unifai-cli
    install:
      - kind: brew
        tap: unifai-network/homebrew-unifai-cli
        formula: ucli
        bins: [ucli]
---

# ucli

A Go CLI for Unifai actions with first-class support for searching and invoking blockchain services across multiple chains.

## What it does

ucli enables you to:

- **Search services**: Find blockchain services and actions using natural language queries
- **Invoke services**: Execute blockchain actions with customizable parameters and retry logic
- **Manage configuration**: Configure API keys and endpoints with multiple priority levels

Default API endpoint: `https://app.uniclaw.ai/api/v1/unifai`

## Installation

### Homebrew

```bash
brew update
brew tap unifai-network/homebrew-unifai-cli
brew install ucli
ucli version
```

## Authentication

API key source priority (highest to lowest):

1. Command-line flag: `--api-key`
2. Environment variable: `UNIFAI_AGENT_API_KEY`
3. Config file: `~/.config/unifai-cli/config.yaml`

### Initialize Configuration

Generate a config template:

```bash
ucli config init
```

Show effective configuration:

```bash
ucli config show
```

## Usage Examples

### Search for Services

Search for blockchain services using natural language:

```bash
# Basic search
ucli search --query "swap usdc to sol"

# Search with pagination
ucli search --query "transfer tokens" --limit 10 --offset 0

# Include specific actions
ucli search --query "defi protocols" --include-actions action1,action2
```

### Invoke Services

Execute blockchain actions with JSON payloads:

```bash
# Inline JSON payload
ucli invoke --action "Meteora--29--swap" --payload '{"amount":100}'

# Read payload from file
ucli invoke --action "MyAction--1--execute" --payload @payload.json

# With custom retries and timeout
ucli invoke --action "MyAction--1--execute" --payload '{"x":1}' --max-retries 3 --timeout 60s
```

### Payload Formats

ucli supports flexible payload handling:

- **auto** (default): Parses valid JSON as object, otherwise treats as string
- **object**: Forces JSON object parsing
- **string**: Sends payload as raw string

```bash
# Force object parsing
ucli invoke --action "MyAction" --payload '{"key":"value"}' --payload-format object

# Force string parsing
ucli invoke --action "MyAction" --payload "raw text" --payload-format string
```

### Output Formats

```bash
# Human-readable output (default)
ucli search --query "swap tokens"

# JSON output for scripting
ucli search --query "swap tokens" --json
```

## Configuration

### Config File Location

`~/.config/unifai-cli/config.yaml`

### Example Configuration

```yaml
api_key: your-unifai-api-key
endpoint: https://app.uniclaw.ai/api/v1/unifai
```

## Retry and Timeout Behavior

- **Default max retries**: 1
- **Retry strategy**: Exponential backoff (1s, 2s, 4s, ...)
- **Default timeout**: 50s
- **Retry conditions**: Network failures and HTTP 5xx errors

## Exit Codes

- **0**: Success
- **1**: API or runtime error
- **2**: Argument or usage error

## Common Use Cases

### Swap Tokens on Solana

```bash
ucli search --query "swap usdc to sol on solana"
ucli invoke --action "Meteora--29--swap" --payload '{"fromToken":"USDC","toToken":"SOL","amount":100}'
```

### Bridge Assets

```bash
ucli search --query "bridge eth to base"
ucli invoke --action "Bridge--1--transfer" --payload '{"chain":"base","amount":0.1}'
```

### Check Service Status

```bash
ucli search --query "protocol health check" --json
```

## When to Use This Skill

Use ucli when you need to:

- Search for blockchain services and actions across multiple chains
- Execute on-chain transactions programmatically
- Integrate Unifai capabilities into scripts and automation workflows
- Quickly test blockchain service invocations with retry logic
- Query available DeFi protocols and their capabilities

## Advanced Features

### Custom Timeouts

```bash
ucli invoke --action "LongRunning--1--process" --payload '{}' --timeout 120s
```

### Retry Configuration

```bash
ucli invoke --action "Unreliable--1--call" --payload '{}' --max-retries 5
```

### API Key Override

```bash
ucli search --query "test" --api-key temporary-key-123
```

## Troubleshooting

### Check Configuration

```bash
ucli config show
```

This displays the effective configuration and shows which source (flag, env, or file) is being used.

### Verify Installation

```bash
ucli version
```

### Test API Connectivity

```bash
ucli search --query "test" --json
```

## Additional Resources

- Config template: `configs/config.example.yaml`
- Homebrew tap: https://github.com/unifai-network/homebrew-unifai-cli
