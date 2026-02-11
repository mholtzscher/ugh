// Code generated from /home/michael/code/ugh/internal/nlp/antlr/UghParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // UghParser
import "github.com/antlr4-go/antlr/v4"

type BaseUghParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseUghParserVisitor) VisitRoot(ctx *RootContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitCommand(ctx *CommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitCreateCommand(ctx *CreateCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitCreateVerb(ctx *CreateVerbContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitCreatePart(ctx *CreatePartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitCreateOp(ctx *CreateOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitCreateText(ctx *CreateTextContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitUpdateCommand(ctx *UpdateCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitUpdateVerb(ctx *UpdateVerbContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitTargetRef(ctx *TargetRefContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitOperation(ctx *OperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterCommand(ctx *FilterCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterVerb(ctx *FilterVerbContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterOrExpr(ctx *FilterOrExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterAndExpr(ctx *FilterAndExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterNotExpr(ctx *FilterNotExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterAtom(ctx *FilterAtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterPredicate(ctx *FilterPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterFieldPredicate(ctx *FilterFieldPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterTagPredicate(ctx *FilterTagPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterTextPredicate(ctx *FilterTextPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterValue(ctx *FilterValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitFilterValueWord(ctx *FilterValueWordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitOrOp(ctx *OrOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitAndOp(ctx *AndOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitNotOp(ctx *NotOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitViewCommand(ctx *ViewCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitViewTarget(ctx *ViewTargetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitContextCommand(ctx *ContextCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitContextArg(ctx *ContextArgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitSetOp(ctx *SetOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitAddFieldOp(ctx *AddFieldOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitRemoveFieldOp(ctx *RemoveFieldOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitClearOp(ctx *ClearOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitTagOp(ctx *TagOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitOpValue(ctx *OpValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitOpValueWord(ctx *OpValueWordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUghParserVisitor) VisitIdent(ctx *IdentContext) interface{} {
	return v.VisitChildren(ctx)
}
