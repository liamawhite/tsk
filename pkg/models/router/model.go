package router

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/tsk/pkg/models/tasks/edit"
	"github.com/liamawhite/tsk/pkg/models/tasks/list"
	"github.com/liamawhite/tsk/pkg/task"
)

const name = "root"

func NewModel(client *task.Client) tea.Model {
	return Model{
		mode:   taskList,
        keys:   keyMap,
        client: client,

        list: list.New(client),
        // edit is populated when the user wants to edit a task
	}
}

type Model struct {
	keys   KeyMap
	mode   mode
	client *task.Client

	list   list.Model
	edit   edit.Model
}

func (m Model) Init() tea.Cmd {
    return m.list.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    slog.Debug("received msg", "model", name, "msg", msg, "msgType", fmt.Sprintf("%T", msg))
	// Stop short if we're quitting
    if msg, ok := msg.(tea.KeyMsg); ok {
        if key.Matches(msg, m.keys.Quit) {
            slog.Info("quitting", "model", name, "mode", m.mode)
            return m, tea.Quit
        }
    }

    // Handle messages from the sub-models
    switch msg := msg.(type) {

    // Handle messages from the task lister
	case list.AddMsg:
		m.mode = taskEdit
		m.edit = edit.New(edit.AddPopulator(), edit.Persister(m.client))
		return m, m.edit.Init()
	case list.EditMsg:
		m.mode = taskEdit
		m.edit = edit.New(edit.EditPopulator(m.client, msg.Id), edit.Persister(m.client))
		return m, m.edit.Init()

    // Handle messages from the task editor
	case edit.CancelMsg, edit.SubmitMsg:
		m.mode = taskList
        return m, nil 
	}

    // Route messages to the appropriate sub-model
    return m.routing(msg)
}

func (m Model) routing(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.mode {

    case taskList:
        m.list.Focus()
        l, cmd := m.list.Update(msg)
        if l, ok := l.(list.Model); ok {
            m.list = l
            return m, cmd
        }

    case taskEdit:
        e, cmd := m.edit.Update(msg)
        if e, ok := e.(edit.Model); ok {
            m.edit = e
            return m, cmd
        }
    }
    return m, nil
}

func (m Model) View() string {
	switch m.mode {
	case taskList:
		return m.list.View()
	
    case taskEdit:
        return m.edit.View()
    }

    return ""
}

