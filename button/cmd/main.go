package main

import (
	"fmt"
	"os"
	_ "time"

	"github.com/azr4e1/polacco/button"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	button1 button.Model
	button2 button.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.button1, cmd = m.button1.Update(msg)
	cmds = append(cmds, cmd)

	m.button2, cmd = m.button2.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	button1 := m.button1.View()
	button2 := m.button2.View()

	return lipgloss.JoinVertical(lipgloss.Center, button1, button2)
}

func main() {
	trigger1 := key.NewBinding(key.WithKeys("a"))
	buttonModel1 := button.New("a", 0, trigger1, button.SetWidth(5), button.SetHeight(3))
	trigger2 := key.NewBinding(key.WithKeys("s"))
	buttonModel2 := button.New("s", 1, trigger2, button.SetWidth(5), button.SetHeight(3))
	model := Model{button1: buttonModel1, button2: buttonModel2}
	// model.button.Static = true
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
