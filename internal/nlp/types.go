package nlp

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
