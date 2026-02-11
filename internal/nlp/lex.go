package nlp

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

// LexToken is a lexer token with symbolic name.
type LexToken struct {
	Name  string
	Value string
	Pos   lexer.Position
}

// Lex tokenizes input with the shared DSL lexer.
func Lex(input string) ([]LexToken, error) {
	lex, err := dslLexer.LexString("", input)
	if err != nil {
		return nil, fmt.Errorf("lex input: %w", err)
	}

	rawTokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, fmt.Errorf("consume tokens: %w", err)
	}

	out := make([]LexToken, 0, len(rawTokens))
	for _, tok := range rawTokens {
		if tok.EOF() {
			continue
		}
		out = append(out, LexToken{
			Name:  symbolName(tok.Type),
			Value: tok.Value,
			Pos:   tok.Pos,
		})
	}

	return out, nil
}
