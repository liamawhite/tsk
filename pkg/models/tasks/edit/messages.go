package edit

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/liamawhite/tsk/pkg/task"
)

func AddPopulator() tea.Cmd {
	return func() tea.Msg {
		slog.Info("editing a new task", "model", name)
		return task.GetMsg{Task: task.Task{Id: uuid.New().String()}, Error: nil}
	}
}

func EditPopulator(client *task.Client, id string) tea.Cmd {
	slog.Info("editing an existing task", "model", name, "id", id)
	return client.Get(id)
}

type CancelMsg struct {
	Error error
}

func (m CancelMsg) String() string {
	return "cancel edit"
}

// Broadcasts the cancel message with an optional error
func Abort(err error) tea.Cmd {
	return func() tea.Msg {
		if err != nil {
			slog.Error("aborting edit", "model", name, "error", err)
		}
		return CancelMsg{Error: err}
	}
}

type SubmitMsg struct{}

func Submit(persister tea.Cmd) tea.Cmd {
	return tea.Batch(persister, func() tea.Msg { return SubmitMsg{} })
}

func Persister(client *task.Client) func(task.Task) tea.Cmd {
	return func(t task.Task) tea.Cmd {
		return client.Put(t)
	}
}
