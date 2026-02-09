package nlp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/domain"
)

type parser struct {
	tokens []token
	idx    int
	opts   ParseOptions
}

const (
	minTagTokenLength   = 2
	predicateStateField = "state"
	predicateDueField   = "due"
	predicateProject    = "project"
	predicateContext    = "context"
	predicateText       = "text"
)

func parseInput(input string, opts ParseOptions) (ParseResult, error) {
	tokens, err := lex(input)
	if err != nil {
		return ParseResult{
			Diagnostics: []Diagnostic{{
				Severity: SeverityError,
				Code:     "E_LEX",
				Message:  err.Error(),
			}},
		}, err
	}

	p := parser{tokens: tokens, opts: opts}
	mode, err := p.resolveMode()
	if err != nil {
		return ParseResult{
			Intent: IntentUnknown,
			Diagnostics: []Diagnostic{{
				Severity: SeverityError,
				Code:     "E_CMD",
				Message:  err.Error(),
			}},
		}, err
	}

	switch mode {
	case ModeCreate:
		cmd, parseErr := p.parseCreateCommand()
		if parseErr != nil {
			return ParseResult{Intent: IntentCreate}, parseErr
		}
		return ParseResult{Intent: IntentCreate, Command: cmd}, nil
	case ModeUpdate:
		cmd, parseErr := p.parseUpdateCommand()
		if parseErr != nil {
			return ParseResult{Intent: IntentUpdate}, parseErr
		}
		return ParseResult{Intent: IntentUpdate, Command: cmd}, nil
	case ModeFilter:
		cmd, parseErr := p.parseFilterCommand()
		if parseErr != nil {
			return ParseResult{Intent: IntentFilter}, parseErr
		}
		return ParseResult{Intent: IntentFilter, Command: cmd}, nil
	case ModeAuto:
		return ParseResult{}, errors.New("mode not resolved")
	default:
		return ParseResult{}, errors.New("unknown parse mode")
	}
}

func (p *parser) resolveMode() (Mode, error) {
	if p.opts.Mode != ModeAuto {
		return p.opts.Mode, nil
	}
	tok := p.current()
	if tok.kind != tokenWord {
		return ModeAuto, errors.New(
			"expected command verb (add, create, new, set, edit, update, find, show, list, filter)",
		)
	}
	word := strings.ToLower(tok.text)
	if isUpdateVerb(word) {
		return ModeUpdate, nil
	}
	if isFilterVerb(word) {
		return ModeFilter, nil
	}
	if isCreateVerb(word) {
		return ModeCreate, nil
	}
	return ModeAuto, fmt.Errorf(
		"unknown command verb: %q (expected: add, create, new, set, edit, update, find, show, list, filter)",
		word,
	)
}

func (p *parser) parseCreateCommand() (CreateCommand, error) {
	start := p.current().span.Start
	if p.current().kind == tokenWord && isCreateVerb(strings.ToLower(p.current().text)) {
		p.advance()
	}

	titleTokens := make([]token, 0)
	ops := make([]Operation, 0)

	for !p.atEOF() {
		if p.isOpStart() {
			op, err := p.parseOperation()
			if err != nil {
				return CreateCommand{}, err
			}
			ops = append(ops, op)
			continue
		}

		if p.current().kind == tokenWord && isRelativeDate(strings.ToLower(p.current().text)) {
			tok := p.current()
			ops = append(ops, SetOp{
				Field: FieldDue,
				Value: Value{Raw: tok.text, Quoted: false, Span: tok.span},
				Span:  tok.span,
			})
			p.advance()
			continue
		}

		titleTokens = append(titleTokens, p.current())
		p.advance()
	}

	title, titleSpan := joinTokens(titleTokens)
	cmd := CreateCommand{
		Title: title,
		Ops:   ops,
		Span:  Span{Start: start, End: p.current().span.End},
	}
	if strings.TrimSpace(title) == "" && !hasFieldOp(ops, FieldTitle) {
		return CreateCommand{}, errors.New("create command requires title or title: field")
	}
	if title != "" && titleSpan.Start.Offset != 0 {
		cmd.Span.Start = titleSpan.Start
	}
	return cmd, nil
}

