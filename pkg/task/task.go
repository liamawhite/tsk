package task

import (
    "time"
    
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

type Task struct {
	Id string

	Name   string
	Notes  string
	Status Status

	Due       time.Time
	Scheduled time.Time
    Recurrence rrule.RRule

    ProjectId *string
}
