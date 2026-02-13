# Codex Multipass (codex-mp) v0.1.0

`codex-mp` is a local CLI for switching Codex accounts by swapping
`auth.json` profiles on disk.

## Quick Start

### 1. Install
Choose one:
- **Homebrew** (Recommended): `brew install BigCactusLabs/tap/codex-mp`
- **Build from Source**: Run `make build` (Requires [Go 1.23+](https://go.dev/doc/install) installed)

### 2. Startup
Initialize your profiles directory:
```bash
codex-mp init
```

### 3. Activate
Save your current login and start switching:
```bash
# Save current session
codex-mp save work

# Switch to profile
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

## Phase 1 usage

Initialize profiles directory:

```bash
codex-mp init
```

Save the current login as a profile:

```bash
codex login
codex-mp save work
```

Switch to a saved profile:

```bash
codex-mp use work
```

List saved profiles:

```bash
codex-mp list
```

Print current auth fingerprint:

```bash
codex-mp who
```

Show resolved paths:

```bash
codex-mp path
```

## Phase 2 usage (Management + UI)

Delete a profile:

```bash
codex-mp delete old-work
```

Rename a profile:

```bash
codex-mp rename personal home
```

Interactive selection (TUI):

```bash
codex-mp pick
# or
codex-mp ui
```

## Release metadata

`VERSION` is the single source of truth for release version.

To update the Homebrew formula metadata from `VERSION`:

```bash
# Update local placeholder and (optionally) the live tap repo
./scripts/update_formula.sh <sha256> [path/to/homebrew-tap]
```
