# Codex Switcher Engineer Handoff

## Overview
`codex-switch` is a Bash-based CLI tool designed to manage and switch between multiple Codex `auth.json` profiles. It operates by copying profile files to and from a central `auth.json` location.

## Architecture
- **Language**: Bash (compatible with macOS default Bash 3.2).
- **Storage**: Profiles are stored as JSON files in `$CODEX_HOME/profiles/`.
- **Active State**: The active profile is always copied to `$CODEX_HOME/auth.json`.
- **Interactive UI**: Uses `fzf` if available; otherwise falls back to a native Bash `select` menu.

## Security
- **Permissions**: 
  - `$CODEX_HOME` and `$CODEX_HOME/profiles/` are locked to `700`.
  - `auth.json` and all profile files are locked to `600`.
- **Information Disclosure**: The `who` command prints only a SHA-256 fingerprint, never the raw token.

## Commands
- `init`: Setup the profiles directory.
- `save <name>`: Backup current `auth.json`.
- `use <name>`: Restore a profile to `auth.json`.
- `list`: Show all saved profiles.
- `who`: Show current auth fingerprint.
- `path`: Show configuration paths.
- `delete <name>`: Remove a profile.
- `rename <old> <new>`: Rename a profile.
- `pick` / `ui`: Interactive selection.

## Release Flow
1. Update version in `codex-switch`.
2. Push a new tag to GitHub.
3. Update the URL and SHA256 in `homebrew/Formula/codex-switch.rb`.
4. Submit PR to the homebrew tap.
