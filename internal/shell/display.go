package shell

import (
	"fmt"
	"io"

	"github.com/mholtzscher/ugh/internal/output"
)

// Display handles output formatting.
type Display struct {
	mode   DisplayMode
	writer output.Writer
}

// DisplayMode defines how to display results.
type DisplayMode int

const (
	DisplayCompact DisplayMode = iota
	DisplayTable
	DisplayDetail
)

// NewDisplay creates a new display handler.
func NewDisplay() *Display {
	return &Display{
		mode:   DisplayCompact,
		writer: output.NewWriter(false, false),
	}
}

// ShowResult displays an execution result.
func (d *Display) ShowResult(result *ExecuteResult) {
	if result == nil {
		return
	}

	switch d.mode {
	case DisplayCompact:
		d.showCompact(result)
	case DisplayTable:
		d.showTable(result)
	case DisplayDetail:
		d.showDetail(result)
	}
}

func (d *Display) showCompact(result *ExecuteResult) {
	if result.Message != "" {
		_ = d.writer.WriteLine(result.Message)
	}
}

func (d *Display) showTable(result *ExecuteResult) {
	if result.Message != "" {
		_ = d.writer.WriteLine(result.Message)
	}
}

func (d *Display) showDetail(result *ExecuteResult) {
	if result.Message != "" {
		_ = d.writer.WriteLine(result.Message)
	}
}

// Clear clears the screen.
func (d *Display) Clear() {
	_, _ = fmt.Fprint(io.Discard, "\033[H\033[2J")
}
