package dashboard

import (
	"github.com/charmbracelet/huh"
	"github.com/liamawhite/tsk/pkg/task"
)


func newEditor(task *task.Task) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("task").
				Key("task"),
		),
	)
}
