# AGENTS.md

## Repository Expectations

- This is a Go CLI for switching Codex file-backed `auth.json` profiles.
- Keep changes minimal and follow the existing package split:
  - `go/internal/app`: Cobra commands and output formatting.
  - `go/internal/profile`: profile storage, switching, locking, diagnostics.
  - `go/internal/config`: Codex path and config parsing.
  - `go/internal/fs`: atomic file and lock primitives.
- Do not print token contents, auth JSON contents, or other credential material.
- Treat `CODEX_HOME` as test-scoped state. Tests must use temporary `CODEX_HOME` values, never the user's real `~/.codex`.
- `codex-mp` currently supports only file-backed Codex auth. Do not add keyring profile storage without a separate design and cross-platform test plan.

## Commands

- Build: `make build`
- Unit tests: `make unit-test`
- Integration tests: `make test`
- Shell lint: `shellcheck bash/codex-switch scripts/*.sh tests/*.sh`
- Go lint: `make lint`
- Race check: from repo root, run `cd go && GOCACHE=$(pwd)/../.gocache GOMODCACHE=$(pwd)/../.gomodcache go test -race ./...`

## Definition Of Done

- Run focused tests for changed behavior first.
- Before claiming completion, run the relevant full verification command set or state exactly what could not be run.
- Do not commit generated artifacts such as `codex-mp`, `.gocache/`, tarballs, or temporary release files.
