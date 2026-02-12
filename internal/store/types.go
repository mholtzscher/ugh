package store

import "time"

type State string

const (
	StateInbox   State = "inbox"
	StateNow     State = "now"
	StateWaiting State = "waiting"
	StateLater   State = "later"
	StateDone    State = "done"
)

const (
	TaskEventKindCreate = "create"
	TaskEventKindUpdate = "update"
	TaskEventKindDone   = "done"
	TaskEventKindUndo   = "undo"
	TaskEventKindDelete = "delete"
)

type Task struct {
	ID          int64
	State       State
	PrevState   *State
	Title       string
	Notes       string
	DueOn       *time.Time
	WaitingFor  string
	CompletedAt *time.Time
	Projects    []string
	Contexts    []string
	Meta        map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListTasksByExprOptions struct {
	ExcludeDone bool
	OnlyDone    bool
}

type NameCount struct {
	Name  string
	Count int64
}

type ShellHistory struct {
	ID            int64
	Timestamp     int64
	Command       string
	Success       bool
	ResultSummary string
	Intent        string
}

type TaskEvent struct {
	ID             int64
	TaskID         int64
	Timestamp      int64
	Kind           string
	Summary        string
	ChangesJSON    string
	Origin         string
	ShellHistoryID *int64
	ShellCommand   string
}
