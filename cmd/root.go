/*
Copyright © 2024 Yagnik Patel <pyagnik409@gmail.com>
*/
package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/yagnik-patel-47/fzf-cli/pkg/algo"
	"github.com/yagnik-patel-47/fzf-cli/pkg/files"
	ui_list "github.com/yagnik-patel-47/fzf-cli/ui/list"
)

type (
	errMsg error
)
type fileMsg []string

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var modeLabelStyle = lipgloss.NewStyle().Bold(true).Margin(0, 1, 0, 0).Padding(0, 1).Foreground(lipgloss.Color("#18181b")).Background(lipgloss.Color("#5eead4"))

type Model struct {
	fileChan    <-chan []string
	textInput   textinput.Model
	list        ui_list.Model
	err         errMsg
	mode        string
	view_height int
	view_width  int
	Item        string
}

func waitForFile(ch <-chan []string) tea.Cmd {
	return func() tea.Msg {
		file, ok := <-ch
		if !ok {
			return nil // Channel closed
		}
		return fileMsg(file)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		waitForFile(m.fileChan),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case fileMsg:
		m.list.ConstItems = make([]string, len(msg))
		m.list.Items = make([]string, len(msg))
		copy(m.list.ConstItems, msg)
		copy(m.list.Items, msg)
		cmds = append(cmds, waitForFile(m.fileChan))
		m.list.UpdatePagesCount(len(m.list.Items))
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			finalitem, err := m.list.GetItem()
			if err != nil {
				return m, nil
			}
			m.Item = finalitem
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.mode == "normal" {
				return m, tea.Quit
			} else {
				m.mode = "normal"
				m.list.SetMode("normal")
				m.textInput.Blur()
				return m, nil
			}
		}

		switch msg.String() {
		case "i":
			if m.mode != "insert" {
				m.mode = "insert"
				m.list.SetMode("insert")
				m.textInput.Cursor.Blink = true
				m.textInput.Focus()
				return m, nil
			}
		}

	case errMsg:
		m.err = msg
		return m, nil

	case tea.BlurMsg:
		m.textInput.Cursor.Blink = false
		log.Println("Blur event received")

	case tea.FocusMsg:
		m.textInput.Cursor.Blink = true
		log.Println("Focus event received")

	case tea.WindowSizeMsg:
		m.view_height = msg.Height
		m.view_width = msg.Width
		m.list.SetListHeight(msg.Height - 13)
	}

	m.list.FilterValue = m.textInput.Value()
	filteredValues := algo.FuzzyFind(m.textInput.Value(), m.list.ConstItems)
	m.list.Items = filteredValues
	m.list.UpdatePagesCount(len(filteredValues))

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	items := m.list.GetSlicedItems()

	var b strings.Builder

	b.WriteString(m.list.View() + "\n\n" + m.textInput.View() + "\n\n")
	b.WriteString(modeLabelStyle.Render(fmt.Sprint(strings.ToUpper(m.mode))))

	if m.mode == "insert" {
		b.WriteString(fmt.Sprintf("%s\n", map[bool]string{true: "esc: normal mode", false: "esc: normal mode"}[len(items) > 0]))
	} else {
		b.WriteString(fmt.Sprintf("%s\n", map[bool]string{true: "←h|j↓|k↑|l→ • i: insert mode • esc: quit", false: "esc: quit"}[len(items) > 0]))
	}

	return docStyle.Render(b.String())
}

func InitializeModel() (Model, error) {
	ti := textinput.New()
	ti.Placeholder = "Search your stuff..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	wd, err := os.Getwd()
	if err != nil {
		return Model{}, err
	}

	target_dir := flag.String("dir", wd, "directory to search in")
	flag.Parse()

	fileChan := files.GetAllFiles(*target_dir)

	const_items := make([]string, 0)
	if files, ok := <-fileChan; ok {
		const_items = append(const_items, files...)
	}
	file_items := make([]string, len(const_items))
	copy(file_items, const_items)

	m := Model{
		textInput: ti,
		list:      ui_list.NewList(file_items, const_items, *target_dir),
		err:       nil,
		mode:      "insert",
		fileChan:  fileChan,
		Item:      "",
	}

	return m, nil
}
