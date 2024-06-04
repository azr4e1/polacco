package readline

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const TAB = "\t"

type InputSentMsg string
type OutputSentMsg string

type REPLUnit struct {
	Input  string
	Output string
}

type Model struct {
	Prompt         string
	MaxHistorySize int
	CursorStyle    lipgloss.Style
	TextStyle      lipgloss.Style
	PromptStyle    lipgloss.Style

	currentInput   string
	history        []REPLUnit
	historyPointer int
	cursor         int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	output := m.PromptStyle.Render(m.Prompt) + m.TextStyle.Render(m.currentInput)

	return output
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
		case key.Matches(msg, DefaultKeyMap.Down):
		case key.Matches(msg, DefaultKeyMap.Start):
			m.cursor = 0
		case key.Matches(msg, DefaultKeyMap.End):
			m.cursor = len(m.currentInput)
		case key.Matches(msg, DefaultKeyMap.Left):
			m.decreaseCursor(1)

		case key.Matches(msg, DefaultKeyMap.Right):
			m.increaseCursor(1)

		case key.Matches(msg, DefaultKeyMap.Tab):
			m.tab()

		case key.Matches(msg, DefaultKeyMap.Delete):
			m.delete()
		case key.Matches(msg, DefaultKeyMap.Backspace):
			m.backspace()

		case key.Matches(msg, DefaultKeyMap.Enter):
			cmd = m.enter()

		case key.Matches(msg, DefaultKeyMap.Esc):
			m.esc()

		case key.Matches(msg, DefaultKeyMap.Truncate):
			m.truncate()

		default:
			m.updateCurrentInput(msg.String())
		}
	}
	return m, cmd
}

func (m *Model) increaseCursor(n int) {
	if m.cursor+n >= len(m.currentInput) {
		m.cursor = len(m.currentInput)
		return
	}

	m.cursor += n
}

func (m *Model) decreaseCursor(n int) {
	if m.cursor-n <= 0 {
		m.cursor = 0
		return
	}

	m.cursor -= n
}

func (m *Model) updateCurrentInput(msg string) {

	prev := m.currentInput[:m.cursor]
	next := m.currentInput[m.cursor:]
	m.currentInput = prev + msg + next
	m.increaseCursor(len(msg))
}

func (m *Model) tab() {
	m.currentInput += TAB
	m.increaseCursor(1)
}

func (m *Model) backspace() {
	prev := m.currentInput[:m.cursor]
	next := m.currentInput[m.cursor:]
	if len(prev) > 0 {
		m.currentInput = prev[:len(prev)-1] + next
		m.decreaseCursor(1)
	}
}

func (m *Model) delete() {
	prev := m.currentInput[:m.cursor]
	next := m.currentInput[m.cursor:]
	if len(next) > 0 {
		m.currentInput = prev + next[1:]
	}
}

func (m *Model) enter() tea.Cmd {
	input := m.currentInput
	cmd := func() tea.Msg { return InputSentMsg(input) }
	m.currentInput = ""
	m.cursor = 0

	return cmd
}

func (m *Model) esc() {
}

func (m *Model) truncate() {
	prev := m.currentInput[:m.cursor]
	m.currentInput = prev
}
