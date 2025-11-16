package memo_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sushichan044/memo-cli/internal/config"
	"github.com/sushichan044/memo-cli/internal/memo"
)

func TestNormalizeFileName(t *testing.T) {
	// normalizeFileName is unexported, so we test it indirectly through Create()
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "with slashes",
			input:    "project/notes",
			contains: "projectnotes",
		},
		{
			name:     "with backslashes",
			input:    "project\\notes",
			contains: "projectnotes",
		},
		{
			name:     "with spaces",
			input:    "my project notes",
			contains: "my-project-notes",
		},
		{
			name:     "with extension",
			input:    "notes.md",
			contains: "notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			cfg := &config.Config{BaseDir: tmpDir}
			creator := memo.New(cfg)

			path, err := creator.Create(tt.input)
			if err != nil {
				t.Fatalf("Create() failed: %v", err)
			}

			filename := filepath.Base(path)
			if !strings.Contains(filename, tt.contains) {
				t.Errorf("filename %q should contain %q", filename, tt.contains)
			}
		})
	}
}

func TestCreate_WithName(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: tmpDir,
	}

	creator := memo.New(cfg)

	// Create memo with custom name
	path, err := creator.Create("test-memo")
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Check that file exists
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Errorf("File %q does not exist", path)
	}

	// Check filename format (should be HH-MM-name.md)
	filename := filepath.Base(path)
	if !strings.HasSuffix(filename, "-test-memo.md") {
		t.Errorf("Filename = %q, want format HH-MM-test-memo.md", filename)
	}

	// Check that timestamp prefix exists (HH-MM format)
	nameWithoutExt := strings.TrimSuffix(filename, ".md")
	parts := strings.SplitN(nameWithoutExt, "-", 3) // HH-MM-name
	const expectedMinParts = 3                      // At least HH-MM-name
	if len(parts) < expectedMinParts {
		t.Errorf("Filename %q should have format HH-MM-name", nameWithoutExt)
	}

	// Check that file is under tmpDir
	if !strings.HasPrefix(path, tmpDir) {
		t.Errorf("Path %q should be under %q", path, tmpDir)
	}

	// Check directory structure (should have YYYYMMDD directory)
	dir := filepath.Dir(path)
	dateDir := filepath.Base(dir)
	const expectedDirLen = 8 // YYYYMMDD
	if len(dateDir) != expectedDirLen {
		t.Errorf("Date directory %q should be 8 characters (YYYYMMDD)", dateDir)
	}
}

func TestCreate_WithoutName(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: tmpDir,
	}

	creator := memo.New(cfg)

	// Create memo without name (should use timestamp)
	path, err := creator.Create("")
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Check that file exists
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Errorf("File %q does not exist", path)
	}

	// Check filename format (should be HH-MM.md)
	filename := filepath.Base(path)
	if !strings.HasSuffix(filename, ".md") {
		t.Errorf("Filename %q should end with .md", filename)
	}

	// Remove .md extension
	nameWithoutExt := strings.TrimSuffix(filename, ".md")

	// Should be in format HH-MM (5 characters with dash)
	parts := strings.Split(nameWithoutExt, "-")
	const expectedParts = 2 // HH-MM
	if len(parts) != expectedParts {
		t.Errorf("Timestamp filename %q should have format HH-MM", nameWithoutExt)
	}
}

func TestCreate_DirectoryCreation(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: filepath.Join(tmpDir, "nonexistent"),
	}

	creator := memo.New(cfg)

	// Create memo (should create directories)
	path, err := creator.Create("test")
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Check that file exists
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Errorf("File %q does not exist", path)
	}

	// Check that base directory was created
	if _, statErr := os.Stat(cfg.BaseDir); os.IsNotExist(statErr) {
		t.Errorf("Base directory %q was not created", cfg.BaseDir)
	}
}

func TestGenerateFilename(t *testing.T) {
	// generateFilename is unexported, so we test it indirectly through Create()
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "with name",
			input: "my-memo",
		},
		{
			name:  "without name",
			input: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			cfg := &config.Config{BaseDir: tmpDir}
			creator := memo.New(cfg)

			path, err := creator.Create(tt.input)
			if err != nil {
				t.Fatalf("Create() failed: %v", err)
			}

			filename := filepath.Base(path)

			if filename == "" {
				t.Error("filename is empty")
			}

			// Should have .md extension
			if !strings.HasSuffix(filename, ".md") {
				t.Errorf("filename %q should have .md extension", filename)
			}
		})
	}
}

func TestCheckGitignore(t *testing.T) {
	// This is a basic test - gitignore checking is difficult to test
	// without setting up a real git repository
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: tmpDir,
	}

	creator := memo.New(cfg)

	// Should not panic
	warning := creator.CheckGitignore()

	// Warning can be empty or non-empty depending on git setup
	// Just check that it doesn't panic
	t.Logf("CheckGitignore() returned: %q", warning)
}
