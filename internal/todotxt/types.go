package todotxt

import "time"

type Parsed struct {
	Done           bool
	Priority       string
	CompletionDate *time.Time
	CreationDate   *time.Time
	Description    string
	Projects       []string
	Contexts       []string
	Meta           map[string]string
	Unknown        []string
}
