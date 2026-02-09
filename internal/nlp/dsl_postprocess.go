package nlp

import (
	"errors"
	"strconv"
	"strings"
)

func (c *CreateCommand) postProcess() error {
	if c == nil {
		return errors.New("nil create command")
	}

	titleTokens := make([]string, 0, len(c.Parts))
	ops := make([]Operation, 0, len(c.Parts))

	for _, part := range c.Parts {
		switch typed := part.(type) {
		case *CreateText:
			text := strings.TrimSpace(string(typed.Text))
			if text != "" {
				titleTokens = append(titleTokens, text)
			}
		case *CreateOpPart:
			if typed.Op == nil {
				continue
			}
			op, ok := typed.Op.(Operation)
			if !ok {
				return errors.New("invalid create op")
			}
			if normalized, keep := normalizeOperation(op); keep {
				ops = append(ops, normalized)
			}
		default:
			// Ignore unknown parts.
		}
	}

	c.Title = strings.TrimSpace(joinTokens(titleTokens))
	c.Ops = ops

	if c.Title == "" && !hasTitleSetOp(c.Ops) {
		return errors.New("create command requires title or title: field")
	}
	return nil
}

func (u *UpdateCommand) postProcess() {
	if u == nil {
		return
	}
	if u.Target == nil {
		u.Target = &TargetRef{Kind: TargetSelected}
	}

	if len(u.Ops) == 0 {
		return
	}
	normalized := make([]Operation, 0, len(u.Ops))
	for _, op := range u.Ops {
		if norm, keep := normalizeOperation(op); keep {
			normalized = append(normalized, norm)
		}
	}
	u.Ops = normalized
}

func (f *FilterCommand) postProcess() error {
	if f == nil {
		return errors.New("nil filter command")
	}
	if f.Chain == nil {
		return errors.New("filter command requires an expression")
	}
	expr := f.Chain.toExpr()
	if expr == nil {
		return errors.New("filter command requires an expression")
	}
	f.Expr = expr
	return nil
}

func (o *FilterOrChain) toExpr() FilterExpr {
	if o == nil || o.Left == nil {
		return nil
	}

	left := o.Left.toExpr()
	if left == nil {
		return nil
	}

	if o.Op == nil || o.Right == nil {
		return left
	}

	right := o.Right.toExpr()
	if right == nil {
		return left
	}

	return FilterBinary{Op: FilterOr, Left: left, Right: right}
}

func (a *FilterAndChain) toExpr() FilterExpr {
	if a == nil || a.Left == nil {
		return nil
	}

	left := a.Left.toExpr()
	if left == nil {
		return nil
	}

	if a.Op == nil || a.Right == nil {
		return left
	}

	right := a.Right.toExpr()
	if right == nil {
		return left
	}

	return FilterBinary{Op: FilterAnd, Left: left, Right: right}
}

func (n *FilterNotExpr) toExpr() FilterExpr {
	if n == nil || n.Atom == nil {
		return nil
	}
	expr := n.Atom.toExpr()
	if expr == nil {
		return nil
	}
	if n.Not != nil {
		return FilterNot{Expr: expr}
	}
	return expr
}

func (a *FilterAtom) toExpr() FilterExpr {
	if a == nil {
		return nil
	}
	if a.Paren != nil {
		return a.Paren.toExpr()
	}
	if a.Pred == nil {
		return nil
	}
	p := a.Pred.toPredicate()
	if p == nil {
		return nil
	}
	return *p
}

func (p *FilterPredicate) toPredicate() *Predicate {
	if p == nil {
		return nil
	}
	if p.Field != nil {
		return p.Field.toPredicate()
	}
	if p.Tag != nil {
		return p.Tag.toPredicate()
	}
	if p.Text != nil {
		return p.Text.toPredicate()
	}
	return nil
}

func (p *FilterFieldPredicate) toPredicate() *Predicate {
	if p == nil {
		return nil
	}
	field := normalizeCapturedField([]string{p.Field})
	value := strings.TrimSpace(string(p.Value))

	switch field {
	case "state":
		return &Predicate{Kind: PredState, Text: value}
	case "due":
		return &Predicate{Kind: PredDue, Text: value}
	case "project", "projects":
		return &Predicate{Kind: PredProject, Text: value}
	case "context", "contexts":
		return &Predicate{Kind: PredContext, Text: value}
	case "text":
		return &Predicate{Kind: PredText, Text: value}
	case "id":
		if id, ok := parsePossibleID(value); ok {
			return &Predicate{Kind: PredID, Text: strconv.FormatInt(id, 10)}
		}
		return &Predicate{Kind: PredID, Text: strings.TrimPrefix(value, "#")}
	default:
		// Unknown field, treat as text search.
		if field == "" {
			return &Predicate{Kind: PredText, Text: value}
		}
		return &Predicate{Kind: PredText, Text: field + ":" + value}
	}
}

func (p *FilterTagPredicate) toPredicate() *Predicate {
	if p == nil {
		return nil
	}
	if p.Project != "" {
		return &Predicate{Kind: PredProject, Text: strings.TrimPrefix(p.Project, "#")}
	}
	if p.Context != "" {
		return &Predicate{Kind: PredContext, Text: strings.TrimPrefix(p.Context, "@")}
	}
	return &Predicate{Kind: PredText, Text: ""}
}

func (p *FilterTextPredicate) toPredicate() *Predicate {
	if p == nil {
		return nil
	}
	value := strings.TrimSpace(string(p.Value))
	if id, ok := parsePossibleID(value); ok {
		return &Predicate{Kind: PredID, Text: strconv.FormatInt(id, 10)}
	}
	return &Predicate{Kind: PredText, Text: value}
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

func hasTitleSetOp(ops []Operation) bool {
	for _, op := range ops {
		setOp, ok := op.(SetOp)
		if ok && setOp.Field == FieldTitle {
			return true
		}
	}
	return false
}

func normalizeOperation(op Operation) (Operation, bool) {
	if op == nil {
		return nil, false
	}

	switch typed := op.(type) {
	case SetOp:
		return typed, true
	case *SetOp:
		if typed == nil {
			return nil, false
		}
		return SetOp{Field: typed.Field, Value: typed.Value}, true
	case AddOp:
		return typed, true
	case *AddOp:
		if typed == nil {
			return nil, false
		}
		return AddOp{Field: typed.Field, Value: typed.Value}, true
	case RemoveOp:
		return typed, true
	case *RemoveOp:
		if typed == nil {
			return nil, false
		}
		return RemoveOp{Field: typed.Field, Value: typed.Value}, true
	case ClearOp:
		return typed, true
	case *ClearOp:
		if typed == nil {
			return nil, false
		}
		return ClearOp{Field: typed.Field}, true
	case TagOp:
		return typed, true
	case *TagOp:
		if typed == nil {
			return nil, false
		}
		return TagOp{Kind: typed.Kind, Value: typed.Value}, true
	case DueShorthandOp:
		return typed, true
	case *DueShorthandOp:
		if typed == nil {
			return nil, false
		}
		return DueShorthandOp{Value: typed.Value}, true
	case *tagOpNode:
		if typed == nil {
			return nil, false
		}
		return TagOp{Kind: typed.Kind, Value: typed.Value}, true
	case *dueShorthandNode:
		if typed == nil {
			return nil, false
		}
		return DueShorthandOp{Value: typed.Value}, true
	default:
		return op, true
	}
}
