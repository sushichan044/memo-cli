package config_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sushichan044/memo-cli/internal/config"
)

func TestNew(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Should use current directory with .{username}/memo
	if cfg.BaseDir == "" {
		t.Error("BaseDir should not be empty")
	}

	// Should contain memo directory
	if !filepath.IsAbs(cfg.BaseDir) {
		t.Errorf("BaseDir should be absolute path, got %q", cfg.BaseDir)
	}
}

func TestGetIgnorePattern(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		wantErr bool
	}{
		{
			name:    "absolute path",
			baseDir: "/tmp/test/.memo/memo",
			wantErr: false,
		},
		{
			name:    "relative path",
			baseDir: ".memo/memo",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{BaseDir: tt.baseDir}
			pattern, err := cfg.GetIgnorePattern()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIgnorePattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if pattern == "" {
					t.Error("GetIgnorePattern() returned empty pattern")
				}
				// Should end with /
				if pattern[len(pattern)-1] != '/' {
					t.Errorf("GetIgnorePattern() = %q, should end with /", pattern)
				}
			}
		})
	}
}

func TestMemoRootDirEnv(t *testing.T) {
	const testDir = "/tmp/memo_test_dir"

	// Set environment variable
	t.Setenv("MEMO_ROOT_DIR", testDir)

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if cfg.BaseDir != testDir {
		t.Errorf("BaseDir = %q; want %q", cfg.BaseDir, testDir)
	}
}

func TestGetUsernameOrDefault(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// BaseDir should contain a username-based directory
	if cfg.BaseDir == "" {
		t.Error("BaseDir should not be empty")
	}

	// Should contain .{username} or .memo in the path
	if !strings.Contains(cfg.BaseDir, "/.") {
		t.Errorf("BaseDir %q should contain a hidden directory (starting with dot)", cfg.BaseDir)
	}
}

func TestNew_Integration(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Change to temporary directory
	t.Chdir(tmpDir)

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// BaseDir should be absolute
	if !filepath.IsAbs(cfg.BaseDir) {
		t.Errorf("BaseDir should be absolute, got %q", cfg.BaseDir)
	}

	// BaseDir should end with "memo"
	if filepath.Base(cfg.BaseDir) != "memo" {
		t.Errorf("BaseDir %q should end with 'memo', got %q", cfg.BaseDir, filepath.Base(cfg.BaseDir))
	}

	// BaseDir's parent should start with a dot (hidden directory like .sushichan044 or .memo)
	parentDir := filepath.Dir(cfg.BaseDir)
	parentName := filepath.Base(parentDir)
	if !strings.HasPrefix(parentName, ".") {
		t.Errorf("BaseDir's parent should be a hidden directory, got %q", parentName)
	}
}
