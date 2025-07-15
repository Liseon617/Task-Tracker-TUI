package model

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

type Task struct {
	status      Status
	title       string
	description string
}

type (
    TaskCreatedMsg Task
    TaskUpdatedMsg struct {
        OldTask Task
        NewTask Task
    }
)

func NewTask(status Status, title, description string) Task {
	return Task{status, title, description}
}

func (t *Task) Next() {
	if t.status < Done {
		t.status++
	}
}

func (t *Task) Prev() {
	if t.status > Todo {
		t.status--
	}
}

func (t *Task) Update (title, description string) {
	t.title = title
	t.description = description
}

func (t Task) FilterValue() string  { return t.title }
func (t Task) Title() string        { return t.title }
func (t Task) Description() string  { return t.description }
func (t Task) Status() Status       { return t.status }
