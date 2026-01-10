package memo_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

			path, err := creator.Create(tt.input, "")
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
	path, err := creator.Create("test-memo", "")
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Check that file exists
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Errorf("File %q does not exist", path)
	}

	// Check filename format (should be HH-MM-SS-test-memo.md)
	filename := filepath.Base(path)
	if !strings.HasSuffix(filename, "-test-memo.md") {
		t.Errorf("Filename %q should end with -test-memo.md", filename)
	}

	// Verify timestamp prefix format (HH-MM-SS-)
	nameWithoutSuffix := strings.TrimSuffix(filename, "-test-memo.md")
	parts := strings.Split(nameWithoutSuffix, "-")
	const expectedParts = 3 // HH-MM-SS
	if len(parts) != expectedParts {
		t.Errorf("Timestamp prefix %q should have format HH-MM-SS", nameWithoutSuffix)
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

func TestCreate_WithExtension(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: tmpDir,
	}

	creator := memo.New(cfg)

	path, err := creator.Create("test-memo", "txt")
	require.NoError(t, err, "Create() should succeed")
	require.FileExists(t, path, "File should be created")

	// Check extension
	filename := filepath.Base(path)
	assert.True(t, strings.HasSuffix(filename, ".txt"),
		"Filename %q should end with .txt", filename)
}

func TestCreate_InvalidExtension(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: tmpDir,
	}

	creator := memo.New(cfg)

	if _, err := creator.Create("test", "md/evil"); err == nil {
		t.Fatal("Create() should fail with invalid extension")
	}
}

