# unifai-cli

A Go CLI for Unifai actions with first-class support for:

- `search_services` via `unifai search`
- `invoke_service` via `unifai invoke`

Default API endpoint:

- `https://app.uniclaw.ai/api/v1/unifai`

## Commands

- `unifai search --query "swap usdc to sol" --limit 10 --offset 0 --include-actions a,b`
- `unifai invoke --action "Meteora--29--..." --payload '{"x":1}' --max-retries 1`
- `unifai config init`
- `unifai config show`
- `unifai version`

## Auth and Config Priority

API key source priority (high to low):

1. `--api-key`
2. `UNIFAI_AGENT_API_KEY`
3. `~/.config/unifai-cli/config.yaml`

Endpoint source priority (high to low):

1. `--endpoint`
2. `UNIFAI_ENDPOINT`
3. `~/.config/unifai-cli/config.yaml`
4. Default: `https://app.uniclaw.ai/api/v1/unifai`

Generate config template:

```bash
unifai config init
```

Show effective config and source:

```bash
unifai config show
```

Reference template: `configs/config.example.yaml`

## Payload Compatibility

`invoke` supports both payload object and payload string to match current TS behavior differences.

- Inline JSON: `--payload '{"a":1}'`
- Read from file: `--payload @payload.json`
- Force parsing mode: `--payload-format auto|object|string`

Behavior:

- `auto` (default): valid JSON -> object; otherwise -> raw string
- `object`: payload must be valid JSON
- `string`: payload is sent as string

## Retries and Timeout

- `--max-retries` (default: `1`)
- Exponential backoff: `1s`, `2s`, `4s`, ...
- `--timeout` (default: `50s`)
- Retry only on network failures and HTTP `5xx`

## Output and Exit Codes

- Default output: concise human-readable
- `--json`: raw API response
- Exit codes:
  - `0`: success
  - `1`: API/runtime error
  - `2`: argument/usage error

`invoke` default rendering normalizes result:

- If response has `payload`, print `payload`
- Otherwise print full response

## Build and Test

```bash
make tidy
make test
make build
```

Binary output:

- `bin/unifai`

## Install

### Homebrew

```bash
brew tap unifai-network/homebrew-unifai-cli
brew install unifai
```

## Release

### Local

```bash
# Snapshot package (no GitHub release)
./scripts/release.sh --snapshot

# Real release
./scripts/release.sh
```

### GitHub Actions

Workflow file: `.github/workflows/release.yml`

Trigger by pushing a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

GoReleaser config: `.goreleaser.yaml`

Release artifacts are published to:

- `https://github.com/unifai-network/homebrew-unifai-cli/releases`

This is required when `unifai-cli` source repository is private, so Homebrew can download assets anonymously.

### Homebrew release setup (maintainer)

1. Create tap repo: `unifai-network/homebrew-unifai-cli` (with `main` branch).
2. Add repository secret `HOMEBREW_TAP_GITHUB_TOKEN` in `unifai-cli` GitHub Actions.
3. Ensure the token can write to `unifai-network/homebrew-unifai-cli`.
4. Push a tag (for example `v0.1.0`) to trigger release.
