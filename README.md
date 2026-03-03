# unifai-cli

A Go CLI for searching and invoking services on the [UnifAI](https://unifai.network) network — DeFi, token data, social media, web search, news, travel, sports, and more.

> **Note:** This Go CLI supports search and invoke only. It **does not support transaction signing**. For full functionality including autonomous transaction signing (Solana, EVM, Polymarket, Hyperliquid, Jito bundles, etc.), use the **recommended** JS CLI from the [unifai-sdk-js](https://github.com/unifai-network/unifai-sdk-js?tab=readme-ov-file#using-the-cli) package:
>
> ```bash
> npx -p unifai-sdk unifai search --query "swap usdc to sol"
> npx -p unifai-sdk unifai invoke --action "Jupiter--5--swap" --payload '{"amount":100}' --sign
> ```
>
> Or install globally:
>
> ```bash
> npm install -g unifai-sdk
> unifai search --query "swap usdc to sol"
> ```
>
> See more at [unifai-sdk-js](https://github.com/unifai-network/unifai-sdk-js?tab=readme-ov-file#using-the-cli).

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

1. Get an API key from [UnifAI Console](https://console.unifai.network/).
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

Execute a blockchain action. Returns a response with an approval link for on-chain transactions. Transaction signing is **not supported** in this Go CLI — use the [JS CLI](https://github.com/unifai-network/unifai-sdk-js?tab=readme-ov-file#using-the-cli) for autonomous signing.

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