func (p *parser) parseUpdateCommand() (UpdateCommand, error) {
	start := p.current().span.Start
	if p.current().kind == tokenWord && isUpdateVerb(strings.ToLower(p.current().text)) {
		p.advance()
	}

	target, err := p.parseTarget()
	if err != nil {
		return UpdateCommand{}, err
	}

	ops := make([]Operation, 0)
	for !p.atEOF() {
		if !p.isOpStart() {
			return UpdateCommand{}, fmt.Errorf("unexpected token %q in update command", p.current().text)
		}
		op, opErr := p.parseOperation()
		if opErr != nil {
			return UpdateCommand{}, opErr
		}
		ops = append(ops, op)
	}

	if len(ops) == 0 {
		return UpdateCommand{}, errors.New("update command requires at least one operation")
	}

	return UpdateCommand{
		Target: target,
		Ops:    ops,
		Span:   Span{Start: start, End: p.current().span.End},
	}, nil
}

func (p *parser) parseFilterCommand() (FilterCommand, error) {
	start := p.current().span.Start
	if p.current().kind == tokenWord && isFilterVerb(strings.ToLower(p.current().text)) {
		p.advance()
	}

	if p.atEOF() {
		return FilterCommand{}, errors.New("filter command requires an expression")
	}

	expr, err := p.parseFilterOrExpr()
	if err != nil {
		return FilterCommand{}, err
	}
	if !p.atEOF() {
		return FilterCommand{}, fmt.Errorf("unexpected token %q in filter expression", p.current().text)
	}

	return FilterCommand{Expr: expr, Span: Span{Start: start, End: p.current().span.End}}, nil
}

func (p *parser) parseTarget() (TargetRef, error) {
	tok := p.current()
	if tok.kind != tokenWord {
		return TargetRef{}, fmt.Errorf("expected update target, got %q", tok.text)
	}

	lower := strings.ToLower(tok.text)
	if lower == "selected" {
		p.advance()
		return TargetRef{Kind: TargetSelected, Span: tok.span}, nil
	}

	if trimmed, ok := strings.CutPrefix(lower, "#"); ok {
		id, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return TargetRef{}, fmt.Errorf("invalid task id target %q", tok.text)
		}
		p.advance()
		return TargetRef{Kind: TargetID, ID: id, Span: tok.span}, nil
	}

	if lower == "id" && p.peek().kind == tokenColon {
		p.advance()
		p.advance()
		if p.current().kind != tokenWord && p.current().kind != tokenQuoted {
			return TargetRef{}, errors.New("expected numeric id after id")
		}
		id, err := strconv.ParseInt(p.current().text, 10, 64)
		if err != nil {
			return TargetRef{}, fmt.Errorf("invalid task id %q", p.current().text)
		}
		span := Span{Start: tok.span.Start, End: p.current().span.End}
		p.advance()
		return TargetRef{Kind: TargetID, ID: id, Span: span}, nil
	}

	return TargetRef{}, fmt.Errorf("invalid update target %q", tok.text)
}

func (p *parser) parseOperation() (Operation, error) {
	tok := p.current()
	switch tok.kind {
	case tokenPlus:
		return p.parseAddRemoveOperation(true)
	case tokenMinus:
		return p.parseAddRemoveOperation(false)
	case tokenBang:
		return p.parseClearOperation()
	case tokenWord:
		if strings.HasPrefix(tok.text, "#") || strings.HasPrefix(tok.text, "@") {
			return p.parseTagOperation()
		}
		if p.peek().kind == tokenColon {
			return p.parseSetOperation()
		}
	case tokenQuoted, tokenColon, tokenLParen, tokenRParen, tokenEOF:
		return nil, fmt.Errorf("expected operation at %q", tok.text)
	}
	return nil, fmt.Errorf("expected operation at %q", tok.text)
}

func (p *parser) parseTagOperation() (Operation, error) {
	tok := p.current()
	if tok.kind != tokenWord {
		return nil, errors.New("expected tag token")
	}
	if len(tok.text) < minTagTokenLength {
		return nil, fmt.Errorf("invalid tag token %q", tok.text)
	}

	value := tok.text[1:]
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("invalid tag token %q", tok.text)
	}

	p.advance()
	if tok.text[0] == '#' {
		return TagOp{Kind: TagProject, Value: value, Span: tok.span}, nil
	}
	return TagOp{Kind: TagContext, Value: value, Span: tok.span}, nil
}

