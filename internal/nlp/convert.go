package nlp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Token represents a single token in the DSL.
type Token struct {
	Kind  string
	Value string
}

// Token kind constants.
const (
	tokenIdent      = "Ident"
	tokenQuoted     = "Quoted"
	tokenColon      = "Colon"
	tokenVerb       = "Verb"
	tokenAddOp      = "AddOp"
	tokenRemoveOp   = "RemoveOp"
	tokenClearOp    = "ClearOp"
	tokenProjectTag = "ProjectTag"
	tokenContextTag = "ContextTag"
	tokenLParen     = "LParen"
	tokenRParen     = "RParen"
	tokenAndOp      = "AndOp"
	tokenOrOp       = "OrOp"
	tokenAnd        = "And"
	tokenOr         = "Or"
	tokenNot        = "Not"
	tokenEOF        = "EOF"
)

func tokenize(g *GInput) []Token {
	tokens := make([]Token, 0, len(g.Tokens))
	for _, t := range g.Tokens {
		tok := extractToken(t)
		if tok.Kind != "" {
			tokens = append(tokens, tok)
		}
	}
	return tokens
}

func extractToken(t *GToken) Token {
	switch {
	case t.Ident != "":
		return Token{Kind: tokenIdent, Value: t.Ident}
	case t.Quoted != "":
		return Token{Kind: tokenQuoted, Value: t.Quoted}
	case t.Colon != "":
		return Token{Kind: tokenColon, Value: t.Colon}
	case t.Verb != "":
		return Token{Kind: tokenVerb, Value: t.Verb}
	case t.AddOp != "":
		return Token{Kind: tokenAddOp, Value: t.AddOp}
	case t.RemoveOp != "":
		return Token{Kind: tokenRemoveOp, Value: t.RemoveOp}
	case t.ClearOp != "":
		return Token{Kind: tokenClearOp, Value: t.ClearOp}
	case t.ProjectTag != "":
		return Token{Kind: tokenProjectTag, Value: t.ProjectTag}
	case t.ContextTag != "":
		return Token{Kind: tokenContextTag, Value: t.ContextTag}
	case t.LParen != "":
		return Token{Kind: tokenLParen, Value: t.LParen}
	case t.RParen != "":
		return Token{Kind: tokenRParen, Value: t.RParen}
	case t.AndOp != "":
		return Token{Kind: tokenAndOp, Value: t.AndOp}
	case t.OrOp != "":
		return Token{Kind: tokenOrOp, Value: t.OrOp}
	case t.And != "":
		return Token{Kind: tokenAnd, Value: t.And}
	case t.Or != "":
		return Token{Kind: tokenOr, Value: t.Or}
	case t.Not != "":
		return Token{Kind: tokenNot, Value: t.Not}
	default:
		return Token{}
	}
}

type tokenStream struct {
	tokens []Token
	idx    int
}

func newTokenStream(tokens []Token) *tokenStream {
	return &tokenStream{tokens: tokens, idx: 0}
}

func (ts *tokenStream) current() Token {
	if ts.idx >= len(ts.tokens) {
		return Token{Kind: tokenEOF}
	}
	return ts.tokens[ts.idx]
}

func (ts *tokenStream) advance() {
	if ts.idx < len(ts.tokens) {
		ts.idx++
	}
}

func (ts *tokenStream) atEOF() bool {
	return ts.idx >= len(ts.tokens)
}

func convertGrammar(g *GInput, opts ParseOptions) (ParseResult, error) {
	verb := strings.ToLower(g.Verb)
	tokens := tokenize(g)
	ts := newTokenStream(tokens)

	if isCreateVerb(verb) {
		cmd, err := parseCreate(ts, opts)
		if err != nil {
			return ParseResult{Intent: IntentCreate}, err
		}
		return ParseResult{Intent: IntentCreate, Command: cmd}, nil
	}

	if isUpdateVerb(verb) {
		cmd, err := parseUpdate(ts)
		if err != nil {
			return ParseResult{Intent: IntentUpdate}, err
		}
		return ParseResult{Intent: IntentUpdate, Command: cmd}, nil
	}

	if isFilterVerb(verb) {
		cmd, err := parseFilter(ts)
		if err != nil {
			return ParseResult{Intent: IntentFilter}, err
		}
		return ParseResult{Intent: IntentFilter, Command: cmd}, nil
	}

	return ParseResult{Intent: IntentUnknown}, fmt.Errorf("unknown command verb: %s", verb)
}

