package list

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
	"github.com/liamawhite/tsk/pkg/models/components/table"
	"github.com/liamawhite/tsk/pkg/task"
	"github.com/samber/lo"
)

var theme = catppuccin.Mocha

func (m Model) View() string {
	return m.table.View()
}

const (
    colStatus = "Status"
    colTask = "Task"
)

func buildTable(tasks []task.Task) table.Model[task.Task] {
	return table.New[task.Task](
		table.WithColumns[task.Task]([]table.Column{
			{Title: colStatus, Width: 1},
			{Title: colTask, Width: 20},
		}),
		table.WithHiddenHeaders[task.Task](true),
		table.WithRows(lo.Map(tasks, func(t task.Task, _ int) table.Row[task.Task] {
			return table.Row[task.Task]{Id: t.Id, Data: t, Renderer: func(t task.Task) []string {
				return []string{statusEmoji(t.Status), t.Name}
			}}
		})),
		table.WithSort(func(a task.Task, b task.Task) bool {
			return a.Status < b.Status
		}),
		table.WithStyles(table.Styles[task.Task]{
			Header:   table.DefaultStyles[task.Task]().Header, // we don't show the header
			Cell:     conditionalFormat,
			Selected: table.DefaultStyles[task.Task]().Selected,
		}),
	)
}

func conditionalFormat(t task.Task, column string) lipgloss.Style {
    common := lipgloss.NewStyle().Padding(0, 1)
    if column == colStatus {
        common = common.Align(lipgloss.Center)
    }

    switch t.Status {
    case task.Blocked:
        return common.Foreground(lipgloss.Color(theme.Red().Hex))
    case task.Paused:
        return common.Foreground(lipgloss.Color(theme.Subtext1().Hex)).Faint(true)
    case task.InProgress:
        return common.Foreground(lipgloss.Color(theme.Green().Hex))
    case task.Done:
        local := common.Faint(true)
        if column == colStatus {
            return local.Foreground(lipgloss.Color(theme.Green().Hex))
        }
        return local.Foreground(lipgloss.Color(theme.Subtext0().Hex)).Strikethrough(true)
    default:
        return common.Foreground(lipgloss.Color(theme.Text().Hex))
}
}

func statusEmoji(s task.Status) string {
	switch s {
	case task.Backlog:
		return ""
	case task.Blocked:
		return "󰹆"
	case task.Paused:
		return "󰏤"
	case task.InProgress:
		return "󰁜"
	case task.Done:
		return "✓"
	default:
		return ""
	}
}
