package shell

import (
	"time"
)

// ExecuteResult contains the result of command execution.
type ExecuteResult struct {
	Intent    string
	Message   string
	TaskIDs   []int64
	Summary   string
	Timestamp time.Time
}

// ExecuteOptions provides options for command execution.
type ExecuteOptions struct {
	Now            time.Time
	SelectedTaskID *int64
}
