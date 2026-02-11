package nlp

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	contextCommandVerb = "context"
)

func parseIdent(lex *lexer.PeekingLexer) (string, error) {
	tok := lex.Peek()
	if tok == nil {
		return "", participle.NextMatch
	}
	if tok.Type != dslSymbols["Ident"] {
		return "", participle.NextMatch
	}
	return strings.ToLower(tok.Value), nil
}

func parseVerb(lex *lexer.PeekingLexer, allowed []string) (string, error) {
	s, err := parseIdent(lex)
	if err != nil {
		return "", err
	}
	if slices.Contains(allowed, s) {
		lex.Next()
		return s, nil
	}
	return "", participle.NextMatch
}

func matchToken(lex *lexer.PeekingLexer, allowed []string) error {
	tok := lex.Peek()
	if tok == nil {
		return participle.NextMatch
	}
	s := strings.ToLower(strings.TrimSpace(tok.Value))
	if slices.Contains(allowed, s) {
		lex.Next()
		return nil
	}
	return participle.NextMatch
}

//nolint:gochecknoglobals // constant lookup table for verb synonyms
var createVerbs = []string{"add", "create", "new"}

func (v *CreateVerb) Parse(lex *lexer.PeekingLexer) error {
	if v == nil {
		return errors.New("nil CreateVerb")
	}
	s, err := parseVerb(lex, createVerbs)
	if err != nil {
		return err
	}
	*v = CreateVerb(s)
	return nil
}

//nolint:gochecknoglobals // constant lookup table for verb synonyms
var updateVerbs = []string{"set", "edit", "update"}

func (v *UpdateVerb) Parse(lex *lexer.PeekingLexer) error {
	if v == nil {
		return errors.New("nil UpdateVerb")
	}
	s, err := parseVerb(lex, updateVerbs)
	if err != nil {
		return err
	}
	*v = UpdateVerb(s)
	return nil
}

//nolint:gochecknoglobals // constant lookup table for verb synonyms
var filterVerbs = []string{"find", "show", "list", "filter"}

func (v *FilterVerb) Parse(lex *lexer.PeekingLexer) error {
	if v == nil {
		return errors.New("nil FilterVerb")
	}
	s, err := parseVerb(lex, filterVerbs)
	if err != nil {
		return err
	}
	*v = FilterVerb(s)
	return nil
}

//nolint:gochecknoglobals // constant lookup table for verb synonyms
var viewVerbs = []string{"view"}

func (v *ViewVerb) Parse(lex *lexer.PeekingLexer) error {
	if v == nil {
		return errors.New("nil ViewVerb")
	}
	s, err := parseVerb(lex, viewVerbs)
	if err != nil {
		return err
	}
	*v = ViewVerb(s)
	return nil
}

//nolint:gochecknoglobals // constant lookup table for verb synonyms
var contextVerbs = []string{contextCommandVerb}

func (v *ContextVerb) Parse(lex *lexer.PeekingLexer) error {
	if v == nil {
		return errors.New("nil ContextVerb")
	}
	s, err := parseVerb(lex, contextVerbs)
	if err != nil {
		return err
	}
	*v = ContextVerb(s)
	return nil
}

func (t *ViewTarget) Parse(lex *lexer.PeekingLexer) error {
	if t == nil {
		return errors.New("nil ViewTarget")
	}
	tok := lex.Peek()
	if tok == nil {
		return participle.NextMatch
	}
	if tok.Type != dslSymbols["Ident"] {
		return fmt.Errorf("invalid view: %s", tok.Value)
	}

	s := strings.ToLower(strings.TrimSpace(tok.Value))
	switch s {
	case "i", viewNameInbox, "n", viewNameNow, "w", viewNameWaiting, "l", viewNameLater, "c", viewNameCalendar, "today":
		lex.Next()
		t.Name = s
		return nil
	default:
		return fmt.Errorf("invalid view: %s", tok.Value)
	}
}

func (a *ContextArg) Parse(lex *lexer.PeekingLexer) error {
	if a == nil {
		return errors.New("nil ContextArg")
	}
	tok := lex.Peek()
	if tok == nil {
		return participle.NextMatch
	}

	if tok.Type == dslSymbols["ProjectTag"] {
		lex.Next()
		a.Project = tok.Value
		return nil
	}
	if tok.Type == dslSymbols["ContextTag"] {
		lex.Next()
		a.Context = tok.Value
		return nil
	}
	if tok.Type == dslSymbols["Ident"] {
		s := strings.ToLower(strings.TrimSpace(tok.Value))
		if s == "clear" {
			lex.Next()
			a.Clear = true
			return nil
		}
	}

	return fmt.Errorf("invalid context argument: %s", tok.Value)
}

