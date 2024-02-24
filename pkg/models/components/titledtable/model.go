package titledtable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/tsk/pkg/models/components/table"
)

type Option[T any] func(*Model[T])

type Model[T any] struct {
    // Deliberately not embedding the table.Model because theres are no good ways
    // to avoid extend the options for the table.Model
    table table.Model[T]

    title string
}

func New[T any](table table.Model[T], opts ...Option[T]) Model[T] {
    m := Model[T]{table: table}
    for _, opt := range opts {
        opt(&m)
    }
    return m
}

func WithTitle[T any](title string) Option[T] {
    return func(m *Model[T]) {
        m.title = title
    }
}

func (m *Model[T]) Focus() {
    m.table.Focus()
}

func (m *Model[T]) Blur() {
    m.table.Blur()
}

func (m Model[T]) SelectedRow() table.Row[T] {
    return m.table.SelectedRow()
}

func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}


func (m Model[T]) View() string {
    if m.title == "" {
        return m.table.View()
    }
    return m.title + "\n" + m.table.View()
}
