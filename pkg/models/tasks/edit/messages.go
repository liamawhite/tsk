package edit

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/liamawhite/tsk/pkg/task"
)


type PopulatorMsg struct {
	Task  task.Task
	Error error
}

func (p PopulatorMsg) String() string {
	return fmt.Sprintf("PopulatorMsg{Task: %v}", p.Task.Name)
}

func NewAddPopulator() tea.Cmd {
	return func() tea.Msg {
		slog.Info("editing a new task", "model", name)
		return PopulatorMsg{Task: task.Task{Id: uuid.New().String()}, Error: nil}
	}
}

func NewEditPopulator(client *task.Client, id string) tea.Cmd {
	return func() tea.Msg {
        slog.Info("editing an existing task", "model", name, "id", id)
		task, err := client.Get(id)
		return PopulatorMsg{Task: task, Error: err}
	}
}

type CancelMsg struct{
    Error error
}

func Abort(err error) tea.Cmd {
    return func() tea.Msg {
        if err != nil {
            slog.Error("aborting edit", "model", name, "error", err)
        }
        return CancelMsg{Error: err}
    }
}

type SubmitMsg struct {
    Task task.Task
    Error error
}

func (s SubmitMsg) String() string {
    if s.Error != nil {
        return fmt.Sprintf("{Error: %v}",s.Error)
    }
    return fmt.Sprintf("{Task: %v}", s.Task.Id)
}

func NewPersister(client *task.Client) func(task.Task) tea.Cmd {
    return func(t task.Task) tea.Cmd {
        return func() tea.Msg {
            slog.Info("persisting task", "model", name, "id", t.Id)
            err := client.Put(t)
            return SubmitMsg{Task: t, Error: err}
        }
    }
}
