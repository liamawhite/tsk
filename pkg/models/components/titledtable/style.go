package titledtable

import "github.com/charmbracelet/lipgloss"

// Styles contains style definitions for this list component. By default, these
// values are generated by DefaultStyles.
type Styles[T any] struct {
	Title   lipgloss.Style
}

// DefaultStyles returns a set of default style definitions for this table.
func DefaultStyles[T any]() Styles[T] {
	return Styles[T]{
	    Title: lipgloss.NewStyle(),
    }
}