func parseCreate(ts *tokenStream, _ ParseOptions) (CreateCommand, error) {
	titleParts := make([]string, 0)
	ops := make([]Operation, 0)

	for !ts.atEOF() {
		tok := ts.current()

		// Check for op-starting tokens
		if ts.isOpStartAtCurrent() {
			op, err := parseOperation(ts)
			if err != nil {
				return CreateCommand{}, err
			}
			ops = append(ops, op)
			continue
		}

		// Check for relative date
		if tok.Kind == tokenIdent && isRelativeDate(strings.ToLower(tok.Value)) {
			ops = append(ops, SetOp{
				Field: FieldDue,
				Value: strings.ToLower(tok.Value),
			})
			ts.advance()
			continue
		}

		// Otherwise, it's part of the title
		titleParts = append(titleParts, getTokenValue(tok))
		ts.advance()
	}

	title := strings.TrimSpace(strings.Join(titleParts, " "))

	// Validate we have a title or title: field
	if title == "" && !hasFieldOp(ops, FieldTitle) {
		return CreateCommand{}, errors.New("create command requires title or title: field")
	}

	return CreateCommand{
		Title: title,
		Ops:   ops,
	}, nil
}

func parseUpdate(ts *tokenStream) (UpdateCommand, error) {
	target, err := parseUpdateTarget(ts)
	if err != nil {
		return UpdateCommand{}, err
	}

	ops, err := parseUpdateOperations(ts)
	if err != nil {
		return UpdateCommand{}, err
	}

	if len(ops) == 0 {
		return UpdateCommand{}, errors.New("update command requires at least one operation")
	}

	return UpdateCommand{
		Target: target,
		Ops:    ops,
	}, nil
}

func parseUpdateTarget(ts *tokenStream) (TargetRef, error) {
	// If token is an operation, default to "selected"
	if ts.atEOF() || ts.isOpStartAtCurrent() {
		return TargetRef{Kind: TargetSelected}, nil
	}

	// Try to parse explicit target
	parsedTarget, err := parseTarget(ts)
	if err == nil {
		return parsedTarget, nil
	}

	// If target parsing fails but we have tag operations, use "selected"
	if ts.current().Kind == tokenProjectTag || ts.current().Kind == tokenContextTag {
		return TargetRef{Kind: TargetSelected}, nil
	}

	return TargetRef{}, err
}

func parseUpdateOperations(ts *tokenStream) ([]Operation, error) {
	ops := make([]Operation, 0)

	for !ts.atEOF() {
		if !ts.isOpStartAtCurrent() {
			return nil, fmt.Errorf("unexpected token %q in update command", ts.current().Value)
		}

		op, err := parseOperation(ts)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	return ops, nil
}

func parseFilter(ts *tokenStream) (FilterCommand, error) {
	if ts.atEOF() {
		return FilterCommand{}, errors.New("filter command requires an expression")
	}

	expr, err := parseFilterOrExpr(ts)
	if err != nil {
		return FilterCommand{}, err
	}

	if !ts.atEOF() {
		return FilterCommand{}, fmt.Errorf("unexpected token %q in filter expression", ts.current().Value)
	}

	return FilterCommand{Expr: expr}, nil
}

func parseTarget(ts *tokenStream) (TargetRef, error) {
	if ts.atEOF() {
		return TargetRef{}, errors.New("expected update target")
	}

	tok := ts.current()

	if tok.Kind == tokenIdent && strings.ToLower(tok.Value) == "selected" {
		ts.advance()
		return TargetRef{Kind: TargetSelected}, nil
	}

	if tok.Kind == tokenIdent {
		lower := strings.ToLower(tok.Value)

		// Check for numeric ID
		if id, err := strconv.ParseInt(lower, 10, 64); err == nil {
			ts.advance()
			return TargetRef{Kind: TargetID, ID: id}, nil
		}

		// Check for #123 format
		if idStr, ok := strings.CutPrefix(lower, "#"); ok {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return TargetRef{}, fmt.Errorf("invalid task id target: %s", tok.Value)
			}
			ts.advance()
			return TargetRef{Kind: TargetID, ID: id}, nil
		}
	}

	return TargetRef{}, fmt.Errorf("invalid update target: %s", tok.Value)
}

