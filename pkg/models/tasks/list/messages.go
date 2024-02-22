package list

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/tsk/pkg/task"
)

type EditMsg struct {
	Id string
}

type AddMsg struct{}

type DeletedTaskMsg struct {
	Error error
}

func NewTaskDeleter(client *task.Client) func(string) tea.Cmd {
	return func(id string) tea.Cmd {
		return func() tea.Msg {
			err := client.Delete(id)
			return DeletedTaskMsg{Error: err}
		}
	}
}

type ListTasksMsg struct {
	tasks []task.Task
	error error
}

func (l ListTasksMsg) String() string {
    if l.error != nil {
        return fmt.Sprintf("{error:%v}", l.error)
    }
    return fmt.Sprintf("{tasks:%v}", len(l.tasks))
}

func NewTaskLister(client *task.Client) tea.Cmd {
	return func() tea.Msg {
		tasks, err := client.List()
		return ListTasksMsg{tasks: tasks, error: err}
	}
}
