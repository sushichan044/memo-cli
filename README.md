# memo-cli

A CLI tool to create and manage markdown memos with interactive selection and gitignore integration.

## Features

- 📝 Create markdown memos with custom names or timestamps
- 🔍 Interactive memo selection with fuzzy finder (fzf-like interface)
- 📂 Organized by date (YYYYMMDD directories)
- ⚠️  Gitignore checking with helpful warnings
- 🎨 Customizable via environment variables
- 🚀 Cross-platform support (Linux, macOS, Windows)

## Installation

### Homebrew (macOS/Linux)

```bash
brew install sushichan044/tap/memo-cli
```

### From Source

```bash
go install github.com/sushichan044/memo-cli/cmd/memo@latest
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/sushichan044/memo-cli/releases)

## Usage

### Create a new memo

```bash
# Create with timestamp (HH-MM-SS.md)
memo

# Create with custom name
memo "project-notes"
memo "meeting/2024"
```

### Select existing memo interactively

```bash
# Launch fuzzy finder to select a memo
memo list
```

### Environment Variables

#### `MEMO_BASE_DIR`

Customize the base directory for memos.

```bash
# Use custom directory
export MEMO_BASE_DIR="$HOME/Documents/memos"
memo

# Default: .{username}/memo in current directory
# If username cannot be determined: .memo/memo
```

### Directory Structure

```
{MEMO_BASE_DIR}/
└── YYYYMMDD/           # Date folder (e.g., 20251031)
    ├── HH-MM-SS.md     # Timestamp memo
    └── custom-name.md  # Named memo
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

# Create a named memo
$ memo "sprint-planning"
✅ Memo created at: /path/to/project/.sushichan044/memo/20251031/sprint-planning.md

# Pipe the path to open in editor
$ vim "$(memo 'quick-note')"

# Select from existing memos
$ memo list
# → Interactive fuzzy finder with preview
```

## Filename Normalization

When creating named memos, the following transformations are applied:

- Slashes (`/`) → Dashes (`-`)
- Spaces (` `) → Dashes (`-`)
- File extensions are removed (`.md` is automatically added)

Examples:
- `project/notes` → `project-notes.md`
- `my meeting notes` → `my-meeting-notes.md`
- `file.txt` → `file.md`

## Gitignore Integration

The tool checks if your memo directory is ignored by git and displays a warning if not.

To suppress the warning, add the memo directory to your `.gitignore`:

```gitignore
# For user-specific memos (default)
.sushichan044/memo/

# Or for custom MEMO_BASE_DIR
/path/to/your/custom/memos/
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
