package shell

import (
	"time"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/store"
)

type ResultLevel int

const (
	ResultLevelInfo ResultLevel = iota
	ResultLevelSuccess
	ResultLevelWarning
	ResultLevelError
)

// ExecuteResult contains the result of command execution.
type ExecuteResult struct {
	Intent    string
	Message   string
	TaskIDs   []int64
	Task      *store.Task
	Tasks     []*store.Task
	Versions  []*store.TaskVersion
	Context   *output.ContextStatus
	ViewHelp  *output.ViewHelp
	Level     ResultLevel
	Summary   string
	Timestamp time.Time
}

// ExecuteOptions provides options for command execution.
type ExecuteOptions struct {
	Now            time.Time
	SelectedTaskID *int64
}
