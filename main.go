package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnik-patel-47/fzf-cli/cmd"
)

func main() {
	f, err := tea.LogToFile("C:\\Codes\\go\\fzf_cli\\debug.log", "debug")
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
	if _, err := p.Run(); err != nil {
		log.Println("Alas, there's been an error:", err)
		os.Exit(1)
	}
}
