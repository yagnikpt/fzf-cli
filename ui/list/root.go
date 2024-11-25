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
	items       []string
	const_items []string
	cursor      int
	paginator   paginator.Model
	mode        string
}

var headerStyle = lipgloss.NewStyle().Padding(1, 4).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color(("#4f46e5")))
var itemStyle = lipgloss.NewStyle()
var activeItemStyle = itemStyle.Background(lipgloss.Color("#f43f5e")).Foreground(lipgloss.Color("#fff")).Bold(true)
var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f43f5e")).Bold(true)

func NewList(items []string, const_items []string) Model {
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
		items:       items,
		const_items: const_items,
		cursor:      0,
		paginator:   p,
		mode:        "insert",
	}
}

func (l *Model) GetItemValues() []string {
	return l.items
}

func (l *Model) GetConstItemValues() []string {
	return l.const_items
}

func (l *Model) SetItemValues(values []string) {
	l.items = values
}

func (l *Model) UpdatePagesCount(total int) {
	if total == 0 || l.paginator.Page*l.paginator.PerPage >= total {
		l.paginator.Page = 0
	}
	l.paginator.SetTotalPages(total)
}

func (l *Model) GetSlicedItems() []string {
	if len(l.items) == 0 {
		return []string{}
	}
	start, end := l.paginator.GetSliceBounds(len(l.items))
	// Ensure bounds are valid
	if start >= len(l.items) {
		start = 0
		l.paginator.Page = 0
	}
	if end > len(l.items) {
		end = len(l.items)
	}
	return l.items[start:end]
}

func (l *Model) SetMode(mode string) {
	l.mode = mode
	if mode == "insert" {
		l.paginator.KeyMap.NextPage.SetEnabled(false)
		l.paginator.KeyMap.PrevPage.SetEnabled(false)
	} else {
		l.paginator.KeyMap.NextPage.SetEnabled(true)
		l.paginator.KeyMap.PrevPage.SetEnabled(true)
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

	if len(l.items) == 0 {
		l.cursor = 0
		l.paginator.Page = 0
		return *l, cmd
	}

	l.paginator, cmd = l.paginator.Update(msg)

	if l.paginator.Page*l.paginator.PerPage >= len(l.items) {
		l.paginator.Page = (len(l.items) - 1) / l.paginator.PerPage
	}

	start, end := l.paginator.GetSliceBounds(len(l.items))
	if start >= len(l.items) {
		start = 0
		l.paginator.Page = 0
	}
	if end > len(l.items) {
		end = len(l.items)
	}
	currentPageItems := l.items[start:end]

	if len(currentPageItems) == 0 {
		l.cursor = 0
	} else if l.cursor >= len(currentPageItems) {
		l.cursor = len(currentPageItems) - 1
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(currentPageItems) > 0 {
				if err := exec.Command("cmd", "/c", "start", currentPageItems[l.cursor]).Run(); err != nil {
					log.Fatal(err)
				}
			}
		}
		switch {
		case key.Matches(msg, listNavigation.Up):
			if l.cursor > 0 {
				l.cursor--
			} else {
				if l.paginator.Page > 0 {
					l.paginator.PrevPage()
					start, end = l.paginator.GetSliceBounds(len(l.items))
					l.cursor = len(l.items[start:end]) - 1
				}
			}
		case key.Matches(msg, listNavigation.Down):
			if l.cursor < len(currentPageItems)-1 {
				l.cursor++
			} else if l.paginator.Page < l.paginator.TotalPages-1 {
				l.paginator.NextPage()
				l.cursor = 0
			}
		}
	}
	return *l, cmd
}

func (l *Model) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("Fuzzy finder with Go :)") + "\n\n")
	start, end := l.paginator.GetSliceBounds(len(l.items))
	if len(l.items[start:end]) > 0 {
		for i, item := range l.items[start:end] {
			cursor := " "
			if l.cursor == i {
				cursor = cursorStyle.Render(">")
				b.WriteString(fmt.Sprintf("%s %s\n", cursor, activeItemStyle.Render(item)))
			} else {
				b.WriteString(fmt.Sprintf("%s %s\n", cursor, itemStyle.Render(item)))
			}
		}
		b.WriteString("\n" + "  " + l.paginator.View())
	} else {
		b.WriteString(fmt.Sprintf("%s\n", "No results found!"))
	}
	return b.String()
}
