#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION_FILE="$ROOT_DIR/VERSION"
FORMULA_FILE="$ROOT_DIR/homebrew/Formula/codex-mp.rb"

usage() {
  cat <<'USAGE'
Usage: scripts/update_formula.sh <sha256> [tap_dir]

Reads version from VERSION and updates the formula in:
1. homebrew/Formula/codex-mp.rb (local placeholder)
2. [tap_dir]/Formula/codex-mp.rb (if provided)
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

update_file() {
  local target="$1"
  local tmp
  tmp="$(mktemp "$target.tmp.XXXXXX")"
  awk -v version="$VERSION" -v sha="$SHA256" '
    /url "https:\/\/github\.com\/BigCactusLabs\/codex-multipass\/archive\/refs\/tags\/v[0-9]+\.[0-9]+\.[0-9]+\.tar\.gz"/ {
      print "  url \"https://github.com/BigCactusLabs/codex-multipass/archive/refs/tags/v" version ".tar.gz\""
      next
    }
    /sha256 "[0-9a-f]+"/ {
      print "  sha256 \"" sha "\""
      next
    }
    { print }
  ' "$target" > "$tmp"
  mv "$tmp" "$target"
}

# Update local copy
update_file "$FORMULA_FILE"
echo "Updated local formula: $FORMULA_FILE"

# Update tap if provided
TAP_DIR="${2:-}"
if [[ -n "$TAP_DIR" ]]; then
  TAP_FORMULA="$TAP_DIR/Formula/codex-mp.rb"
  if [[ -f "$TAP_FORMULA" ]]; then
    update_file "$TAP_FORMULA"
    echo "Updated tap formula: $TAP_FORMULA"
  else
    echo "Warning: Tap formula not found at $TAP_FORMULA" >&2
  fi
fi

echo "Success: Version $VERSION, SHA256 $SHA256"
