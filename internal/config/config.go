package config

import (
	"os"
	"os/user"
	"path/filepath"
)

const (
	// EnvMemoBaseDir is the environment variable name for the memo base directory.
	EnvMemoBaseDir = "MEMO_BASE_DIR"

	// defaultMemoDirName is the default directory name when username cannot be determined.
	defaultMemoDirName = "memo"
)

// Config holds the configuration for the memo CLI.
type Config struct {
	// BaseDir is the base directory where memos are stored.
	BaseDir string
}

// New creates a new Config instance.
// It reads the MEMO_BASE_DIR environment variable.
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
	fromEnv := os.Getenv(EnvMemoBaseDir)
	if fromEnv != "" {
		return fromEnv, nil
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
		return defaultMemoDirName
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
