# memo-cli

A CLI tool to create and manage markdown memos with interactive selection and gitignore integration.

## Features

- 📝 Create markdown memos with custom names or timestamps
- 📂 Organized by date (YYYYMMDD directories)
- ⚠️  Gitignore checking with helpful warnings

## Installation

```bash
go install github.com/sushichan044/memo-cli/cmd/memo@latest
```

## Usage

### Create a new memo

```bash
# Create with timestamp (HH-MM-SS.md)
memo

# Create with custom name
memo "project-notes"
memo "meeting/2024"
```

### Directory Structure

```
.{$USER}/
└── YYYYMMDD/                  # Date folder (e.g., 20251031)
    ├── HH-MM-SS.md            # Timestamp memo (no name provided)
    └── HH-MM-SS-custom-name.md  # Named memo (with timestamp prefix)
```

### Examples

```bash
# Create a memo for today
$ memo
⚠️  Warning: Memo directory is not in .gitignore
    Please add the following line to your .gitignore:
    .sushichan044/memo/
✅ Memo created at: /path/to/project/.sushichan044/memo/20251031/14-30-45.md
/path/to/project/.sushichan044/memo/20251031/14-30-45.md

# Create a named memo (timestamp is always added as prefix)
$ memo "sprint-planning"
✅ Memo created at: /path/to/project/.sushichan044/memo/20251031/14-30-45-sprint-planning.md

# Pipe the path to open in editor
$ vim "$(memo 'quick-note')"

# Select from existing memos (Coming Soon)
# $ memo list
# Interactive fuzzy finder with preview
```

## Configuration

### Custom Memo Base Directory

You can set a custom base directory for your memos by defining the `MEMO_ROOT_DIR` environment variable. Make sure to use an absolute path.

```bash
MEMO_ROOT_DIR="/path/to/your/custom/memos" memo new "custom-note"
# This will create the memo in /path/to/your/custom/memos/YYYYMMDD/
```

## Gitignore Integration

The tool checks if your memo directory is ignored by git and displays a warning if not.

To suppress the warning, add the memo directory to your `.gitignore`:

```gitignore
# For user-specific memos (default)
.sushichan044/memo/

# Or for custom MEMO_BASE_DIR (Planned Feature)
# /path/to/your/custom/memos/
```

## Development

### Prerequisites

- Go 1.24.0 or later
- [mise](https://mise.jdx.dev/) (optional, for task runner)

### Build

```bash
# Using mise
mise run build-snapshot

# Using go
go build -o memo ./cmd/memo
```

### Test

```bash
# Using mise
mise run test

# Using go
go test ./...
```

### Lint

```bash
# Using mise
mise run lint

# Or
golangci-lint run
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
