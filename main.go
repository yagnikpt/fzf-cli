package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/fzf-cli/cmd"
)

func get_log_path() string {
	homepath, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		dirPath := filepath.Join(homepath, "AppData", "Local", "fzf_cli")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				log.Fatalf("Failed to create directory %s: %v", dirPath, err)
			}
		}
		path := filepath.Join(dirPath, "debug.log")
		return path
	}
	dirPath := filepath.Join(homepath, ".local", "state", "fzf_cli")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dirPath, err)
		}
	}
	path := filepath.Join(dirPath, "debug.log")
	return path
}

func main() {
	f, err := tea.LogToFile(get_log_path(), "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m, err := cmd.InitializeModel()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	rawModel, err := p.Run()
	if err != nil {
		log.Fatalln("Alas, there's been an error:", err)
	}

	final, ok := rawModel.(cmd.Model)
	if !ok {
		log.Fatalln("Alas, final model is not of type cmd.Model")
	}

	fmt.Println(final.Item)
	os.Exit(0)
}
