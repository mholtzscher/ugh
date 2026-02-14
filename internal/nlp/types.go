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

type DiagnosticError struct {
	err         error
	diagnostics []Diagnostic
}

func NewDiagnosticError(err error, diagnostics []Diagnostic) error {
	if err == nil {
		return nil
	}
	if len(diagnostics) == 0 {
		return err
	}
	return DiagnosticError{err: err, diagnostics: diagnostics}
}

func (e DiagnosticError) Error() string {
	if e.err == nil {
		if len(e.diagnostics) == 0 {
			return ""
		}
		return e.diagnostics[0].Message
	}
	return e.err.Error()
}

func (e DiagnosticError) Unwrap() error {
	return e.err
}

func (e DiagnosticError) Diagnostics() []Diagnostic {
	if len(e.diagnostics) == 0 {
		return nil
	}
	items := make([]Diagnostic, len(e.diagnostics))
	copy(items, e.diagnostics)
	return items
}
