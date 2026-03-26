#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

TARGET_TAG="v0.1.6"
TARGET_TAP_REPO="BigCactusLabs/tap"
TARGET_FORMULA="BigCactusLabs/tap/codex-mp"
FIXED_HEAD_SHA="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
FIXED_ASSET_SHA256="abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"

write_stub() {
  local name="$1"
  cat >"$FAKE_BIN/$name"
  chmod +x "$FAKE_BIN/$name"
}

assert_log_sequence() {
  local file="$1"
  shift

  local previous_line=0
  local needle line

  for needle in "$@"; do
    line="$(grep -nF "$needle" "$file" | head -n1 | cut -d: -f1 || true)"
    [[ -n "$line" ]] || {
      echo "FAIL: expected log to contain: $needle" >&2
      cat "$file" >&2
      exit 1
    }
    [[ "$line" -gt "$previous_line" ]] || {
      echo "FAIL: expected log entry after line $previous_line: $needle" >&2
      cat "$file" >&2
      exit 1
    }
    previous_line="$line"
  done
}

ROOT="$(mktemp -d)"
FAKE_BIN="$ROOT/bin"
LOG_FILE="$ROOT/invocations.log"
TAP_DIR="$ROOT/tap"

export LOG_FILE FIXED_HEAD_SHA FIXED_ASSET_SHA256 TAP_DIR

trap 'rm -rf "$ROOT"' EXIT

mkdir -p "$FAKE_BIN" "$TAP_DIR/Formula" "$TAP_DIR/.git"
: >"$LOG_FILE"

write_stub git <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="${LOG_FILE:?}"
FIXED_HEAD_SHA="${FIXED_HEAD_SHA:?}"

printf 'git %s\n' "$*" >>"$LOG_FILE"

args=("$@")
if [[ "${args[0]:-}" == "-C" ]]; then
  shift 2
  args=("$@")
fi

case "${args[0]:-}" in
  status)
    [[ "${args[1]:-}" == "--short" ]] || {
      echo "unexpected git status invocation: $*" >&2
      exit 1
    }
    exit 0
    ;;
  branch)
    [[ "${args[1]:-}" == "--show-current" ]] || {
      echo "unexpected git branch invocation: $*" >&2
      exit 1
    }
    printf 'main\n'
    exit 0
    ;;
  rev-parse)
    case "${args[1]:-}" in
      HEAD|origin/main)
        printf '%s\n' "$FIXED_HEAD_SHA"
        exit 0
        ;;
      *)
        echo "unexpected git rev-parse invocation: $*" >&2
        exit 1
        ;;
    esac
    ;;
  remote)
    [[ "${args[1]:-}" == "get-url" && "${args[2]:-}" == "origin" ]] || {
      echo "unexpected git remote invocation: $*" >&2
      exit 1
    }
    printf 'https://github.com/BigCactusLabs/codex-multipass.git\n'
    exit 0
    ;;
  fetch|checkout|pull|reset)
    exit 0
    ;;
  *)
    echo "unexpected git invocation: $*" >&2
    exit 1
    ;;
esac
EOF

write_stub gh <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="${LOG_FILE:?}"

printf 'gh %s\n' "$*" >>"$LOG_FILE"

if [[ "$1" == "release" && "${2:-}" == "view" && "${3:-}" == "v0.1.6" ]]; then
  exit 1
fi

if [[ "$1" == "release" && "${2:-}" == "create" ]]; then
  exit 0
fi

echo "unexpected gh invocation: $*" >&2
exit 1
EOF

write_stub brew <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="${LOG_FILE:?}"
TAP_DIR="${TAP_DIR:?}"

printf 'brew %s\n' "$*" >>"$LOG_FILE"

if [[ "${1:-}" == "--repository" && "${2:-}" == "BigCactusLabs/tap" ]]; then
  printf '%s\n' "$TAP_DIR"
  exit 0
fi

if [[ "${1:-}" == "info" && "${2:-}" == "BigCactusLabs/tap/codex-mp" ]]; then
  printf 'codex-mp: stable 0.1.6\n'
  exit 0
fi

if [[ "${1:-}" == "fetch" && "${2:-}" == "--force" && "${3:-}" == "--build-from-source" && "${4:-}" == "BigCactusLabs/tap/codex-mp" ]]; then
  exit 0
fi

echo "unexpected brew invocation: $*" >&2
exit 1
EOF

write_stub shasum <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="${LOG_FILE:?}"
FIXED_ASSET_SHA256="${FIXED_ASSET_SHA256:?}"

printf 'shasum %s\n' "$*" >>"$LOG_FILE"

if [[ "${1:-}" == "-a" && "${2:-}" == "256" ]]; then
  printf '%s  -\n' "$FIXED_ASSET_SHA256"
  exit 0
fi

echo "unexpected shasum invocation: $*" >&2
exit 1
EOF

set +e
PATH="$FAKE_BIN:$PATH" "$ROOT_DIR/scripts/publish_homebrew_release.sh" --dry-run >>"$LOG_FILE" 2>&1
status=$?
set -e

[[ $status -eq 0 ]] || {
  echo "FAIL: scripts/publish_homebrew_release.sh exited $status" >&2
  cat "$LOG_FILE" >&2
  exit 1
}

assert_log_sequence "$LOG_FILE" \
  "git status --short" \
  "git branch --show-current" \
  "git rev-parse HEAD" \
  "git rev-parse origin/main" \
  "git remote get-url origin" \
  "gh release view $TARGET_TAG" \
  "gh release create $TARGET_TAG" \
  "brew --repository $TARGET_TAP_REPO" \
  "git -C $TAP_DIR status --short" \
  "git -C $TAP_DIR branch --show-current" \
  "shasum -a 256" \
  "scripts/update_formula.sh" \
  "brew info $TARGET_FORMULA" \
  "brew fetch --force --build-from-source $TARGET_FORMULA"
