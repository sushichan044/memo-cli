# AGENTS.md

This file provides guidance to Coding Agents when working with code in this repository.

---

## Quick Commands

```bash
mise run dev                # Run in development mode
mise run test               # Run tests using gotestsum
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

## Project Context

memo-cli is a Go CLI tool for creating and managing markdown memos.

## Sources of Truth

Keep this file light. For implementation details, refer to:

- Product and usage overview: `README.md`
- CLI entry point: `cmd/memo/main.go`
- Package layout and behavior: `internal/`
- Dependencies and versions: `go.mod`, `go.sum`
- Task runner and scripts: `mise.toml`
- Lint/format rules: `.golangci.yml`
- Release/build configuration: `.goreleaser.yaml`
