#!/usr/bin/env bash
set -euo pipefail

# Corruption/Edge Case test for codex-switch
# Usage: ./tests/corrupt_storage_test.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CODEX_MP="${CODEX_MP:-$SCRIPT_DIR/../codex-mp}"

BASE_TEMP=$(mktemp -d)
trap 'rm -rf "$BASE_TEMP"' EXIT

export CODEX_HOME="$BASE_TEMP"
"$CODEX_MP" init

echo "Testing: Corrupted auth.json (Binary data)"
dd if=/dev/urandom of="$CODEX_HOME/auth.json" bs=1024 count=1 2>/dev/null
"$CODEX_MP" save corrupt-test
[[ -f "$CODEX_HOME/profiles/corrupt-test.json" ]] || exit 1

echo "Testing: Missing profiles directory"
rm -rf "$CODEX_HOME/profiles"
set +e
"$CODEX_MP" save missing-dir 2>/dev/null
RET=$?
set -e
[[ $RET -ne 0 ]] || { echo "FAIL: Should fail when profiles dir is missing"; exit 1; }

echo "Testing: Profile name with path traversal"
set +e
"$CODEX_MP" save "../hacker" 2>/dev/null
RET=$?
set -e
[[ $RET -ne 0 ]] || { echo "FAIL: Should fail with path traversal in name"; exit 1; }

echo "Testing: Profile name with null byte (if possible in bash/go)"
set +e
"$CODEX_MP" save "name"$'\0'"suffix" 2>/dev/null
RET=$?
set -e
# Most CLI parsers will stop at null, so "name" might be saved or it might fail
# Success is either failing gracefully or saving "name" but NOT being insecure

echo "Testing: auth.json is a directory"
rm "$CODEX_HOME/auth.json"
mkdir "$CODEX_HOME/auth.json"
set +e
"$CODEX_MP" save dir-auth 2>/dev/null
RET=$?
set -e
[[ $RET -ne 0 ]] || { echo "FAIL: Should fail when auth.json is a directory"; exit 1; }

echo "✓ Corruption and edge case tests completed"
