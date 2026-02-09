package nlp

import "github.com/alecthomas/participle/v2/lexer"

//nolint:gochecknoglobals // Used by Parseable implementations.
var dslSymbols map[string]lexer.TokenType = dslLexer.Symbols()
