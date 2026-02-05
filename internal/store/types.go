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
	Priority    string
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

type Filters struct {
	All        bool
	DoneOnly   bool
	TodoOnly   bool
	State      string
	Project    string
	Context    string
	Priority   string
	Search     string
	DueSetOnly bool
}

type NameCount struct {
	Name  string
	Count int64
}
