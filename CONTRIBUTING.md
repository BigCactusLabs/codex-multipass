# Contributing

We welcome contributions! The project is written in **Go**.

## Development Setup

1.  **Install Go**: 1.23 or later.
2.  **Install Dependencies**: `make tidy`
3.  **Build**: `make build`
4.  **Run Integration Tests**: `make test`
5.  **Run Unit Tests**: `cd go && go test ./internal/app ./internal/profile`

## Project Structure

- `go/cmd/codex-mp`: Main entry point.

- `go/internal/app`: CLI commands and wiring.
- `go/internal/profile`: Core profile management logic.
- `go/internal/config`: Configuration handling.
- `go/internal/ui`: User interface components.
- `go/internal/fs`: Atomic file system operations.
- `bash/codex-switch`: Compatibility wrapper that delegates to `codex-mp`.
- `tests/`: Integration tests (Bash scripts).

## Testing Notes

- Integration test scripts read the binary path from `CODEX_MP`.
- `make test` builds `./codex-mp` and runs smoke, battle, concurrency, and corrupt-storage scripts against it.
- Do not commit generated artifacts (for example `__pycache__/` and `*.pyc`).
- You can run scripts directly, for example:
  - `CODEX_MP=./codex-mp ./tests/smoke.sh`
  - `CODEX_MP=./codex-mp ./tests/battle.sh`
  - `CODEX_MP=./codex-mp ./tests/concurrency_test.sh`
  - `CODEX_MP=./codex-mp ./tests/corrupt_storage_test.sh`

## CI

- Test job should build `codex-mp`, then run smoke, battle, concurrency, and corrupt-storage scripts.
- Lint job runs `shellcheck bash/codex-switch scripts/*.sh tests/*.sh`.

## Release Workflow

1.  Update `VERSION` and commit the release changes.
2.  Run `./scripts/publish_homebrew_release.sh` from a clean worktree.
3.  Verify `brew info BigCactusLabs/tap/codex-mp` and `brew fetch --force --build-from-source BigCactusLabs/tap/codex-mp`.

Do not try to retrofit an old release from a branch with newer unreleased changes; publish from the intended release version or from a detached worktree at the existing tag.
