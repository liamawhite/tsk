package list

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/tsk/pkg/models/tasks/edit"
	"github.com/liamawhite/tsk/pkg/task"
	"github.com/samber/lo"
)

const name = "tasks/list"

func New(lister tea.Cmd, deleter func(string) tea.Cmd) Model {
	return Model{
        keys: keyMap,
		tasks:  []task.Task{},
		table:  buildTable([]task.Task{}),
		lister: lister,
        deleter: deleter,
	}
}

type Model struct {
	keys KeyMap

	tasks []task.Task
	table table.Model

	lister tea.Cmd
    deleter func(string) tea.Cmd
}

func (m Model) Init() tea.Cmd {
	return m.lister
}

func (m *Model) Focus() {
	m.table.Focus()
}

func (m *Model) Blur() {
	m.table.Blur()
}

func (m Model) Refresh() tea.Cmd {
    return m.lister
}

func (m Model) SelectedTask() task.Task {
	row := m.table.SelectedRow()
	return task.Task{
		Id:   row[1],
		Name: row[0],
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("received msg", "model", name, "msg", msg)

	// Handle list wide navigation
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keys.Add):
			return m, func() tea.Msg { return AddMsg{} }
        case key.Matches(msg, m.keys.Edit):
            return m, func() tea.Msg { return EditMsg{Id: m.SelectedTask().Id} }
        case key.Matches(msg, m.keys.Delete):
            return m, m.deleter(m.SelectedTask().Id)    
		}
	}

	// Handle the messages we care about
    switch msg := msg.(type) {
    
    // Indicates that we have a new list of tasks to display
    case ListTasksMsg:
		m.tasks = msg.tasks
		m.table = buildTable(m.tasks)
        return m, nil

    // If we have deleted a task or the editor submitted a new one then fetch the new list of tasks
    case DeletedTaskMsg, edit.SubmitMsg:
         return m, m.Refresh()
	
    }

	return m.routing(msg)
}

func (m Model) routing(msg tea.Msg) (tea.Model, tea.Cmd) {
    slog.Debug("routing message", "model", name, "msg", msg, "msgType", fmt.Sprintf("%T", msg))

    t, cmd := m.table.Update(msg)
    m.table = t
    return m, cmd
}


func (m Model) View() string {
	return m.table.View()
}

func buildTable(tasks []task.Task) table.Model {
	return table.New(
		table.WithColumns([]table.Column{
			{"Task", 20},
			{"Id", 0},
		}),
		table.WithRows(lo.Map(tasks, func(t task.Task, _ int) table.Row {
			return table.Row{t.Name, t.Id}
		})),
	)
}
