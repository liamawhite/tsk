package list

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
	"github.com/liamawhite/tsk/pkg/models/components/table"
	"github.com/liamawhite/tsk/pkg/models/components/titledtable"
	"github.com/liamawhite/tsk/pkg/task"
	"github.com/samber/lo"
)

var theme = catppuccin.Mocha

func (m Model) View() string {
	return m.table.View()
}

const (
	colStatus = "Status"
	colTask   = "Task"
)

func (m Model) buildTable() titledtable.Model[task.Task] {
	tbl := table.New[task.Task](
		table.WithColumns[task.Task]([]table.Column{
			{Title: colStatus, Width: 1},
			{Title: colTask, Width: 20},
		}),
		table.WithHiddenHeaders[task.Task](true),
		table.WithRows(lo.Map(m.tasks, func(t task.Task, _ int) table.Row[task.Task] {
			return table.Row[task.Task]{Id: t.Id, Data: t, Renderer: func(t task.Task) []string {
				return []string{statusIcon(t.Status), t.Name}
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
        table.WithWidth[task.Task](m.width),
        table.WithHeight[task.Task](m.height),
	)
	return titledtable.New[task.Task](tbl,
		titledtable.WithTitle[task.Task]("Tasks"),
		titledtable.WithTitleStyle[task.Task](titledtable.Styles[task.Task]{
			Title: lipgloss.NewStyle().
				Margin(0, 1, 1, 1).
                Padding(0, 1).
                Background(lipgloss.Color(theme.Blue().Hex)).
                Foreground(lipgloss.Color(theme.Crust().Hex)),
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
	case task.Active:
		return common.Foreground(lipgloss.Color(theme.Green().Hex))
	case task.Complete:
		local := common.Faint(true).Foreground(lipgloss.Color(theme.Subtext0().Hex))
		if column == colStatus {
			return local
		}
		return local.Strikethrough(true)
    case task.Abandoned:
        local := common.Foreground(lipgloss.Color(theme.Overlay1().Hex)).Faint(true)
        if column == colStatus {
            return local
        }
        return local.Strikethrough(true)
	default:
		return common.Foreground(lipgloss.Color(theme.Text().Hex))
	}
}

func statusIcon(s task.Status) string {
	switch s {
	case task.Todo:
		return "󰁍"
	case task.Blocked:
		return "󰹆"
	case task.Paused:
		return "󰏤"
	case task.Active:
		return "󰁔"
	case task.Complete:
		return "✓"
    case task.Abandoned:
        return "x"
	default:
		return ""
	}
}
