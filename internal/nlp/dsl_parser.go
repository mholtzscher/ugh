package nlp

import (
	"github.com/alecthomas/participle/v2"
)

//nolint:gochecknoglobals // Parser is a package-level singleton for participle.
var dslParser = participle.MustBuild[Root](
	participle.Lexer(dslLexer),
	participle.Elide("Whitespace"),
	participle.Unquote("Quoted"),
	participle.Union[Command](
		&CreateCommand{},
		&UpdateCommand{},
		&FilterCommand{},
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
		&dueShorthandNode{},
	),
)
