# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`unifai-cli` (binary: `unifai`) is a Go CLI tool for interacting with the Unifai blockchain services API. It provides two primary operations: searching for services and invoking blockchain actions across multiple chains (Solana, Base, Ethereum).

## Development Commands

### Build and Test
```bash
# Tidy dependencies
make tidy

# Format code
make fmt

# Run all tests
make test

# Build binary (outputs to bin/ucli)
make build

# Run directly without building
make run
```

### Release
```bash
# Snapshot release (local, no GitHub)
make snapshot-release

# Production release (requires goreleaser)
make release

# GitHub release via tag
git tag v0.1.0
git push origin v0.1.0
```

## Architecture

### Entry Point and Command Flow
- **Entry**: `cmd/unifai/main.go` → `internal/app/run.go`
- **Command Router**: `internal/command/root.go` sets up cobra commands with global options
- **Exit Code Handling**: `internal/app/run.go` distinguishes between usage errors (exit 2) and runtime errors (exit 1)

### Core Packages

#### `internal/command/`
Command implementations using cobra framework:
- `root.go`: Main command router with `GlobalOptions` struct (ConfigPath, APIKey, Endpoint, Timeout)
- `search.go`: Implements `unifai search` with query, limit, offset, include-actions
- `invoke.go`: Implements `unifai invoke` with action, payload, max-retries, payload-format
- `config.go`: Implements `unifai config init` and `unifai config show`
- `version.go`: Implements `unifai version`
- `common.go`: Shared command utilities

#### `internal/config/`
Configuration loading with priority hierarchy (flag > env > file):
- `loader.go`: `Resolve()` function implements 3-tier priority system
- `config.go`: Type definitions for `FileConfig` and `EffectiveConfig`
- API Key Priority: `--api-key` flag > `UNIFAI_AGENT_API_KEY` env > `~/.config/unifai-cli/config.yaml`
- Endpoint Priority: `--endpoint` flag > `UNIFAI_ENDPOINT` env > `~/.config/unifai-cli/config.yaml` > default (set in `internal/config/config.go`)
- Default timeout: 50s

#### `internal/unifai/`
API client implementation:
- `client.go`: `Client` struct with `Search()` and `Invoke()` methods
- `types.go`: Request/response types (`SearchRequest`, `InvokeRequest`, `APIError`)
- Uses standard `net/http` with context support
- Authorization via header: `Authorization: {apiKey}`

#### `internal/retry/`
Exponential backoff retry logic:
- Retries on network failures and HTTP 5xx errors
- Default max retries: 1
- Backoff sequence: 1s, 2s, 4s, 8s, ...

#### `internal/output/`
Output formatting:
- `print.go`: Handles both human-readable and `--json` output modes
- `invoke` normalizes output by extracting `payload` field if present

#### `internal/errors/`
Custom error types:
- `UsageError`: Maps to exit code 2
- Exit codes defined: `ExitOK (0)`, `ExitError (1)`, `ExitUsage (2)`

#### `internal/version/`
Version info injected at build time via ldflags:
- `Version`, `Commit`, `BuildDate` set by Makefile/GoReleaser

### Key Design Patterns

1. **Config Resolution**: Three-tier priority (flag > env > file) implemented in `config.Resolve()`
2. **Error Classification**: Usage vs runtime errors determine exit codes
3. **Payload Flexibility**: `invoke` supports `--payload '{"x":1}'`, `--payload @file.json`, and `--payload-format` (auto/object/string)
4. **Retry Strategy**: Exponential backoff with configurable max retries
5. **Output Normalization**: `invoke` extracts `payload` field from response if present for cleaner output

## Testing

Run tests with:
```bash
make test
# or directly:
go test ./...
```

No test files currently present in codebase.

## Release Process

### GoReleaser Configuration
- Config: `.goreleaser.yaml`
- Builds for: linux/darwin/windows on amd64/arm64
- Archives: tar.gz (zip for Windows)
- Publishes to: `unifai-network/homebrew-unifai-cli` repo (not source repo)
- Requires: `HOMEBREW_TAP_GITHUB_TOKEN` secret for GitHub Actions

### Why Separate Homebrew Repo?
Release artifacts are published to `homebrew-unifai-cli` repo to enable anonymous downloads when the source repo is private.

### Homebrew Installation
```bash
brew tap unifai-network/homebrew-unifai-cli
brew install ucli
```

## Configuration

### Config File
- Location: `~/.config/unifai-cli/config.yaml`
- Template: `configs/config.example.yaml`
- Fields: `apiKey`, `endpoint`, `timeoutSeconds`

### Generate Config
```bash
ucli config init
```

### View Effective Config
```bash
ucli config show
# Shows resolved values and their sources (flag/env/file/default)
```

## Module and Dependencies

- Module: `unifai`
- Go version: 1.20
- Primary dependencies:
  - `github.com/spf13/cobra` - CLI framework
  - `gopkg.in/yaml.v3` - YAML config parsing

## Build Details

### Version Injection
The Makefile injects version info via ldflags:
```
-X unifai/internal/version.Version=$(VERSION)
-X unifai/internal/version.Commit=$(COMMIT)
-X unifai/internal/version.BuildDate=$(BUILD_DATE)
```

### Binary Output
- Local builds: `bin/unifai`
- GoReleaser builds: Multi-platform archives in `dist/`

## Common Development Patterns

When adding new commands:
1. Create command constructor in `internal/command/`
2. Accept `*GlobalOptions` parameter for shared flags
3. Register in `root.go` via `root.AddCommand()`
4. Return appropriate error types (`UsageError` for bad input)

When modifying API client:
1. Update types in `internal/unifai/types.go`
2. Add/modify methods in `internal/unifai/client.go`
3. Use `doJSON()` helper for consistent HTTP handling
