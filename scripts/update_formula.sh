#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION_FILE="$ROOT_DIR/VERSION"
FORMULA_FILE="$ROOT_DIR/homebrew/Formula/codex-switch.rb"

usage() {
  cat <<'USAGE'
Usage: scripts/update_formula.sh <sha256>

Reads version from VERSION and updates:
- url ".../refs/tags/v<version>.tar.gz"
- sha256 "<sha256>"
USAGE
}

[[ -f "$VERSION_FILE" ]] || {
  echo "Error: missing VERSION file at $VERSION_FILE" >&2
  exit 1
}
[[ -f "$FORMULA_FILE" ]] || {
  echo "Error: missing formula file at $FORMULA_FILE" >&2
  exit 1
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

SHA256="${1:-}"
[[ -n "$SHA256" ]] || {
  usage
  exit 2
}
if [[ ! "$SHA256" =~ ^[0-9a-f]{64}$ ]]; then
  echo "Error: sha256 must be 64 lowercase hex characters." >&2
  exit 2
fi

VERSION="$(tr -d '[:space:]' < "$VERSION_FILE")"
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: VERSION must match semantic version format (e.g. 0.1.0)." >&2
  exit 2
fi

TMP_FILE="$(mktemp "$FORMULA_FILE.tmp.XXXXXX")"
awk -v version="$VERSION" -v sha="$SHA256" '
  /url "https:\/\/github\.com\/quinn\/multidex\/archive\/refs\/tags\/v[0-9]+\.[0-9]+\.[0-9]+\.tar\.gz"/ {
    print "  url \"https://github.com/quinn/multidex/archive/refs/tags/v" version ".tar.gz\""
    next
  }
  /sha256 "[0-9a-f]+"/ {
    print "  sha256 \"" sha "\""
    next
  }
  { print }
' "$FORMULA_FILE" > "$TMP_FILE"

mv "$TMP_FILE" "$FORMULA_FILE"

echo "Updated formula to version $VERSION with sha256 $SHA256"
