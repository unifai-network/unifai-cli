#!/usr/bin/env bash
set -euo pipefail

if ! command -v goreleaser >/dev/null 2>&1; then
  echo "goreleaser is required. Install from https://goreleaser.com/install/" >&2
  exit 1
fi

if [[ "${1:-}" == "--snapshot" ]]; then
  goreleaser release --snapshot --clean
  exit 0
fi

if [[ -z "${HOMEBREW_TAP_GITHUB_TOKEN:-}" ]]; then
  echo "HOMEBREW_TAP_GITHUB_TOKEN is required for publishing Homebrew formula." >&2
  exit 1
fi

goreleaser release --clean
