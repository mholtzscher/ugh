//nolint:testpackage // Tests validate internal table helpers directly.
package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"

	"github.com/mholtzscher/ugh/internal/store"
)

func TestTaskTableRowsMatchColumns(t *testing.T) {
	tasks := []*store.Task{
		{ID: 1, State: store.StateInbox, Title: "first"},
		{ID: 2, State: store.StateDone, Title: "second"},
	}

	tests := []struct {
		name      string
		showState bool
	}{
		{name: "all tab includes state", showState: true},
		{name: "state tab omits state", showState: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := taskTableColumns(tt.showState, 80)
			rows := taskTableRows(tasks, tt.showState)
			if len(rows) != len(tasks) {
				t.Fatalf("expected %d rows, got %d", len(tasks), len(rows))
			}
			for i, row := range rows {
				if len(row) != len(cols) {
					t.Fatalf("row %d has %d cells, want %d", i, len(row), len(cols))
				}
			}
		})
	}
}

func TestSetTaskTableDataHandlesColumnCountChange(t *testing.T) {
	taskTable := table.New()

	setTaskTableData(
		&taskTable,
		taskTableColumns(true, 80),
		taskTableRows([]*store.Task{{ID: 1, State: store.StateInbox, Title: "first"}}, true),
		0,
	)

	setTaskTableData(
		&taskTable,
		taskTableColumns(false, 80),
		taskTableRows([]*store.Task{{ID: 1, State: store.StateInbox, Title: "first"}}, false),
		0,
	)

	if got := len(taskTable.Columns()); got != 3 {
		t.Fatalf("expected 3 columns after switch, got %d", got)
	}
	if got := len(taskTable.Rows()); got != 1 {
		t.Fatalf("expected 1 row after switch, got %d", got)
	}
}
