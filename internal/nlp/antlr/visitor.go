// Package antlr provides an ANTLR-based parser for the ugh DSL.
//
// This is a proof-of-concept alternative to the participle-based parser
// in the parent nlp package. It produces identical nlp.* AST types.
//
// Type assertions on ANTLR context types are safe because the parser grammar
// guarantees the concrete types at each position in the parse tree.
//
//nolint:errcheck // type assertions on ANTLR contexts are grammar-guaranteed.
package antlr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/antlr/parser"
)

// astBuilder walks the ANTLR parse tree and produces nlp.* AST nodes.
type astBuilder struct {
	parser.BaseUghParserVisitor

	err error // first error encountered during traversal
}

// ─── Root / Command dispatch ────────────────────────────────────────────────

//nolint:ireturn,modernize // implements generated UghParserVisitor interface.
func (v *astBuilder) VisitRoot(ctx *parser.RootContext) interface{} {
	if ctx.Command() == nil {
		v.err = errors.New("empty parse result")
		return nil
	}
	return ctx.Command().Accept(v)
}

//nolint:ireturn,modernize // implements generated UghParserVisitor interface.
func (v *astBuilder) VisitCommand(ctx *parser.CommandContext) interface{} {
	if ctx.CreateCommand() != nil {
		return ctx.CreateCommand().Accept(v)
	}
	if ctx.UpdateCommand() != nil {
		return ctx.UpdateCommand().Accept(v)
	}
	if ctx.FilterCommand() != nil {
		return ctx.FilterCommand().Accept(v)
	}
	if ctx.ViewCommand() != nil {
		return ctx.ViewCommand().Accept(v)
	}
	if ctx.ContextCommand() != nil {
		return ctx.ContextCommand().Accept(v)
	}
	v.err = errors.New("unrecognized command")
	return nil
}

// ─── CREATE ─────────────────────────────────────────────────────────────────

//nolint:ireturn,modernize // implements generated UghParserVisitor interface.
func (v *astBuilder) VisitCreateCommand(ctx *parser.CreateCommandContext) interface{} {
	verbText := strings.ToLower(ctx.CreateVerb().GetText())

	titleTokens := make([]string, 0)
	ops := make([]nlp.Operation, 0)

	for _, part := range ctx.AllCreatePart() {
		partCtx := part.(*parser.CreatePartContext)
		if partCtx.CreateOp() != nil {
			op := v.visitCreateOp(partCtx.CreateOp().(*parser.CreateOpContext))
			if op != nil {
				ops = append(ops, op)
			}
		} else if partCtx.CreateText() != nil {
			text := v.visitCreateText(partCtx.CreateText().(*parser.CreateTextContext))
			if text != "" {
				titleTokens = append(titleTokens, text)
			}
		}
	}

	title := strings.TrimSpace(joinTokens(titleTokens))
	if title == "" && !hasTitleSetOp(ops) {
		v.err = errors.New("create command requires title or title: field")
		return nil
	}

	return &nlp.CreateCommand{
		Verb:  nlp.CreateVerb(verbText),
		Title: title,
		Ops:   ops,
	}
}

func (v *astBuilder) visitCreateText(ctx *parser.CreateTextContext) string {
	if ctx.QUOTED() != nil {
		return unquote(ctx.QUOTED().GetText())
	}
	// For ident, HASH_NUMBER, COMMA - just get the text
	return ctx.GetText()
}

func (v *astBuilder) visitCreateOp(ctx *parser.CreateOpContext) nlp.Operation {
	if ctx.SetOp() != nil {
		return v.visitSetOp(ctx.SetOp().(*parser.SetOpContext))
	}
	if ctx.AddFieldOp() != nil {
		return v.visitAddFieldOp(ctx.AddFieldOp().(*parser.AddFieldOpContext))
	}
	if ctx.RemoveFieldOp() != nil {
		return v.visitRemoveFieldOp(ctx.RemoveFieldOp().(*parser.RemoveFieldOpContext))
	}
	if ctx.ClearOp() != nil {
		return v.visitClearOp(ctx.ClearOp().(*parser.ClearOpContext))
	}
	if ctx.TagOp() != nil {
		return v.visitTagOp(ctx.TagOp().(*parser.TagOpContext))
	}
	return nil
}

// ─── UPDATE ─────────────────────────────────────────────────────────────────