func TestCreate_ExtensionNormalization(t *testing.T) {
	tests := []struct {
		name    string
		ext     string
		want    string
		wantErr bool
	}{
		{"default on empty", "", "md", false},
		{"strips leading dot", ".txt", "txt", false},
		{"lowercases", "MD", "md", false},
		{"trims whitespace", "  txt  ", "txt", false},
		{"allows dash", "my-ext", "my-ext", false},
		{"allows underscore", "my_ext", "my_ext", false},
		{"allows numbers", "ext2", "ext2", false},
		{"rejects slash", "md/evil", "", true},
		{"rejects backslash", "md\\evil", "", true},
		{"rejects semicolon", "md;sh", "", true},
		{"rejects null byte", "md\x00", "", true},
		{"rejects space", "md sh", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			cfg := &config.Config{
				BaseDir: tmpDir,
			}
			creator := memo.New(cfg)

			path, err := creator.Create("test", tt.ext)

			if tt.wantErr {
				require.Error(t, err, "Create() should fail with invalid extension")
				assert.Contains(t, err.Error(), "invalid extension",
					"Error should mention invalid extension")
			} else {
				require.NoError(t, err, "Create() should succeed")
				require.FileExists(t, path, "File should be created")

				filename := filepath.Base(path)
				expectedExt := "." + tt.want
				assert.True(t, strings.HasSuffix(filename, expectedExt),
					"Filename %q should end with %q", filename, expectedExt)
			}
		})
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
	path, err := creator.Create("", "")
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Check that file exists
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Errorf("File %q does not exist", path)
	}

	// Check filename format (should be HH-MM-SS.md)
	filename := filepath.Base(path)
	if !strings.HasSuffix(filename, ".md") {
		t.Errorf("Filename %q should end with .md", filename)
	}

	// Remove .md extension
	nameWithoutExt := strings.TrimSuffix(filename, ".md")

	// Should be in format HH-MM-SS (8 characters with dashes)
	parts := strings.Split(nameWithoutExt, "-")
	const expectedParts = 3 // HH-MM-SS
	if len(parts) != expectedParts {
		t.Errorf("Timestamp filename %q should have format HH-MM-SS", nameWithoutExt)
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
	path, err := creator.Create("test", "")
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

// verifyTimestampPrefix checks that the first three parts of dash-separated name are 2-digit timestamps.
func verifyTimestampPrefix(t *testing.T, parts []string) {
	t.Helper()

	if len(parts) < 3 {
		t.Error("not enough parts for timestamp prefix")

		return
	}

	for i := range 3 {
		if len(parts[i]) != 2 {
			t.Errorf("timestamp part %d %q should be 2 digits", i, parts[i])
		}
	}
}

func TestGenerateFilename(t *testing.T) {
	// generateFilename is unexported, so we test it indirectly through Create()
	tests := []struct {
		name           string
		input          string
		expectedSuffix string
		minParts       int // minimum number of dash-separated parts (HH-MM-SS = 3, HH-MM-SS-name = 4+)
	}{
		{
			name:           "with name",
			input:          "my-memo",
			expectedSuffix: "-my-memo.md",
			minParts:       4, // HH-MM-SS-my-memo
		},
		{
			name:           "without name",
			input:          "",
			expectedSuffix: ".md",
			minParts:       3, // HH-MM-SS
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			cfg := &config.Config{BaseDir: tmpDir}
			creator := memo.New(cfg)

			path, err := creator.Create(tt.input, "")
			if err != nil {
				t.Fatalf("Create() failed: %v", err)
			}

			filename := filepath.Base(path)

			if filename == "" {
				t.Error("filename is empty")
			}

			// Should have expected suffix
			if !strings.HasSuffix(filename, tt.expectedSuffix) {
				t.Errorf("filename %q should have suffix %q", filename, tt.expectedSuffix)
			}

			// Should have .md extension
			if !strings.HasSuffix(filename, ".md") {
				t.Errorf("filename %q should have .md extension", filename)
			}

			// Verify timestamp is always present
			nameWithoutExt := strings.TrimSuffix(filename, ".md")
			parts := strings.Split(nameWithoutExt, "-")
			if len(parts) < tt.minParts {
				t.Errorf(
					"filename %q should have at least %d dash-separated parts, got %d",
					nameWithoutExt,
					tt.minParts,
					len(parts),
				)
			}

			// First three parts should be timestamp (HH-MM-SS)
			verifyTimestampPrefix(t, parts)
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

func TestCreateDirectory_WithName(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: tmpDir,
	}

	creator := memo.New(cfg)

	// Create directory with custom name
	path, err := creator.CreateDirectory("test-directory")
	require.NoError(t, err, "CreateDirectory() should succeed")

	// Check that directory exists
	info, err := os.Stat(path)
	require.NoError(t, err, "Directory should exist")
	assert.True(t, info.IsDir(), "Path should be a directory")

	// Check directory name format (should be HH-MM-SS-test-directory)
	dirname := filepath.Base(path)
	assert.True(t, strings.HasSuffix(dirname, "-test-directory"),
		"Directory name %q should end with -test-directory", dirname)

	// Verify timestamp prefix format (HH-MM-SS-)
	nameWithoutSuffix := strings.TrimSuffix(dirname, "-test-directory")
	parts := strings.Split(nameWithoutSuffix, "-")
	const expectedParts = 3 // HH-MM-SS
	assert.Len(t, parts, expectedParts,
		"Timestamp prefix %q should have format HH-MM-SS", nameWithoutSuffix)

	// Check that directory is under tmpDir
	assert.True(t, strings.HasPrefix(path, tmpDir),
		"Path %q should be under %q", path, tmpDir)

	// Check directory structure (should have YYYYMMDD parent directory)
	parentDir := filepath.Dir(path)
	dateDir := filepath.Base(parentDir)
	const expectedDirLen = 8 // YYYYMMDD
	assert.Len(t, dateDir, expectedDirLen,
		"Date directory %q should be 8 characters (YYYYMMDD)", dateDir)
}

func TestCreateDirectory_WithoutName(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		BaseDir: tmpDir,
	}

	creator := memo.New(cfg)

	// Create directory without name (should use timestamp)
	path, err := creator.CreateDirectory("")
	require.NoError(t, err, "CreateDirectory() should succeed")

	// Check that directory exists
	info, err := os.Stat(path)
	require.NoError(t, err, "Directory should exist")
	assert.True(t, info.IsDir(), "Path should be a directory")

	// Check directory name format (should be HH-MM-SS)
	dirname := filepath.Base(path)

	// Should be in format HH-MM-SS (8 characters with dashes)
	parts := strings.Split(dirname, "-")
	const expectedParts = 3 // HH-MM-SS
	assert.Len(t, parts, expectedParts,
		"Timestamp directory name %q should have format HH-MM-SS", dirname)
}
