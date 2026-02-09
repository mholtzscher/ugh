package nlp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// convertGrammar converts the participle grammar types to AST types.
func convertGrammar(g *GCommand, opts ParseOptions) (ParseResult, error) {
	if g.Create != nil {
		cmd, err := convertCreate(g.Create, opts)
		if err != nil {
			return ParseResult{Intent: IntentCreate}, err
		}
		return ParseResult{Intent: IntentCreate, Command: cmd}, nil
	}

	if g.Update != nil {
		cmd, err := convertUpdate(g.Update)
		if err != nil {
			return ParseResult{Intent: IntentUpdate}, err
		}
		return ParseResult{Intent: IntentUpdate, Command: cmd}, nil
	}

	if g.Filter != nil {
		cmd, err := convertFilter(g.Filter)
		if err != nil {
			return ParseResult{Intent: IntentFilter}, err
		}
		return ParseResult{Intent: IntentFilter, Command: cmd}, nil
	}

	return ParseResult{Intent: IntentUnknown}, errors.New("unknown command type")
}

func convertCreate(g *GCreateCmd, _ ParseOptions) (CreateCommand, error) {
	titleParts := make([]string, 0)
	ops := make([]Operation, 0)

	for _, item := range g.Items {
		switch v := item.(type) {
		case GCreateTitle:
			titleParts = append(titleParts, v.Words...)
		case *GCreateTitle:
			titleParts = append(titleParts, v.Words...)
		case GCreateOp:
			op, err := convertOp(&v.Op)
			if err != nil {
				return CreateCommand{}, err
			}
			ops = append(ops, op)
		case *GCreateOp:
			op, err := convertOp(&v.Op)
			if err != nil {
				return CreateCommand{}, err
			}
			ops = append(ops, op)
		case GCreateRelativeDate:
			// Relative date as due date shorthand
			ops = append(ops, SetOp{
				Field: FieldDue,
				Value: v.Date,
			})
		case *GCreateRelativeDate:
			// Relative date as due date shorthand
			ops = append(ops, SetOp{
				Field: FieldDue,
				Value: v.Date,
			})
		}
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

func convertUpdate(g *GUpdateCmd) (UpdateCommand, error) {
	target, err := convertTarget(g.Target)
	if err != nil {
		return UpdateCommand{}, err
	}

	ops := make([]Operation, 0, len(g.Ops))
	for _, gOp := range g.Ops {
		op, opErr := convertOp(&gOp)
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
	}, nil
}

func convertTarget(target string) (TargetRef, error) {
	if target == "" {
		return TargetRef{Kind: TargetSelected}, nil
	}

	lower := strings.ToLower(target)

	switch lower {
	case "selected", "it", "this":
		return TargetRef{Kind: TargetSelected}, nil
	case "that":
		// "that" refers to second-to-last, handled at higher level
		return TargetRef{Kind: TargetSelected}, nil
	}

	// Try numeric ID
	if id, err := strconv.ParseInt(lower, 10, 64); err == nil {
		return TargetRef{Kind: TargetID, ID: id}, nil
	}

	// Try #123 format
	if strings.HasPrefix(lower, "#") {
		if id, err := strconv.ParseInt(lower[1:], 10, 64); err == nil {
			return TargetRef{Kind: TargetID, ID: id}, nil
		}
	}

	// Invalid target
	return TargetRef{}, fmt.Errorf("invalid update target: %s", target)
}

func convertRawValue(v GRawValue) string {
	if v.Quoted != "" {
		return v.Quoted
	}
	return strings.TrimSpace(strings.Join(v.Tokens, " "))
}

func convertOp(g *GOp) (Operation, error) {
	if g.Set != nil {
		field, err := parseFieldFromToken(g.Set.Field)
		if err != nil {
			return nil, err
		}
		return SetOp{
			Field: field,
			Value: convertRawValue(g.Set.Value),
		}, nil
	}

	if g.Add != nil {
		field, err := parseFieldFromToken(g.Add.Field)
		if err != nil {
			return nil, err
		}
		return AddOp{
			Field: field,
			Value: convertRawValue(g.Add.Value),
		}, nil
	}

	if g.Remove != nil {
		field, err := parseFieldFromToken(g.Remove.Field)
		if err != nil {
			return nil, err
		}
		return RemoveOp{
			Field: field,
			Value: convertRawValue(g.Remove.Value),
		}, nil
	}

	if g.Clear != nil {
		field, err := parseFieldFromToken(g.Clear.Field)
		if err != nil {
			return nil, err
		}
		return ClearOp{Field: field}, nil
	}

	if g.Tag != nil {
		if g.Tag.Project != "" {
			return TagOp{
				Kind:  TagProject,
				Value: strings.TrimPrefix(g.Tag.Project, "#"),
			}, nil
		}
		if g.Tag.Context != "" {
			return TagOp{
				Kind:  TagContext,
				Value: strings.TrimPrefix(g.Tag.Context, "@"),
			}, nil
		}
	}

	return nil, errors.New("unknown operation type")
}

func convertFilter(g *GFilterCmd) (FilterCommand, error) {
	if g.Expr == nil {
		return FilterCommand{}, errors.New("filter command requires an expression")
	}
	expr := convertOrChain(g.Expr)
	if expr == nil {
		return FilterCommand{}, errors.New("filter command requires an expression")
	}
	return FilterCommand{Expr: expr}, nil
}

func convertOrChain(g *GOrChain) FilterExpr {
	if g == nil || g.Left == nil {
		return nil
	}

	left := convertAndChain(g.Left)
	if left == nil {
		return nil
	}

	// If no right side or no operator, just return left
	if g.Right == nil || g.Op == "" {
		return left
	}

	// Build OR expression
	right := convertOrChain(g.Right)
	if right == nil {
		return left
	}

	return FilterBinary{
		Op:    FilterOr,
		Left:  left,
		Right: right,
	}
}

func convertAndChain(g *GAndChain) FilterExpr {
	if g == nil || g.Left == nil {
		return nil
	}

	left := convertNotExpr(g.Left)
	if left == nil {
		return nil
	}

	// If no right side or no operator, just return left
	if g.Right == nil || g.Op == "" {
		return left
	}

	// Build AND expression
	right := convertAndChain(g.Right)
	if right == nil {
		return left
	}

	return FilterBinary{
		Op:    FilterAnd,
		Left:  left,
		Right: right,
	}
}

func convertNotExpr(g *GNotExpr) FilterExpr {
	if g == nil || g.Atom == nil {
		return nil
	}

	expr := convertFilterAtom(g.Atom)
	if expr == nil {
		return nil
	}

	// Check if there's a NOT operator
	if g.Not != "" && (g.Not == "not" || g.Not == "!") {
		return FilterNot{Expr: expr}
	}

	return expr
}

func convertFilterAtom(g *GFilterAtom) FilterExpr {
	if g.Paren != nil {
		return convertOrChain(g.Paren)
	}

	if g.Pred != nil {
		return convertPredicate(g.Pred)
	}

	return nil
}

func convertPredicate(g *GPredicate) FilterExpr {
	if g.FieldPred != nil {
		return convertFieldPredicate(g.FieldPred)
	}

	if g.TagPred != nil {
		return convertTagPredicate(g.TagPred)
	}

	if g.TextPred != nil {
		return convertTextPredicate(g.TextPred)
	}

	if g.IDPred != nil {
		return Predicate{
			Kind: PredID,
			Text: g.IDPred.ID,
		}
	}

	return nil
}

func convertFieldPredicate(g *GFieldPredicate) Predicate {
	field := extractFieldName(g.Field)
	value := convertRawValue(g.Value)

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

func convertTagPredicate(g *GTagPredicate) Predicate {
	if g.Project != "" {
		return Predicate{
			Kind: PredProject,
			Text: strings.TrimPrefix(g.Project, "#"),
		}
	}
	if g.Context != "" {
		return Predicate{
			Kind: PredContext,
			Text: strings.TrimPrefix(g.Context, "@"),
		}
	}
	return Predicate{Kind: PredText, Text: ""}
}

func convertTextPredicate(g *GTextPredicate) Predicate {
	value := convertRawValue(g.Value)

	// Check if value is purely numeric - treat as ID lookup
	if isNumeric(value) {
		return Predicate{Kind: PredID, Text: value}
	}

	return Predicate{Kind: PredText, Text: value}
}

// Helper functions

func parseFieldFromToken(token string) (Field, error) {
	// Token includes the colon, e.g., "title:" or "+projects:"
	field := extractFieldName(token)
	return parseFieldName(field)
}

func extractFieldName(token string) string {
	// Remove leading + or - and trailing colon
	token = strings.TrimSpace(token)
	token = strings.TrimPrefix(token, "+")
	token = strings.TrimPrefix(token, "-")
	token = strings.TrimPrefix(token, "!")
	token = strings.TrimSuffix(token, ":")
	return strings.ToLower(strings.TrimSpace(token))
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
	case "id":
		// ID is special - not a real field for SetOp
		return 0, errors.New("id cannot be set directly")
	case "text":
		// Text is special - for filtering only
		return 0, errors.New("text is not a settable field")
	default:
		return 0, fmt.Errorf("unknown field %q", value)
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

func hasFieldOp(ops []Operation, field Field) bool {
	for _, op := range ops {
		if setOp, ok := op.(SetOp); ok && setOp.Field == field {
			return true
		}
	}
	return false
}
