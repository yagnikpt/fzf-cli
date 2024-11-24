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
	selected    int
	paginator   paginator.Model
}

var headerStyle = lipgloss.NewStyle().Padding(1, 1).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color(("#4f46e5")))
var itemStyle = lipgloss.NewStyle()
var activeItemStyle = itemStyle.Underline(true)

func NewList(items []string, const_items []string) Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.KeyMap = paginator.KeyMap{
		PrevPage: key.NewBinding(key.WithKeys("pgup", "left"), key.WithHelp("←", "page left")),
		NextPage: key.NewBinding(key.WithKeys("pgdown", "right"), key.WithHelp("→", "page right")),
	}
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(items))

	return Model{
		items:       items,
		const_items: const_items,
		cursor:      0,
		selected:    -1,
		paginator:   p,
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
	l.paginator.SetTotalPages(total)
}

func (l *Model) GetSlicedItems() []string {
	start, end := l.paginator.GetSliceBounds(len(l.items))
	return l.items[start:end]
}

func (l *Model) Init() tea.Cmd {
	return nil
}

func (l *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	start, end := l.paginator.GetSliceBounds(len(l.items))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if l.cursor > 0 {
				l.cursor--
			} else {
				if l.paginator.Page > 0 {
					l.paginator.PrevPage()
					l.cursor = l.paginator.PerPage - 1
				}
			}
		case tea.KeyDown:
			if l.cursor < l.paginator.PerPage-1 {
				if l.cursor < len(l.items[start:end])-1 {
					l.cursor++
				}
			} else {
				l.paginator.NextPage()
				l.cursor = 0
			}
		case tea.KeyRight, tea.KeyLeft, tea.KeyPgUp, tea.KeyPgDown:
			if l.cursor > len(l.items[start:end])-1 || l.cursor < 0 {
				l.cursor = len(l.items[start:end]) - 1
			}
		case tea.KeyEnter:
			l.selected = l.cursor
			if err := exec.Command("cmd", "/c", "start", l.items[l.cursor]).Run(); err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}
		}
	}
	l.paginator, cmd = l.paginator.Update(msg)
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
				cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4e94")).Render(">")
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
