package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"

	"github.com/mholtzscher/ugh/internal/output"
)

// Display handles output formatting.
type Display struct {
	mode    DisplayMode
	writer  output.Writer
	noColor bool
}

// DisplayMode defines how to display results.
type DisplayMode int

const (
	DisplayCompact DisplayMode = iota
	DisplayTable
	DisplayDetail
)

// NewDisplay creates a new display handler.
func NewDisplay(noColor bool) *Display {
	return &Display{
		mode:    DisplayCompact,
		writer:  output.NewWriter(false, noColor),
		noColor: noColor,
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
	if result.Message == "" {
		return
	}

	if d.noColor {
		_ = d.writer.WriteLine(result.Message)
		return
	}

	// Map intents to appropriate pterm styles
	intent := strings.ToLower(result.Intent)
	switch {
	case strings.Contains(intent, "add"), strings.Contains(intent, "create"), strings.Contains(intent, "new"):
		_ = pterm.Success.Println(result.Message)
	case strings.Contains(intent, "done"), strings.Contains(intent, "complete"):
		_ = pterm.Success.Println(result.Message)
	case strings.Contains(intent, "undo"), strings.Contains(intent, "revert"):
		_ = pterm.Success.Println(result.Message)
	case strings.Contains(intent, "delete"), strings.Contains(intent, "rm"), strings.Contains(intent, "remove"):
		_ = pterm.Warning.Println(result.Message)
	case strings.Contains(intent, "error"), strings.Contains(intent, "fail"):
		_ = pterm.Error.Println(result.Message)
	case strings.Contains(intent, "show"),
		strings.Contains(intent, "list"),
		strings.Contains(intent, "find"),
		strings.Contains(intent, "filter"):
		_ = d.writer.WriteLine(result.Message)
	case strings.Contains(intent, "context"), strings.Contains(intent, "help"):
		_ = d.writer.WriteLine(result.Message)
	default:
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
	_, _ = fmt.Fprint(os.Stdout, "\033[H\033[2J")
}
