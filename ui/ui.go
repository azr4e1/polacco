package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/azr4e1/polacco/button"
	"github.com/azr4e1/polacco/readline"
	"github.com/azr4e1/polacco/rpn"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const TAB = "\t"
const BUTNWIDTH = 9
const BUTNHEIGHT = 3
const ROWLEN = 4
const TOTALWIDTH = (2 + BUTNWIDTH) * ROWLEN

// 4 for the readline, 3 for the stack
const TOTALHEIGHT = 4 + 3 + 4*(BUTNHEIGHT+2)

var Help = `pop:   pop last element from stack
list:  show stack
reset: reset stack
quit:  quit
`

var HelpStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#71797E"))

type model struct {
	rl            readline.Model
	stack         *rpn.RPNStack
	currentOutput string
	outputStyle   lipgloss.Style
	borderStyle   lipgloss.Style
	history       string
	quitting      bool
	buttons       []button.Model
	error         error
}

func (m model) Init() tea.Cmd {
	return m.rl.Blink()
}

func paddingRight(output, padChar string, padLength int) string {
	lines := strings.Split(output, "\n")
	result := []string{}
	for _, line := range lines {
		diff := padLength - len(line)
		if diff > 0 {
			line += strings.Repeat(padChar, diff)
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func (m model) View() string {
	if m.quitting {
		return m.rl.TextStyle.Render("Bye!\n")
	}
	if m.error != nil {
		return m.rl.TextStyle.Render(m.error.Error() + "\n")
	}
	// readline
	output := m.rl.View()
	resultOutput := m.currentOutput
	if len(resultOutput) > m.rl.Width+len(m.rl.Prompt)-2 {
		resultOutput = resultOutput[len(resultOutput)-m.rl.Width-len(m.rl.Prompt)+2:]
	}
	resultOutput = m.outputStyle.Render(resultOutput)

	output = lipgloss.JoinVertical(lipgloss.Left, output, resultOutput)
	output = m.borderStyle.Render(output)

	// keyboard
	rows := []string{}
	row := []string{}
	for _, b := range m.buttons {
		row = append(row, b.View())
		if len(row) == ROWLEN {
			rowString := lipgloss.JoinHorizontal(lipgloss.Center, row...)
			rows = append(rows, rowString)
			row = []string{}
		}
	}
	if len(row) != 0 {
		rowString := lipgloss.JoinHorizontal(lipgloss.Center, row...)
		rows = append(rows, rowString)
		row = []string{}
	}
	keyboard := lipgloss.JoinVertical(lipgloss.Center, rows...)

	output = lipgloss.JoinVertical(lipgloss.Left, output, keyboard)

	// stack
	stack := m.stack.GetValues()
	stackLength := 0
	stackEls := []string{}
	truncated := false
	for i := len(stack) - 1; i >= 0; i-- {
		stackEl := fmt.Sprintf("%.2f", stack[i])
		stackLength += len(stackEl) + 2
		if stackLength+5 > TOTALWIDTH {
			truncated = true
			break
		}
		stackEls = append([]string{m.borderStyle.Render(m.outputStyle.Render(stackEl))}, stackEls...)
	}
	if truncated {
		stackEls = append([]string{lipgloss.NewStyle().Height(3).AlignVertical(lipgloss.Center).Render("...")}, stackEls...)
	}
	if len(stackEls) == 0 {
		stackEls = []string{m.borderStyle.Render("EMPTY")}
	}
	stackString := lipgloss.JoinHorizontal(lipgloss.Center, stackEls...)

	output = lipgloss.JoinVertical(lipgloss.Left, output, stackString)

	// help
	output += fmt.Sprintf("\n%s", HelpStyle.Render(Help))

	return output
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width-1 < TOTALWIDTH {
			m.error = errors.New("Window width is too small.")
			return m, tea.Quit
		}
		totalHeight := TOTALHEIGHT + len(strings.Split(Help, "\n"))
		if msg.Height-1 < totalHeight {
			m.error = errors.New("Window height is too small.")
			return m, tea.Quit
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			m.quitting = true
			return m, tea.Quit
		}

	case readline.ReadlineMsg:
		m.actionParse(string(msg))
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.rl, cmd = m.rl.Update(msg)
	cmds = append(cmds, cmd)

	for i, btn := range m.buttons {
		btn, cmd = btn.Update(msg)
		m.buttons[i] = btn
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func initialModel() tea.Model {
	stack := rpn.NewStack()
	buttons := []button.Model{}
	for id, label := range []string{"7", "8", "9", "+", "4", "5", "6", "-", "1", "2", "3", "*", "0", ".", "^", "/"} {
		trigger := key.NewBinding(key.WithKeys(label))
		btn := button.New(label, id, trigger, button.SetWidth(BUTNWIDTH), button.SetHeight(BUTNHEIGHT))
		buttons = append(buttons, btn)
	}
	// 2 accounts for the border width
	rl := readline.New(readline.SetWidth(TOTALWIDTH - 2))
	return model{
		stack:       stack,
		rl:          rl,
		outputStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#870087")),
		borderStyle: lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#3C3C3C")),
		buttons:     buttons,
	}
}

func (m *model) actionParse(input string) {
	cleanExpr := strings.ToLower(strings.TrimSpace(input))
	switch cleanExpr {
	case "q", "qu", "qui", "quit":
		m.quitting = true
	case "l", "ls", "li", "lis", "list":
		m.currentOutput = fmt.Sprintf("%v", m.stack.GetValues())
	case "p", "po", "pop":
		val, err := m.stack.Pop()
		if err != nil {
			m.currentOutput = fmt.Sprint("error: ", err)
			return
		}
		m.currentOutput = fmt.Sprintf("%f", val)
	case "r", "re", "res", "rese", "reset":
		m.stack = rpn.NewStack()
		m.currentOutput = ""
	default:
		err := rpn.StringParser(m.stack, cleanExpr)
		if err != nil {
			m.currentOutput = fmt.Sprint("error: ", err)
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
