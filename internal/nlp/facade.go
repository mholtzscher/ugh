package nlp

import (
	"context"
	"time"
)

type ParseOptions struct {
	Mode Mode
	Now  time.Time
}

type Parser interface {
	Parse(ctx context.Context, input string, opts ParseOptions) (ParseResult, error)
}

type defaultParser struct{}

func NewParser() Parser {
	return defaultParser{}
}

func Parse(input string, opts ParseOptions) (ParseResult, error) {
	return defaultParser{}.Parse(context.Background(), input, opts)
}

func (defaultParser) Parse(_ context.Context, input string, opts ParseOptions) (ParseResult, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}
	return parseInput(input, opts)
}
