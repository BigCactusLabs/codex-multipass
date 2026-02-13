#!/usr/bin/env bash
set -euo pipefail

# Smoke test for codex-mp Phase 1
# Usage: ./tests/smoke.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CODEX_HOME="$(mktemp -d)"
export CODEX_HOME
CODEX_MP="${CODEX_MP:-$SCRIPT_DIR/../codex-mp}"

echo "Using temporary CODEX_HOME: $CODEX_HOME"

# Clean up on exit
trap 'rm -rf "$CODEX_HOME"' EXIT

# 1. Init
echo "Testing: init"
"$CODEX_MP" init
[[ -d "$CODEX_HOME" ]]
[[ -d "$CODEX_HOME/profiles" ]]

# Verify permissions (stat -f %Lp on macOS)
if [[ "$OSTYPE" == "darwin"* ]]; then
    [[ "$(stat -f %Lp "$CODEX_HOME")" == "700" ]]
    [[ "$(stat -f %Lp "$CODEX_HOME/profiles")" == "700" ]]
fi

# 2. Path
echo "Testing: path"
"$CODEX_MP" path | grep "CODEX_HOME=$CODEX_HOME" > /dev/null

# 3. Save (without auth.json should fail)
echo "Testing: save (failing case)"
set +e
"$CODEX_MP" save test-profile 2>/dev/null
[[ $? -eq 1 ]]
set -e

# 4. Save (with auth.json)
echo "Testing: save"
echo '{"token": "dummy-token"}' > "$CODEX_HOME/auth.json"
chmod 600 "$CODEX_HOME/auth.json"
"$CODEX_MP" save work
[[ -f "$CODEX_HOME/profiles/work.json" ]]
if [[ "$OSTYPE" == "darwin"* ]]; then
    [[ "$(stat -f %Lp "$CODEX_HOME/profiles/work.json")" == "600" ]]
fi

# 5. List
echo "Testing: list"
"$CODEX_MP" list | grep "work" > /dev/null

# 6. Who
echo "Testing: who"
FINGERPRINT=$("$CODEX_MP" who)
[[ -n "$FINGERPRINT" ]]
[[ ! "$FINGERPRINT" =~ "token" ]] # Ensure token is not leaked

# 7. Use
echo "Testing: use"
# Create another profile
echo '{"token": "other-token"}' > "$CODEX_HOME/auth.json"
"$CODEX_MP" save personal
"$CODEX_MP" use work
grep "dummy-token" "$CODEX_HOME/auth.json" > /dev/null

# 8. Delete
echo "Testing: delete"
"$CODEX_MP" delete personal
[[ ! -f "$CODEX_HOME/profiles/personal.json" ]]

# 9. Delete (failing case)
echo "Testing: delete (failing case)"
set +e
"$CODEX_MP" delete non-existent 2>/dev/null
[[ $? -eq 1 ]]
set -e

# 10. Rename
echo "Testing: rename"
"$CODEX_MP" rename work job
[[ -f "$CODEX_HOME/profiles/job.json" ]]
[[ ! -f "$CODEX_HOME/profiles/work.json" ]]

# 11. Rename (failing case - collision)
echo "Testing: rename (collision)"
"$CODEX_MP" save hobby
set +e
"$CODEX_MP" rename hobby job 2>/dev/null
[[ $? -eq 1 ]]
set -e

# 12. Name validation (new commands)
echo "Testing: name validation (rename)"
set +e
"$CODEX_MP" rename job "invalid name!" 2>/dev/null
[[ $? -eq 2 ]]
set -e

# 13. JSON output for mutating commands
echo "Testing: --json save/use/rename/delete"
echo '{"token": "json-token"}' > "$CODEX_HOME/auth.json"
SAVE_JSON=$("$CODEX_MP" --json save json-profile)
echo "$SAVE_JSON" | grep '"ok":true' > /dev/null
echo "$SAVE_JSON" | grep '"action":"save"' > /dev/null

USE_JSON=$("$CODEX_MP" --json use json-profile)
echo "$USE_JSON" | grep '"action":"use"' > /dev/null
grep "json-token" "$CODEX_HOME/auth.json" > /dev/null

RENAME_JSON=$("$CODEX_MP" --json rename json-profile json-profile-2)
echo "$RENAME_JSON" | grep '"action":"rename"' > /dev/null
[[ -f "$CODEX_HOME/profiles/json-profile-2.json" ]]

DELETE_JSON=$("$CODEX_MP" --json delete json-profile-2)
echo "$DELETE_JSON" | grep '"action":"delete"' > /dev/null
[[ ! -f "$CODEX_HOME/profiles/json-profile-2.json" ]]

# 14. version command uses VERSION file
echo "Testing: version"
EXPECTED_VERSION=$(tr -d '[:space:]' < "$SCRIPT_DIR/../VERSION")
ACTUAL_VERSION=$("$CODEX_MP" version)
[[ "$ACTUAL_VERSION" == "$EXPECTED_VERSION" ]]

echo "Phase 2 smoke tests PASSED!"
