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

var Help = `pop: pop last element from stack
list: show stack
reset: reset stack
quit: quit
`

var HelpStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#71797E"))

type model struct {
	rl            readline.Model
	stack         *rpn.RPNStack
	currentOutput string
	outputStyle   lipgloss.Style
	history       string
	quitting      bool
	help          bool
}

func (m model) Init() tea.Cmd {
	return m.rl.Blink()
}

func (m model) View() string {
	if m.quitting {
		return m.rl.TextStyle.Render("Bye!\n")
	}
	output := m.rl.View()
	output += m.outputStyle.Render(fmt.Sprintf("\n%s\n", m.currentOutput))

	if m.help {
		output += fmt.Sprintf("\n%s", HelpStyle.Render(Help))
	}

	return output
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}
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
		outputStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#85c1e9")),
	}
}

func (m *model) actionParse(input string) {
	cleanExpr := strings.ToLower(strings.TrimSpace(input))
	switch cleanExpr {
	case "h", "he", "hel", "help":
		m.help = !m.help
	case "q", "qu", "qui", "quit":
		m.quitting = true
	case "l", "ls", "li", "lis", "list":
		m.currentOutput = fmt.Sprintf("%v", m.stack.GetValues())
	case "p", "po", "pop":
		val, err := m.stack.Pop()
		if err != nil {
			m.currentOutput = fmt.Sprint("error:", err)
			return
		}
		m.currentOutput = fmt.Sprintf("%f", val)
	case "r", "re", "res", "rese", "reset":
		m.stack = rpn.NewStack()
		m.currentOutput = ""
	default:
		err := rpn.StringParser(m.stack, cleanExpr)
		if err != nil {
			m.currentOutput = fmt.Sprint("error:", err)
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
