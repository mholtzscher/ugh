package nlp

import (
	"context"
	"errors"
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

	root, err := dslParser.ParseString("", input)
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

	if root == nil || root.Cmd == nil {
		return ParseResult{Intent: IntentUnknown}, errors.New("empty parse result")
	}

	intent, cmdResult, postErr := postProcess(root.Cmd)
	if postErr != nil {
		return ParseResult{Intent: intent}, postErr
	}

	if opts.Mode != ModeAuto {
		want := IntentUnknown
		switch opts.Mode {
		case ModeAuto:
			want = IntentUnknown
		case ModeCreate:
			want = IntentCreate
		case ModeUpdate:
			want = IntentUpdate
		case ModeFilter:
			want = IntentFilter
		case ModeView:
			want = IntentView
		case ModeContext:
			want = IntentContext
		}
		if want != IntentUnknown && intent != want {
			return ParseResult{Intent: intent, Command: cmdResult}, errors.New("command does not match parse mode")
		}
	}

	return ParseResult{Intent: intent, Command: cmdResult}, nil
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

func postProcess(cmd Command) (Intent, Command, error) {
	switch typed := cmd.(type) {
	case *CreateCommand:
		if err := typed.postProcess(); err != nil {
			return IntentCreate, typed, err
		}
		return IntentCreate, typed, nil
	case *UpdateCommand:
		typed.postProcess()
		return IntentUpdate, typed, nil
	case *FilterCommand:
		if err := typed.postProcess(); err != nil {
			return IntentFilter, typed, err
		}
		return IntentFilter, typed, nil
	case *ViewCommand:
		if err := typed.postProcess(); err != nil {
			return IntentView, typed, err
		}
		return IntentView, typed, nil
	case *ContextCommand:
		if err := typed.postProcess(); err != nil {
			return IntentContext, typed, err
		}
		return IntentContext, typed, nil
	default:
		return IntentUnknown, cmd, errors.New("unknown command type")
	}
}