func (p *parser) parseSetOperation() (Operation, error) {
	fieldToken := p.current()
	field, err := parseFieldName(fieldToken.text)
	if err != nil {
		return nil, err
	}
	p.advance()
	if p.current().kind != tokenColon {
		return nil, errors.New("expected : after field")
	}
	p.advance()

	value, span := p.parseValueUntilOpBoundary()
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("empty value for %s", fieldToken.text)
	}

	return SetOp{
		Field: field,
		Value: Value{Raw: value, Quoted: false, Span: span},
		Span:  Span{Start: fieldToken.span.Start, End: span.End},
	}, nil
}

func (p *parser) parseAddRemoveOperation(add bool) (Operation, error) {
	opToken := p.current()
	p.advance()
	if p.current().kind != tokenWord {
		return nil, fmt.Errorf("expected field name after %s", opToken.text)
	}

	fieldToken := p.current()
	field, err := parseFieldName(fieldToken.text)
	if err != nil {
		return nil, err
	}
	if !isListField(field) {
		return nil, fmt.Errorf("%s only supports list fields", opToken.text)
	}
	p.advance()

	if p.current().kind != tokenColon {
		return nil, fmt.Errorf("expected : after %s", fieldToken.text)
	}
	p.advance()

	value, span := p.parseValueUntilOpBoundary()
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("empty value for %s", fieldToken.text)
	}

	if add {
		return AddOp{
			Field: field,
			Value: Value{Raw: value, Quoted: false, Span: span},
			Span:  Span{Start: opToken.span.Start, End: span.End},
		}, nil
	}

	return RemoveOp{
		Field: field,
		Value: Value{Raw: value, Quoted: false, Span: span},
		Span:  Span{Start: opToken.span.Start, End: span.End},
	}, nil
}

func (p *parser) parseClearOperation() (Operation, error) {
	bang := p.current()
	p.advance()
	if p.current().kind != tokenWord {
		return nil, errors.New("expected clear field after bang")
	}
	field, err := parseFieldName(p.current().text)
	if err != nil {
		return nil, err
	}
	if !isClearableField(field) {
		return nil, fmt.Errorf("%q is not clearable", p.current().text)
	}
	fieldSpan := p.current().span
	p.advance()
	return ClearOp{Field: field, Span: Span{Start: bang.span.Start, End: fieldSpan.End}}, nil
}

func (p *parser) parseValueUntilOpBoundary() (string, Span) {
	parts := make([]token, 0)
	for !p.atEOF() {
		if p.isOpStart() {
			break
		}
		parts = append(parts, p.current())
		p.advance()
	}
	value, span := joinTokens(parts)
	return strings.TrimSpace(value), span
}

func (p *parser) parseFilterOrExpr() (FilterExpr, error) {
	left, err := p.parseFilterAndExpr()
	if err != nil {
		return nil, err
	}

	for isLogicalOr(p.current()) {
		p.advance()
		right, parseErr := p.parseFilterAndExpr()
		if parseErr != nil {
			return nil, parseErr
		}
		left = FilterBinary{
			Op:    FilterOr,
			Left:  left,
			Right: right,
			Span:  Span{Start: left.NodeSpan().Start, End: right.NodeSpan().End},
		}
	}
	return left, nil
}

func (p *parser) parseFilterAndExpr() (FilterExpr, error) {
	left, err := p.parseFilterNotExpr()
	if err != nil {
		return nil, err
	}

	for isLogicalAnd(p.current()) {
		p.advance()
		right, parseErr := p.parseFilterNotExpr()
		if parseErr != nil {
			return nil, parseErr
		}
		left = FilterBinary{
			Op:    FilterAnd,
			Left:  left,
			Right: right,
			Span:  Span{Start: left.NodeSpan().Start, End: right.NodeSpan().End},
		}
	}
	return left, nil
}

