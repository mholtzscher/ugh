package shell

import (
	"context"
	"sort"
	"strings"
	"unicode"

	"github.com/chzyer/readline"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

const (
	ansiReset   = "\033[0m"
	ansiBlue    = "\033[34m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"

	identTokenName   = "Ident"
	maxViewTokenArgs = 2
)

func commandSuggestions() []string {
	return []string{
		"add", "create", "new",
		"set", "edit", "update",
		"find", "show", "list", "filter",
		"context", "view", "help", "clear", "quit", "exit",
	}
}

func genericSuggestions() []string {
	return []string{
		"title:", "notes:", "due:", "waiting:", "state:",
		"project:", "projects:", "context:", "contexts:",
		"+project:", "+context:", "-project:", "-context:",
		"!due", "!waiting", "!notes",
		"and", "or", "not", "&&", "||",
		"today", "tomorrow",
	}
}

func stateSuggestions() []string {
	return []string{"state:inbox", "state:now", "state:waiting", "state:later", "state:done"}
}

func dueSuggestions() []string {
	return []string{"due:today", "due:tomorrow"}
}

func viewSuggestions() []string {
	return []string{
		"i", "inbox",
		"n", "now",
		"w", "waiting",
		"l", "later",
		"c", "calendar", "today",
	}
}

type shellCompleter struct {
	listProjects func(context.Context) ([]string, error)
	listContexts func(context.Context) ([]string, error)
}

var _ readline.AutoCompleter = (*shellCompleter)(nil)

func newShellCompleter(svc service.Service) *shellCompleter {
	return &shellCompleter{
		listProjects: func(ctx context.Context) ([]string, error) {
			rows, err := svc.ListProjects(ctx, service.ListTagsRequest{})
			if err != nil {
				return nil, err
			}
			return extractNames(rows), nil
		},
		listContexts: func(ctx context.Context) ([]string, error) {
			rows, err := svc.ListContexts(ctx, service.ListTagsRequest{})
			if err != nil {
				return nil, err
			}
			return extractNames(rows), nil
		},
	}
}

func (c *shellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	if pos < 0 {
		pos = 0
	}
	if pos > len(line) {
		pos = len(line)
	}

	prefixRunes := line[:pos]
	prefix := string(prefixRunes)
	tokens, err := nlp.Lex(prefix)
	if err != nil {
		return nil, 0
	}
	if hasOpenQuote(tokens) {
		return nil, 0
	}

	_, fragment := splitFragment(prefixRunes)
	fragment = strings.TrimSpace(fragment)
	fragmentLower := strings.ToLower(fragment)

	suggestions := c.suggest(tokens, fragment, fragmentLower)
	if len(suggestions) == 0 {
		return nil, 0
	}

	return toCompletionSuffixes(fragment, suggestions), len([]rune(fragment))
}

func (c *shellCompleter) suggest(tokens []nlp.LexToken, fragment string, fragmentLower string) []string {
	nonWhitespace := filterNonWhitespace(tokens)
	if len(nonWhitespace) == 0 {
		return filterCandidates(fragment, commandSuggestions())
	}

	if suggestions, handled := c.viewCommandSuggestions(nonWhitespace, fragment); handled {
		return suggestions
	}

	if suggestions, handled := c.contextCommandSuggestions(nonWhitespace, fragment); handled {
		return suggestions
	}

	if strings.HasPrefix(fragmentLower, "#") {
		return filterCandidates(fragment, prefixed(c.projectNames(), "#"))
	}
	if strings.HasPrefix(fragmentLower, "@") {
		return filterCandidates(fragment, prefixed(c.contextNames(), "@"))
	}

	if strings.HasPrefix(fragmentLower, "state:") {
		return filterCandidates(fragment, stateSuggestions())
	}
	if strings.HasPrefix(fragmentLower, "due:") {
		return filterCandidates(fragment, dueSuggestions())
	}

	if fieldPrefix, valuePrefix, ok := splitFieldValuePrefix(fragmentLower); ok {
		switch fieldPrefix {
		case "project:", "projects:", "+project:", "+projects:", "-project:", "-projects:":
			return filterCandidates(fragment, withValuePrefix(fieldPrefix, valuePrefix, c.projectNames()))
		case "context:", "contexts:", "+context:", "+contexts:", "-context:", "-contexts:":
			return filterCandidates(fragment, withValuePrefix(fieldPrefix, valuePrefix, c.contextNames()))
		}
	}

	if fragment == "" {
		candidates := append([]string{}, genericSuggestions()...)
		candidates = append(candidates, prefixed(c.projectNames(), "#")...)
		candidates = append(candidates, prefixed(c.contextNames(), "@")...)
		return dedupe(candidates)
	}

	candidates := append([]string{}, commandSuggestions()...)
	candidates = append(candidates, genericSuggestions()...)
	return filterCandidates(fragment, dedupe(candidates))
}

func (c *shellCompleter) viewCommandSuggestions(nonWhitespace []nlp.LexToken, fragment string) ([]string, bool) {
	if len(nonWhitespace) == 0 ||
		nonWhitespace[0].Name != identTokenName ||
		!strings.EqualFold(nonWhitespace[0].Value, "view") {
		return nil, false
	}

	if fragment != "" && len(nonWhitespace) == 1 {
		return nil, false
	}

	if fragment == "" {
		if len(nonWhitespace) == 1 {
			return filterCandidates(fragment, viewSuggestions()), true
		}
		return nil, true
	}

	if len(nonWhitespace) > maxViewTokenArgs {
		return nil, true
	}

	return filterCandidates(fragment, viewSuggestions()), true
}

func (c *shellCompleter) contextCommandSuggestions(nonWhitespace []nlp.LexToken, fragment string) ([]string, bool) {
	if fragment != "" {
		return nil, false
	}
	if len(nonWhitespace) == 0 ||
		nonWhitespace[0].Name != identTokenName ||
		!strings.EqualFold(nonWhitespace[0].Value, "context") {
		return nil, false
	}

	candidates := []string{"clear"}
	candidates = append(candidates, prefixed(c.projectNames(), "#")...)
	candidates = append(candidates, prefixed(c.contextNames(), "@")...)
	return filterCandidates(fragment, dedupe(candidates)), true
}

func (c *shellCompleter) projectNames() []string {
	if c.listProjects == nil {
		return nil
	}
	names, err := c.listProjects(context.Background())
	if err != nil {
		return nil
	}
	return sortedUnique(names)
}

func (c *shellCompleter) contextNames() []string {
	if c.listContexts == nil {
		return nil
	}
	names, err := c.listContexts(context.Background())
	if err != nil {
		return nil
	}
	return sortedUnique(names)
}

func extractNames(rows []store.NameCount) []string {
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Name) == "" {
			continue
		}
		names = append(names, row.Name)
	}
	return names
}

