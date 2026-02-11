// Code generated from /home/michael/code/ugh/internal/nlp/antlr/UghParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // UghParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by UghParser.
type UghParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by UghParser#root.
	VisitRoot(ctx *RootContext) interface{}

	// Visit a parse tree produced by UghParser#command.
	VisitCommand(ctx *CommandContext) interface{}

	// Visit a parse tree produced by UghParser#createCommand.
	VisitCreateCommand(ctx *CreateCommandContext) interface{}

	// Visit a parse tree produced by UghParser#createVerb.
	VisitCreateVerb(ctx *CreateVerbContext) interface{}

	// Visit a parse tree produced by UghParser#createPart.
	VisitCreatePart(ctx *CreatePartContext) interface{}

	// Visit a parse tree produced by UghParser#createOp.
	VisitCreateOp(ctx *CreateOpContext) interface{}

	// Visit a parse tree produced by UghParser#createText.
	VisitCreateText(ctx *CreateTextContext) interface{}

	// Visit a parse tree produced by UghParser#updateCommand.
	VisitUpdateCommand(ctx *UpdateCommandContext) interface{}

	// Visit a parse tree produced by UghParser#updateVerb.
	VisitUpdateVerb(ctx *UpdateVerbContext) interface{}

	// Visit a parse tree produced by UghParser#targetRef.
	VisitTargetRef(ctx *TargetRefContext) interface{}

	// Visit a parse tree produced by UghParser#operation.
	VisitOperation(ctx *OperationContext) interface{}

	// Visit a parse tree produced by UghParser#filterCommand.
	VisitFilterCommand(ctx *FilterCommandContext) interface{}

	// Visit a parse tree produced by UghParser#filterVerb.
	VisitFilterVerb(ctx *FilterVerbContext) interface{}

	// Visit a parse tree produced by UghParser#filterOrExpr.
	VisitFilterOrExpr(ctx *FilterOrExprContext) interface{}

	// Visit a parse tree produced by UghParser#filterAndExpr.
	VisitFilterAndExpr(ctx *FilterAndExprContext) interface{}

	// Visit a parse tree produced by UghParser#filterNotExpr.
	VisitFilterNotExpr(ctx *FilterNotExprContext) interface{}

	// Visit a parse tree produced by UghParser#filterAtom.
	VisitFilterAtom(ctx *FilterAtomContext) interface{}

	// Visit a parse tree produced by UghParser#filterPredicate.
	VisitFilterPredicate(ctx *FilterPredicateContext) interface{}

	// Visit a parse tree produced by UghParser#filterFieldPredicate.
	VisitFilterFieldPredicate(ctx *FilterFieldPredicateContext) interface{}

	// Visit a parse tree produced by UghParser#filterTagPredicate.
	VisitFilterTagPredicate(ctx *FilterTagPredicateContext) interface{}

	// Visit a parse tree produced by UghParser#filterTextPredicate.
	VisitFilterTextPredicate(ctx *FilterTextPredicateContext) interface{}

	// Visit a parse tree produced by UghParser#filterValue.
	VisitFilterValue(ctx *FilterValueContext) interface{}

	// Visit a parse tree produced by UghParser#filterValueWord.
	VisitFilterValueWord(ctx *FilterValueWordContext) interface{}

	// Visit a parse tree produced by UghParser#orOp.
	VisitOrOp(ctx *OrOpContext) interface{}

	// Visit a parse tree produced by UghParser#andOp.
	VisitAndOp(ctx *AndOpContext) interface{}

	// Visit a parse tree produced by UghParser#notOp.
	VisitNotOp(ctx *NotOpContext) interface{}

	// Visit a parse tree produced by UghParser#viewCommand.
	VisitViewCommand(ctx *ViewCommandContext) interface{}

	// Visit a parse tree produced by UghParser#viewTarget.
	VisitViewTarget(ctx *ViewTargetContext) interface{}

	// Visit a parse tree produced by UghParser#contextCommand.
	VisitContextCommand(ctx *ContextCommandContext) interface{}

	// Visit a parse tree produced by UghParser#contextArg.
	VisitContextArg(ctx *ContextArgContext) interface{}

	// Visit a parse tree produced by UghParser#setOp.
	VisitSetOp(ctx *SetOpContext) interface{}

	// Visit a parse tree produced by UghParser#addFieldOp.
	VisitAddFieldOp(ctx *AddFieldOpContext) interface{}

	// Visit a parse tree produced by UghParser#removeFieldOp.
	VisitRemoveFieldOp(ctx *RemoveFieldOpContext) interface{}

	// Visit a parse tree produced by UghParser#clearOp.
	VisitClearOp(ctx *ClearOpContext) interface{}

	// Visit a parse tree produced by UghParser#tagOp.
	VisitTagOp(ctx *TagOpContext) interface{}

	// Visit a parse tree produced by UghParser#opValue.
	VisitOpValue(ctx *OpValueContext) interface{}

	// Visit a parse tree produced by UghParser#opValueWord.
	VisitOpValueWord(ctx *OpValueWordContext) interface{}

	// Visit a parse tree produced by UghParser#ident.
	VisitIdent(ctx *IdentContext) interface{}
}
