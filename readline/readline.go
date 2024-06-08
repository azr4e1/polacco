package readline

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const TAB = "\t"
const EmptyChar = " "

type ReadlineMsg string

type option func(*Model) error

type Model struct {
	Prompt         string
	MaxHistorySize int
	TextStyle      lipgloss.Style
	PromptStyle    lipgloss.Style

	currentPrompt       string
	history             []string
	historyPointer      int
	historyPromptCached string
	cursorPointer       int
	cursor              cursor.Model
}

func New(opts ...option) Model {
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color("ff0000"))
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	m := &Model{
		Prompt:         "> ",
		MaxHistorySize: 100,
		TextStyle:      textStyle,
		PromptStyle:    textStyle,
		cursor:         cursor.New(),
	}

	m.cursor.Style = cursorStyle
	m.cursor.TextStyle = textStyle
	m.cursor.SetChar(EmptyChar)
	m.cursor.Focus()

	for _, o := range opts {
		if err := o(m); err != nil {
			continue
		}
	}

	return *m
}

func SetPrompt(prompt string) option {
	return func(m *Model) error {
		m.Prompt = prompt
		return nil
	}
}

func SetMaxHistory(mh int) option {
	return func(m *Model) error {
		m.MaxHistorySize = mh
		return nil
	}
}

func SetCursorStyle(style lipgloss.Style) option {
	return func(m *Model) error {
		m.cursor.Style = style
		return nil
	}
}

func SetTextStyle(style lipgloss.Style) option {
	return func(m *Model) error {
		m.TextStyle = style
		m.cursor.TextStyle = style
		return nil
	}
}

func SetPromptStyle(style lipgloss.Style) option {
	return func(m *Model) error {
		m.PromptStyle = style
		return nil
	}
}

func (m Model) View() string {
	output := m.PromptStyle.Inline(true).Render(m.Prompt)
	switch m.cursorPointer {
	case len(m.currentPrompt):
		output += m.TextStyle.Inline(true).Render(m.currentPrompt) + m.cursor.View()

	case 0:
		output += m.cursor.View() + m.TextStyle.Inline(true).Render(m.currentPrompt[1:])

	case len(m.currentPrompt) - 1:
		output += m.TextStyle.Inline(true).Render(m.currentPrompt[:len(m.currentPrompt)-1]) + m.cursor.View()

	default:
		prev := m.currentPrompt[:m.cursorPointer]
		next := m.currentPrompt[m.cursorPointer+1:]
		output += m.TextStyle.Inline(true).Render(prev) + m.cursor.View() + m.TextStyle.Inline(true).Render(next)
	}

	return output
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	oldPos := m.cursorPointer //nolint
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
			m.decreaseCursor(len(m.currentPrompt))

		case key.Matches(msg, DefaultKeyMap.End):
			m.increaseCursor(len(m.currentPrompt))

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
			cmds = append(cmds, m.enter())

		case key.Matches(msg, DefaultKeyMap.Esc):
			m.esc()

		case key.Matches(msg, DefaultKeyMap.Truncate):
			m.deleteAfterCursor()

		case msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace:
			m.updateCurrentInput(msg.String())
		}
	}

	var cmd tea.Cmd
	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	if oldPos != m.cursorPointer && m.cursor.Mode() == cursor.CursorBlink {
		m.cursor.Blink = false
		cmds = append(cmds, m.cursor.BlinkCmd())
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) increaseCursor(n int) {
	if m.cursorPointer+n >= len(m.currentPrompt) {
		m.cursorPointer = len(m.currentPrompt)
		m.cursor.SetChar(EmptyChar)
		return
	}

	m.cursorPointer += n
	m.cursor.SetChar(string(m.currentPrompt[m.cursorPointer]))
}

func (m *Model) decreaseCursor(n int) {
	if m.cursorPointer-n <= 0 {
		m.cursorPointer = 0
		if len(m.currentPrompt) == 0 {
			m.cursor.SetChar(EmptyChar)
			return
		}
		m.cursor.SetChar(string(m.currentPrompt[m.cursorPointer]))
		return
	}

	m.cursorPointer -= n
	if len(m.currentPrompt) == m.cursorPointer {
		m.cursor.SetChar(EmptyChar)
		return
	}
	m.cursor.SetChar(string(m.currentPrompt[m.cursorPointer]))
}

func (m *Model) updateCurrentInput(msg string) {

	prev := m.currentPrompt[:m.cursorPointer]
	next := m.currentPrompt[m.cursorPointer:]
	m.currentPrompt = prev + msg + next
	m.increaseCursor(len(msg))
	m.cacheHistory()
}

func (m *Model) cacheHistory() {
	cachedPrompt := m.currentPrompt
	m.historyPromptCached = cachedPrompt
	m.historyPointer = len(m.history)
}

func (m *Model) tab() {
	m.currentPrompt += TAB
	m.increaseCursor(1)
}

func (m *Model) backspace() {
	prev := m.currentPrompt[:m.cursorPointer]
	next := m.currentPrompt[m.cursorPointer:]
	if len(prev) > 0 {
		m.currentPrompt = prev[:len(prev)-1] + next
		m.decreaseCursor(1)
	}
	m.cacheHistory()
}

func (m *Model) delete() {
	prev := m.currentPrompt[:m.cursorPointer]
	next := m.currentPrompt[m.cursorPointer:]
	if len(next) > 0 {
		m.currentPrompt = prev + next[1:]
		m.cacheHistory()

		if len(next[1:]) > 0 {
			m.cursor.SetChar(string(next[1]))
			return
		}
	}
	m.cursor.SetChar(EmptyChar)
}

func (m *Model) enter() tea.Cmd {
	input := m.currentPrompt
	cmd := func() tea.Msg { return ReadlineMsg(input) }
	m.updateHistory(input)
	m.currentPrompt = ""
	m.cursorPointer = 0

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
	if strings.TrimSpace(input) == "" {
		return
	}
	if m.MaxHistorySize <= 0 {
		m.history = []string{}
		return
	} else if len(m.history) >= m.MaxHistorySize {
		m.history = m.history[len(m.history)-m.MaxHistorySize+1:]
	}
	m.history = append(m.history, input)
	m.historyPointer = len(m.history)
	m.historyPromptCached = ""
}

func (m *Model) setHistoryPrompt() {
	if m.historyPointer >= len(m.history) {
		m.currentPrompt = m.historyPromptCached
		m.cursorPointer = len(m.currentPrompt)
		return
	}
	m.currentPrompt = m.history[m.historyPointer]
	m.cursorPointer = len(m.currentPrompt)
}

func (m *Model) esc() {
}

func (m *Model) deleteAfterCursor() {
	prev := m.currentPrompt[:m.cursorPointer]
	m.currentPrompt = prev
	m.cursor.SetChar(EmptyChar)
	m.cacheHistory()
}

func (m *Model) deleteBeforeCursor() {
	next := m.currentPrompt[m.cursorPointer:]
	m.currentPrompt = next
	m.cacheHistory()
	m.cursor.SetChar(EmptyChar)
	if len(next) > 0 {
		m.cursor.SetChar(string(next[0]))
	}
}

func (m *Model) Blink() tea.Cmd {
	return cursor.Blink
}
