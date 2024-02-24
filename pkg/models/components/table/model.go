package table

import (
	"sort"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model defines a state for the table widget.
type Model[T any] struct {
	KeyMap KeyMap

    hideHeaders bool
	cols   []Column
	rows   []Row[T]
	cursor int
	focus  bool
	styles Styles[T]

	viewport viewport.Model
	start    int
	end      int

    sorter func(a, b T) bool
    sorted bool
}


// Column defines the table structure.
type Column struct {
	Title string
	Width int
}

// SetStyles sets the table styles.
func (m *Model[T]) SetStyles(s Styles[T]) {
	m.styles = s
	m.UpdateViewport()
}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))
type Option[T any] func(*Model[T])

// New creates a new model for the table widget.
func New[T any](opts ...Option[T]) Model[T] {
	m := Model[T]{
		cursor:   0,
		viewport: viewport.New(0, 20),

		KeyMap: DefaultKeyMap(),
		styles: DefaultStyles[T](),
	}

	for _, opt := range opts {
		opt(&m)
	}
    
    m.UpdateViewport()

	return m
}

// WithColumns sets the table columns (headers).
func WithColumns[T any](cols []Column) Option[T] {
	return func(m *Model[T]) {
		m.cols = cols
	}
}

func WithHiddenHeaders[T any](h bool) Option[T] {
    return func(m *Model[T]) {
        m.hideHeaders = h
    }
}

// WithRows sets the table rows (data).
func WithRows[T any](rows []Row[T]) Option[T] {
	return func(m *Model[T]) {
		m.rows = rows
	}
}

// WithHeight sets the height of the table.
func WithHeight[T any](h int) Option[T] {
	return func(m *Model[T]) {
		m.viewport.Height = h
	}
}

// WithWidth sets the width of the table.
func WithWidth[T any](w int) Option[T] {
	return func(m *Model[T]) {
		m.viewport.Width = w
	}
}

// WithFocused sets the focus state of the table.
func WithFocused[T any](f bool) Option[T] {
	return func(m *Model[T]) {
		m.focus = f
	}
}

// WithStyles sets the table styles.
func WithStyles[T any](s Styles[T]) Option[T]{
	return func(m *Model[T]) {
		m.styles = s
	}
}

// WithKeyMap sets the key map.
func WithKeyMap[T any](km KeyMap) Option[T] {
	return func(m *Model[T]) {
		m.KeyMap = km
	}
}

func WithSort[T any](sorter func(a, b T) bool) Option[T] {
    return func(m *Model[T]) {
        m.sorter = sorter
    }
}

// Update is the Bubble Tea update loop.
func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.PageUp):
			m.MoveUp(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		}
	}

	return m, nil
}

// Focused returns the focus state of the table.
func (m Model[T]) Focused() bool {
	return m.focus
}

// Focus focuses the table, allowing the user to move around the rows and
// interact.
func (m *Model[T]) Focus() {
	m.focus = true
	m.UpdateViewport()
}

// Blur blurs the table, preventing selection or movement.
func (m *Model[T]) Blur() {
	m.focus = false
	m.UpdateViewport()
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m Model[T]) SelectedRow() Row[T] {
	return m.rows[m.cursor]
}

// Rows returns the current rows.
func (m Model[T]) Rows() []Row[T] {
	return m.rows
}

// SetRows sets a new rows state.
func (m *Model[T]) SetRows(r []Row[T]) {
	m.rows = r
    m.sorted = false
	m.UpdateViewport()
}

// SetColumns sets a new columns state.
func (m *Model[T]) SetColumns(c []Column) {
	m.cols = c
	m.UpdateViewport()
}

// SetWidth sets the width of the viewport of the table.
func (m *Model[T]) SetWidth(w int) {
	m.viewport.Width = w
	m.UpdateViewport()
}

// SetHeight sets the height of the viewport of the table.
func (m *Model[T]) SetHeight(h int) {
	m.viewport.Height = h
	m.UpdateViewport()
}

// Height returns the viewport height of the table.
func (m Model[T]) Height() int {
	return m.viewport.Height
}

// Width returns the viewport width of the table.
func (m Model[T]) Width() int {
	return m.viewport.Width
}

// Cursor returns the index of the selected row.
func (m Model[T]) Cursor() int {
	return m.cursor
}

// SetCursor sets the cursor position in the table.
func (m *Model[T]) SetCursor(n int) {
	m.cursor = clamp(n, 0, len(m.rows)-1)
	m.UpdateViewport()
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *Model[T]) MoveUp(n int) {
	m.cursor = clamp(m.cursor-n, 0, len(m.rows)-1)
	switch {
	case m.start == 0:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset, 0, m.cursor))
	case m.start < m.viewport.Height:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset+n, 0, m.cursor))
	case m.viewport.YOffset >= 1:
		m.viewport.YOffset = clamp(m.viewport.YOffset+n, 1, m.viewport.Height)
	}
	m.UpdateViewport()
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *Model[T]) MoveDown(n int) {
	m.cursor = clamp(m.cursor+n, 0, len(m.rows)-1)
	m.UpdateViewport()

	switch {
	case m.end == len(m.rows):
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.viewport.Height))
	case m.cursor > (m.end-m.start)/2:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.cursor))
	case m.viewport.YOffset > 1:
	case m.cursor > m.viewport.YOffset+m.viewport.Height-1:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset+1, 0, 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *Model[T]) GotoTop() {
	m.MoveUp(m.cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Model[T]) GotoBottom() {
	m.MoveDown(len(m.rows))
}

// Sort sorts the table based on the current sorter function.
func (m *Model[T]) Sort() {
    if m.sorted || m.sorter == nil {
        return
    }
    sort.Slice(m.rows, func(i, j int) bool { return m.sorter(m.rows[i].Data, m.rows[j].Data) })
    m.sorted = true
}
