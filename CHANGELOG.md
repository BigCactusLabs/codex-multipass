# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