func (t *TargetRef) Parse(lex *lexer.PeekingLexer) error {
	if t == nil {
		return errors.New("nil TargetRef")
	}
	tok := lex.Peek()
	if tok == nil {
		return participle.NextMatch
	}
	// If the next token is not a plausible target, treat it as absent.
	if tok.Type != dslSymbols["Ident"] && tok.Type != dslSymbols["HashNumber"] {
		return participle.NextMatch
	}

	text := strings.ToLower(strings.TrimSpace(tok.Value))
	if tok.Type == dslSymbols["HashNumber"] {
		id, err := strconv.ParseInt(strings.TrimPrefix(text, "#"), 10, 64)
		if err != nil || id <= 0 {
			return fmt.Errorf("invalid update target: %s", tok.Value)
		}
		lex.Next()
		t.Kind = TargetID
		t.ID = id
		return nil
	}

	switch text {
	case "selected", "it", "this", "that":
		lex.Next()
		t.Kind = TargetSelected
		t.ID = 0
		return nil
	}

	if isDigits(text) {
		id, err := strconv.ParseInt(text, 10, 64)
		if err != nil || id <= 0 {
			return fmt.Errorf("invalid update target: %s", tok.Value)
		}
		lex.Next()
		t.Kind = TargetID
		t.ID = id
		return nil
	}

	return fmt.Errorf("invalid update target: %s", tok.Value)
}

//nolint:gochecknoglobals // constant lookup table for operator symbols and keywords
var orTokens = []string{"||", "or"}

func (o *OrOperator) Parse(lex *lexer.PeekingLexer) error {
	if o == nil {
		return errors.New("nil OrOperator")
	}
	return matchToken(lex, orTokens)
}

//nolint:gochecknoglobals // constant lookup table for operator symbols and keywords
var andTokens = []string{"&&", "and"}

func (a *AndOperator) Parse(lex *lexer.PeekingLexer) error {
	if a == nil {
		return errors.New("nil AndOperator")
	}
	return matchToken(lex, andTokens)
}

//nolint:gochecknoglobals // constant lookup table for operator symbols and keywords
var notTokens = []string{"!", "not"}

func (n *NotOperator) Parse(lex *lexer.PeekingLexer) error {
	if n == nil {
		return errors.New("nil NotOperator")
	}
	return matchToken(lex, notTokens)
}

func (f *Field) Capture(values []string) error {
	name := normalizeCapturedField(values)
	switch name {
	case "title":
		*f = FieldTitle
		return nil
	case "notes":
		*f = FieldNotes
		return nil
	case "due":
		*f = FieldDue
		return nil
	case "waiting", "waiting-for", "waiting_for":
		*f = FieldWaiting
		return nil
	case "state":
		*f = FieldState
		return nil
	case "project", "projects":
		*f = FieldProjects
		return nil
	case "context", "contexts":
		*f = FieldContexts
		return nil
	case "meta":
		*f = FieldMeta
		return nil
	case "id":
		return errors.New("id cannot be set directly")
	case "text":
		return errors.New("text is not a settable field")
	default:
		return fmt.Errorf("unknown field %q", name)
	}
}

func (v *OpValue) Capture(values []string) error {
	*v = OpValue(joinTokens(values))
	return nil
}

func (v *FilterValue) Parse(lex *lexer.PeekingLexer) error {
	if v == nil {
		return errors.New("nil FilterValue")
	}
	peek := lex.Peek()
	if peek == nil {
		return participle.NextMatch
	}

	if peek.Type == dslSymbols["Quoted"] {
		lex.Next()
		*v = FilterValue(peek.Value)
		return nil
	}

	values := make([]string, 0)
	for {
		tok := lex.Peek()
		if tok == nil || isFilterValueDelimiter(tok) {
			break
		}
		if !isFilterValueToken(tok) {
			break
		}
		lex.Next()
		values = append(values, tok.Value)
	}

	if len(values) == 0 {
		return errors.New("expected value")
	}
	*v = FilterValue(joinTokens(values))
	return nil
}

func isFilterValueDelimiter(tok *lexer.Token) bool {
	if tok.Type == dslSymbols["RParen"] || tok.Type == dslSymbols["AndOp"] || tok.Type == dslSymbols["OrOp"] {
		return true
	}
	if tok.Type == dslSymbols["Ident"] {
		s := strings.ToLower(tok.Value)
		return s == "and" || s == "or" || s == "not"
	}
	return false
}

func isFilterValueToken(tok *lexer.Token) bool {
	return tok.Type == dslSymbols["Ident"] ||
		tok.Type == dslSymbols["HashNumber"] ||
		tok.Type == dslSymbols["Colon"] ||
		tok.Type == dslSymbols["Comma"]
}

func normalizeCapturedField(values []string) string {
	joined := strings.Join(values, "")
	joined = strings.TrimSpace(joined)
	joined = strings.TrimPrefix(joined, "+")
	joined = strings.TrimPrefix(joined, "-")
	joined = strings.TrimPrefix(joined, "!")
	joined = strings.TrimSuffix(joined, ":")
	joined = strings.TrimSpace(joined)
	return strings.ToLower(joined)
}

func joinTokens(values []string) string {
	values = trimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	var b strings.Builder
	for i, tok := range values {
		if i == 0 {
			b.WriteString(tok)
			continue
		}
		if tok == ":" || tok == "," {
			b.WriteString(tok)
			continue
		}
		prev := values[i-1]
		if prev == ":" || prev == "," {
			b.WriteString(tok)
			continue
		}
		b.WriteByte(' ')
		b.WriteString(tok)
	}

	return strings.TrimSpace(b.String())
}

func trimEmpty(in []string) []string {
	out := in[:0]
	for _, s := range in {
		if strings.TrimSpace(s) == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
