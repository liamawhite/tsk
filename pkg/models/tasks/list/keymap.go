package list

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
    Add    key.Binding
    Edit   key.Binding
    Delete key.Binding
}

var keyMap = KeyMap{ 
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
    Edit: key.NewBinding(
        key.WithKeys("e"),
        key.WithHelp("e", "edit"),
    ),
    Delete: key.NewBinding(
        key.WithKeys("d", "delete"),
        key.WithHelp("d/del", "delete"),
    ),
}