func parseOperation(ts *tokenStream) (Operation, error) {
	tok := ts.current()

	switch tok.Kind {
	case "AddOp":
		return parseAddRemoveOperation(ts, true)
	case "RemoveOp":
		return parseAddRemoveOperation(ts, false)
	case "ClearOp":
		return parseClearOperation(ts)
	case "ProjectTag":
		return parseTagOperation(ts, TagProject)
	case "ContextTag":
		return parseTagOperation(ts, TagContext)
	case tokenIdent:
		if ts.idx+1 < len(ts.tokens) && ts.tokens[ts.idx+1].Kind == tokenColon {
			return parseSetOperation(ts)
		}
		return nil, fmt.Errorf("expected operation at %q", tok.Value)
	default:
		return nil, fmt.Errorf("expected operation at %q", tok.Value)
	}
}

func parseTagOperation(ts *tokenStream, kind TagKind) (Operation, error) {
	tok := ts.current()
	value := tok.Value

	if kind == TagProject {
		value = strings.TrimPrefix(value, "#")
	} else {
		value = strings.TrimPrefix(value, "@")
	}

	if len(value) < 1 {
		return nil, fmt.Errorf("invalid tag token: %s", tok.Value)
	}

	ts.advance()
	return TagOp{Kind: kind, Value: value}, nil
}

func parseSetOperation(ts *tokenStream) (Operation, error) {
	fieldTok := ts.current()
	field, err := parseFieldName(fieldTok.Value)
	if err != nil {
		return nil, err
	}
	ts.advance()

	// Expect colon
	if ts.current().Kind != "Colon" {
		return nil, errors.New("expected : after field")
	}
	ts.advance()

	// Parse value (can be multiple tokens until op boundary)
	valueParts := make([]string, 0)
	for !ts.atEOF() && !ts.isOpStartAtCurrent() {
		valueParts = append(valueParts, getTokenValue(ts.current()))
		ts.advance()
	}

	value := strings.TrimSpace(strings.Join(valueParts, " "))
	if value == "" {
		return nil, fmt.Errorf("empty value for %s", fieldTok.Value)
	}

	return SetOp{Field: field, Value: value}, nil
}

func parseAddRemoveOperation(ts *tokenStream, add bool) (Operation, error) {
	ts.advance() // consume + or -

	if ts.atEOF() {
		return nil, errors.New("expected field name after + or -")
	}

	fieldTok := ts.current()
	field, err := parseFieldName(fieldTok.Value)
	if err != nil {
		return nil, err
	}

	if !isListField(field) {
		if add {
			return nil, errors.New("+ only supports list fields (projects, contexts, meta)")
		}
		return nil, errors.New("- only supports list fields (projects, contexts, meta)")
	}
	ts.advance()

	// Expect colon
	if ts.current().Kind != "Colon" {
		return nil, fmt.Errorf("expected : after %s", fieldTok.Value)
	}
	ts.advance()

	// Parse value
	valueParts := make([]string, 0)
	for !ts.atEOF() && !ts.isOpStartAtCurrent() {
		valueParts = append(valueParts, getTokenValue(ts.current()))
		ts.advance()
	}

	value := strings.TrimSpace(strings.Join(valueParts, " "))
	if value == "" {
		return nil, fmt.Errorf("empty value for %s", fieldTok.Value)
	}

	if add {
		return AddOp{Field: field, Value: value}, nil
	}
	return RemoveOp{Field: field, Value: value}, nil
}

func parseClearOperation(ts *tokenStream) (Operation, error) {
	ts.advance() // consume !

	if ts.atEOF() {
		return nil, errors.New("expected field name after clear operator")
	}

	fieldTok := ts.current()
	field, err := parseFieldName(fieldTok.Value)
	if err != nil {
		return nil, err
	}

	if !isClearableField(field) {
		return nil, fmt.Errorf("!%s is not clearable", fieldTok.Value)
	}

	ts.advance()
	return ClearOp{Field: field}, nil
}

func parseFilterOrExpr(ts *tokenStream) (FilterExpr, error) {
	left, err := parseFilterAndExpr(ts)
	if err != nil {
		return nil, err
	}

	for isLogicalOr(ts.current()) {
		ts.advance()
		right, rightErr := parseFilterAndExpr(ts)
		if rightErr != nil {
			return nil, rightErr
		}
		left = FilterBinary{Op: FilterOr, Left: left, Right: right}
	}

	return left, nil
}

func parseFilterAndExpr(ts *tokenStream) (FilterExpr, error) {
	left, err := parseFilterNotExpr(ts)
	if err != nil {
		return nil, err
	}

	for isLogicalAnd(ts.current()) {
		ts.advance()
		right, rightErr := parseFilterNotExpr(ts)
		if rightErr != nil {
			return nil, rightErr
		}
		left = FilterBinary{Op: FilterAnd, Left: left, Right: right}
	}

	return left, nil
}

