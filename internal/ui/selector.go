package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

const (
	previewMaxLines     = 20
	previewWindowOffset = 2 // Offset for preview window height calculation
)

// SelectMemo shows an interactive fuzzy finder for selecting a memo file.
// Returns the absolute path to the selected file.
func SelectMemo(baseDir string) (string, error) {
	// Find all memo files
	memoFiles, err := listMemoFiles(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to list memo files: %w", err)
	}

	if len(memoFiles) == 0 {
		return "", fmt.Errorf("no memo files found in %s", baseDir)
	}

	// Show fuzzy finder
	idx, err := fuzzyfinder.Find(
		memoFiles,
		func(i int) string {
			// Display relative path from baseDir
			rel, relErr := filepath.Rel(baseDir, memoFiles[i])
			if relErr != nil {
				return filepath.Base(memoFiles[i])
			}
			return rel
		},
		fuzzyfinder.WithPreviewWindow(func(i, _, h int) string {
			if i == -1 {
				return ""
			}
			return previewFile(memoFiles[i], h-previewWindowOffset)
		}),
		fuzzyfinder.WithHeader(fmt.Sprintf("Select a memo file (%d total)", len(memoFiles))),
	)
	if err != nil {
		return "", fmt.Errorf("selection cancelled or failed: %w", err)
	}

	return memoFiles[idx], nil
}

// listMemoFiles recursively finds all .md files in the given directory.
func listMemoFiles(baseDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip directories that we can't read
			if os.IsPermission(err) {
				return nil
			}
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include .md files
		if filepath.Ext(path) == ".md" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// previewFile returns the first N lines of a file for preview.
// If the file cannot be read, returns an error message.
// The maxLines parameter limits the number of lines shown, defaulting to previewMaxLines if <= 0.
func previewFile(path string, maxLines int) string {
	if maxLines <= 0 {
		maxLines = previewMaxLines
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	lineCount := 0
	for scanner.Scan() && lineCount < maxLines {
		lines = append(lines, scanner.Text())
		lineCount++
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return fmt.Sprintf("Error reading file: %v", scanErr)
	}

	if lineCount >= maxLines {
		lines = append(lines, "", "... (more lines)")
	}

	return strings.Join(lines, "\n")
}