//nolint:ireturn,modernize // implements generated UghParserVisitor interface.
func (v *astBuilder) VisitUpdateCommand(ctx *parser.UpdateCommandContext) interface{} {
	verbText := strings.ToLower(ctx.UpdateVerb().GetText())

	target := &nlp.TargetRef{Kind: nlp.TargetSelected}
	if ctx.TargetRef() != nil {
		t, err := v.visitTargetRef(ctx.TargetRef().(*parser.TargetRefContext))
		if err != nil {
			v.err = err
			return nil
		}
		target = t
	}

	ops := make([]nlp.Operation, 0, len(ctx.AllOperation()))
	for _, opCtx := range ctx.AllOperation() {
		op := v.visitOperation(opCtx.(*parser.OperationContext))
		if op != nil {
			ops = append(ops, op)
		}
	}

	return &nlp.UpdateCommand{
		Verb:   nlp.UpdateVerb(verbText),
		Target: target,
		Ops:    ops,
	}
}

func (v *astBuilder) visitTargetRef(ctx *parser.TargetRefContext) (*nlp.TargetRef, error) {
	if ctx.HASH_NUMBER() != nil {
		text := ctx.HASH_NUMBER().GetText()
		idStr := strings.TrimPrefix(text, "#")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("invalid update target: %s", text)
		}
		return &nlp.TargetRef{Kind: nlp.TargetID, ID: id}, nil
	}

	if ctx.IDENT() != nil {
		text := strings.ToLower(strings.TrimSpace(ctx.IDENT().GetText()))
		switch text {
		case "selected", "it", "this", "that":
			return &nlp.TargetRef{Kind: nlp.TargetSelected}, nil
		}
		if isDigits(text) {
			id, err := strconv.ParseInt(text, 10, 64)
			if err != nil || id <= 0 {
				return nil, fmt.Errorf("invalid update target: %s", text)
			}
			return &nlp.TargetRef{Kind: nlp.TargetID, ID: id}, nil
		}
		return nil, fmt.Errorf("invalid update target: %s", text)
	}

	return &nlp.TargetRef{Kind: nlp.TargetSelected}, nil
}

func (v *astBuilder) visitOperation(ctx *parser.OperationContext) nlp.Operation {
	if ctx.SetOp() != nil {
		return v.visitSetOp(ctx.SetOp().(*parser.SetOpContext))
	}
	if ctx.AddFieldOp() != nil {
		return v.visitAddFieldOp(ctx.AddFieldOp().(*parser.AddFieldOpContext))
	}
	if ctx.RemoveFieldOp() != nil {
		return v.visitRemoveFieldOp(ctx.RemoveFieldOp().(*parser.RemoveFieldOpContext))
	}
	if ctx.ClearOp() != nil {
		return v.visitClearOp(ctx.ClearOp().(*parser.ClearOpContext))
	}
	if ctx.TagOp() != nil {
		return v.visitTagOp(ctx.TagOp().(*parser.TagOpContext))
	}
	return nil
}

// ─── FILTER ─────────────────────────────────────────────────────────────────

//nolint:ireturn,modernize // implements generated UghParserVisitor interface.
func (v *astBuilder) VisitFilterCommand(ctx *parser.FilterCommandContext) interface{} {
	verbText := strings.ToLower(ctx.FilterVerb().GetText())

	expr := v.visitFilterOrExpr(ctx.FilterOrExpr().(*parser.FilterOrExprContext))
	if expr == nil {
		v.err = errors.New("filter command requires an expression")
		return nil
	}

	return &nlp.FilterCommand{
		Verb: nlp.FilterVerb(verbText),
		Expr: expr,
	}
}

func (v *astBuilder) visitFilterOrExpr(ctx *parser.FilterOrExprContext) nlp.FilterExpr {
	andExprs := ctx.AllFilterAndExpr()
	if len(andExprs) == 0 {
		return nil
	}

	left := v.visitFilterAndExpr(andExprs[0].(*parser.FilterAndExprContext))
	if left == nil {
		return nil
	}

	for i := 1; i < len(andExprs); i++ {
		right := v.visitFilterAndExpr(andExprs[i].(*parser.FilterAndExprContext))
		if right == nil {
			return left
		}
		left = nlp.FilterBinary{Op: nlp.FilterOr, Left: left, Right: right}
	}

	return left
}

func (v *astBuilder) visitFilterAndExpr(ctx *parser.FilterAndExprContext) nlp.FilterExpr {
	notExprs := ctx.AllFilterNotExpr()
	if len(notExprs) == 0 {
		return nil
	}

	left := v.visitFilterNotExpr(notExprs[0].(*parser.FilterNotExprContext))
	if left == nil {
		return nil
	}

	for i := 1; i < len(notExprs); i++ {
		right := v.visitFilterNotExpr(notExprs[i].(*parser.FilterNotExprContext))
		if right == nil {
			return left
		}
		left = nlp.FilterBinary{Op: nlp.FilterAnd, Left: left, Right: right}
	}

	return left
}

func (v *astBuilder) visitFilterNotExpr(ctx *parser.FilterNotExprContext) nlp.FilterExpr {
	expr := v.visitFilterAtom(ctx.FilterAtom().(*parser.FilterAtomContext))
	if expr == nil {
		return nil
	}

	if ctx.NotOp() != nil {
		return nlp.FilterNot{Expr: expr}
	}
	return expr
}

