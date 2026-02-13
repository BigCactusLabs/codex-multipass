#!/usr/bin/env bash
set -euo pipefail

# Concurrency test for codex-mp
# Usage: ./tests/concurrency_test.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CODEX_MP="${CODEX_MP:-$SCRIPT_DIR/../codex-mp}"

BASE_TEMP=$(mktemp -d)
trap 'rm -rf "$BASE_TEMP"' EXIT

export CODEX_HOME="$BASE_TEMP"
"$CODEX_MP" init

echo "Testing concurrency (multiple saves)..."

# Create a dummy auth.json
echo '{"token": "initial-token"}' > "$CODEX_HOME/auth.json"

# Launch multiple save operations in parallel
NUM_CONCURRENT=20
for i in $(seq 1 $NUM_CONCURRENT); do
    "$CODEX_MP" save "profile-$i" &
done

wait

# Check if all profiles exist
for i in $(seq 1 $NUM_CONCURRENT); do
    if [[ ! -f "$CODEX_HOME/profiles/profile-$i.json" ]]; then
        echo "FAIL: profile-$i.json was not saved properly"
        exit 1
    fi
done

echo "Testing concurrency (alternating use)..."
# This is more likely to hit locking if they are all trying to modify auth.json
for i in $(seq 1 $NUM_CONCURRENT); do
    # Create the profile first
    echo "{\"token\": \"token-$i\"}" > "$CODEX_HOME/profiles/p$i.json"
done

for i in $(seq 1 $NUM_CONCURRENT); do
    "$CODEX_MP" use "p$i" &
done

wait

# The result should be one of the tokens
FINAL_TOKEN=$(grep -o "token-[0-9]*" "$CODEX_HOME/auth.json" || true)
if [[ -n "$FINAL_TOKEN" ]]; then
    echo "Final token: $FINAL_TOKEN"
else
    echo "FAIL: auth.json content is corrupted or missing token"
    exit 1
fi

echo "✓ Concurrency tests passed"
