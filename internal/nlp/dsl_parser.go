package nlp

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

//nolint:gochecknoglobals // Parser is a package-level singleton for participle.
var dslParser = participle.MustBuild[Root](
	participle.Lexer(dslLexer),
	participle.Elide("Whitespace"),
	participle.Unquote("Quoted"),
	participle.CaseInsensitive("Ident"),
	participle.Map(trimPrefixTokenMapper("#"), "ProjectTag"),
	participle.Map(trimPrefixTokenMapper("@"), "ContextTag"),
	participle.Union[Command](
		&CreateCommand{},
		&UpdateCommand{},
		&FilterCommand{},
		&ViewCommand{},
		&ContextCommand{},
		&LogCommand{},
	),
	participle.Union[CreatePart](
		&CreateOpPart{},
		&CreateText{},
	),
	participle.Union[Operation](
		&SetOp{},
		&AddOp{},
		&RemoveOp{},
		&ClearOp{},
		&tagOpNode{},
	),
	participle.Union[CreateOp](
		&SetOp{},
		&AddOp{},
		&RemoveOp{},
		&ClearOp{},
		&tagOpNode{},
	),
)

func trimPrefixTokenMapper(prefix string) func(lexer.Token) (lexer.Token, error) {
	return func(tok lexer.Token) (lexer.Token, error) {
		tok.Value = strings.TrimPrefix(tok.Value, prefix)
		return tok, nil
	}
}
