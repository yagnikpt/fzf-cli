package files

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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

var (
	maxWorkers = runtime.NumCPU() * 4
)

func GetAllFiles(root string) <-chan []string {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Fatalf("Error occurred while traversing: %v\n", err)
	}

	files := make([]string, 0, 1000)
	done := make(chan struct{})
	fileChan := make(chan []string)
	errorChan := make(chan error, 10)

	go func() {
		startTime := time.Now()

		var fileCount int
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				errorChan <- err
				return nil
			}

			if d.IsDir() {
				if shouldSkip(d.Name()) {
					return filepath.SkipDir
				}
				return nil
			}
			relPath, _ := filepath.Rel(root, path)
			files = append(files, relPath)
			fileCount++
			return nil
		})

		if err != nil {
			log.Printf("Walk error: %v\n", err)
		}

		done <- struct{}{}
		duration := time.Since(startTime)
		log.Printf("Processed %d files in %v (%.2f files/sec)\n",
			fileCount,
			duration,
			float64(fileCount)/duration.Seconds())
		close(done)
	}()

	go func() {
		sendTicker := time.NewTicker(5 * time.Millisecond)
		defer sendTicker.Stop()

		for {
			select {
			case <-done:
				if len(files) > 0 {
					fileChan <- files
				}
				close(fileChan)
				close(errorChan)
				return
			case <-sendTicker.C:
				if len(files) > 0 {
					select {
					case fileChan <- files:
						log.Printf("Sent %d files\n", len(files))
					default:
					}
				}
			}
		}
	}()

	go func() {
		for err := range errorChan {
			if err != nil {
				log.Printf("Error occurred while traversing: %v\n", err)
			}
		}
	}()

	return fileChan
}
