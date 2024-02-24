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

func editMsg(id string) tea.Cmd {
    return func() tea.Msg {
        return EditMsg{Id: id}
    }
}

type NewMsg struct {}

func (m NewMsg) String() string {
    return "add task"
}

func newMsg() tea.Cmd {
    return func() tea.Msg {
        return NewMsg{}
    }
}

