#!/usr/bin/env bash
set -euo pipefail

# Battle test for codex-switch
# Usage: ./tests/battle.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CODEX_SWITCH="$SCRIPT_DIR/../cli/codex-switch"

# Use a base temp dir for all tests
BASE_TEMP=$(mktemp -d)
trap 'rm -rf "$BASE_TEMP"' EXIT

run_test() {
    local name="$1"
    echo "Running test: $name..."
    ( cd "$BASE_TEMP" && "$@" )
}

# 1. Paths with spaces
echo "--- Test: Paths with spaces ---"
SPACE_DIR="$BASE_TEMP/my codex dir"
mkdir -p "$SPACE_DIR"
export CODEX_HOME="$SPACE_DIR"

"$CODEX_SWITCH" init
[[ -d "$SPACE_DIR/profiles" ]] || exit 1

echo '{"token": "space-token"}' > "$SPACE_DIR/auth.json"
"$CODEX_SWITCH" save "myprofile" # Note: cli validates name, spaces are NOT allowed in profile names per script
[[ -f "$SPACE_DIR/profiles/myprofile.json" ]] || exit 1

"$CODEX_SWITCH" use "myprofile"
grep "space-token" "$SPACE_DIR/auth.json" > /dev/null || exit 1
echo "✓ Paths with spaces passed"

# 2. Symlinked auth.json
echo "--- Test: Symlinked auth.json ---"
SYM_DIR="$BASE_TEMP/symlink_test"
mkdir -p "$SYM_DIR/real_storage"
mkdir -p "$SYM_DIR/codex_home"
export CODEX_HOME="$SYM_DIR/codex_home"

echo '{"token": "link-token"}' > "$SYM_DIR/real_storage/auth.json"
ln -s "$SYM_DIR/real_storage/auth.json" "$CODEX_HOME/auth.json"

"$CODEX_SWITCH" init
"$CODEX_SWITCH" save linked
# cp -f should have copied the content to profiles/linked.json
grep "link-token" "$CODEX_HOME/profiles/linked.json" > /dev/null || exit 1

# Now use a different profile and see if it breaks the link
echo '{"token": "new-token"}' > "$CODEX_HOME/profiles/new.json"
"$CODEX_SWITCH" use new
grep "new-token" "$CODEX_HOME/auth.json" > /dev/null || exit 1

if [[ -L "$CODEX_HOME/auth.json" ]]; then
    echo "Note: auth.json is still a symlink (cp -f behavior on this system)"
    grep "new-token" "$SYM_DIR/real_storage/auth.json" > /dev/null || exit 1
else
    echo "Note: auth.json is no longer a symlink (cp -f replaced it)"
fi
echo "✓ Symlink test passed"

# 3. Invalid Secret contents (Large file)
echo "--- Test: Large auth.json ---"
LARGE_DIR="$BASE_TEMP/large_test"
mkdir -p "$LARGE_DIR"
export CODEX_HOME="$LARGE_DIR"
dd if=/dev/urandom of="$CODEX_HOME/auth.json" bs=1M count=1 2>/dev/null
"$CODEX_SWITCH" init
"$CODEX_SWITCH" save huge
[[ -f "$CODEX_HOME/profiles/huge.json" ]] || exit 1
# Fingerprint should still work
FP=$("$CODEX_SWITCH" who)
[[ -n "$FP" ]] || exit 1
echo "✓ Large file test passed"

# 4. Read-only profiles
echo "--- Test: Read-only profiles ---"
RO_DIR="$BASE_TEMP/ro_test"
mkdir -p "$RO_DIR"
export CODEX_HOME="$RO_DIR"
"$CODEX_SWITCH" init
echo '{"token": "ro-token"}' > "$RO_DIR/auth.json"
"$CODEX_SWITCH" save ro-profile
chmod 400 "$RO_DIR/profiles/ro-profile.json"

# Switch away
echo '{"token": "temp"}' > "$RO_DIR/auth.json"
# Switch back - should work as cp can read it
"$CODEX_SWITCH" use ro-profile
grep "ro-token" "$RO_DIR/auth.json" > /dev/null || exit 1
echo "✓ Read-only profile test passed"

# 5. Overwriting profile
echo "--- Test: Overwriting profile ---"
OVR_DIR="$BASE_TEMP/ovr_test"
mkdir -p "$OVR_DIR"
export CODEX_HOME="$OVR_DIR"
"$CODEX_SWITCH" init
echo '{"token": "v1"}' > "$OVR_DIR/auth.json"
"$CODEX_SWITCH" save p1
echo '{"token": "v2"}' > "$OVR_DIR/auth.json"
"$CODEX_SWITCH" save p1
grep "v2" "$OVR_DIR/profiles/p1.json" > /dev/null || exit 1
echo "✓ Overwrite test passed"

# 6. Malformed JSON in list (TTY mode)
echo "--- Test: Malformed JSON list ---"
MAL_DIR="$BASE_TEMP/mal_test"
mkdir -p "$MAL_DIR"
export CODEX_HOME="$MAL_DIR"
"$CODEX_SWITCH" init
echo 'not json but valid token' > "$MAL_DIR/auth.json"
"$CODEX_SWITCH" save plain
# List should show it without crashing
IS_TTY=true "$CODEX_SWITCH" list > /dev/null
echo "✓ Malformed JSON list passed"

# 7. Directory permissions error
echo "--- Test: Directory permissions error ---"
PERM_DIR="$BASE_TEMP/perm_error_test"
mkdir -p "$PERM_DIR"
chmod 000 "$PERM_DIR"
export CODEX_HOME="$PERM_DIR"
set +e
"$CODEX_SWITCH" init 2>/dev/null
[[ $? -ne 0 ]]
set -e
# Restore permissions for cleanup
chmod 700 "$PERM_DIR"
echo "✓ Directory permissions error passed"

# 8. Delete non-existent (covered by smoke, but let's double check exit code)
echo "--- Test: Delete non-existent ---"
DN_DIR="$BASE_TEMP/dn_test"
mkdir -p "$DN_DIR"
export CODEX_HOME="$DN_DIR"
"$CODEX_SWITCH" init
set +e
"$CODEX_SWITCH" delete "ghost" 2>/dev/null
RET=$?
set -e
[[ $RET -eq 1 ]]
echo "✓ Delete non-existent passed"

# 9. Rename non-existent
echo "--- Test: Rename non-existent ---"
RN_DIR="$BASE_TEMP/rn_test"
mkdir -p "$RN_DIR"
export CODEX_HOME="$RN_DIR"
"$CODEX_SWITCH" init
set +e
"$CODEX_SWITCH" rename "ghost" "real" 2>/dev/null
RET=$?
set -e
[[ $RET -eq 1 ]]
echo "✓ Rename non-existent passed"

echo "ALL BATTLE TESTS PASSED!"
