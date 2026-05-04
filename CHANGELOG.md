# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-05-04

### Added
- Added repo-local `AGENTS.md` with Codex-specific build, test, lint, and credential-handling guidance.
- Added `codex-mp doctor` and JSON doctor output for checking Codex auth/profile compatibility.

### Changed
- `make test` now runs smoke, battle, concurrency, and corrupt-storage integration scripts.
- README and contributor docs now explain that `codex-mp` only manages file-backed `auth.json` sessions when Codex is configured to use `cli_auth_credentials_store = "file"`.
- JSON command output now uses structured encoding instead of hand-built strings.
- Go module dependencies are refreshed, the module now targets Go 1.25+, and CI uses current checkout/setup-go/golangci-lint actions.
- `make lint` now runs ShellCheck and a pinned `golangci-lint` v2.12.1.

### Fixed
- `save`, `use`, and `who` now fail with actionable guidance when Codex is configured for `keyring`, `auto`, or `ephemeral` credential storage instead of explicit file-backed `auth.json`.
- Smoke test negative cases now enforce exit codes instead of silently passing.
- `make clean` no longer deletes the tracked `go/go.sum` file.

## [0.1.6] - 2026-02-25

### Changed
- CI now builds `codex-mp` before running smoke and battle test scripts.
- CI shell lint now targets `bash/codex-switch` directly.
- `bash/codex-switch` now acts as a compatibility wrapper that delegates to `codex-mp`.
- `make test` wiring and script docs now consistently use `CODEX_MP`.

### Fixed
- Corrupt storage test expectation now matches runtime behavior when `profiles/` is missing.
- CLI failure behavior is now testable in unit tests via injected exit handling.

### Removed
- Removed tracked compiled Python cache artifact from docs tooling.

## [0.1.3] - 2026-02-13

### Fixed
- Fixed version mismatch where v0.1.2 tag reported version 0.1.1.
- Updated Homebrew formula to sync with correct release.

## [0.1.2] - 2026-02-13
### Yanked
- Release created with incorrect version string (0.1.1).

## [0.1.1] - 2026-02-13
### Changed
- Maintenance updates.

## [0.1.0] - 2026-02-13

### Added
- Initial public release of `codex-mp` CLI.
- Commands: `init`, `save`, `use`, `list`, `who`, `path`, `delete`, `rename`, `version`.
- Interactive TUI mode (`pick`, `ui`).
- Atomic profile switching with permission hardening.
- Homebrew formula for easy installation.
