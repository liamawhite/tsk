package list

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding

	Todo      key.Binding
	Blocked   key.Binding
	Paused    key.Binding
	Active    key.Binding
	Complete  key.Binding
	Abandoned key.Binding
}

var keyMap = KeyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "delete"),
		key.WithHelp("d/del", "delete"),
	),
	Todo: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "mark todo"),
	),
	Blocked: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "mark blocked"),
	),
	Paused: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "mark paused"),
	),
	Active: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "mark active"),
	),
	Complete: key.NewBinding(
		key.WithKeys("c", "enter"),
		key.WithHelp("c/enter", "mark complete"),
	),
    Abandoned: key.NewBinding(
        key.WithKeys("x"),
        key.WithHelp("x", "mark abandoned"),
    ),
}
