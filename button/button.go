package button

import (
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActivateMsg struct {
	button Model
	tag    int
}

type DeactivateMsg struct {
	button Model
	tag    int
}

type Model struct {
	Label         string
	Border        bool
	InactiveStyle lipgloss.Style
	ActiveStyle   lipgloss.Style
	BorderStyle   lipgloss.Style
	Static        bool
	Trigger       key.Binding

	active     bool
	height     int
	width      int
	delay      time.Duration
	id         int
	msgCounter int
	blank      bool
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

func New(label string, id int, trigger key.Binding, opts ...option) Model {
	m := &Model{
		Label:         label,
		Border:        true,
		InactiveStyle: lipgloss.NewStyle().UnsetBackground().UnsetForeground(),
		ActiveStyle:   lipgloss.NewStyle().Background(lipgloss.Color("#870087")).Foreground(lipgloss.Color("#000000")),
		BorderStyle:   lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#3C3C3C")),
		Static:        false,
		Trigger:       trigger,

		height: 1,
		width:  len(label),
		delay:  100 * time.Millisecond,
		id:     id,
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
	if m.blank {
		return ""
	}

	style := lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center).AlignVertical(lipgloss.Center)
	button := style.Render(m.Label)
	if m.active {
		button = m.ActiveStyle.Render(button)
	} else {
		button = m.InactiveStyle.Render(button)
	}
	if m.Border {
		button = m.BorderStyle.Render(button)
	}
	return button
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ActivateMsg:
		if msg.button.id == m.id {
			cmd := m.activate()
			return m, cmd
		}
	case DeactivateMsg:
		if msg.button.id == m.id && msg.tag == m.msgCounter {
			m.deactivate()
		}
	case tea.WindowSizeMsg:
		if msg.Width+1 < m.width {
			errorMessage := errors.New("Window width is too small.")
			m.blank = true
			return m, tea.Sequence(tea.Println(errorMessage), tea.Quit)
		}
		if msg.Height+1 < m.height {
			errorMessage := errors.New("Window height is too small.")
			m.blank = true
			return m, tea.Sequence(tea.Println(errorMessage), tea.Quit)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Trigger):
			if m.Static {
				m.active = !m.active
				return m, nil
			}
			return m, SendActivateMsg(&m)

		}
	}
	return m, nil
}

func (m *Model) activate() tea.Cmd {
	m.active = true
	if m.Static {
		return nil
	}
	return tea.Tick(m.delay, func(_ time.Time) tea.Msg { return DeactivateMsg{button: *m, tag: m.msgCounter} })
}

func (m *Model) deactivate() {
	m.active = false
}

func SendActivateMsg(b *Model) tea.Cmd {
	b.msgCounter++
	return func() tea.Msg {
		return ActivateMsg{button: *b, tag: b.msgCounter}
	}
}
