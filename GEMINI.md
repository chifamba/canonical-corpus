# GEMINI.md

## Project Overview
`canonical-corpus` is a Go-based tool designed to fetch, parse, normalize, and export biblical and related texts from public-domain sources (e.g., Wikisource, Sacred-Texts, Project Gutenberg) into a structured Markdown corpus. It supports multiple categories (canonical, extra-canonical, dead-sea-scrolls) and multiple languages/translations (NIV,ESV, SNA2002, SUB, WLC, SBLG, LXX,  WEB, ASV, HMT, etc.).

### Main Technologies
- **Go 1.26**
- **Cobra**: CLI framework for commands and flags.
- **Zap**: High-performance logging.
- **YAML**: Configuration management.
- **Custom Components**: Downloader (with retries), Rate Limiter, Parser, Normalizer, and Compiler.

### Architecture
- `cmd/corpus-builder/`: Entry point for the CLI application.
- `internal/compiler/`: The core engine that orchestrates the fetch-parse-normalize-write lifecycle.
- `internal/downloader/`: Handles robust HTTP requests with rate limiting and retries.
- `internal/parser/`: Extracts content from various formats (HTML, TXT).
- `internal/normalizer/`: Cleans text and deduplicates passages.
- `internal/progress/`: Manages state to allow resumable builds.
- `sources/`: The "source of truth" for book metadata and their remote URLs.
- `pkg/markdown/`: Utilities for generating consistent Markdown output.

## Building and Running

### Prerequisites
- Go 1.26 or later.

### Key Commands
- **Build the CLI**:
  ```bash
  go build -o corpus-builder ./cmd/corpus-builder
  ```
- **Run the full build**:
  ```bash
  go run cmd/corpus-builder/main.go build
  ```
  Options:
  - `-k, --category`: Filter by category (`canonical`, `extra-canonical`, `dead-sea-scrolls`).
  - `-l, --language`: Filter by BCP-47 language code (e.g., `en`, `el`, `he`, `sn`).
  - `-t, --translation`: Filter by translation ID (e.g., `kjv`, `web`, `asv`).
  - `-f, --force`: Re-compile books even if they are marked as complete.
- **Fetch raw sources only**:
  ```bash
  go run cmd/corpus-builder/main.go fetch
  ```
- **Verify corpus integrity**:
  ```bash
  go run cmd/corpus-builder/main.go verify
  ```
- **Export the corpus**:
  ```bash
  go run cmd/corpus-builder/main.go export --dest ./dist
  ```
- **Run tests**:
  ```bash
  go test ./...
  ```

## Development Conventions

### Project Structure
- Follows standard Go layout. logic specific to the application resides in `internal/`, while reusable utilities are in `pkg/`.
- Book metadata and source definitions must be added to the `sources/` package.

### Coding Style
- Use `zap` for logging. Inject loggers into components during initialization.
- Configuration is loaded from `configs/config.yaml` and should be passed down where needed.
- Ensure new parsers or normalizers are placed in their respective `internal/` packages.

### Testing
- Add tests for new packages.
- Use `go test ./...` to ensure no regressions in parsing or normalization logic.
- Progress state is tracked in the output directory (default `./corpus/.progress.json`).