func prefixed(values []string, prefix string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, prefix+value)
	}
	return out
}

func withValuePrefix(fieldPrefix string, valuePrefix string, values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if valuePrefix != "" && !strings.HasPrefix(strings.ToLower(value), valuePrefix) {
			continue
		}
		out = append(out, fieldPrefix+value)
	}
	return out
}

func filterCandidates(fragment string, candidates []string) []string {
	if len(candidates) == 0 {
		return nil
	}
	if fragment == "" {
		return dedupe(candidates)
	}

	fragmentLower := strings.ToLower(fragment)
	out := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if strings.HasPrefix(strings.ToLower(candidate), fragmentLower) {
			out = append(out, candidate)
		}
	}
	return dedupe(out)
}

func dedupe(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func sortedUnique(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	copyValues := dedupe(values)
	sort.Strings(copyValues)
	return copyValues
}

func toCompletionSuffixes(fragment string, fullCandidates []string) [][]rune {
	fragmentRunes := []rune(fragment)
	fragLen := len(fragmentRunes)
	out := make([][]rune, 0, len(fullCandidates))
	for _, candidate := range fullCandidates {
		candidateRunes := []rune(candidate)
		if fragLen > len(candidateRunes) {
			continue
		}
		suffix := candidateRunes[fragLen:]
		out = append(out, suffix)
	}
	return out
}

func splitFieldValuePrefix(fragmentLower string) (string, string, bool) {
	idx := strings.Index(fragmentLower, ":")
	if idx < 0 {
		return "", "", false
	}
	field := fragmentLower[:idx+1]
	valuePrefix := ""
	if idx+1 < len(fragmentLower) {
		valuePrefix = fragmentLower[idx+1:]
	}
	return field, valuePrefix, true
}

func splitFragment(prefix []rune) (int, string) {
	i := len(prefix)
	for i > 0 && !isFragmentDelimiter(prefix[i-1]) {
		i--
	}
	return i, string(prefix[i:])
}

func isFragmentDelimiter(r rune) bool {
	return unicode.IsSpace(r) || r == '(' || r == ')' || r == ','
}

func hasOpenQuote(tokens []nlp.LexToken) bool {
	open := 0
	for _, tok := range tokens {
		switch tok.Name {
		case "QuoteStart":
			open++
		case "QuoteEnd":
			if open > 0 {
				open--
			}
		}
	}
	return open > 0
}

func filterNonWhitespace(tokens []nlp.LexToken) []nlp.LexToken {
	out := make([]nlp.LexToken, 0, len(tokens))
	for _, tok := range tokens {
		if tok.Name == "Whitespace" {
			continue
		}
		out = append(out, tok)
	}
	return out
}

type shellPainter struct{}

var _ readline.Painter = (*shellPainter)(nil)

func newShellPainter() *shellPainter {
	return &shellPainter{}
}

func (*shellPainter) Paint(line []rune, _ int) []rune {
	if len(line) == 0 {
		return line
	}

	input := string(line)
	tokens, err := nlp.Lex(input)
	if err != nil {
		return line
	}
	if len(tokens) == 0 {
		return line
	}

	var b strings.Builder
	cursor := 0
	for _, tok := range tokens {
		start := tok.Pos.Offset
		if start < cursor || start > len(input) {
			continue
		}
		end := min(start+len(tok.Value), len(input))

		if cursor < start {
			b.WriteString(input[cursor:start])
		}

		color := colorForToken(tok)
		if color == "" {
			b.WriteString(input[start:end])
		} else {
			b.WriteString(color)
			b.WriteString(input[start:end])
			b.WriteString(ansiReset)
		}

		cursor = end
	}
	if cursor < len(input) {
		b.WriteString(input[cursor:])
	}

	return []rune(b.String())
}

func colorForToken(tok nlp.LexToken) string {
	switch tok.Name {
	case "Quoted", "QuoteStart", "QuoteEnd", "StringText", "StringEscape", "StringBackslash":
		return ansiYellow
	case "ProjectTag", "ProjectTagPrefix":
		return ansiBlue
	case "ContextTag", "ContextTagPrefix":
		return ansiGreen
	case "SetField", "AddField", "RemoveField", "ClearField", "ClearOp", "AddOp", "RemoveOp":
		return ansiMagenta
	case "AndOp", "OrOp":
		return ansiCyan
	case "HashNumber":
		return ansiCyan
	case "Ident":
		lower := strings.ToLower(tok.Value)
		if lower == "and" || lower == "or" || lower == "not" {
			return ansiCyan
		}
		if lower == "add" || lower == "create" || lower == "new" ||
			lower == "set" || lower == "edit" || lower == "update" ||
			lower == "find" || lower == "show" || lower == "list" || lower == "filter" ||
			lower == "view" || lower == "context" {
			return ansiYellow
		}
	}
	return ""
}
