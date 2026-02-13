#!/usr/bin/env bash
set -euo pipefail

# Smoke test for codex-switch Phase 1
# Usage: ./tests/smoke.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export CODEX_HOME="$(mktemp -d)"
CODEX_SWITCH="$SCRIPT_DIR/../cli/codex-switch"

echo "Using temporary CODEX_HOME: $CODEX_HOME"

# Clean up on exit
trap 'rm -rf "$CODEX_HOME"' EXIT

# 1. Init
echo "Testing: init"
"$CODEX_SWITCH" init
[[ -d "$CODEX_HOME" ]]
[[ -d "$CODEX_HOME/profiles" ]]

# Verify permissions (stat -f %Lp on macOS)
if [[ "$OSTYPE" == "darwin"* ]]; then
    [[ "$(stat -f %Lp "$CODEX_HOME")" == "700" ]]
    [[ "$(stat -f %Lp "$CODEX_HOME/profiles")" == "700" ]]
fi

# 2. Path
echo "Testing: path"
"$CODEX_SWITCH" path | grep "CODEX_HOME=$CODEX_HOME" > /dev/null

# 3. Save (without auth.json should fail)
echo "Testing: save (failing case)"
set +e
"$CODEX_SWITCH" save test-profile 2>/dev/null
[[ $? -eq 1 ]]
set -e

# 4. Save (with auth.json)
echo "Testing: save"
echo '{"token": "dummy-token"}' > "$CODEX_HOME/auth.json"
chmod 600 "$CODEX_HOME/auth.json"
"$CODEX_SWITCH" save work
[[ -f "$CODEX_HOME/profiles/work.json" ]]
if [[ "$OSTYPE" == "darwin"* ]]; then
    [[ "$(stat -f %Lp "$CODEX_HOME/profiles/work.json")" == "600" ]]
fi

# 5. List
echo "Testing: list"
"$CODEX_SWITCH" list | grep "work" > /dev/null

# 6. Who
echo "Testing: who"
FINGERPRINT=$("$CODEX_SWITCH" who)
[[ -n "$FINGERPRINT" ]]
[[ ! "$FINGERPRINT" =~ "token" ]] # Ensure token is not leaked

# 7. Use
echo "Testing: use"
# Create another profile
echo '{"token": "other-token"}' > "$CODEX_HOME/auth.json"
"$CODEX_SWITCH" save personal
"$CODEX_SWITCH" use work
grep "dummy-token" "$CODEX_HOME/auth.json" > /dev/null

# 8. Delete
echo "Testing: delete"
"$CODEX_SWITCH" delete personal
[[ ! -f "$CODEX_HOME/profiles/personal.json" ]]

# 9. Delete (failing case)
echo "Testing: delete (failing case)"
set +e
"$CODEX_SWITCH" delete non-existent 2>/dev/null
[[ $? -eq 1 ]]
set -e

# 10. Rename
echo "Testing: rename"
"$CODEX_SWITCH" rename work job
[[ -f "$CODEX_HOME/profiles/job.json" ]]
[[ ! -f "$CODEX_HOME/profiles/work.json" ]]

# 11. Rename (failing case - collision)
echo "Testing: rename (collision)"
"$CODEX_SWITCH" save hobby
set +e
"$CODEX_SWITCH" rename hobby job 2>/dev/null
[[ $? -eq 1 ]]
set -e

# 12. Name validation (new commands)
echo "Testing: name validation (rename)"
set +e
"$CODEX_SWITCH" rename job "invalid name!" 2>/dev/null
[[ $? -eq 2 ]]
set -e

echo "Phase 2 smoke tests PASSED!"
