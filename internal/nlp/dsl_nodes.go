package nlp

import (
	"errors"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type tagOpNode struct {
	Kind  TagKind
	Value string
}

func (*tagOpNode) operation() {}
func (*tagOpNode) createOp()  {}

func (t *tagOpNode) Parse(lex *lexer.PeekingLexer) error {
	if t == nil {
		return errors.New("nil tagOpNode")
	}
	tok := lex.Peek()
	if tok == nil {
		return participle.NextMatch
	}
	if tok.Type == dslSymbols["ProjectTag"] {
		lex.Next()
		t.Kind = TagProject
		t.Value = tok.Value
		return nil
	}
	if tok.Type == dslSymbols["ContextTag"] {
		lex.Next()
		t.Kind = TagContext
		t.Value = tok.Value
		return nil
	}
	return participle.NextMatch
}
