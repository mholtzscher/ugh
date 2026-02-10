package shell_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	require.NoError(t, scanner.Err(), "scan error")
	assert.Equal(t, []string{"line1", "line2", "line3"}, lines, "lines mismatch")
}

func TestScriptScannerLineNumber(t *testing.T) {
	t.Parallel()

	input := "first\nsecond\nthird"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	expectedLineNums := []int{1, 2, 3}
	i := 0
	for scanner.Scan() {
		assert.Equal(t, expectedLineNums[i], scanner.LineNumber(), "line number mismatch at index %d", i)
		i++
	}
}

func TestScriptScannerEmptyInput(t *testing.T) {
	t.Parallel()

	input := ""
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	assert.False(t, scanner.Scan(), "expected no lines from empty input")
	require.NoError(t, scanner.Err(), "scan error")
}

func TestScriptScannerSingleLine(t *testing.T) {
	t.Parallel()

	input := "only line"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	require.True(t, scanner.Scan(), "expected one line")
	assert.Equal(t, "only line", scanner.Text(), "line text mismatch")
	assert.Equal(t, 1, scanner.LineNumber(), "line number mismatch")
	assert.False(t, scanner.Scan(), "expected no more lines")
}

func TestScriptScannerWithComments(t *testing.T) {
	t.Parallel()

	input := "# this is a comment\nactual command\n  # indented comment\n"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.Equal(t, []string{"# this is a comment", "actual command", "  # indented comment"}, lines, "lines mismatch")
}

func TestScriptScannerWithBlankLines(t *testing.T) {
	t.Parallel()

	input := "line1\n\nline2\n\n\nline3"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.Equal(t, []string{"line1", "", "line2", "", "", "line3"}, lines, "lines mismatch")
}

func TestScriptScannerWithQuitCommand(t *testing.T) {
	t.Parallel()

	input := "command1\nquit\ncommand2"
	scanner := shell.NewScriptScanner(strings.NewReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.Equal(t, []string{"command1", "quit", "command2"}, lines, "lines mismatch")
}
