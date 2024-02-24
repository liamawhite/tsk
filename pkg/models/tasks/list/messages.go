package list

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type EditMsg struct {
    Id string
}

func (m EditMsg) String() string {
    return fmt.Sprintf("edit task %s", m.Id)
}

func newEditMsg(id string) tea.Cmd {
    return func() tea.Msg {
        return EditMsg{Id: id}
    }
}

type AddMsg struct {}

func (m AddMsg) String() string {
    return "add task"
}

func newAddMsg() tea.Cmd {
    return func() tea.Msg {
        return AddMsg{}
    }
}

