package readline

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const TAB = "\t"

type ReadlineMsg string

type Model struct {
	Prompt         string
	MaxHistorySize int
	CursorStyle    lipgloss.Style
	TextStyle      lipgloss.Style
	PromptStyle    lipgloss.Style

	currentPrompt       string
	history             []string
	historyPointer      int
	historyPromptCached string
	cursor              int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	output := m.PromptStyle.Render(m.Prompt) + m.TextStyle.Render(m.currentPrompt)

	return output
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
			m.decreaseHistoryPointer(1)
			m.setHistoryPrompt()

		case key.Matches(msg, DefaultKeyMap.Down):
			m.increaseHistoryPointer(1)
			m.setHistoryPrompt()

		case key.Matches(msg, DefaultKeyMap.Start):
			m.cursor = 0

		case key.Matches(msg, DefaultKeyMap.End):
			m.cursor = len(m.currentPrompt)

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

		case msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace:
			m.updateCurrentInput(msg.String())
		}
	}
	return m, cmd
}

func (m *Model) increaseCursor(n int) {
	if m.cursor+n >= len(m.currentPrompt) {
		m.cursor = len(m.currentPrompt)
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

	prev := m.currentPrompt[:m.cursor]
	next := m.currentPrompt[m.cursor:]
	m.currentPrompt = prev + msg + next
	m.increaseCursor(len(msg))
	cachedPrompt := m.currentPrompt
	m.historyPromptCached = cachedPrompt
}

func (m *Model) tab() {
	m.currentPrompt += TAB
	m.increaseCursor(1)
}

func (m *Model) backspace() {
	prev := m.currentPrompt[:m.cursor]
	next := m.currentPrompt[m.cursor:]
	if len(prev) > 0 {
		m.currentPrompt = prev[:len(prev)-1] + next
		m.decreaseCursor(1)
	}
}

func (m *Model) delete() {
	prev := m.currentPrompt[:m.cursor]
	next := m.currentPrompt[m.cursor:]
	if len(next) > 0 {
		m.currentPrompt = prev + next[1:]
	}
}

func (m *Model) enter() tea.Cmd {
	input := m.currentPrompt
	cmd := func() tea.Msg { return ReadlineMsg(input) }
	m.updateHistory(input)
	m.currentPrompt = ""
	m.historyPromptCached = ""
	m.cursor = 0

	return cmd
}

func (m *Model) decreaseHistoryPointer(n int) {
	if m.historyPointer-n <= 0 {
		m.historyPointer = 0
		return
	}
	m.historyPointer -= n
}

func (m *Model) increaseHistoryPointer(n int) {
	if m.historyPointer+n >= len(m.history) {
		m.historyPointer = len(m.history)
		return
	}
	m.historyPointer += n
}

func (m *Model) updateHistory(input string) {
	if m.MaxHistorySize <= 0 {
		m.history = []string{}
		return
	} else if len(m.history) >= m.MaxHistorySize {
		m.history = m.history[len(m.history)-m.MaxHistorySize+1:]
	}
	m.history = append(m.history, input)
	m.historyPointer = len(m.history)
}

func (m *Model) setHistoryPrompt() {
	if m.historyPointer >= len(m.history) {
		m.currentPrompt = m.historyPromptCached
		m.cursor = len(m.currentPrompt)
		return
	}
	m.currentPrompt = m.history[m.historyPointer]
	m.cursor = len(m.currentPrompt)
}

func (m *Model) esc() {
}

func (m *Model) truncate() {
	prev := m.currentPrompt[:m.cursor]
	m.currentPrompt = prev
}
