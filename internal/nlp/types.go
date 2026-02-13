package nlp

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=Mode,Intent,Severity -output=types_string.go

type Mode int

const (
	ModeAuto Mode = iota
	ModeCreate
	ModeUpdate
	ModeFilter
	ModeView
	ModeContext
)

type Intent int

const (
	IntentUnknown Intent = iota
	IntentCreate
	IntentUpdate
	IntentFilter
	IntentView
	IntentContext
	IntentLog
)

type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
)

type Diagnostic struct {
	Severity Severity
	Code     string
	Message  string
	Hint     string
}

type ParseResult struct {
	Intent      Intent
	Command     Command
	Diagnostics []Diagnostic
}