func (v *astBuilder) visitFilterAtom(ctx *parser.FilterAtomContext) nlp.FilterExpr {
	if ctx.FilterOrExpr() != nil {
		return v.visitFilterOrExpr(ctx.FilterOrExpr().(*parser.FilterOrExprContext))
	}
	if ctx.FilterPredicate() != nil {
		return v.visitFilterPredicate(ctx.FilterPredicate().(*parser.FilterPredicateContext))
	}
	return nil
}

func (v *astBuilder) visitFilterPredicate(ctx *parser.FilterPredicateContext) nlp.FilterExpr {
	if ctx.FilterFieldPredicate() != nil {
		return v.visitFilterFieldPredicate(ctx.FilterFieldPredicate().(*parser.FilterFieldPredicateContext))
	}
	if ctx.FilterTagPredicate() != nil {
		return v.visitFilterTagPredicate(ctx.FilterTagPredicate().(*parser.FilterTagPredicateContext))
	}
	if ctx.FilterTextPredicate() != nil {
		return v.visitFilterTextPredicate(ctx.FilterTextPredicate().(*parser.FilterTextPredicateContext))
	}
	return nil
}

func (v *astBuilder) visitFilterFieldPredicate(ctx *parser.FilterFieldPredicateContext) nlp.FilterExpr {
	rawField := ctx.SET_FIELD().GetText()
	field := normalizeFieldName(rawField)
	value := strings.TrimSpace(v.collectFilterValue(ctx.FilterValue().(*parser.FilterValueContext)))

	switch field {
	case "state":
		return nlp.Predicate{Kind: nlp.PredState, Text: value}
	case "due":
		return nlp.Predicate{Kind: nlp.PredDue, Text: value}
	case "project", "projects":
		return nlp.Predicate{Kind: nlp.PredProject, Text: value}
	case "context", "contexts":
		return nlp.Predicate{Kind: nlp.PredContext, Text: value}
	case "text":
		return nlp.Predicate{Kind: nlp.PredText, Text: value}
	case "id":
		if id, ok := parsePossibleID(value); ok {
			return nlp.Predicate{Kind: nlp.PredID, Text: strconv.FormatInt(id, 10)}
		}
		return nlp.Predicate{Kind: nlp.PredID, Text: strings.TrimPrefix(value, "#")}
	default:
		if field == "" {
			return nlp.Predicate{Kind: nlp.PredText, Text: value}
		}
		return nlp.Predicate{Kind: nlp.PredText, Text: field + ":" + value}
	}
}

func (v *astBuilder) visitFilterTagPredicate(ctx *parser.FilterTagPredicateContext) nlp.FilterExpr {
	if ctx.PROJECT_TAG() != nil {
		name := strings.TrimPrefix(ctx.PROJECT_TAG().GetText(), "#")
		return nlp.Predicate{Kind: nlp.PredProject, Text: name}
	}
	if ctx.CONTEXT_TAG() != nil {
		name := strings.TrimPrefix(ctx.CONTEXT_TAG().GetText(), "@")
		return nlp.Predicate{Kind: nlp.PredContext, Text: name}
	}
	return nlp.Predicate{Kind: nlp.PredText, Text: ""}
}

func (v *astBuilder) visitFilterTextPredicate(ctx *parser.FilterTextPredicateContext) nlp.FilterExpr {
	value := strings.TrimSpace(v.collectFilterValue(ctx.FilterValue().(*parser.FilterValueContext)))
	if id, ok := parsePossibleID(value); ok {
		return nlp.Predicate{Kind: nlp.PredID, Text: strconv.FormatInt(id, 10)}
	}
	return nlp.Predicate{Kind: nlp.PredText, Text: value}
}

func (v *astBuilder) collectFilterValue(ctx *parser.FilterValueContext) string {
	if ctx.QUOTED() != nil {
		return unquote(ctx.QUOTED().GetText())
	}
	words := ctx.AllFilterValueWord()
	parts := make([]string, 0, len(words))
	for _, w := range words {
		parts = append(parts, w.GetText())
	}
	return joinTokens(parts)
}

// ─── VIEW ───────────────────────────────────────────────────────────────────

//nolint:ireturn,modernize // implements generated UghParserVisitor interface.
func (v *astBuilder) VisitViewCommand(ctx *parser.ViewCommandContext) interface{} {
	cmd := &nlp.ViewCommand{
		Verb: nlp.ViewVerb("view"),
	}

	if ctx.ViewTarget() != nil {
		targetText := strings.ToLower(strings.TrimSpace(ctx.ViewTarget().GetText()))
		canonical := canonicalViewName(targetText)
		if canonical == "" {
			v.err = fmt.Errorf("invalid view: %s", targetText)
			return nil
		}
		cmd.Target = &nlp.ViewTarget{Name: canonical}
	}

	return cmd
}

