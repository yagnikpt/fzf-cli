package files

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ignoredPaths = []string{
	"node_modules",
	".git",
	"dist",
	"build",
	".next",
	"target",
	"bin",
	"obj",
	"vendor",
	".idea",
	".vscode",
	"__pycache__",
}

func shouldSkip(path string) bool {
	for _, ignored := range ignoredPaths {
		if strings.Contains(path, ignored) {
			return true
		}
	}
	return false
}

func GetAllFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldSkip(path) || info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		files = append(files, relPath)
		return nil
	})

	sort.Strings(files)
	return files, err
}
