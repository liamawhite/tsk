package task

import (
	"time"

	"github.com/liamawhite/jsondb"
	"github.com/teambition/rrule-go"
)

type Status int

const (
	Backlog Status = iota
	Blocked
	Paused
	InProgress
	Done
)

var Statuses = map[Status]string{
    Backlog:    "Backlog",
    Blocked:    "Blocked",
    Paused:     "Paused",
    InProgress: "In Progress",
    Done:       "Done",
}


func (s Status) String() string {
	return Statuses[s]
}

type Task struct {
	Id string

	Name   string
	Notes  string
	Status Status

	Due        time.Time
	Scheduled  time.Time
	Recurrence rrule.RRule

	ProjectId *string
}

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

func (c *Client) Get(id string) (Task, error) {
    return c.db.Read(id)
}

// Idempotent based on the task's ID.
func (c *Client) Put(task Task) error {
    return c.db.Write(task.Id, task)
}

func (c *Client) Delete(id string) error {
    return c.db.Delete(id)
}

func (c *Client) List() ([]Task, error) {
    return c.db.List()
}
