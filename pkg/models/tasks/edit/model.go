package edit

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/liamawhite/tsk/pkg/task"
)

const name = "tasks/edit"

func New(populator tea.Cmd, persister func(task.Task) tea.Cmd) Model {
	return Model{keys: keyMap, populator: populator, persister: persister}
}

type Model struct {
	keys KeyMap
	form *huh.Form

	populator tea.Cmd
	persister func(task.Task) tea.Cmd

	id string
}

func (m Model) Init() tea.Cmd {
	slog.Debug("initializing model", "model", name)
	return m.populator
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("received msg", "model", name, "msg", msg, "msgType", fmt.Sprintf("%T", msg))

	// Populate task is only ever called on Init so we can safely call form init here
	if msg, ok := msg.(PopulatorMsg); ok {
		if msg.Error != nil {
			return m, Abort(msg.Error)
		}
		m.id = msg.Task.Id
		m.form = buildForm(msg.Task)
		return m, m.form.Init()
	}

	// If the form is done, we can write to the database, then broadcast 
	if m.form.State == huh.StateCompleted {
		slog.Info("form completed", "model", name)    
		return m, m.persister(m.task())
	}

    // If the form is aborted, we broadcast the abort message
	if m.form.State == huh.StateAborted {
		slog.Info("form aborted", "model", name)
        return m, Abort(nil)
	}

    return m.routing(msg)
}

func (m Model) task() task.Task {
    return task.Task{
        Id:   m.id,
        Name: m.form.GetString(taskKey),
    }
}

func (m Model) routing(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("routing message", "model", name, "msg", msg, "msgType", fmt.Sprintf("%T", msg))

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.form == nil {
		return ""
	}
	return m.form.View()
}

var (
	idKey   = "id"
	taskKey = "task"
)

func buildForm(t task.Task) *huh.Form {
	nameInput := huh.NewInput().Key(taskKey).Title("Task").Value(&t.Name)

	return huh.NewForm(
		huh.NewGroup(nameInput),
	)
}
