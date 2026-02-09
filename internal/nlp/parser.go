package nlp

import (
	"context"
	"time"
)

// ParseOptions configures the parser behavior.
type ParseOptions struct {
	Mode Mode
	Now  time.Time
}

// Parse parses the input string and returns a ParseResult.
func Parse(input string, opts ParseOptions) (ParseResult, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}

	// Use participle to parse the input into grammar structures
	gInput, err := dslParser.ParseString("", input)
	if err != nil {
		// Try to determine what kind of error this is
		return ParseResult{
			Intent: IntentUnknown,
			Diagnostics: []Diagnostic{{
				Severity: SeverityError,
				Code:     "E_PARSE",
				Message:  err.Error(),
			}},
		}, err
	}

	// Convert grammar to AST
	return convertGrammar(gInput, opts)
}

// Parser interface for dependency injection.
type Parser interface {
	Parse(ctx context.Context, input string, opts ParseOptions) (ParseResult, error)
}

type defaultParser struct{}

// NewParser creates a new Parser instance.
func NewParser() Parser {
	return defaultParser{}
}

// Parse implements the Parser interface.
func (defaultParser) Parse(_ context.Context, input string, opts ParseOptions) (ParseResult, error) {
	return Parse(input, opts)
}
