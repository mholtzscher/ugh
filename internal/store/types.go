package store

import "time"

type Task struct {
	ID             int64
	Done           bool
	Priority       string
	CompletionDate *time.Time
	CreationDate   *time.Time
	Description    string
	Projects       []string
	Contexts       []string
	Meta           map[string]string
	Unknown        []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Filters struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
	Project  string
	Context  string
	Priority string
	Search   string
	Sort     string
}
