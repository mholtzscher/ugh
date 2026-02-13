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
	if err := validateFilterChain(f.Chain); err != nil {
		return err
	}
	expr := f.Chain.toExpr()
	if expr == nil {
		return errors.New("filter command requires an expression")
	}
	f.Expr = expr
	return nil
}

func (v *ViewCommand) postProcess() error {
	if v == nil {
		return errors.New("nil view command")
	}
	if v.Target == nil {
		return nil
	}

	v.Target.Name = canonicalViewName(v.Target.Name)
	if v.Target.Name == "" {
		return errors.New("view command requires a valid view name")
	}
	return nil
}

func (c *ContextCommand) postProcess() error {
	if c == nil {
		return errors.New("nil context command")
	}
	if c.Arg == nil {
		return nil
	}

	nonEmpty := 0
	if c.Arg.Clear {
		nonEmpty++
	}
	if strings.TrimSpace(c.Arg.Project) != "" {
		nonEmpty++
	}
	if strings.TrimSpace(c.Arg.Context) != "" {
		nonEmpty++
	}
	if nonEmpty != 1 {
		return errors.New("context command requires exactly one argument")
	}

	c.Arg.Project = strings.TrimSpace(c.Arg.Project)
	c.Arg.Context = strings.TrimSpace(c.Arg.Context)
	return nil
}

func (l *LogCommand) postProcess() error {
	if l == nil {
		return errors.New("nil log command")
	}
	if l.Target == nil {
		return errors.New("log command requires a task id")
	}
	if l.Target.Kind != TargetID {
		return errors.New("log command requires a numeric task id")
	}
	return nil
}

func canonicalViewName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "i", viewNameInbox:
		return viewNameInbox
	case "n", viewNameNow:
		return viewNameNow
	case "w", viewNameWaiting:
		return viewNameWaiting
	case "l", viewNameLater:
		return viewNameLater
	case "c", viewNameCalendar, "today":
		return viewNameCalendar
	default:
		return ""
	}
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
	value := ""
	if p.Value != nil {
		value = strings.TrimSpace(string(*p.Value))
	}

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

func validateFilterChain(chain *FilterOrChain) error {
	if chain == nil {
		return nil
	}
	if err := validateFilterAndChain(chain.Left); err != nil {
		return err
	}
	return validateFilterChain(chain.Right)
}

func validateFilterAndChain(chain *FilterAndChain) error {
	if chain == nil {
		return nil
	}
	if err := validateFilterNotExpr(chain.Left); err != nil {
		return err
	}
	return validateFilterAndChain(chain.Right)
}

func validateFilterNotExpr(expr *FilterNotExpr) error {
	if expr == nil || expr.Atom == nil {
		return nil
	}
	return validateFilterAtom(expr.Atom)
}

func validateFilterAtom(atom *FilterAtom) error {
	if atom == nil {
		return nil
	}
	if atom.Paren != nil {
		return validateFilterChain(atom.Paren)
	}
	if atom.Pred == nil {
		return nil
	}
	return validateFilterPredicate(atom.Pred)
}

func validateFilterPredicate(pred *FilterPredicate) error {
	if pred == nil || pred.Field == nil {
		return nil
	}
	return validateFieldPredicate(pred.Field)
}

func validateFieldPredicate(pred *FilterFieldPredicate) error {
	if pred == nil {
		return nil
	}
	if pred.Value == nil || strings.TrimSpace(string(*pred.Value)) == "" {
		return errors.New("expected value")
	}
	return nil
}

func (p *FilterTagPredicate) toPredicate() *Predicate {
	if p == nil {
		return nil
	}
	if p.Project != "" {
		return &Predicate{Kind: PredProject, Text: p.Project}
	}
	if p.Context != "" {
		return &Predicate{Kind: PredContext, Text: p.Context}
	}
	return &Predicate{Kind: PredText, Text: ""}
}

func (p *FilterTextPredicate) toPredicate() *Predicate {
	if p == nil {
		return nil
	}
	value := strings.TrimSpace(string(p.Value))
	if strings.EqualFold(value, "recent") {
		return &Predicate{Kind: PredRecent, Text: ""}
	}
	if recentLimit, ok := strings.CutPrefix(strings.ToLower(value), "recent:"); ok {
		return &Predicate{Kind: PredRecent, Text: strings.TrimSpace(recentLimit)}
	}
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
	case *SetOp:
		if typed == nil {
			return nil, false
		}
		return *typed, true
	case *AddOp:
		if typed == nil {
			return nil, false
		}
		return *typed, true
	case *RemoveOp:
		if typed == nil {
			return nil, false
		}
		return *typed, true
	case *ClearOp:
		if typed == nil {
			return nil, false
		}
		return *typed, true
	case *tagOpNode:
		if typed == nil {
			return nil, false
		}
		return TagOp{Kind: typed.Kind, Value: typed.Value}, true
	default:
		return op, true
	}
}
