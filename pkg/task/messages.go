package task

import "fmt"

type ModifyMsg struct {
    Task Task
    Error error
}

func (s ModifyMsg) String() string {
    if s.Error != nil {
        return fmt.Sprintf("{Error: %v}",s.Error)
    }
    return fmt.Sprintf("{Task: %v}", s.Task.Id)
}

type DeleteMsg struct {
    Id string
    Error error
}

func (s DeleteMsg) String() string {
    if s.Error != nil {
        return fmt.Sprintf("{Error: %v}",s.Error)
    }
    return fmt.Sprintf("{Task: %v}", s.Id)
}

type GetMsg struct {
    Task Task
    Error error
}

func (s GetMsg) String() string {
    if s.Error != nil {
        return fmt.Sprintf("{Error: %v}",s.Error)
    }
    return fmt.Sprintf("{Task: %v}", s.Task.Id)
}

type ListMsg struct {
    Tasks []Task
    Error error
}

func (s ListMsg) String() string {
    if s.Error != nil {
        return fmt.Sprintf("{Error: %v}",s.Error)
    }
    return fmt.Sprintf("{Tasks: %v}", len(s.Tasks))
}
