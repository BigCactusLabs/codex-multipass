# Contributing

We welcome contributions! The project is written in **Go**.

## Development Setup

1.  **Install Go**: 1.23 or later.
2.  **Install Dependencies**: `make tidy`
3.  **Run Tests**: `make test`

## Project Structure

- `go/cmd/codex-mp`: Main entry point.

- `go/internal/app`: CLI commands and wiring.
- `go/internal/profile`: Core profile management logic.
- `go/internal/config`: Configuration handling.
- `go/internal/ui`: User interface components.
- `go/internal/fs`: Atomic file system operations.
- `bash/`: Original Bash implementation (legacy/reference).
- `tests/`: Integration tests (Bash scripts).

## Release Workflow

1.  Update `VERSION` file.
2.  Commit and tag.
3.  CI will build and release.
