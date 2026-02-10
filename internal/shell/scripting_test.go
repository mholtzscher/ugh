package shell_test

import (
	"strings"
	"testing"

	"github.com/mholtzscher/ugh/internal/shell"
)

func TestScriptScannerScanLines(t *testing.T) {
	t.Parallel()

	input := "line1\nline2\nline3"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("scan error: %v", err)
	}

	want := []string{"line1", "line2", "line3"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}
	for i, line := range lines {
		if line != want[i] {
			t.Errorf("line %d: got %q, want %q", i, line, want[i])
		}
	}
}

func TestScriptScannerLineNumber(t *testing.T) {
	t.Parallel()

	input := "first\nsecond\nthird"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	expectedLineNums := []int{1, 2, 3}
	i := 0
	for scanner.Scan() {
		if scanner.LineNumber() != expectedLineNums[i] {
			t.Errorf("line number: got %d, want %d", scanner.LineNumber(), expectedLineNums[i])
		}
		i++
	}
}

func TestScriptScannerEmptyInput(t *testing.T) {
	t.Parallel()

	input := ""
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	if scanner.Scan() {
		t.Error("expected no lines from empty input")
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("scan error: %v", err)
	}
}

func TestScriptScannerSingleLine(t *testing.T) {
	t.Parallel()

	input := "only line"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	if !scanner.Scan() {
		t.Fatal("expected one line")
	}

	if scanner.Text() != "only line" {
		t.Errorf("got %q, want %q", scanner.Text(), "only line")
	}

	if scanner.LineNumber() != 1 {
		t.Errorf("line number: got %d, want 1", scanner.LineNumber())
	}

	if scanner.Scan() {
		t.Error("expected no more lines")
	}
}

func TestScriptScannerWithComments(t *testing.T) {
	t.Parallel()

	input := "# this is a comment\nactual command\n  # indented comment\n"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	want := []string{"# this is a comment", "actual command", "  # indented comment"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}
	for i, line := range lines {
		if line != want[i] {
			t.Errorf("line %d: got %q, want %q", i, line, want[i])
		}
	}
}

func TestScriptScannerWithBlankLines(t *testing.T) {
	t.Parallel()

	input := "line1\n\nline2\n\n\nline3"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	want := []string{"line1", "", "line2", "", "", "line3"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}
	for i, line := range lines {
		if line != want[i] {
			t.Errorf("line %d: got %q, want %q", i, line, want[i])
		}
	}
}

func TestScriptScannerWithQuitCommand(t *testing.T) {
	t.Parallel()

	input := "command1\nquit\ncommand2"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	want := []string{"command1", "quit", "command2"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}
	for i, line := range lines {
		if line != want[i] {
			t.Errorf("line %d: got %q, want %q", i, line, want[i])
		}
	}
}
