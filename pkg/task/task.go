package task

import (
	"time"

	"github.com/teambition/rrule-go"
)

type Status int

const (
	Todo Status = iota
	Blocked
	Paused
	Active
	Complete
	Abandoned
)

var Statuses = map[Status]string{
	Todo:      "Todo",
	Blocked:   "Blocked",
	Paused:    "Paused",
	Active:    "Active",
	Complete:  "Complete",
	Abandoned: "Abandoned",
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
