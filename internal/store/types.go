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
	Recent      bool
	Limit       int64
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

type TaskVersion struct {
	VersionID   int64
	TaskID      int64
	State       State
	PrevState   *State
	Title       string
	Notes       string
	DueOn       *time.Time
	WaitingFor  string
	CompletedAt *time.Time
	UpdatedAt   time.Time
	Deleted     bool
	Projects    []string
	Contexts    []string
	Meta        map[string]string
}