func parseFilterNotExpr(ts *tokenStream) (FilterExpr, error) {
	if isLogicalNot(ts.current()) {
		ts.advance()
		expr, err := parseFilterNotExpr(ts)
		if err != nil {
			return nil, err
		}
		return FilterNot{Expr: expr}, nil
	}
	return parseFilterAtom(ts)
}

func parseFilterAtom(ts *tokenStream) (FilterExpr, error) {
	if ts.current().Kind == "LParen" {
		ts.advance()
		expr, err := parseFilterOrExpr(ts)
		if err != nil {
			return nil, err
		}
		if ts.current().Kind != "RParen" {
			return nil, errors.New("expected )")
		}
		ts.advance()
		return expr, nil
	}
	return parseFilterPredicate(ts)
}

func parseFilterPredicate(ts *tokenStream) (FilterExpr, error) {
	// Check for field:value format
	if ts.current().Kind == tokenIdent && ts.idx+1 < len(ts.tokens) && ts.tokens[ts.idx+1].Kind == tokenColon {
		field := strings.ToLower(ts.current().Value)
		ts.advance()
		ts.advance() // consume colon

		valueParts := make([]string, 0)
		for !ts.atEOF() && !isFilterBoundary(ts.current()) {
			valueParts = append(valueParts, getTokenValue(ts.current()))
			ts.advance()
		}

		value := strings.TrimSpace(strings.Join(valueParts, " "))
		if value == "" {
			return nil, fmt.Errorf("missing predicate value for %s", field)
		}

		return parsePredicateField(field, value), nil
	}

	// Otherwise, treat as text search (or ID if numeric)
	valueParts := make([]string, 0)
	for !ts.atEOF() && !isFilterBoundary(ts.current()) {
		valueParts = append(valueParts, getTokenValue(ts.current()))
		ts.advance()
	}

	value := strings.TrimSpace(strings.Join(valueParts, " "))
	if value == "" {
		return nil, errors.New("expected predicate")
	}

	// Check if value is purely numeric - treat as ID lookup
	if isNumeric(value) {
		return Predicate{Kind: PredID, Text: value}, nil
	}

	return Predicate{Kind: PredText, Text: value}, nil
}

func parsePredicateField(field, value string) Predicate {
	switch field {
	case "state":
		return Predicate{Kind: PredState, Text: value}
	case "due":
		return Predicate{Kind: PredDue, Text: value}
	case "project":
		return Predicate{Kind: PredProject, Text: value}
	case "context":
		return Predicate{Kind: PredContext, Text: value}
	case "text":
		return Predicate{Kind: PredText, Text: value}
	case "id":
		return Predicate{Kind: PredID, Text: value}
	default:
		// Unknown field, treat as text search
		return Predicate{Kind: PredText, Text: field + ":" + value}
	}
}

// Helper functions

func getTokenValue(tok Token) string {
	if tok.Kind == tokenQuoted {
		// Remove quotes
		val := tok.Value
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			return val[1 : len(val)-1]
		}
	}
	return tok.Value
}

func (ts *tokenStream) isOpStartAtCurrent() bool {
	tok := ts.current()
	switch tok.Kind {
	case tokenAddOp, tokenRemoveOp, tokenClearOp, tokenProjectTag, tokenContextTag:
		return true
	case tokenIdent:
		// Check if next token is Colon (field:value) and it's a known field
		if ts.idx+1 < len(ts.tokens) && ts.tokens[ts.idx+1].Kind == tokenColon {
			// Only treat as op start if it's a known field name
			if isKnownFieldName(tok.Value) {
				return true
			}
		}
		return false
	}
	return false
}

func isKnownFieldName(value string) bool {
	_, err := parseFieldName(value)
	return err == nil
}

func isLogicalOr(tok Token) bool {
	return tok.Kind == tokenOr || tok.Kind == tokenOrOp
}

func isLogicalAnd(tok Token) bool {
	return tok.Kind == tokenAnd || tok.Kind == tokenAndOp
}

func isLogicalNot(tok Token) bool {
	return tok.Kind == tokenNot || tok.Kind == tokenClearOp
}

func isFilterBoundary(tok Token) bool {
	return tok.Kind == tokenRParen || isLogicalAnd(tok) || isLogicalOr(tok) || tok.Kind == tokenEOF
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

func isNumeric(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(value) > 0
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

func hasFieldOp(ops []Operation, field Field) bool {
	for _, op := range ops {
		if setOp, ok := op.(SetOp); ok && setOp.Field == field {
			return true
		}
	}
	return false
}
