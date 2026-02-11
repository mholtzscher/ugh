package antlr

import (
	"context"
	"errors"
	"fmt"
	"time"

	plexer "github.com/alecthomas/participle/v2/lexer"
	antlrrt "github.com/antlr4-go/antlr/v4"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/antlr/parser"
)

// Parse parses the input string using the ANTLR-generated parser
// and produces the same nlp.ParseResult as the participle-based parser.
func Parse(input string, opts nlp.ParseOptions) (nlp.ParseResult, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}

	// Set up the ANTLR lexer and parser.
	is := antlrrt.NewInputStream(input)
	lexer := parser.NewUghLexer(is)

	errListener := &syntaxError{}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errListener)

	stream := antlrrt.NewCommonTokenStream(lexer, antlrrt.TokenDefaultChannel)
	p := parser.NewUghParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)

	// Parse the input.
	tree := p.Root()

	if errListener.hasErrors() {
		return nlp.ParseResult{
			Intent: nlp.IntentUnknown,
			Diagnostics: []nlp.Diagnostic{{
				Severity: nlp.SeverityError,
				Code:     "E_PARSE",
				Message:  errListener.Error(),
			}},
		}, fmt.Errorf("%s", errListener.Error())
	}

	// Walk the parse tree with the AST builder visitor.
	visitor := &astBuilder{}
	result := tree.Accept(visitor)

	if visitor.err != nil {
		intent := nlp.IntentUnknown
		if result != nil {
			intent = intentForCommand(result)
		}
		return nlp.ParseResult{Intent: intent}, visitor.err
	}

	if result == nil {
		return nlp.ParseResult{Intent: nlp.IntentUnknown}, errors.New("empty parse result")
	}

	cmd, ok := result.(nlp.Command)
	if !ok {
		return nlp.ParseResult{Intent: nlp.IntentUnknown}, errors.New("unexpected parse result type")
	}

	intent := intentForCommand(cmd)

	// Mode enforcement
	want := nlp.IntentUnknown
	switch opts.Mode {
	case nlp.ModeAuto:
	case nlp.ModeCreate:
		want = nlp.IntentCreate
	case nlp.ModeUpdate:
		want = nlp.IntentUpdate
	case nlp.ModeFilter:
		want = nlp.IntentFilter
	case nlp.ModeView:
		want = nlp.IntentView
	case nlp.ModeContext:
		want = nlp.IntentContext
	}
	if want != nlp.IntentUnknown && intent != want {
		return nlp.ParseResult{Intent: intent, Command: cmd}, errors.New("command does not match parse mode")
	}

	return nlp.ParseResult{Intent: intent, Command: cmd}, nil
}

// intentForCommand derives the Intent from a Command.
func intentForCommand(cmd any) nlp.Intent {
	switch cmd.(type) {
	case *nlp.CreateCommand:
		return nlp.IntentCreate
	case *nlp.UpdateCommand:
		return nlp.IntentUpdate
	case *nlp.FilterCommand:
		return nlp.IntentFilter
	case *nlp.ViewCommand:
		return nlp.IntentView
	case *nlp.ContextCommand:
		return nlp.IntentContext
	default:
		return nlp.IntentUnknown
	}
}

// Parser implements the nlp.Parser interface using the ANTLR backend.
type Parser struct{}

// NewParser creates a new ANTLR-based Parser.
func NewParser() nlp.Parser {
	return Parser{}
}

// Parse implements the nlp.Parser interface.
func (Parser) Parse(_ context.Context, input string, opts nlp.ParseOptions) (nlp.ParseResult, error) {
	return Parse(input, opts)
}

// Lex tokenizes input using the ANTLR lexer and returns LexTokens
// compatible with the existing nlp.LexToken type used for syntax highlighting.
func Lex(input string) ([]nlp.LexToken, error) {
	is := antlrrt.NewInputStream(input)
	lexer := parser.NewUghLexer(is)

	errListener := &syntaxError{}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errListener)

	tokens := make([]nlp.LexToken, 0)
	for {
		tok := lexer.NextToken()
		if tok.GetTokenType() == antlrrt.TokenEOF {
			break
		}
		name := antlrTokenName(tok.GetTokenType())
		if name == "WS" {
			continue
		}
		tokens = append(tokens, nlp.LexToken{
			Name:  name,
			Value: tok.GetText(),
			Pos: plexer.Position{
				Offset: tok.GetStart(),
				Line:   tok.GetLine(),
				Column: tok.GetColumn(),
			},
		})
	}

	if errListener.hasErrors() {
		return tokens, fmt.Errorf("lex errors: %s", errListener.Error())
	}

	return tokens, nil
}

// antlrTokenName maps ANTLR token type IDs to symbolic names matching the
// participle lexer's token names (for compatibility with the syntax highlighter).
func antlrTokenName(tokenType int) string {
	// The symbolic names are available from the lexer, but for direct mapping
	// we use the same names as the participle lexer.
	switch tokenType {
	case parser.UghLexerQUOTED:
		return "Quoted"
	case parser.UghLexerHASH_NUMBER:
		return "HashNumber"
	case parser.UghLexerPROJECT_TAG:
		return "ProjectTag"
	case parser.UghLexerCONTEXT_TAG:
		return "ContextTag"
	case parser.UghLexerPROJECT_TAG_PREFIX:
		return "ProjectTagPrefix"
	case parser.UghLexerCONTEXT_TAG_PREFIX:
		return "ContextTagPrefix"
	case parser.UghLexerSET_FIELD:
		return "SetField"
	case parser.UghLexerADD_FIELD:
		return "AddField"
	case parser.UghLexerREMOVE_FIELD:
		return "RemoveField"
	case parser.UghLexerCLEAR_FIELD:
		return "ClearField"
	case parser.UghLexerCLEAR_OP:
		return "ClearOp"
	case parser.UghLexerADD_OP:
		return "AddOp"
	case parser.UghLexerREMOVE_OP:
		return "RemoveOp"
	case parser.UghLexerCOLON:
		return "Colon"
	case parser.UghLexerCOMMA:
		return "Comma"
	case parser.UghLexerLPAREN:
		return "LParen"
	case parser.UghLexerRPAREN:
		return "RParen"
	case parser.UghLexerAND_OP:
		return "AndOp"
	case parser.UghLexerOR_OP:
		return "OrOp"
	case parser.UghLexerIDENT:
		return "Ident"
	// Keyword tokens map to "Ident" for compatibility with the participle
	// lexer's syntax highlighting (which treats all words as Ident).
	case parser.UghLexerKW_ADD, parser.UghLexerKW_CREATE, parser.UghLexerKW_NEW,
		parser.UghLexerKW_SET, parser.UghLexerKW_EDIT, parser.UghLexerKW_UPDATE,
		parser.UghLexerKW_FIND, parser.UghLexerKW_SHOW, parser.UghLexerKW_LIST,
		parser.UghLexerKW_FILTER, parser.UghLexerKW_VIEW, parser.UghLexerKW_CONTEXT,
		parser.UghLexerKW_AND, parser.UghLexerKW_OR, parser.UghLexerKW_NOT:
		return "Ident"
	case parser.UghLexerWS:
		return "Whitespace"
	default:
		return fmt.Sprintf("Unknown(%d)", tokenType)
	}
}
