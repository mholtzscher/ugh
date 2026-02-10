package nlp

import (
	"errors"
	"strings"

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

type dueShorthandNode struct {
	Value string
}

func (*dueShorthandNode) operation() {}
func (*dueShorthandNode) createOp()  {}

func (d *dueShorthandNode) Parse(lex *lexer.PeekingLexer) error {
	if d == nil {
		return errors.New("nil dueShorthandNode")
	}
	tok := lex.Peek()
	if tok == nil {
		return participle.NextMatch
	}
	if tok.Type != dslSymbols["Ident"] {
		return participle.NextMatch
	}
	s := strings.ToLower(strings.TrimSpace(tok.Value))
	switch s {
	case "today", "tomorrow", "next-week":
		lex.Next()
		d.Value = s
		return nil
	default:
		return participle.NextMatch
	}
}
