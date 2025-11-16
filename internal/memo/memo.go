package memo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pathologize"

	"github.com/sushichan044/memo-cli/internal/config"
	"github.com/sushichan044/memo-cli/internal/gitignore"
)

// Creator handles memo creation logic.
type Creator struct {
	config *config.Config
}

// New creates a new Creator instance.
func New(cfg *config.Config) *Creator {
	return &Creator{config: cfg}
}

// Create creates a new memo file with the given name.
// If name is empty, uses timestamp (HH-MM-SS) as filename.
// Returns the absolute path to the created file.
func (c *Creator) Create(name string) (string, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(c.config.BaseDir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create base directory: %w", err)
	}

	// Generate filename
	filename := normalizeFileName(c.generateFilename(name))

	// Create date directory (YYYYMMDD)
	now := time.Now()
	dateDir := now.Format("20060102")
	fullDir := filepath.Join(c.config.BaseDir, dateDir)

	if err := os.MkdirAll(fullDir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create date directory: %w", err)
	}

	// Create file path
	filePath := filepath.Join(fullDir, filename+".md")

	// Create empty file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create memo file: %w", err)
	}
	defer file.Close()

	return filePath, nil
}

// generateFilename creates a normalized filename from user input.
// All filenames are prefixed with timestamp (HH-MM).
// If name is empty, uses only the timestamp.
// Otherwise, generates "HH-MM-name" format by:
// - Prefixing with timestamp.
// - Removing extension if provided.
// - Replacing slashes with dashes.
// - Replacing spaces with dashes.
func (c *Creator) generateFilename(name string) string {
	timestamp := time.Now().Format("15-04")
	if name == "" {
		// Use timestamp as default
		return timestamp
	}

	return timestamp + "-" + pathologize.Clean(name)
}

// CheckGitignore checks if the memo base directory is ignored by git.
// Returns a warning message if not ignored, empty string otherwise.
// Silently returns empty string if gitignore checking fails (e.g., not a git repository).
func (c *Creator) CheckGitignore() string {
	matcher, err := gitignore.NewFromCWD()
	if err != nil {
		// Not a git repository or error reading gitignore - skip check silently
		return ""
	}

	if matcher.IsIgnored(c.config.BaseDir) {
		// Directory is already ignored
		return ""
	}

	// Generate warning message
	pattern, err := c.config.GetIgnorePattern()
	if err != nil {
		// Can't generate pattern - skip warning
		return ""
	}

	return fmt.Sprintf(
		"⚠️  Warning: Memo directory is not in .gitignore\n"+
			"    Please add the following line to your .gitignore:\n"+
			"    %s",
		pattern,
	)
}

func normalizeFileName(name string) string {
	normalized := name

	rep := strings.NewReplacer(
		" ", "-",
		// add more replacements if needed
	)
	normalized = rep.Replace(normalized)

	return normalized
}
