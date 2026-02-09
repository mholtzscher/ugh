package shell

import (
	"bufio"
	"io"
)

// ScriptScanner scans input line by line for scripting mode.
type ScriptScanner struct {
	scanner *bufio.Scanner
	lineNum int
	text    string
}

// NewScriptScanner creates a scanner from an [io.Reader].
func NewScriptScanner(r io.Reader) *ScriptScanner {
	return &ScriptScanner{
		scanner: bufio.NewScanner(r),
		lineNum: 0,
	}
}

// Scan advances to the next line.
func (s *ScriptScanner) Scan() bool {
	if s.scanner.Scan() {
		s.lineNum++
		s.text = s.scanner.Text()
		return true
	}
	return false
}

// Text returns the current line.
func (s *ScriptScanner) Text() string {
	return s.text
}

// LineNumber returns the current line number.
func (s *ScriptScanner) LineNumber() int {
	return s.lineNum
}

// Err returns any scanning error.
func (s *ScriptScanner) Err() error {
	return s.scanner.Err()
}
