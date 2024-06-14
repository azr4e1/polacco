package button

import (
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActivateMsg struct {
	button Model
}

type DeactivateMsg struct {
	button Model
}

type Model struct {
	Label         string
	Border        bool
	InactiveStyle lipgloss.Style
	ActiveStyle   lipgloss.Style
	BorderStyle   lipgloss.Style

	active       bool
	height       int
	width        int
	padding      string
	delay        time.Duration
	id           int
	errorMessage error
}

type option func(*Model) error

func SetHeight(height int) option {
	return func(m *Model) error {
		if height <= 0 {
			return errors.New("cannot set height to less than 0")
		}
		m.height = height
		return nil
	}
}

func SetWidth(width int) option {
	return func(m *Model) error {
		if width <= 0 {
			return errors.New("cannot set width to less than 0")
		}
		m.width = width
		return nil
	}
}

func SetDelay(delay time.Duration) option {
	return func(m *Model) error {
		if delay <= 0 {
			return errors.New("cannot set delay to less than 0")
		}
		m.delay = delay
		return nil
	}
}

func New(label string, id int, opts ...option) Model {
	m := &Model{
		Label:         label,
		Border:        true,
		InactiveStyle: lipgloss.NewStyle().UnsetBackground().UnsetForeground(),
		ActiveStyle:   lipgloss.NewStyle().Background(lipgloss.Color("blue")).Foreground(lipgloss.Color("black")),
		BorderStyle:   lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("grey")),

		height:  1,
		width:   len(label),
		padding: " ",
		delay:   100 * time.Millisecond,
		id:      id,
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
	if err := m.errorMessage; err != nil {
		return err.Error()
	}
	if m.active {
		return m.Label
	}
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ActivateMsg:
		if msg.button.id == m.id {
			cmd := m.activate()
			return m, cmd
		}
	case DeactivateMsg:
		if msg.button.id == m.id {
			m.deactivate()
		}
	case tea.WindowSizeMsg:
		if msg.Width < m.width {
			m.errorMessage = errors.New("Window width is too small.")
			return m, tea.Quit
		}
		if msg.Height < m.height {
			m.errorMessage = errors.New("Window height is too small.")
			return m, tea.Quit
		}
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC:
			return m, tea.Quit
		case msg.String() == "q":
			return m, SendActivateMsg(m)

		}
	}
	return m, nil
}

func (m *Model) activate() tea.Cmd {
	m.active = true
	return tea.Tick(m.delay, func(_ time.Time) tea.Msg { return DeactivateMsg{button: *m} })
}

func (m *Model) deactivate() {
	m.active = false
}

func SendActivateMsg(b Model) tea.Cmd {
	return func() tea.Msg {
		return ActivateMsg{button: b}
	}
}
