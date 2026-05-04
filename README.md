<div align="center">

# Codex Multipass

**Switch Codex accounts in seconds. No logout required.**

[![Version](https://img.shields.io/badge/version-0.2.0-blue)]()
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)]()
[![License](https://img.shields.io/badge/license-MIT-green)]()

</div>

---

Juggling multiple Codex accounts? `codex-mp` lets you save, name, and hot-swap file-backed login sessions from the command line — like browser profiles, but for your terminal.

> **Heads up:** Codex supports multiple credential storage modes: `file`, `keyring`, `auto`, and `ephemeral`. `codex-mp` manages file-backed `auth.json` sessions only. Set `cli_auth_credentials_store = "file"` in `~/.codex/config.toml`, then sign in again before saving profiles. Keyring-backed profile storage is intentionally out of scope for this release.

## Quick Start

**Install:**

```bash
brew install BigCactusLabs/tap/codex-mp
```

Or build from source with [Go 1.25+](https://go.dev/doc/install): `make build`

**Use it:**

```bash
codex-mp init          # one-time setup
codex-mp save work     # snapshot current session as "work"
codex-mp save personal # log into another account, save it too
codex-mp use work      # swap back instantly
codex-mp doctor        # verify Codex is using file-backed auth
codex-mp ui            # or pick from an interactive list
```

That's it. No repos touched, no tools reconfigured — just swaps files under your Codex state directory.

## Commands

| Command | What it does |
|---------|-------------|
| `codex-mp init` | Set up the profiles directory |
| `codex-mp save <name>` | Snapshot current session as a named profile |
| `codex-mp use <name>` | Switch to a saved profile |
| `codex-mp ui` / `pick` | Interactive profile selector (TUI) |
| `codex-mp list` | Show all saved profiles |
| `codex-mp who` | Print SHA-256 fingerprint of current auth |
| `codex-mp rename <old> <new>` | Rename a profile |
| `codex-mp delete <name>` | Delete a profile |
| `codex-mp path` | Show resolved Codex state paths |
| `codex-mp doctor` | Check Codex auth/profile compatibility |
| `codex-mp completion <shell>` | Generate shell completions (bash/zsh/fish/powershell) |
| `codex-mp --json <cmd>` | JSON output for scripting |

## Security

Your tokens never leave your machine and never get printed to the terminal.

- **Local only** — all data stays on disk, nothing phones home
- **Atomic writes** — temp file + rename, so a crash can't corrupt your auth
- **Process locking** — no concurrent-write races
- **Strict permissions** — directories `700`, auth files `600`, fails closed if hardening fails
- **Fingerprints, not secrets** — `codex-mp who` shows a SHA-256 hash, never the token itself

## Paths

By default Codex state lives at `~/.codex`. Override with `CODEX_HOME`:

```bash
CODEX_HOME=/custom/path/.codex codex-mp path
```

```
~/.codex/
  auth.json          # active session
  config.toml        # set cli_auth_credentials_store = "file"
  profiles/
    work.json        # saved profiles
    personal.json
```

## How It Works

`codex-mp` tracks which profile is active and syncs the latest `auth.json` back before switching. This means rotated refresh tokens are preserved automatically — no stale-token surprises.

If Codex is configured with `keyring`, `auto`, or `ephemeral`, profile commands fail with remediation instead of guessing where credentials live. `auto` is rejected because Codex may use the OS keyring when available and fall back to `auth.json` only on failure; that ambiguity is not safe for a profile switcher.

## Installation

### Homebrew (recommended)

```bash
brew install BigCactusLabs/tap/codex-mp
```

### Maintainer Release Flow

Homebrew publishing is maintained locally, not through GitHub Actions.

For a release, bump `VERSION`, commit the release changes, then run `./scripts/publish_homebrew_release.sh`.
The script tags the release if needed, uploads a stable tarball asset, updates `BigCactusLabs/tap`, and validates the formula with `brew`.

### From Source

```bash
git clone https://github.com/BigCactusLabs/codex-multipass.git
cd codex-multipass
make build    # requires Go 1.25+
# binary lands at ./codex-mp
```

> **Legacy users:** If you have automation calling `codex-switch`, point it at `bash/codex-switch` — it delegates to the Go binary.

## Development

```bash
make build                                        # build the binary
make test                                         # integration tests against local binary
make unit-test                                    # Go unit tests
shellcheck bash/codex-switch scripts/*.sh tests/*.sh
```

CI runs Go unit tests, smoke, battle, concurrency, and corrupt-storage tests plus shell and Go linting on every push.

---

<div align="center">

Made by [Big Cactus Labs](https://github.com/BigCactusLabs) 🌵

</div>
