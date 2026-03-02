# Releasing

Instructions for maintainers on publishing new releases.

## Prerequisites

- [GoReleaser](https://goreleaser.com/install/) installed locally (for local releases)
- `HOMEBREW_TAP_GITHUB_TOKEN` GitHub secret configured (for CI releases)
- The tap repo `unifai-network/homebrew-unifai-cli` exists with a `main` branch

## Local release

### Snapshot (no GitHub publish)

```bash
./scripts/release.sh --snapshot
```

### Full release

```bash
export HOMEBREW_TAP_GITHUB_TOKEN=ghp_...
./scripts/release.sh
```

## GitHub Actions release

Push a semver tag to trigger the release workflow (`.github/workflows/release.yml`):

```bash
git tag v0.2.0
git push origin v0.2.0
```

GoReleaser config: `.goreleaser.yaml`

## Where artifacts go

Release artifacts are published to `unifai-network/homebrew-unifai-cli`, not this repository.
This allows Homebrew to download assets anonymously even when the source repo is private.

## Homebrew tap setup (one-time)

1. Create the tap repo `unifai-network/homebrew-unifai-cli` with a `main` branch.
2. Add the repository secret `HOMEBREW_TAP_GITHUB_TOKEN` in `unifai-cli` GitHub Actions settings.
3. Ensure the token has write access to `unifai-network/homebrew-unifai-cli`.
4. Push a tag to trigger the first release.
