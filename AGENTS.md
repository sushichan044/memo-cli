# AGENTS.md

This file provides guidance to Coding Agents when working with code in this repository.

---

## Project Context

### Overview

memo-cli is a Go CLI tool for creating and managing markdown memos with interactive selection capabilities. It provides a simple way to create timestamped or named markdown files organized by date, with built-in gitignore checking to help keep memos out of version control.

### Architecture

#### Core Structure

```
memo-cli/
├── cmd/
│   └── memo/           # CLI entry point (kong-based)
│       └── main.go
├── internal/
│   ├── config/         # Configuration management
│   │   ├── config.go   # Environment variable handling, default paths
│   │   └── config_test.go
│   ├── memo/           # Memo creation logic
│   │   ├── memo.go     # File creation, normalization, gitignore checking
│   │   └── memo_test.go
│   ├── ui/             # Interactive selection
│   │   └── selector.go # Fuzzyfinder integration
│   └── gitignore/      # Gitignore pattern matching
│       ├── matcher.go  # Pattern matching engine
│       └── path.go     # Path resolution for gitignore files
```

### Key Features

1. **Memo Creation**: Create markdown files with timestamp or custom names
2. **Date Organization**: Automatic YYYYMMDD directory structure
3. **Filename Normalization**: Safe filename generation (slash/space → dash, extension removal)
4. **Interactive Selection**: Fuzzy finder (go-fuzzyfinder) with file preview
5. **Gitignore Integration**: Checks if memo directory is ignored, shows helpful warnings
6. **Environment Customization**: `MEMO_BASE_DIR` for custom base directory
7. **Cross-platform**: Supports Linux, macOS, and Windows

### Design Decisions

#### User-specific Directory Structure

- **Decision**: Use `.{username}/memo` as default path
- **Rationale**: Avoids hardcoding usernames, supports multiple users
- **Fallback**: `.memo/memo` when username cannot be determined

#### Gitignore Warning (Non-blocking)

- **Decision**: Warn but don't block when memo directory is not ignored
- **Rationale**: Better UX, lets users decide their own workflow
- **Implementation**: Uses existing `internal/gitignore.Matcher`

#### System Timezone

- **Decision**: Use system timezone for all timestamps
- **Rationale**: Simplicity, matches user's local time
- **Alternative Considered**: `MEMO_TIMEZONE` environment variable (deferred)

#### Output Format

- **stdout**: File path (pipeable to other commands)
- **stderr**: Human-readable messages and warnings
- **Rationale**: Follows Unix convention for tool composition

#### CLI Parser

- **Decision**: kong (github.com/alecthomas/kong)
- **Rationale**: Type-safe, minimal boilerplate, excellent help generation

#### Interactive Selection

- **Decision**: go-fuzzyfinder (github.com/ktr0731/go-fuzzyfinder)
- **Rationale**: Pure Go implementation, no external dependencies, preview window support

---

## Development Environment

### Runtime & Package Manager

- **Go Version**: 1.24.0 or later (toolchain: go1.25.1)
- **Task Runner**: mise

### Code Quality Tools

#### Linter

- **Tool**: golangci-lint v2.1.5
- **Configuration**: `.golangci.yml`
- **Settings**:
  - Timeout: 5 minutes
  - Maximum line length: 120 characters
  - 80+ linters enabled including security checks (gosec), code complexity analysis, and strict formatting rules

#### Formatter

- **Tools**: goimports, golines, gci (via golangci-lint)
- **Auto-fix**: Available via `mise run lint-fix`

### Testing

- **Test Runner**: gotestsum v1.12.2 (enhanced output wrapper around standard Go testing)
- **Assertion Library**: github.com/stretchr/testify v1.11.1
- **Standard Commands**: `go test ./...` also supported

### Build

- **Build System**: goreleaser v2
- **Configuration**: `.goreleaser.yaml`
- **Supported Platforms**:
  - Linux (x86_64, i386, ARM variants)
  - Windows (via ZIP archives)
  - macOS (universal binaries with ARM64 support)
- **Packaging**: tar.gz for Unix, ZIP for Windows
- **Pre-build Hooks**: `go mod tidy`, `go generate ./...`

### Package Publishing

- **Release Platform**: GitHub Releases
- **Distribution**: Homebrew Tap (sushichan044/homebrew-tap)
- **Checksum**: Automatically generated for all artifacts
- **Version Management**: svu (semantic versioning utility)

### Available Scripts

```bash
mise run dev                # Run in development mode
mise run test               # Run tests using gotestsum
mise run test-coverage      # Run tests with coverage reporting
mise run lint               # Run golangci-lint for code quality checks
mise run lint-fix           # Auto-fix linting issues
mise run fmt                # Format code
mise run build-snapshot     # Build cross-platform binaries with goreleaser
mise run clean              # Remove generated files

# Standard Go commands
go run ./cmd/memo           # Run CLI in development
go run ./cmd/memo "note"    # Create a memo named "note"
go run ./cmd/memo list      # Launch interactive selector
go test ./...               # Run all tests
go mod tidy                 # Clean up dependencies
```

### Development Workflow

1. **Code**: Write Go code following golangci-lint rules (120 char line limit)
2. **Test**: Run `mise run test` or `go test ./...`
3. **Lint**: Ensure code passes `mise run lint` (80+ linters)
4. **Format**: Auto-format with `mise run fmt` or `mise run lint-fix`
5. **Build**: Test cross-platform builds with `mise run build-snapshot`
6. **Release**: Managed via goreleaser with GitHub Releases and Homebrew Tap

### Key Constraints

- **Minimum Go Version**: 1.24.0 required
- **Gitignore Compliance**: All file operations respect Git ignore rules (global and local)
- **Error Isolation**: Individual file processing errors don't halt entire operations
- **Concurrency Safety**: All concurrent operations use semaphore-based control
- **Cross-platform**: Must support Linux, Windows, and macOS
