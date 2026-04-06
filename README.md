<div align="center">

# Codex Multipass

**Switch Codex accounts in seconds. No logout required.**

[![Version](https://img.shields.io/badge/version-0.1.6-blue)]()
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)]()
[![License](https://img.shields.io/badge/license-MIT-green)]()

</div>

---

Juggling multiple Codex accounts? `codex-mp` lets you save, name, and hot-swap login sessions from the command line — like browser profiles, but for your terminal.

> **Heads up:** Codex can cache credentials in the OS keyring instead of `auth.json`. `codex-mp` manages file-backed sessions only. If `codex login status` shows you're logged in but `auth.json` is missing, set `cli_auth_credentials_store = "file"` in `~/.codex/config.toml` and sign in again.

## Quick Start

**Install:**

```bash
brew install BigCactusLabs/tap/codex-mp
```

Or build from source with [Go 1.23+](https://go.dev/doc/install): `make build`

**Use it:**

```bash
codex-mp init          # one-time setup
codex-mp save work     # snapshot current session as "work"
codex-mp save personal # log into another account, save it too
codex-mp use work      # swap back instantly
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
  profiles/
    work.json        # saved profiles
    personal.json
```

## How It Works

`codex-mp` tracks which profile is active and syncs the latest `auth.json` back before switching. This means rotated refresh tokens are preserved automatically — no stale-token surprises.

## Installation

### Homebrew (recommended)

```bash
brew install BigCactusLabs/tap/codex-mp
```

### From Source

```bash
git clone https://github.com/BigCactusLabs/codex-multipass.git
cd codex-multipass
make build    # requires Go 1.23+
# binary lands at ./codex-mp
```

> **Legacy users:** If you have automation calling `codex-switch`, point it at `bash/codex-switch` — it delegates to the Go binary.

## Development

```bash
make build                                        # build the binary
make test                                         # integration tests against local binary
cd go && go test ./internal/app ./internal/profile # unit tests
```

CI runs smoke, battle, concurrency, and corrupt-storage tests plus shell linting on every push.

---

<div align="center">

Made by [Big Cactus Labs](https://github.com/BigCactusLabs) 🌵

</div>
