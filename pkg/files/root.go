package files

import (
	"log"
	"os"
	"path/filepath"
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
	".pnpm-store",
}

func shouldSkip(path string) bool {
	for _, ignored := range ignoredPaths {
		if strings.Compare(path, ignored) == 0 {
			return true
		}
	}
	return false
}

func traverseDir(dir string, fileChan chan<- string, errorChan chan<- error, wg *sync.WaitGroup, rootdir string) {
	defer wg.Done()

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Error reading directory %s: %v\n", dir, err)
		errorChan <- err
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if !shouldSkip(entry.Name()) {
				wg.Add(1)
				go traverseDir(fullPath, fileChan, errorChan, wg, rootdir)
			}
		} else {
			path, err := filepath.Rel(rootdir, fullPath)
			if err != nil {
				log.Printf("Error getting relative path for %s: %v\n", fullPath, err)
				continue
			}
			if path != "" {
				fileChan <- path
			}
		}
	}
}

func GetAllFiles(root string) <-chan string {
	fileChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		var wg sync.WaitGroup

		if _, err := os.Stat(root); os.IsNotExist(err) {
			log.Printf("Error occurred while traversing: %v\n", err)
		}

		go func() {
			for err := range errorChan {
				if err != nil {
					log.Printf("Error occurred while traversing: %v\n", err)
				}
			}
		}()

		wg.Add(1)
		go traverseDir(root, fileChan, errorChan, &wg, root)

		wg.Wait()
		close(fileChan)
		close(errorChan)
	}()

	return fileChan
}