// ─── CONTEXT ────────────────────────────────────────────────────────────────

//nolint:ireturn,modernize // implements generated UghParserVisitor interface.
func (v *astBuilder) VisitContextCommand(ctx *parser.ContextCommandContext) interface{} {
	cmd := &nlp.ContextCommand{
		Verb: nlp.ContextVerb("context"),
	}

	if ctx.ContextArg() != nil {
		argCtx := ctx.ContextArg().(*parser.ContextArgContext)
		arg, err := v.visitContextArg(argCtx)
		if err != nil {
			v.err = err
			return nil
		}
		cmd.Arg = arg
	}

	return cmd
}

func (v *astBuilder) visitContextArg(ctx *parser.ContextArgContext) (*nlp.ContextArg, error) {
	if ctx.PROJECT_TAG() != nil {
		name := strings.TrimPrefix(ctx.PROJECT_TAG().GetText(), "#")
		return &nlp.ContextArg{Project: name}, nil
	}
	if ctx.CONTEXT_TAG() != nil {
		name := strings.TrimPrefix(ctx.CONTEXT_TAG().GetText(), "@")
		return &nlp.ContextArg{Context: name}, nil
	}
	if ctx.Ident() != nil {
		text := strings.ToLower(strings.TrimSpace(ctx.Ident().GetText()))
		if text == "clear" {
			return &nlp.ContextArg{Clear: true}, nil
		}
		return nil, fmt.Errorf("invalid context argument: %s", text)
	}
	return nil, errors.New("invalid context argument")
}

// ─── Shared operations ──────────────────────────────────────────────────────

func (v *astBuilder) visitSetOp(ctx *parser.SetOpContext) nlp.Operation {
	rawField := ctx.SET_FIELD().GetText()
	field, err := captureField(rawField)
	if err != nil {
		v.err = err
		return nil
	}
	value := v.collectOpValue(ctx.OpValue().(*parser.OpValueContext))
	return nlp.SetOp{Field: field, Value: nlp.OpValue(value)}
}

func (v *astBuilder) visitAddFieldOp(ctx *parser.AddFieldOpContext) nlp.Operation {
	rawField := ctx.ADD_FIELD().GetText()
	field, err := captureField(rawField)
	if err != nil {
		v.err = err
		return nil
	}
	value := v.collectOpValue(ctx.OpValue().(*parser.OpValueContext))
	return nlp.AddOp{Field: field, Value: nlp.OpValue(value)}
}

func (v *astBuilder) visitRemoveFieldOp(ctx *parser.RemoveFieldOpContext) nlp.Operation {
	rawField := ctx.REMOVE_FIELD().GetText()
	field, err := captureField(rawField)
	if err != nil {
		v.err = err
		return nil
	}
	value := v.collectOpValue(ctx.OpValue().(*parser.OpValueContext))
	return nlp.RemoveOp{Field: field, Value: nlp.OpValue(value)}
}

func (v *astBuilder) visitClearOp(ctx *parser.ClearOpContext) nlp.Operation {
	if ctx.CLEAR_FIELD() != nil {
		rawField := ctx.CLEAR_FIELD().GetText()
		field, err := captureClearField(rawField)
		if err != nil {
			v.err = err
			return nil
		}
		return nlp.ClearOp{Field: field}
	}
	// ! ident form
	if ctx.Ident() != nil {
		rawField := ctx.Ident().GetText()
		field, err := captureClearField("!" + rawField)
		if err != nil {
			v.err = err
			return nil
		}
		return nlp.ClearOp{Field: field}
	}
	v.err = errors.New("invalid clear operation")
	return nil
}

func (v *astBuilder) visitTagOp(ctx *parser.TagOpContext) nlp.Operation {
	if ctx.PROJECT_TAG() != nil {
		name := strings.TrimPrefix(ctx.PROJECT_TAG().GetText(), "#")
		return nlp.TagOp{Kind: nlp.TagProject, Value: name}
	}
	if ctx.CONTEXT_TAG() != nil {
		name := strings.TrimPrefix(ctx.CONTEXT_TAG().GetText(), "@")
		return nlp.TagOp{Kind: nlp.TagContext, Value: name}
	}
	return nil
}

func (v *astBuilder) collectOpValue(ctx *parser.OpValueContext) string {
	if ctx.QUOTED() != nil {
		return unquote(ctx.QUOTED().GetText())
	}
	words := ctx.AllOpValueWord()
	parts := make([]string, 0, len(words))
	for _, w := range words {
		parts = append(parts, w.GetText())
	}
	return joinTokens(parts)
}
