package router



type mode int

const (
	taskList mode = iota
    taskEdit
)

func (m mode) String() string {
    switch m {
    case taskList:
        return "taskList"
    case taskEdit:
        return "taskEdit"
    }


    return "unknown"
}
