# Contributing to Codex Switcher

Thank you for your interest in contributing! We welcome bug reports, feature requests, and pull requests.

## Development Setup

1. **Fork and Clone** the repo.
2. **Install Dependencies**:
   - `bash` (3.2+ is fine)
   - `shellcheck` (for linting)
   - `jq` (optional, for JSON debugging)

## Running Tests

We have two levels of tests:

1. **Smoke Tests**: Quick check of core functionality.
   ```bash
   ./tests/smoke.sh
   # PASS output means basic commands work
   ```

2. **Battle Tests**: Edge cases, tricky filenames, permission errors, and concurrency.
   ```bash
   ./tests/battle.sh
   # PASS output means deeper logic is solid
   ```

## Code Style

- Write POSIX-compatible Bash where possible, but `bash` 3.2+ features (arrays, `[[ ]]`) are allowed.
- Use `snake_case` for variables and functions.
- All scripts must pass **ShellCheck**:
  ```bash
  shellcheck cli/codex-switch scripts/*.sh tests/*.sh
  ```

## Release Flow (Maintainers)

1. Update `VERSION` file with the new semantic version.
2. Update `CHANGELOG.md` with the new release details.
3. Commit and push the changes.
4. Create a new GitHub Release tag matching `VERSION` (e.g., `v0.1.1`).
5. Wait for the tarball to be available.
6. Run the update script to refresh the Homebrew formula:
   ```bash
   # Get the sha256 of the new tarball first
   shasum -a 256 codex-switch-0.1.1.tar.gz
   
   # Update the formula
   ./scripts/update_formula.sh <new_sha256>
   ```
7. Commit the updated formula and push.
