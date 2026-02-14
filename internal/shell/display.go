package shell

import (
	"errors"
	"fmt"
	"os"

	"github.com/mholtzscher/ugh/internal/output"
)

// Display handles output formatting.
type Display struct {
	writer output.Writer
}

// NewDisplay creates a new display handler.
func NewDisplay(tty bool, writer output.Writer) *Display {
	writer.TTY = tty

	return &Display{
		writer: writer,
	}
}

// ShowResult displays an execution result.
func (d *Display) ShowResult(result *ExecuteResult) {
	if result == nil {
		return
	}

	if d.showPayload(result) {
		return
	}
	d.writeMessage(result.Message, result.Level)
}

func (d *Display) ShowError(err error) {
	_ = d.writer.WriteErr(err)
}

func (d *Display) showPayload(result *ExecuteResult) bool {
	if result.Context != nil {
		_ = d.writer.WriteContextStatus(*result.Context)
		return true
	}
	if result.ViewHelp != nil {
		_ = d.writer.WriteViewHelp(*result.ViewHelp)
		return true
	}
	if result.Tasks != nil {
		_ = d.writer.WriteTasks(result.Tasks)
		return true
	}
	if result.Task != nil {
		_ = d.writer.WriteTask(result.Task)
		return true
	}
	if result.Versions != nil {
		_ = d.writer.WriteTaskVersionDiff(result.Versions)
		return true
	}
	return false
}

func (d *Display) writeMessage(message string, level ResultLevel) {
	if message == "" {
		return
	}

	switch level {
	case ResultLevelInfo:
		_ = d.writer.WriteLine(message)
	case ResultLevelError:
		_ = d.writer.WriteErr(errors.New(message))
	case ResultLevelSuccess:
		_ = d.writer.WriteSuccess(message)
	case ResultLevelWarning:
		_ = d.writer.WriteWarning(message)
	}
}

// Clear clears the screen.
func (d *Display) Clear() {
	if !d.writer.TTY {
		return
	}
	_, _ = fmt.Fprint(os.Stdout, "\033[H\033[2J")
}
