package files

import (
	"os"
	"path/filepath"
	"sort"
)

func GetAllFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's a directory
		if info.IsDir() {
			return nil
		}

		// Convert to relative path for files only
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
