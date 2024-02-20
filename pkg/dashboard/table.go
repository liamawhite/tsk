package dashboard

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/liamawhite/tsk/pkg/task"
	"github.com/samber/lo"
)


func taskRow(t task.Task, _ int) table.Row {
	return table.Row{t.Name}
}

func newTable(tasks []task.Task) table.Model {
	return table.New(
		table.WithColumns([]table.Column{
			{Title: "Task", Width: 15},
		}),
		table.WithRows(lo.Map(tasks, taskRow)),
	)
}
