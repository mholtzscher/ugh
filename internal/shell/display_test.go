package shell_test

import (
	"testing"

	"github.com/mholtzscher/ugh/internal/shell"
)

func TestNewDisplay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		noColor bool
	}{
		{
			name:    "with color",
			noColor: false,
		},
		{
			name:    "no color",
			noColor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			display := shell.NewDisplay(tt.noColor)
			if display == nil {
				t.Fatal("NewDisplay returned nil")
			}
		})
	}
}

func TestDisplayShowResultNil(t *testing.T) {
	t.Parallel()

	display := shell.NewDisplay(true)

	// Should not panic when result is nil
	display.ShowResult(nil)
}

func TestDisplayShowResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		result  *shell.ExecuteResult
		noColor bool
	}{
		{
			name: "create intent",
			result: &shell.ExecuteResult{
				Intent:  "create",
				Message: "Created task #1",
			},
			noColor: true,
		},
		{
			name: "update intent",
			result: &shell.ExecuteResult{
				Intent:  "update",
				Message: "Updated task #1",
			},
			noColor: true,
		},
		{
			name: "filter intent",
			result: &shell.ExecuteResult{
				Intent:  "filter",
				Message: "Found 5 tasks",
			},
			noColor: true,
		},
		{
			name: "context intent",
			result: &shell.ExecuteResult{
				Intent:  "context",
				Message: "Current context: #work",
			},
			noColor: true,
		},
		{
			name: "help intent",
			result: &shell.ExecuteResult{
				Intent:  "help",
				Message: "Available commands...",
			},
			noColor: true,
		},
		{
			name: "done intent",
			result: &shell.ExecuteResult{
				Intent:  "done",
				Message: "Marked 3 tasks as done",
			},
			noColor: true,
		},
		{
			name: "delete intent",
			result: &shell.ExecuteResult{
				Intent:  "delete",
				Message: "Deleted task #1",
			},
			noColor: true,
		},
		{
			name: "error intent",
			result: &shell.ExecuteResult{
				Intent:  "error",
				Message: "Something went wrong",
			},
			noColor: true,
		},
		{
			name: "empty message",
			result: &shell.ExecuteResult{
				Intent:  "create",
				Message: "",
			},
			noColor: true,
		},
		{
			name: "with color",
			result: &shell.ExecuteResult{
				Intent:  "create",
				Message: "Created task #1",
			},
			noColor: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			display := shell.NewDisplay(tt.noColor)

			// Should not panic
			display.ShowResult(tt.result)
		})
	}
}

func TestDisplayClear(t *testing.T) {
	t.Parallel()

	display := shell.NewDisplay(true)

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

	display := shell.NewDisplay(true)

	// Should not panic
	display.ShowResult(result)
}
