package gitignore

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Songmu/gitconfig"

	"github.com/sushichan044/memo-cli/internal/xdg"
)

func getGlobalGitIgnorePath() (string, error) {
	val, err := gitconfig.Get("core.excludesFile")
	if err != nil && !gitconfig.IsNotFound(err) {
		return "", err
	}

	if val != "" {
		// TODO: resolve ~/ to home directory
		if val[:2] == "~/" {
			home, homeErr := os.UserHomeDir()
			if homeErr != nil {
				return "", homeErr
			}
			val = filepath.Join(home, val[2:])
		}
		return val, nil
	}

	return getDefaultExcludesFilePath()
}

func getDefaultExcludesFilePath() (string, error) {
	configHome, err := xdg.ConfigHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(configHome, "git", "ignore"), nil
}

func getLocalGitIgnorePath() (string, error) {
	gitDir := os.Getenv("GIT_DIR")
	if gitDir == "" {
		if wd, err := os.Getwd(); err == nil {
			gitDir = filepath.Join(wd, ".git")
		}
	}

	if gitDir == "" {
		return "", errors.New("cannot resolve git directory")
	}

	return filepath.Join(gitDir, "info", "exclude"), nil
}