func (p *parser) parseFilterNotExpr() (FilterExpr, error) {
	if p.current().kind == tokenBang || tokenIsWord(p.current(), "not") {
		start := p.current().span.Start
		p.advance()
		expr, err := p.parseFilterNotExpr()
		if err != nil {
			return nil, err
		}
		return FilterNot{Expr: expr, Span: Span{Start: start, End: expr.NodeSpan().End}}, nil
	}
	return p.parseFilterAtom()
}

func (p *parser) parseFilterAtom() (FilterExpr, error) {
	if p.current().kind == tokenLParen {
		start := p.current().span.Start
		p.advance()
		expr, err := p.parseFilterOrExpr()
		if err != nil {
			return nil, err
		}
		if p.current().kind != tokenRParen {
			return nil, errors.New("expected )")
		}
		end := p.current().span.End
		p.advance()
		switch node := expr.(type) {
		case Predicate:
			node.Span = Span{Start: start, End: end}
			return node, nil
		default:
			return expr, nil
		}
	}
	return p.parseFilterPredicate()
}

func (p *parser) parseFilterPredicate() (FilterExpr, error) {
	tok := p.current()
	if tok.kind == tokenWord && p.peek().kind == tokenColon {
		field := strings.ToLower(tok.text)
		if isPredicateField(field) {
			p.advance()
			p.advance()
			value, span := p.parseValueUntilFilterBoundary()
			if value == "" {
				return nil, fmt.Errorf("missing predicate value for %s", field)
			}
			return parsePredicate(field, value, Span{Start: tok.span.Start, End: span.End})
		}
	}

	value, span := p.parseValueUntilFilterBoundary()
	if value == "" {
		return nil, errors.New("expected predicate")
	}
	return Predicate{Kind: PredText, Text: value, Span: span}, nil
}

func (p *parser) parseValueUntilFilterBoundary() (string, Span) {
	parts := make([]token, 0)
	for !p.atEOF() {
		if p.current().kind == tokenRParen || isLogicalAnd(p.current()) || isLogicalOr(p.current()) {
			break
		}
		parts = append(parts, p.current())
		p.advance()
	}
	value, span := joinTokens(parts)
	return strings.TrimSpace(value), span
}

func parsePredicate(field string, raw string, span Span) (Predicate, error) {
	value := strings.TrimSpace(raw)
	switch field {
	case predicateStateField:
		state, err := normalizeStateToken(value)
		if err != nil {
			return Predicate{}, err
		}
		return Predicate{Kind: PredState, Text: state, Span: span}, nil
	case predicateProject:
		return Predicate{Kind: PredProject, Text: value, Span: span}, nil
	case predicateContext:
		return Predicate{Kind: PredContext, Text: value, Span: span}, nil
	case predicateText:
		return Predicate{Kind: PredText, Text: value, Span: span}, nil
	case predicateDueField:
		cmp, dateValue := splitDateCmp(value)
		kind, text, err := parseDateValue(dateValue)
		if err != nil {
			return Predicate{}, err
		}
		return Predicate{Kind: PredDue, DateCmp: cmp, DateKind: kind, DateText: text, Span: span}, nil
	default:
		return Predicate{}, fmt.Errorf("unsupported predicate field %q", field)
	}
}

func splitDateCmp(value string) (DateCmp, string) {
	switch {
	case strings.HasPrefix(value, ">="):
		return DateGTE, strings.TrimSpace(strings.TrimPrefix(value, ">="))
	case strings.HasPrefix(value, "<="):
		return DateLTE, strings.TrimSpace(strings.TrimPrefix(value, "<="))
	case strings.HasPrefix(value, ">"):
		return DateGT, strings.TrimSpace(strings.TrimPrefix(value, ">"))
	case strings.HasPrefix(value, "<"):
		return DateLT, strings.TrimSpace(strings.TrimPrefix(value, "<"))
	default:
		return DateEq, strings.TrimSpace(value)
	}
}

func parseDateValue(value string) (DateValueKind, string, error) {
	lower := strings.ToLower(strings.TrimSpace(value))
	switch lower {
	case "today":
		return DateToday, "", nil
	case "tomorrow":
		return DateTomorrow, "", nil
	case "next-week":
		return DateNextWeek, "", nil
	}
	if _, err := time.Parse(domain.DateLayoutYYYYMMDD, lower); err != nil {
		return 0, "", domain.InvalidDateFormatError(value)
	}
	return DateAbsolute, lower, nil
}

