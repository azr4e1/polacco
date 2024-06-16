package main

import (
	"fmt"
	"os"
	"time"

	"github.com/azr4e1/polacco/button"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	button button.Model
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
	var cmd tea.Cmd
	m.button, cmd = m.button.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.button.View()
}

func main() {
	trigger := key.NewBinding(key.WithKeys("q"))
	buttonModel := button.New("ok", 0, trigger, button.SetDelay(5*time.Second), button.SetWidth(100), button.SetHeight(20))
	model := Model{button: buttonModel}
	model.button.Static = true
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
