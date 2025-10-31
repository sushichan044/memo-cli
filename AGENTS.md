# AGENTS.md

This file provides guidance to Coding Agents when working with code in this repository.

---

## Project Context

### Overview

mdfm is a Go CLI tool / Library that finds Markdown files using glob patterns and extracts their frontmatter metadata while respecting Git ignore rules.

### Architecture

#### Core Structure

<!-- placeholder -->

### Key Features

<!-- placeholder -->

### Design Decisions

<!-- placeholder -->

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
mise run dev                # Run in development mode (e.g., mise run dev "**/*.md")
mise run test               # Run tests using gotestsum
mise run test-coverage      # Run tests with coverage reporting
mise run lint               # Run golangci-lint for code quality checks
mise run lint-fix           # Auto-fix linting issues
mise run fmt                # Format code
mise run build-snapshot     # Build cross-platform binaries with goreleaser
mise run clean              # Remove generated files

# Standard Go commands
go run ./cmd/cli "**/*.md"  # Output results as JSON
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
