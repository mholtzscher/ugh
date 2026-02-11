package antlr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	antlrrt "github.com/antlr4-go/antlr/v4"

	"github.com/mholtzscher/ugh/internal/nlp"
)

// ─── Error listener ─────────────────────────────────────────────────────────

// syntaxError captures syntax errors from the ANTLR lexer/parser.
type syntaxError struct {
	antlrrt.DefaultErrorListener

	errors []string
}

func (l *syntaxError) SyntaxError(
	_ antlrrt.Recognizer,
	_ any,
	line, column int,
	msg string,
	_ antlrrt.RecognitionException,
) {
	l.errors = append(l.errors, fmt.Sprintf("%d:%d: %s", line, column, msg))
}

func (l *syntaxError) hasErrors() bool {
	return len(l.errors) > 0
}

func (l *syntaxError) Error() string {
	return strings.Join(l.errors, "; ")
}

// ─── View name canonicalization ─────────────────────────────────────────────

//nolint:goconst // these are short literal values in a switch, not reusable constants.
func canonicalViewName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "i", "inbox":
		return "inbox"
	case "n", "now":
		return "now"
	case "w", "waiting":
		return "waiting"
	case "l", "later":
		return "later"
	case "c", "calendar", "today":
		return "calendar"
	default:
		return ""
	}
}

// ─── Field capture ──────────────────────────────────────────────────────────

// captureField parses a raw field token (e.g. "title:", "+project:", "-context:") into a Field.
func captureField(raw string) (nlp.Field, error) {
	name := normalizeFieldName(raw)
	return fieldFromName(name)
}

// captureClearField parses a raw clear field token (e.g. "!due", "! notes") into a Field.
func captureClearField(raw string) (nlp.Field, error) {
	name := normalizeFieldName(raw)
	return fieldFromName(name)
}

// normalizeFieldName strips prefixes (+, -, !), suffixes (:), and whitespace.
func normalizeFieldName(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "+")
	s = strings.TrimPrefix(s, "-")
	s = strings.TrimPrefix(s, "!")
	s = strings.TrimSuffix(s, ":")
	s = strings.TrimSpace(s)
	return strings.ToLower(s)
}

// fieldFromName converts a normalized field name string to a nlp.Field.
func fieldFromName(name string) (nlp.Field, error) {
	switch name {
	case "title":
		return nlp.FieldTitle, nil
	case "notes":
		return nlp.FieldNotes, nil
	case "due":
		return nlp.FieldDue, nil
	case "waiting", "waiting-for", "waiting_for":
		return nlp.FieldWaiting, nil
	case "state":
		return nlp.FieldState, nil
	case "project", "projects":
		return nlp.FieldProjects, nil
	case "context", "contexts":
		return nlp.FieldContexts, nil
	case "meta":
		return nlp.FieldMeta, nil
	case "id":
		return 0, errors.New("id cannot be set directly")
	case "text":
		return 0, errors.New("text is not a settable field")
	default:
		return 0, fmt.Errorf("unknown field %q", name)
	}
}

// ─── Token joining ──────────────────────────────────────────────────────────

// joinTokens joins token strings with minimal normalization.
// Colons and commas are joined without spaces; other tokens get spaces between them.
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

// ─── String utilities ───────────────────────────────────────────────────────

// unquote removes surrounding double quotes and unescapes.
func unquote(s string) string {
	if len(s) < 2 { //nolint:mnd // minimum length for a quoted string (open + close quote)
		return s
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	// Unescape backslash sequences
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
}

func isDigits(s string) bool {
	return s != "" && !strings.ContainsFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	})
}

func parsePossibleID(value string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if rest, ok := strings.CutPrefix(value, "#"); ok {
		if !isDigits(rest) {
			return 0, false
		}
		id, err := strconv.ParseInt(rest, 10, 64)
		return id, err == nil && id > 0
	}
	if !isDigits(value) {
		return 0, false
	}
	id, err := strconv.ParseInt(value, 10, 64)
	return id, err == nil && id > 0
}

func hasTitleSetOp(ops []nlp.Operation) bool {
	for _, op := range ops {
		setOp, ok := op.(nlp.SetOp)
		if ok && setOp.Field == nlp.FieldTitle {
			return true
		}
	}
	return false
}
