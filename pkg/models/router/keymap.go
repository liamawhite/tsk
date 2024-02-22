package router

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit        key.Binding
}

var keyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
