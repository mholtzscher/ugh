package nlp

type Mode int

const (
	ModeAuto Mode = iota
	ModeCreate
	ModeUpdate
	ModeFilter
)

type Intent int

const (
	IntentUnknown Intent = iota
	IntentCreate
	IntentUpdate
	IntentFilter
)

type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
)

type Position struct {
	Offset int
	Line   int
	Column int
}

type Span struct {
	Start Position
	End   Position
}

type Diagnostic struct {
	Severity Severity
	Code     string
	Message  string
	Span     Span
	Hint     string
}

type ParseResult struct {
	Intent      Intent
	Command     CommandAST
	Diagnostics []Diagnostic
	Canonical   string
}
