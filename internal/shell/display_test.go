package shell_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/shell"
)

func TestNewDisplay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			display := shell.NewDisplay(false, config.Display{})
			assert.NotNil(t, display, "NewDisplay returned nil")
		})
	}
}

func TestDisplayShowResultNil(t *testing.T) {
	t.Parallel()

	display := shell.NewDisplay(false, config.Display{})

	// Should not panic when result is nil
	display.ShowResult(nil)
}

func TestDisplayShowResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		result *shell.ExecuteResult
	}{
		{
			name: "create intent",
			result: &shell.ExecuteResult{
				Intent:  "create",
				Message: "Created task #1",
			},
		},
		{
			name: "update intent",
			result: &shell.ExecuteResult{
				Intent:  "update",
				Message: "Updated task #1",
			},
		},
		{
			name: "filter intent",
			result: &shell.ExecuteResult{
				Intent:  "filter",
				Message: "Found 5 tasks",
			},
		},
		{
			name: "context intent",
			result: &shell.ExecuteResult{
				Intent:  "context",
				Message: "Current context: #work",
			},
		},
		{
			name: "help intent",
			result: &shell.ExecuteResult{
				Intent:  "help",
				Message: "Available commands...",
			},
		},
		{
			name: "done intent",
			result: &shell.ExecuteResult{
				Intent:  "done",
				Message: "Marked 3 tasks as done",
			},
		},
		{
			name: "delete intent",
			result: &shell.ExecuteResult{
				Intent:  "delete",
				Message: "Deleted task #1",
			},
		},
		{
			name: "error intent",
			result: &shell.ExecuteResult{
				Intent:  "error",
				Message: "Something went wrong",
			},
		},
		{
			name: "empty message",
			result: &shell.ExecuteResult{
				Intent:  "create",
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			display := shell.NewDisplay(false, config.Display{})

			// Should not panic
			display.ShowResult(tt.result)
		})
	}
}

func TestDisplayClear(t *testing.T) {
	t.Parallel()

	display := shell.NewDisplay(false, config.Display{})

	// Should not panic - outputs ANSI escape codes
	display.Clear()
}

func TestDisplayShowResultWithSummary(t *testing.T) {
	t.Parallel()

	result := &shell.ExecuteResult{
		Intent:  "create",
		Message: "Created task #1: Buy milk",
		Summary: "created task #1",
		TaskIDs: []int64{1},
	}

	display := shell.NewDisplay(false, config.Display{})

	// Should not panic
	display.ShowResult(result)
}