func parseFieldName(value string) (Field, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "title":
		return FieldTitle, nil
	case "notes":
		return FieldNotes, nil
	case "due":
		return FieldDue, nil
	case "waiting", "waiting-for", "waiting_for":
		return FieldWaiting, nil
	case "state":
		return FieldState, nil
	case "project", "projects":
		return FieldProjects, nil
	case "context", "contexts":
		return FieldContexts, nil
	case "meta":
		return FieldMeta, nil
	default:
		return 0, fmt.Errorf("unknown field %q", value)
	}
}

func normalizeStateToken(value string) (string, error) {
	state := strings.ToLower(strings.TrimSpace(value))
	if state == "todo" {
		state = "inbox"
	}
	if !domain.IsTaskState(state) {
		return "", domain.InvalidStateExpectedError(value)
	}
	return state, nil
}

func (p *parser) isOpStart() bool {
	tok := p.current()
	if tok.kind == tokenPlus || tok.kind == tokenMinus || tok.kind == tokenBang {
		return true
	}
	if tok.kind == tokenWord {
		if strings.HasPrefix(tok.text, "#") || strings.HasPrefix(tok.text, "@") {
			return true
		}
		return p.peek().kind == tokenColon && isKnownFieldName(tok.text)
	}
	return false
}

func isKnownFieldName(value string) bool {
	_, err := parseFieldName(value)
	return err == nil
}

func tokenIsWord(tok token, text string) bool {
	return tok.kind == tokenWord && strings.EqualFold(tok.text, text)
}

func isLogicalOr(tok token) bool {
	if tok.kind != tokenWord {
		return false
	}
	value := strings.ToLower(tok.text)
	return value == "or" || value == "||"
}

func isLogicalAnd(tok token) bool {
	if tok.kind != tokenWord {
		return false
	}
	value := strings.ToLower(tok.text)
	return value == "and" || value == "&&"
}

func isPredicateField(value string) bool {
	switch strings.ToLower(value) {
	case predicateStateField, predicateDueField, predicateProject, predicateContext, predicateText:
		return true
	default:
		return false
	}
}

func isCreateVerb(value string) bool {
	switch value {
	case "add", "create", "new":
		return true
	default:
		return false
	}
}

func isUpdateVerb(value string) bool {
	switch value {
	case "edit", "set", "update":
		return true
	default:
		return false
	}
}

func isFilterVerb(value string) bool {
	switch value {
	case "find", "show", "filter", "list":
		return true
	default:
		return false
	}
}

func isRelativeDate(value string) bool {
	switch value {
	case "today", "tomorrow", "next-week":
		return true
	default:
		return false
	}
}

func isListField(field Field) bool {
	return field == FieldProjects || field == FieldContexts || field == FieldMeta
}

func isClearableField(field Field) bool {
	switch field {
	case FieldTitle, FieldState:
		return false
	case FieldNotes, FieldDue, FieldWaiting, FieldProjects, FieldContexts, FieldMeta:
		return true
	default:
		return false
	}
}

func joinTokens(tokens []token) (string, Span) {
	if len(tokens) == 0 {
		return "", Span{}
	}
	parts := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		parts = append(parts, tok.text)
	}
	return strings.Join(parts, " "), Span{Start: tokens[0].span.Start, End: tokens[len(tokens)-1].span.End}
}

func hasFieldOp(ops []Operation, field Field) bool {
	for _, op := range ops {
		typed, ok := op.(SetOp)
		if ok && typed.Field == field {
			return true
		}
	}
	return false
}

func (p *parser) current() token {
	if p.idx >= len(p.tokens) {
		return token{kind: tokenEOF}
	}
	return p.tokens[p.idx]
}

func (p *parser) peek() token {
	if p.idx+1 >= len(p.tokens) {
		return token{kind: tokenEOF}
	}
	return p.tokens[p.idx+1]
}

func (p *parser) advance() {
	if p.idx < len(p.tokens)-1 {
		p.idx++
	}
}

func (p *parser) atEOF() bool {
	return p.current().kind == tokenEOF
}
