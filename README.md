# unifai-cli

A Go CLI for searching and invoking services on the [Unifai](https://unifai.network) network — DeFi, token data, social media, web search, news, travel, sports, and more.

## Installation

### Homebrew

```bash
brew tap unifai-network/homebrew-unifai-cli
brew install unifai
```

### From source

```bash
git clone https://github.com/unifai-network/unifai-cli.git
cd unifai-cli
make build
./bin/unifai version
```

## Quick Start

1. Get an API key from [Unifai](https://unifai.network).
2. Initialize configuration:

```bash
unifai config init
```

3. Edit `~/.config/unifai-cli/config.yaml` and add your API key.
4. Run your first search:

```bash
unifai search --query "swap usdc to sol"
```

## Commands

### search

Search available blockchain services and actions.

```bash
unifai search --query "swap usdc to sol"
unifai search --query "transfer tokens" --limit 10 --offset 0
unifai search --query "defi protocols" --include-actions action1,action2
```

### invoke

Execute a blockchain action.

```bash
unifai invoke --action "Meteora--29--swap" --payload '{"amount":100}'
unifai invoke --action "MyAction--1--execute" --payload @payload.json
unifai invoke --action "MyAction--1--execute" --payload '{"x":1}' --max-retries 3 --timeout 60s
```

### config

Manage configuration.

```bash
unifai config init    # Create config template
unifai config show    # Show effective config and sources
```

### version

Print build version.

```bash
unifai version
```

## Configuration

### Config file

Location: `~/.config/unifai-cli/config.yaml`

```yaml
apiKey: "your-api-key"
endpoint: ""
timeoutSeconds: 50
```

Reference template: [`configs/config.example.yaml`](configs/config.example.yaml)

### Priority order (highest to lowest)

| Setting  | Flag         | Environment variable   | Config file      |
|----------|--------------|------------------------|------------------|
| API key  | `--api-key`  | `UNIFAI_AGENT_API_KEY` | `apiKey`         |
| Endpoint | `--endpoint` | `UNIFAI_ENDPOINT`      | `endpoint`       |
| Timeout  | `--timeout`  | —                      | `timeoutSeconds` |

## Payload formats

`invoke` supports flexible payload handling:

- Inline JSON: `--payload '{"a":1}'`
- From file: `--payload @payload.json`
- Force mode: `--payload-format auto|object|string`

| Mode     | Behavior                                         |
|----------|--------------------------------------------------|
| `auto`   | Valid JSON parsed as object; otherwise raw string |
| `object` | Must be valid JSON; error if not                  |
| `string` | Sent as raw string                                |

## Output and exit codes

- Default: concise human-readable output
- `--json`: raw API JSON response
- `invoke` extracts `payload` field from response when present

| Exit code | Meaning              |
|-----------|----------------------|
| 0         | Success              |
| 1         | API / runtime error  |
| 2         | Argument/usage error |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
