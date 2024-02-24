package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// View renders the component.
func (m Model[T]) View() string {
    if m.hideHeaders {
        return m.viewport.View()
    }
	return m.headersView() + "\n" + m.viewport.View()
}

// UpdateViewport updates the list content based on the previously defined
// columns and rows.
func (m *Model[T]) UpdateViewport() {
    m.Sort()
	renderedRows := make([]string, 0, len(m.rows))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.cursor >= 0 {
		m.start = clamp(m.cursor-m.viewport.Height, 0, m.cursor)
	} else {
		m.start = 0
	}
	m.end = clamp(m.cursor+m.viewport.Height, m.cursor, len(m.rows))
	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m Model[T]) headersView() string {
	var s = make([]string, 0, len(m.cols))
	for _, col := range m.cols {
		style := lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
		renderedCell := style.Render(runewidth.Truncate(col.Title, col.Width, "…"))
		s = append(s, m.styles.Header.Render(renderedCell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *Model[T]) renderRow(rowIndex int) string {
    s := make([]string, 0, len(m.cols))
    row := m.rows[rowIndex]
	for i, value := range row.Render() {
		style := lipgloss.NewStyle().Width(m.cols[i].Width).MaxWidth(m.cols[i].Width).Inline(true)
		renderedCell := m.styles.Cell(row.Data, m.cols[i].Title).Render(style.Render(runewidth.Truncate(value, m.cols[i].Width, "…")))
		s = append(s, renderedCell)
	}

	styled := lipgloss.JoinHorizontal(lipgloss.Left, s...)

	if rowIndex == m.cursor {
		return m.styles.Selected.Render(styled)
	}

	return styled
}
