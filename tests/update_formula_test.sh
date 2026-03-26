#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

TARGET_URL="https://github.com/BigCactusLabs/codex-multipass/releases/download/v0.1.6/codex-multipass-v0.1.6.tar.gz"
TARGET_SHA256="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

make_fixture_root() {
  local root="$1"
  local formula_contents="$2"

  mkdir -p "$root/scripts" "$root/homebrew/Formula" "$root/tap/Formula"
  cp "$REPO_ROOT/scripts/update_formula.sh" "$root/scripts/update_formula.sh"
  cp "$REPO_ROOT/homebrew/Formula/codex-mp.rb" "$root/homebrew/Formula/codex-mp.rb"
  cp "$REPO_ROOT/homebrew/Formula/codex-mp.rb" "$root/tap/Formula/codex-mp.rb"
  printf '0.1.6\n' > "$root/VERSION"

  if [[ -n "$formula_contents" ]]; then
    printf '%s\n' "$formula_contents" > "$root/homebrew/Formula/codex-mp.rb"
    printf '%s\n' "$formula_contents" > "$root/tap/Formula/codex-mp.rb"
  fi
}

assert_contains() {
  local file="$1"
  local needle="$2"

  grep -F "$needle" "$file" >/dev/null || {
    echo "FAIL: expected $file to contain: $needle" >&2
    exit 1
  }
}

scenario_updates_local_and_tap_formula() {
  local root tap_dir
  root="$(mktemp -d)"
  tap_dir="$root/tap"
  trap 'rm -rf "$root"' RETURN

  make_fixture_root "$root" ""

  (
    cd "$root"
    scripts/update_formula.sh --url "$TARGET_URL" --sha256 "$TARGET_SHA256" --tap-dir "$tap_dir"
  )

  assert_contains "$root/homebrew/Formula/codex-mp.rb" "$TARGET_URL"
  assert_contains "$root/homebrew/Formula/codex-mp.rb" "$TARGET_SHA256"
  assert_contains "$tap_dir/Formula/codex-mp.rb" "$TARGET_URL"
  assert_contains "$tap_dir/Formula/codex-mp.rb" "$TARGET_SHA256"
}

scenario_rejects_formula_missing_sha256() {
  local root tap_dir malformed_formula
  root="$(mktemp -d)"
  tap_dir="$root/tap"
  trap 'rm -rf "$root"' RETURN

  malformed_formula='class CodexMp < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/BigCactusLabs/codex-multipass"
  url "https://github.com/BigCactusLabs/codex-multipass/archive/refs/tags/v0.1.6.tar.gz"
  license "MIT"
end'

  make_fixture_root "$root" "$malformed_formula"

  set +e
  (
    cd "$root"
    scripts/update_formula.sh --url "$TARGET_URL" --sha256 "$TARGET_SHA256" --tap-dir "$tap_dir"
  )
  status=$?
  set -e

  [[ $status -ne 0 ]] || {
    echo "FAIL: expected updater to exit non-zero for malformed formula missing sha256" >&2
    exit 1
  }
}

scenario_updates_local_and_tap_formula
scenario_rejects_formula_missing_sha256
