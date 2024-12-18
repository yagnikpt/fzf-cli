package ui_list

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Items       []string
	ConstItems  []string
	Cursor      int
	Paginator   paginator.Model
	Mode        string
	FilterValue string
	Dir         string
}

var headerStyle = lipgloss.NewStyle().Padding(0, 1).Bold(true).Italic(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color(("#4f46e5")))
var itemLengthStyle = lipgloss.NewStyle().Margin(1, 0).Foreground(lipgloss.Color("#5e5e5e"))
var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f43f5e")).Bold(true)
var searchFilterBg = lipgloss.NewStyle().Background(lipgloss.Color("#f43f5e"))
var listContainerStyle = lipgloss.NewStyle()

func NewList(items []string, const_items []string, dir string) Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.KeyMap = paginator.KeyMap{
		PrevPage: key.NewBinding(key.WithKeys("pgup", "left", "a"), key.WithHelp("←", "page left")),
		NextPage: key.NewBinding(key.WithKeys("pgdown", "right", "d"), key.WithHelp("→", "page right")),
	}
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(items))
	p.KeyMap.NextPage.SetEnabled(false)
	p.KeyMap.PrevPage.SetEnabled(false)

	return Model{
		Items:       items,
		ConstItems:  const_items,
		Cursor:      0,
		Paginator:   p,
		Mode:        "insert",
		FilterValue: "",
		Dir:         dir,
	}
}

func (l *Model) SetListHeight(height int) {
	l.Paginator.PerPage = height
	listContainerStyle = listContainerStyle.Height(height)
}

func (l *Model) UpdatePagesCount(total int) {
	if total == 0 || l.Paginator.Page*l.Paginator.PerPage >= total {
		l.Paginator.Page = 0
	}
	l.Paginator.SetTotalPages(total)
}

func (l *Model) GetSlicedItems() []string {
	if len(l.Items) == 0 {
		return []string{}
	}
	start, end := l.Paginator.GetSliceBounds(len(l.Items))
	// Ensure bounds are valid
	if start >= len(l.Items) {
		start = 0
		l.Paginator.Page = 0
	}
	if end > len(l.Items) {
		end = len(l.Items)
	}
	return l.Items[start:end]
}

func (l *Model) SetMode(mode string) {
	l.Mode = mode
	if mode == "insert" {
		l.Paginator.KeyMap.NextPage.SetEnabled(false)
		l.Paginator.KeyMap.PrevPage.SetEnabled(false)
	} else {
		l.Paginator.KeyMap.NextPage.SetEnabled(true)
		l.Paginator.KeyMap.PrevPage.SetEnabled(true)
	}
}

func (l *Model) Init() tea.Cmd {
	return nil
}

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

var listNavigation = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("w", "up"),
		key.WithHelp("↑/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("s", "down"),
		key.WithHelp("↓/s", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("a", "left", "pgup"),
		key.WithHelp("←/a/pgdn", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("d", "right", "pgdn"),
		key.WithHelp("→/d/pgup", "move right"),
	),
}

func (l *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if len(l.Items) == 0 {
		l.Cursor = 0
		l.Paginator.Page = 0
		return *l, cmd
	}

	l.Paginator, cmd = l.Paginator.Update(msg)

	if l.Paginator.Page*l.Paginator.PerPage >= len(l.Items) {
		l.Paginator.Page = (len(l.Items) - 1) / l.Paginator.PerPage
	}

	start, end := l.Paginator.GetSliceBounds(len(l.Items))
	if start >= len(l.Items) {
		start = 0
		l.Paginator.Page = 0
	}
	if end > len(l.Items) {
		end = len(l.Items)
	}
	currentPageItems := l.Items[start:end]

	if len(currentPageItems) == 0 {
		l.Cursor = 0
	} else if l.Cursor >= len(currentPageItems) {
		l.Cursor = len(currentPageItems) - 1
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(currentPageItems) > 0 {
				includesTrailingSlash := strings.HasSuffix(l.Dir, "/") || strings.HasSuffix(l.Dir, "\\")
				if err := exec.Command("cmd", "/c", "start", map[bool]string{true: l.Dir + currentPageItems[l.Cursor], false: l.Dir + "\\" + currentPageItems[l.Cursor]}[includesTrailingSlash]).Run(); err != nil {
					log.Print(err)
				}
			}
		}
		switch {
		case key.Matches(msg, listNavigation.Up):
			if l.Cursor > 0 {
				l.Cursor--
			} else {
				if l.Paginator.Page > 0 {
					l.Paginator.PrevPage()
					start, end = l.Paginator.GetSliceBounds(len(l.Items))
					l.Cursor = len(l.Items[start:end]) - 1
				}
			}
		case key.Matches(msg, listNavigation.Down):
			if l.Cursor < len(currentPageItems)-1 {
				l.Cursor++
			} else if l.Paginator.Page < l.Paginator.TotalPages-1 {
				l.Paginator.NextPage()
				l.Cursor = 0
			}
		}
	}
	return *l, cmd
}

func (l *Model) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("Fuzzy finder with Go :)") + "\n")
	b.WriteString(itemLengthStyle.Render(fmt.Sprintf("%d items", len(l.Items))) + "\n")
	start, end := l.Paginator.GetSliceBounds(len(l.Items))
	var liView strings.Builder
	if len(l.Items[start:end]) > 0 {
		for i, item := range l.Items[start:end] {
			cursor := " "
			displayItem := item
			// searchItem := strings.ToLower(item)
			// searchFilter := strings.ToLower(l.filter_value)
			// lastIndex := 0

			// for _, letter := range searchFilter {
			// 	if idx := strings.Index(searchItem[lastIndex:], string(letter)); idx != -1 {
			// 		actualIdx := lastIndex + idx
			// 		displayItem = displayItem[:actualIdx] + searchFilterBg.Render(string(displayItem[actualIdx])) + displayItem[actualIdx+1:]
			// 		lastIndex = actualIdx + 1
			// 	}
			// }

			if l.Cursor == i {
				cursor = cursorStyle.Render(">")
			}
			liView.WriteString(fmt.Sprintf("%s %s\n", cursor, displayItem))
		}
		b.WriteString(listContainerStyle.Render(liView.String()))
		b.WriteString("\n" + "  " + l.Paginator.View())
	} else {
		b.WriteString(fmt.Sprintf("%s\n", "No results found!"))
	}
	return b.String()
}
