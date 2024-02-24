package edit

import (
	"fmt"
	"log/slog"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/liamawhite/tsk/pkg/task"
	"github.com/samber/lo"
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

    // If a task is retrieved, and we get routed the message, assume
    // it's a request to add/edit a task
	if msg, ok := msg.(task.GetMsg); ok {
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
		return m, Submit(m.persister(m.task()))
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
		Id:     m.id,
		Name:   m.form.GetString(taskKey),
		Notes:  m.form.GetString(notesKey),
		Status: m.form.Get(statusKey).(task.Status),
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
	taskKey   = "task"
	statusKey = "status"
	notesKey  = "notes"
)

func buildForm(t task.Task) *huh.Form {
	nameInput := huh.NewInput().Key(taskKey).Title("Task").Value(&t.Name)
	notesInput := huh.NewText().Key(notesKey).Title("Notes").Value(&t.Notes)

    // Sort statuses so they are in the correct order every time
	statuses := lo.Entries(task.Statuses)
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Key < statuses[j].Key
	})

	statusInput := huh.NewSelect[task.Status]().
		Key(statusKey).
		Title("Status").
		Options(
			lo.Map(statuses, func(e lo.Entry[task.Status, string], _ int) huh.Option[task.Status] {
				return huh.NewOption(e.Value, e.Key)
			})...,
		).Value(&t.Status)

	return huh.NewForm(
		huh.NewGroup(nameInput, statusInput, notesInput),
	).WithShowErrors(true)
}
