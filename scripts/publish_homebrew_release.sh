#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION="$(tr -d '[:space:]' < "$ROOT_DIR/VERSION")"
TAG="v$VERSION"
OWNER="BigCactusLabs"
REPO="codex-multipass"
TAP_SLUG="$OWNER/tap"
FORMULA_NAME="codex-mp"
FORMULA_SLUG="$TAP_SLUG/$FORMULA_NAME"
ASSET_NAME="codex-multipass-$TAG.tar.gz"
REPO_REMOTE="https://github.com/$OWNER/$REPO.git"
TAP_REMOTE_HTTPS="https://github.com/$OWNER/homebrew-tap.git"
TAP_REMOTE_SSH="git@github.com:$OWNER/homebrew-tap.git"

DRY_RUN=false
SKIP_SMOKE_TEST=false
TAP_DIR_OVERRIDE=""

usage() {
  cat <<'USAGE'
Usage: scripts/publish_homebrew_release.sh [--tap-dir <path>] [--skip-smoke-test] [--dry-run]

Creates or updates the Homebrew release asset, refreshes the tap formula, and
validates the formula locally before pushing the tap update.
USAGE
}

die() {
  echo "Error: $*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"
}

log_cmd() {
  echo "$*"
}

run_logged() {
  log_cmd "$*"
  if [[ "$DRY_RUN" != "true" ]]; then
    "$@"
  fi
}

release_asset_url() {
  printf 'https://github.com/%s/%s/releases/download/%s/%s\n' "$OWNER" "$REPO" "$TAG" "$ASSET_NAME"
}

parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --tap-dir)
        [[ $# -ge 2 ]] || die "--tap-dir requires a path"
        TAP_DIR_OVERRIDE="$2"
        shift 2
        ;;
      --skip-smoke-test)
        SKIP_SMOKE_TEST=true
        shift
        ;;
      --dry-run)
        DRY_RUN=true
        shift
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        usage
        exit 2
        ;;
    esac
  done
}

