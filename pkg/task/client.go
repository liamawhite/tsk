package task

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/jsondb"
)

// This is pass-through for now, but expect this to also support caching and additional list filtering in the future.
func NewClient(persistenceDir string) (*Client, error) {
	db, err := jsondb.NewFS[Task](persistenceDir)
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

type Client struct {
	db jsondb.Client[Task]
}

func (c *Client) Get(id string) tea.Cmd {
    return func() tea.Msg {
        slog.Info("retrieving task", "id", id)
        task, err := c.db.Read(id)
        return GetMsg{Task: task, Error: err}
    }
}

// Idempotent based on the task's ID.
func (c *Client) Put(task Task) tea.Cmd {
    return func() tea.Msg {
        slog.Info("persisting task", "id", task.Id)
        err := c.db.Write(task.Id, task)
        return ModifyMsg{Task: task, Error: err}
    }
}

func (c *Client) Delete(id string) tea.Cmd {
    return func() tea.Msg {
        slog.Info("deleting task", "id", id)
        err := c.db.Delete(id)
        return DeleteMsg{Id: id, Error: err}
    }
}

func (c *Client) List() tea.Cmd {
    return func() tea.Msg {
        slog.Info("retrieving all tasks")
        tasks, err := c.db.List()
        return ListMsg{Tasks: tasks, Error: err}
    }
}


