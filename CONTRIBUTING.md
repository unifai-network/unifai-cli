# Contributing

Thanks for your interest in contributing to unifai-cli.

## Development setup

```bash
git clone https://github.com/unifai-network/unifai-cli.git
cd unifai-cli
make build
```

## Common commands

```bash
make tidy       # Tidy Go module dependencies
make fmt        # Format code
make test       # Run all tests
make build      # Build binary to bin/unifai
```

## Making changes

1. Fork the repository and create a feature branch.
2. Make your changes.
3. Run `make fmt && make test` to ensure code is formatted and tests pass.
4. Open a pull request against `main`.

## Reporting issues

Open a GitHub issue with steps to reproduce the problem, expected behavior, and actual behavior.

## Releasing

See [RELEASING.md](RELEASING.md) for maintainer release instructions.
