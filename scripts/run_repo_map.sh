#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-generate}"
SUMMARY="${2:-Regenerated folder map snapshot}"
shift $(( $# > 0 ? 1 : 0 ))
shift $(( $# > 0 ? 1 : 0 ))
EXTRA_ARGS=("$@")

if command -v git >/dev/null 2>&1; then
  REPO_ROOT="${REPO_ROOT:-$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
else
  REPO_ROOT="${REPO_ROOT:-$(pwd)}"
fi

GENERATOR="${REPO_MAP_GENERATOR:-tools/docs/generate_folder_map.py}"
if [[ "$GENERATOR" != /* ]]; then
  GENERATOR="$REPO_ROOT/$GENERATOR"
fi

if [[ ! -f "$GENERATOR" ]]; then
  echo "Missing generator: $GENERATOR" >&2
  echo "Set REPO_MAP_GENERATOR to the repository's generator path." >&2
  exit 1
fi

case "$MODE" in
  check)
    python3 "$GENERATOR" --repo-root "$REPO_ROOT" --check --bump none "${EXTRA_ARGS[@]}"
    ;;
  generate|patch)
    python3 "$GENERATOR" --repo-root "$REPO_ROOT" --bump patch --summary "$SUMMARY" "${EXTRA_ARGS[@]}"
    ;;
  minor)
    python3 "$GENERATOR" --repo-root "$REPO_ROOT" --bump minor --summary "$SUMMARY" "${EXTRA_ARGS[@]}"
    ;;
  major)
    python3 "$GENERATOR" --repo-root "$REPO_ROOT" --bump major --summary "$SUMMARY" "${EXTRA_ARGS[@]}"
    ;;
  *)
    echo "Unsupported mode: $MODE" >&2
    echo "Usage: run_repo_map.sh [generate|check|patch|minor|major] [summary] [extra generator args...]" >&2
    exit 2
    ;;
esac
