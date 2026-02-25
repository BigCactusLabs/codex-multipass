# Codex Multipass (codex-mp) v0.1.6

`codex-mp` is a local CLI for switching Codex accounts by swapping
`auth.json` profiles on disk.

`bash/codex-switch` is a compatibility wrapper that delegates to `codex-mp`.
All profile-switching behavior is implemented in the Go CLI.

## Quick Start

### 1. Install
Choose one:
- **Homebrew** (Recommended): `brew install BigCactusLabs/tap/codex-mp`
- **Build from Source**: Run `make build` (Requires [Go 1.23+](https://go.dev/doc/install) installed)

### 2. Usage
```bash
# Initialize
codex-mp init

# Save current session
codex-mp save work

# Switch (Interactive)
codex-mp ui

# Switch (Command)
codex-mp use work
```

It does not change your repos or tools. It only copies files under your Codex
state directory.

## Security model

- All data stays local.
- The CLI never prints token contents.
- `codex-mp who` prints only a SHA-256 fingerprint of `auth.json`.
- Atomic writes for `save` and `use` (temporary file + rename).
- Process lock for profile mutations to avoid concurrent-write races.
- Permission hardening on every write (fails closed if hardening fails):
  - `CODEX_DIR` mode `700`
  - `profiles/` mode `700`
  - `auth.json` mode `600`
  - `profiles/*.json` mode `600`

## Paths

By default, Codex state is resolved from:

- `CODEX_DIR=$HOME/.codex`
- `AUTH=$CODEX_DIR/auth.json`
- `PROFILES_DIR=$CODEX_DIR/profiles`

If `CODEX_HOME` is set, it overrides the base directory:

```bash
CODEX_HOME=/custom/path/.codex codex-mp path
```

## Installation

### From Source (Go)

Prerequisites: Go 1.23+

```bash
git clone https://github.com/BigCactusLabs/codex-multipass.git
cd codex-multipass
make build
# Binary is at ./codex-mp
```

If you use legacy automation that calls `codex-switch`, point it at
`bash/codex-switch`. The wrapper executes `./codex-mp` (building it if
needed and Go is available) or falls back to an installed `codex-mp`.

### Homebrew

```bash
brew install BigCactusLabs/tap/codex-mp
```

## Commands

```bash
codex-mp init
codex-mp save <name>
codex-mp use <name>
codex-mp list
codex-mp who
codex-mp path
codex-mp delete <name>
codex-mp rename <old> <new>
codex-mp pick
codex-mp ui
codex-mp version
codex-mp help
```

Global output flags:

```bash
codex-mp --plain <command>
codex-mp --json <command>
```

## Usage

### 1. Initialize
Set up the profiles directory:
```bash
codex-mp init
```

### 2. Save Profile
Save your current cached login as a profile:
```bash
codex login
codex-mp save work
```

### 3. Switch Profile
Switch to a saved profile:
```bash
codex-mp use work
```

`codex-mp` tracks the active profile and syncs the latest `auth.json`
back to that profile before switching. This preserves rotated refresh
tokens and avoids stale-token switch failures.

### 4. Interactive Selection (TUI)
Select a profile from a list:
```bash
codex-mp ui
# or
codex-mp pick
```

### 5. Manage Profiles
List, delete, or rename profiles:
```bash
codex-mp list
codex-mp delete old-work
codex-mp rename personal home
```

### 6. Inspect
Check current auth fingerprint or resolved paths:
```bash
codex-mp who
codex-mp path
```

### 7. Shell Completion
Generate completion script for your shell (bash, zsh, fish, powershell):
```bash
codex-mp completion zsh > /usr/local/share/zsh/site-functions/_codex-mp
```

## Release metadata

`VERSION` is the single source of truth for release version.

To update the Homebrew formula metadata from `VERSION`:

```bash
# Update local placeholder and (optionally) the live tap repo
./scripts/update_formula.sh <sha256> [path/to/homebrew-tap]
```

## Development

```bash
# Build
make build

# Run integration test scripts against local binary
make test

# Run targeted unit tests
cd go && go test ./internal/app ./internal/profile
```

CI builds `codex-mp` first, runs smoke + battle tests with
`CODEX_MP=./codex-mp`, and runs shell linting on `bash/codex-switch`.