validate_version() {
  [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "VERSION must match semantic version format (e.g. 0.1.0)."
}

check_repo_state() {
  local git_status current_branch head_sha origin_main remote_url

  git_status="$(git status --short)"
  [[ -z "$git_status" ]] || die "Worktree must be clean before publishing"

  current_branch="$(git branch --show-current)"
  head_sha="$(git rev-parse HEAD)"
  origin_main="$(git rev-parse origin/main)"
  remote_url="$(git remote get-url origin)"

  [[ "$remote_url" == "$REPO_REMOTE" ]] || die "origin must point to $REPO_REMOTE"

  CURRENT_BRANCH="$current_branch"
  HEAD_SHA="$head_sha"
  ORIGIN_MAIN_SHA="$origin_main"
}

resolve_tap_dirs() {
  CANONICAL_TAP_DIR="$(brew --repository "$TAP_SLUG")"
  TAP_DIR="${TAP_DIR_OVERRIDE:-$CANONICAL_TAP_DIR}"
}

prepare_tap_checkout() {
  [[ -d "$TAP_DIR" ]] || die "Tap checkout not found at $TAP_DIR"

  if [[ "$DRY_RUN" != "true" ]]; then
    local tap_remote tap_status

    tap_remote="$(git -C "$TAP_DIR" remote get-url origin)"
    [[ "$tap_remote" == "$TAP_REMOTE_HTTPS" || "$tap_remote" == "$TAP_REMOTE_SSH" ]] || {
      die "Tap checkout origin must point to $OWNER/homebrew-tap"
    }

    tap_status="$(git -C "$TAP_DIR" status --short)"
    [[ -z "$tap_status" ]] || die "Tap checkout must be clean before publishing"
  else
    git -C "$TAP_DIR" status --short >/dev/null
    git -C "$TAP_DIR" branch --show-current >/dev/null
  fi

  if [[ "$DRY_RUN" != "true" ]]; then
    git -C "$TAP_DIR" status --short >/dev/null
    git -C "$TAP_DIR" branch --show-current >/dev/null
    run_logged git -C "$TAP_DIR" fetch origin
    run_logged git -C "$TAP_DIR" checkout main
    run_logged git -C "$TAP_DIR" pull --ff-only origin main
  fi

  if [[ -n "$TAP_DIR_OVERRIDE" && "$TAP_DIR_OVERRIDE" != "$CANONICAL_TAP_DIR" ]]; then
    run_logged brew tap --custom-remote "$TAP_SLUG" "$TAP_DIR"
  fi
}

ensure_release_is_publishable() {
  if gh release view "$TAG" --repo "$OWNER/$REPO" >/dev/null 2>&1; then
    RELEASE_EXISTS=true
    return
  fi

  RELEASE_EXISTS=false
  [[ "$CURRENT_BRANCH" == "main" ]] || die "Missing release requires publishing from main"
  [[ "$HEAD_SHA" == "$ORIGIN_MAIN_SHA" ]] || die "Missing release requires HEAD to match origin/main"
}

prepare_tarball() {
  local tmp_base

  tmp_base="$(mktemp "${TMPDIR:-/tmp}/codex-mp-release.XXXXXX")"
  TMP_TARBALL="$tmp_base.tar.gz"
  rm -f "$tmp_base"
  trap 'rm -f "$TMP_TARBALL"' EXIT

  if [[ "$DRY_RUN" == "true" ]]; then
    : >"$TMP_TARBALL"
  else
    git archive --format=tar --prefix="codex-multipass-$VERSION/" "$TAG" | gzip -n >"$TMP_TARBALL"
  fi
}

calculate_asset_sha() {
  ASSET_SHA="$(shasum -a 256 "$TMP_TARBALL" | awk '{print $1}')"
}

publish_release_asset() {
  if [[ "$RELEASE_EXISTS" == "false" ]]; then
    run_logged gh release create "$TAG" "$TMP_TARBALL" --repo "$OWNER/$REPO" --title "$TAG" --notes ""
    return
  fi

  if [[ "$DRY_RUN" == "true" ]]; then
    run_logged gh release view "$TAG" --repo "$OWNER/$REPO" --json assets,targetCommitish
    run_logged gh release upload "$TAG" "$TMP_TARBALL#$ASSET_NAME" --repo "$OWNER/$REPO"
    return
  fi

  gh release view "$TAG" --repo "$OWNER/$REPO" --json assets,targetCommitish >/dev/null
  die "Existing release asset reconciliation is not implemented yet; inspect the release manually before rerunning."
}

update_formulae() {
  # The release asset must exist before Homebrew can fetch and verify it.
  # If validation fails after release creation, exit loudly and do not push the tap.
  run_logged "$ROOT_DIR/scripts/update_formula.sh" \
    --url "$(release_asset_url)" \
    --sha256 "$ASSET_SHA" \
    --tap-dir "$TAP_DIR"
}

run_validation_or_die() {
  local validation_failed=false

  if ! run_logged brew info "$FORMULA_SLUG"; then
    validation_failed=true
  fi

  if ! run_logged brew fetch --force --build-from-source "$FORMULA_SLUG"; then
    validation_failed=true
  fi

  if [[ "$SKIP_SMOKE_TEST" != "true" ]]; then
    if [[ "$DRY_RUN" == "true" ]]; then
      run_logged HOMEBREW_NO_INSTALL_FROM_API=1 brew install --build-from-source "$FORMULA_SLUG"
      run_logged brew test "$FORMULA_NAME"
    else
      if brew list "$FORMULA_NAME" >/dev/null 2>&1; then
        if ! HOMEBREW_NO_INSTALL_FROM_API=1 brew reinstall --build-from-source "$FORMULA_SLUG"; then
          validation_failed=true
        fi
      else
        if ! HOMEBREW_NO_INSTALL_FROM_API=1 brew install --build-from-source "$FORMULA_SLUG"; then
          validation_failed=true
        fi
      fi

      if ! brew test "$FORMULA_NAME"; then
        validation_failed=true
      fi
    fi
  fi

  if [[ "$validation_failed" == "true" ]]; then
    echo "Release asset uploaded, but tap update was not pushed." >&2
    echo "Fix the validation failure, then rerun the publisher." >&2
    echo "If the release asset itself is wrong, delete the asset manually before rerunning." >&2
    exit 1
  fi
}

finalize_tap() {
  run_logged git -C "$TAP_DIR" add "Formula/$FORMULA_NAME.rb"
  run_logged git -C "$TAP_DIR" commit -m "$FORMULA_NAME $TAG"
  run_logged git -C "$TAP_DIR" push origin HEAD
}

main() {
  parse_args "$@"
  require_cmd git
  require_cmd gh
  require_cmd brew
  require_cmd shasum
  validate_version
  check_repo_state
  ensure_release_is_publishable
  prepare_tarball
  publish_release_asset
  resolve_tap_dirs
  prepare_tap_checkout
  calculate_asset_sha
  update_formulae
  run_validation_or_die
  finalize_tap
}

main "$@"
