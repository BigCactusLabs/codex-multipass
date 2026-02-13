# Codex Switcher (Phase 1)

`codex-switch` is a local CLI for switching Codex accounts by swapping
`auth.json` profiles on disk.

It does not change your repos or tools. It only copies files under your Codex
state directory.

## Security model

- All data stays local.
- The CLI never prints token contents.
- `codex-switch who` prints only a SHA-256 fingerprint of `auth.json`.
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
CODEX_HOME=/custom/path/.codex codex-switch path
```

## Installation

Install via Homebrew:

```bash
brew install quinn/tap/codex-switch
```

Or manually:

1. Copy `cli/codex-switch` to a directory in your `PATH`.
2. Make it executable: `chmod +x codex-switch`.

## Commands

```bash
codex-switch init
codex-switch save <name>
codex-switch use <name>
codex-switch list
codex-switch who
codex-switch path
codex-switch delete <name>
codex-switch rename <old> <new>
codex-switch pick
codex-switch ui
codex-switch version
codex-switch help
```

Global output flags:

```bash
codex-switch --plain <command>
codex-switch --json <command>
```

## Phase 1 usage

Initialize profiles directory:

```bash
codex-switch init
```

Save the current login as a profile:

```bash
codex login
codex-switch save work
```

Switch to a saved profile:

```bash
codex-switch use work
```

List saved profiles:

```bash
codex-switch list
```

Print current auth fingerprint:

```bash
codex-switch who
```

Show resolved paths:

```bash
codex-switch path
```

## Phase 2 usage (Management + UI)

Delete a profile:

```bash
codex-switch delete old-work
```

Rename a profile:

```bash
codex-switch rename personal home
```

Interactive selection (TUI):

```bash
codex-switch pick
# or
codex-switch ui
```

## Release metadata

`VERSION` is the single source of truth for release version.

To update the Homebrew formula metadata from `VERSION`:

```bash
./scripts/update_formula.sh <sha256>
```
