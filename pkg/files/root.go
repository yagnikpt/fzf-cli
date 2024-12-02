package files

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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
	".astro",
	".cache",
	".vercel",
	".netlify",
	".github",
	".wrangler",
	".svelte-kit",
}

func shouldSkip(path string) bool {
	for _, ignored := range ignoredPaths {
		if strings.Compare(path, ignored) == 0 {
			return true
		}
	}
	return false
}

func traverseDir(dir string, fileChan chan<- string, wg *sync.WaitGroup, rootdir string) {
	defer wg.Done()

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", dir, err)
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if !shouldSkip(entry.Name()) {
				wg.Add(1)
				go traverseDir(fullPath, fileChan, wg, rootdir)
			}
		} else {
			path, err := filepath.Rel(rootdir, fullPath)
			if err != nil {
				log.Fatalf("Error getting relative path for %s: %v\n", fullPath, err)
				continue
			}
			fileChan <- path
		}
	}
}

func GetAllFiles(root string) ([]string, error) {
	var files []string
	fileChan := make(chan string)
	var wg sync.WaitGroup

	go func() {
		for file := range fileChan {
			files = append(files, file)
		}
	}()

	wg.Add(1)
	go traverseDir(root, fileChan, &wg, root)

	wg.Wait()
	close(fileChan)

	sort.Strings(files)
	return files, nil
}
