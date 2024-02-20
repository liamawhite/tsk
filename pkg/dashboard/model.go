package dashboard

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/liamawhite/tsk/pkg/task"
	"github.com/samber/lo"
)

func NewModel() tea.Model {
	tasks := []task.Task{
		{Name: "Task 1"},
	}
	return Model{
		mode:   base,
		tasks:  tasks,
		editor: newEditor(nil),
		table:  newTable(tasks),
	}
}

type mode int

const (
	base mode = iota
	edit
)

type Model struct {
	mode   mode
    tasks  []task.Task
	editor tea.Model
	table  table.Model
}

func (m Model) Init() tea.Cmd {
	return m.editor.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quitting at the top level.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.mode {
	case edit:
		form, cmd := m.editor.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.editor = f
		}
		return m, cmd
	default:
		// Handle navigation key presses when in the base mode only, otherwise delegate to the submodels
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "a":
				m.mode = edit
				m.editor = newEditor(nil)
				return m, nil
			}
		}

	}

	return m, nil
}

func (m Model) View() string {
	switch m.mode {
	case edit:
		return m.editor.View()
	default:
		return m.table.View()
	}
}

func taskRow(t task.Task, _ int) table.Row {
	return table.Row{t.Name}
}

func newTable(tasks []task.Task) table.Model {
	return table.New(
		table.WithColumns([]table.Column{
			{Title: "Task", Width: 15},
		}),
		table.WithRows(lo.Map(tasks, taskRow)),
	)
}

func newEditor(task *task.Task) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("task").
				Key("task"),
		),
	)
}
