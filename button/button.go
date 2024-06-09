package button

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Label         string
	Border        bool
	InactiveStyle lipgloss.Style
	ActiveStyle   lipgloss.Style
	BorderStyle   lipgloss.Style

	height  int
	width   int
	padding string
	delay   time.Duration
}

type option func(*Model) error

func New(label string, opts ...option) Model {
	m := &Model{
		Label:   label,
		height:  1,
		width:   len(label),
		padding: " ",
		Border:  false,
		delay:   100 * time.Millisecond,
	}

	for _, o := range opts {
		if err := o(m); err != nil {
			continue
		}
	}

	return *m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return ""
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}
