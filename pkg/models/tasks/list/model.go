package list

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/tsk/pkg/models/components/titledtable"
	"github.com/liamawhite/tsk/pkg/task"
)

const name = "tasks/list"

func New(client *task.Client) Model {
    m := Model{
        // hardcode these for now
        width:  80,
        height: 20,

        client:  client,

		keys:    keyMap,
		tasks:   []task.Task{},
	}
    m.table = m.buildTable()
    return m
}

type Model struct {
    width int
    height int
	keys KeyMap

    client *task.Client

	tasks []task.Task
	table titledtable.Model[task.Task]
    

}

func (m Model) Init() tea.Cmd {
	return m.client.List()
}

func (m *Model) Focus() {
	m.table.Focus()
}

func (m *Model) Blur() {
	m.table.Blur()
}

func (m Model) Refresh() tea.Cmd {
	return m.client.List()
}

func (m Model) SelectedTask() task.Task {
	return m.table.SelectedRow().Data
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("received msg", "model", name, "msg", msg)

	// Handle list wide navigation
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keys.New):
			return m, newMsg() 
		case key.Matches(msg, m.keys.Edit):
            return m, editMsg(m.SelectedTask().Id)  
		case key.Matches(msg, m.keys.Delete):
			return m, m.client.Delete(m.SelectedTask().Id)
        
        // Handle status changes
        case key.Matches(msg, m.keys.Todo):
            return m, m.setSelectedStatus(task.Todo)
        case key.Matches(msg, m.keys.Blocked):
            return m, m.setSelectedStatus(task.Blocked)
        case key.Matches(msg, m.keys.Paused):
            return m, m.setSelectedStatus(task.Paused)
        case key.Matches(msg, m.keys.Active):   
            return m, m.setSelectedStatus(task.Active)
        case key.Matches(msg, m.keys.Complete): 
            return m, m.setSelectedStatus(task.Complete)
        case key.Matches(msg, m.keys.Abandoned):
            return m, m.setSelectedStatus(task.Abandoned)
            
            
            
		}
	}

	// Handle the messages we care about
	switch msg := msg.(type) {

	// Indicates that we have a new list of tasks to display
	case task.ListMsg:
		m.tasks = msg.Tasks
		m.table = m.buildTable()
		return m, nil

	// If we have deleted a task or the editor submitted a new one then fetch the new list of tasks
	case task.ModifyMsg, task.DeleteMsg:
		return m, m.Refresh()

	}

	return m.routing(msg)
}

func (m Model) setSelectedStatus(status task.Status) tea.Cmd {
    t := m.SelectedTask()
    t.Status = status
    return m.client.Put(t)
}

func (m Model) routing(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("routing message", "model", name, "msg", msg, "msgType", fmt.Sprintf("%T", msg))

	t, cmd := m.table.Update(msg)
	m.table = t
	return m, cmd
}

