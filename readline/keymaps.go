package readline

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Start     key.Binding
	End       key.Binding
	Left      key.Binding
	Right     key.Binding
	Quit      key.Binding
	Enter     key.Binding
	Tab       key.Binding
	Delete    key.Binding
	Backspace key.Binding
	Esc       key.Binding
	Truncate  key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "pgup"),
		key.WithHelp("↑", "previous item"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "pgdown"),
		key.WithHelp("down", "next item"),
	),
	Start: key.NewBinding(
		key.WithKeys("ctrl+a", "home"),
		key.WithHelp("C-a", "go to start of line"),
	),
	End: key.NewBinding(
		key.WithKeys("ctrl+e", "end"),
		key.WithHelp("C-e", "go to end of line"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("left", "move to left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("right", "move to right"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("C-c", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Truncate: key.NewBinding(
		key.WithKeys("ctrl+k"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
	),
	Backspace: key.NewBinding(
		key.WithKeys("backspace"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
	),
}
