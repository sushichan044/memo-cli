package config

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

const (
	// fallbackMemoDir is the default directory name when username cannot be determined.
	fallbackMemoDir = "memo"

	memoRootDirEnv = "MEMO_ROOT_DIR"
)

// Config holds the configuration for the memo CLI.
type Config struct {
	// BaseDir is the base directory where memos are stored.
	BaseDir string
}

// New creates a new Config instance.
// It checks the MEMO_ROOT_DIR environment variable for custom base directory.
// If not set, it uses the current directory with .{username}/memo structure.
// If username cannot be determined, it falls back to .memo/memo.
func New() (*Config, error) {
	baseDir, err := getBaseDir()
	if err != nil {
		return nil, err
	}

	cleanBaseDir := filepath.Clean(baseDir)
	return &Config{BaseDir: cleanBaseDir}, nil
}

func getBaseDir() (string, error) {
	if envDir := os.Getenv(memoRootDirEnv); envDir != "" {
		if !filepath.IsAbs(envDir) {
			return "", errors.New(memoRootDirEnv + " must be an absolute path")
		}
		return envDir, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	username := getUsernameOrDefault()
	return filepath.Join(cwd, "."+username, "memo"), nil
}

// getUsernameOrDefault returns the current user's username.
// If it cannot be determined, it returns the default directory name without the leading dot.
func getUsernameOrDefault() string {
	currentUser, err := user.Current()
	if err != nil || currentUser.Username == "" {
		return fallbackMemoDir
	}
	return currentUser.Username
}

// GetIgnorePattern returns the gitignore pattern for the memo directory.
// This is relative to the current working directory.
func (c *Config) GetIgnorePattern() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(cwd, c.BaseDir)
	if err != nil {
		// If we can't get a relative path, use the base directory as-is
		relPath = c.BaseDir
	}

	return filepath.ToSlash(relPath) + "/", nil
}
