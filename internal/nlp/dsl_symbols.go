package nlp

import "github.com/alecthomas/participle/v2/lexer"

var (
	//nolint:gochecknoglobals // Used by Parseable implementations.
	dslSymbols = dslLexer.Symbols()

	//nolint:gochecknoglobals // Used for debug/introspection helpers.
	dslSymbolNames = invertSymbolMap(dslSymbols)
)

func invertSymbolMap(symbols map[string]lexer.TokenType) map[lexer.TokenType]string {
	out := make(map[lexer.TokenType]string, len(symbols))
	for name, tt := range symbols {
		out[tt] = name
	}
	return out
}

func symbolName(tt lexer.TokenType) string {
	name, ok := dslSymbolNames[tt]
	if !ok {
		return ""
	}
	return name
}
