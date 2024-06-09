package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/azr4e1/polacco/readline"
	"github.com/azr4e1/polacco/rpn"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const TAB = "\t"

var Help = ""

type model struct {
	rl            readline.Model
	stack         *rpn.RPNStack
	currentOutput string
	outputStyle   lipgloss.Style
	help          string
	history       string
}

func (m model) Init() tea.Cmd {
	return m.rl.Blink()
}

func (m model) View() string {
	output := m.rl.View() + "\n"
	output += m.outputStyle.Inline(true).Render(m.currentOutput)

	return output
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		}

	case readline.ReadlineMsg:
		m.actionParse(string(msg))
	}

	var cmd tea.Cmd
	m.rl, cmd = m.rl.Update(msg)

	return m, cmd
}

func initialModel() tea.Model {
	stack := rpn.NewStack()
	rl := readline.New()
	return model{
		stack:       stack,
		rl:          rl,
		outputStyle: lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#71797E")),
	}
}

func (m *model) actionParse(input string) {
	cleanExpr := strings.ToLower(strings.TrimSpace(input))
	switch cleanExpr {
	// case "h", "he", "hel", "help":
	// 	s.Help()
	case "l", "ls", "li", "lis", "list":
		m.currentOutput = fmt.Sprintf("%v\n", m.stack.GetValues())
	case "p", "po", "pop":
		val, err := m.stack.Pop()
		if err != nil {
			m.currentOutput = fmt.Sprintln("error:", err)
			return
		}
		m.currentOutput = fmt.Sprintf("%f\n", val)
	case "r", "re", "res", "rese", "reset":
		m.stack = rpn.NewStack()
		m.currentOutput = ""
	default:
		err := rpn.StringParser(m.stack, cleanExpr)
		if err != nil {
			m.currentOutput = fmt.Sprintln("error:", err)
			return
		}
		m.currentOutput = ""
	}
}

func Main() int {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return 1
	}

	return 0
}
