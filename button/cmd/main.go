package main

import (
	"fmt"
	"os"
	"time"

	"github.com/azr4e1/polacco/button"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	model := button.New("ok", 0, button.SetDelay(2*time.Second))
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
