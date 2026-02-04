package store

import "time"

type Status string

const (
	StatusInbox   Status = "inbox"
	StatusNext    Status = "next"
	StatusWaiting Status = "waiting"
	StatusSomeday Status = "someday"
)

type Task struct {
	ID          int64
	Done        bool
	Status      Status
	Priority    string
	Title       string
	Notes       string
	DueOn       *time.Time
	DeferUntil  *time.Time
	WaitingFor  string
	CompletedAt *time.Time
	Projects    []string
	Contexts    []string
	Meta        map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Filters struct {
	All             bool
	DoneOnly        bool
	TodoOnly        bool
	Status          string
	Project         string
	Context         string
	Priority        string
	Search          string
	DueSetOnly      bool
	DeferAfter      string // YYYY-MM-DD
	DeferOnOrBefore string // YYYY-MM-DD
}

type NameCount struct {
	Name  string
	Count int64
}
