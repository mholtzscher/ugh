// Code generated from /home/michael/code/ugh/internal/nlp/antlr/UghParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // UghParser
import "github.com/antlr4-go/antlr/v4"

// UghParserListener is a complete listener for a parse tree produced by UghParser.
type UghParserListener interface {
	antlr.ParseTreeListener

	// EnterRoot is called when entering the root production.
	EnterRoot(c *RootContext)

	// EnterCommand is called when entering the command production.
	EnterCommand(c *CommandContext)

	// EnterCreateCommand is called when entering the createCommand production.
	EnterCreateCommand(c *CreateCommandContext)

	// EnterCreateVerb is called when entering the createVerb production.
	EnterCreateVerb(c *CreateVerbContext)

	// EnterCreatePart is called when entering the createPart production.
	EnterCreatePart(c *CreatePartContext)

	// EnterCreateOp is called when entering the createOp production.
	EnterCreateOp(c *CreateOpContext)

	// EnterCreateText is called when entering the createText production.
	EnterCreateText(c *CreateTextContext)

	// EnterUpdateCommand is called when entering the updateCommand production.
	EnterUpdateCommand(c *UpdateCommandContext)

	// EnterUpdateVerb is called when entering the updateVerb production.
	EnterUpdateVerb(c *UpdateVerbContext)

	// EnterTargetRef is called when entering the targetRef production.
	EnterTargetRef(c *TargetRefContext)

	// EnterOperation is called when entering the operation production.
	EnterOperation(c *OperationContext)

	// EnterFilterCommand is called when entering the filterCommand production.
	EnterFilterCommand(c *FilterCommandContext)

	// EnterFilterVerb is called when entering the filterVerb production.
	EnterFilterVerb(c *FilterVerbContext)

	// EnterFilterOrExpr is called when entering the filterOrExpr production.
	EnterFilterOrExpr(c *FilterOrExprContext)

	// EnterFilterAndExpr is called when entering the filterAndExpr production.
	EnterFilterAndExpr(c *FilterAndExprContext)

	// EnterFilterNotExpr is called when entering the filterNotExpr production.
	EnterFilterNotExpr(c *FilterNotExprContext)

	// EnterFilterAtom is called when entering the filterAtom production.
	EnterFilterAtom(c *FilterAtomContext)

	// EnterFilterPredicate is called when entering the filterPredicate production.
	EnterFilterPredicate(c *FilterPredicateContext)

	// EnterFilterFieldPredicate is called when entering the filterFieldPredicate production.
	EnterFilterFieldPredicate(c *FilterFieldPredicateContext)

	// EnterFilterTagPredicate is called when entering the filterTagPredicate production.
	EnterFilterTagPredicate(c *FilterTagPredicateContext)

	// EnterFilterTextPredicate is called when entering the filterTextPredicate production.
	EnterFilterTextPredicate(c *FilterTextPredicateContext)

	// EnterFilterValue is called when entering the filterValue production.
	EnterFilterValue(c *FilterValueContext)

	// EnterFilterValueWord is called when entering the filterValueWord production.
	EnterFilterValueWord(c *FilterValueWordContext)

	// EnterOrOp is called when entering the orOp production.
	EnterOrOp(c *OrOpContext)

	// EnterAndOp is called when entering the andOp production.
	EnterAndOp(c *AndOpContext)

	// EnterNotOp is called when entering the notOp production.
	EnterNotOp(c *NotOpContext)

	// EnterViewCommand is called when entering the viewCommand production.
	EnterViewCommand(c *ViewCommandContext)

	// EnterViewTarget is called when entering the viewTarget production.
	EnterViewTarget(c *ViewTargetContext)

	// EnterContextCommand is called when entering the contextCommand production.
	EnterContextCommand(c *ContextCommandContext)

	// EnterContextArg is called when entering the contextArg production.
	EnterContextArg(c *ContextArgContext)

	// EnterSetOp is called when entering the setOp production.
	EnterSetOp(c *SetOpContext)

	// EnterAddFieldOp is called when entering the addFieldOp production.
	EnterAddFieldOp(c *AddFieldOpContext)

	// EnterRemoveFieldOp is called when entering the removeFieldOp production.
	EnterRemoveFieldOp(c *RemoveFieldOpContext)

	// EnterClearOp is called when entering the clearOp production.
	EnterClearOp(c *ClearOpContext)

	// EnterTagOp is called when entering the tagOp production.
	EnterTagOp(c *TagOpContext)

	// EnterOpValue is called when entering the opValue production.
	EnterOpValue(c *OpValueContext)

	// EnterOpValueWord is called when entering the opValueWord production.
	EnterOpValueWord(c *OpValueWordContext)

	// EnterIdent is called when entering the ident production.
	EnterIdent(c *IdentContext)

	// ExitRoot is called when exiting the root production.
	ExitRoot(c *RootContext)

	// ExitCommand is called when exiting the command production.
	ExitCommand(c *CommandContext)

	// ExitCreateCommand is called when exiting the createCommand production.
	ExitCreateCommand(c *CreateCommandContext)

	// ExitCreateVerb is called when exiting the createVerb production.
	ExitCreateVerb(c *CreateVerbContext)

	// ExitCreatePart is called when exiting the createPart production.
	ExitCreatePart(c *CreatePartContext)

	// ExitCreateOp is called when exiting the createOp production.
	ExitCreateOp(c *CreateOpContext)

	// ExitCreateText is called when exiting the createText production.
	ExitCreateText(c *CreateTextContext)

	// ExitUpdateCommand is called when exiting the updateCommand production.
	ExitUpdateCommand(c *UpdateCommandContext)

	// ExitUpdateVerb is called when exiting the updateVerb production.
	ExitUpdateVerb(c *UpdateVerbContext)

	// ExitTargetRef is called when exiting the targetRef production.
	ExitTargetRef(c *TargetRefContext)

	// ExitOperation is called when exiting the operation production.
	ExitOperation(c *OperationContext)

	// ExitFilterCommand is called when exiting the filterCommand production.
	ExitFilterCommand(c *FilterCommandContext)

	// ExitFilterVerb is called when exiting the filterVerb production.
	ExitFilterVerb(c *FilterVerbContext)

	// ExitFilterOrExpr is called when exiting the filterOrExpr production.
	ExitFilterOrExpr(c *FilterOrExprContext)

	// ExitFilterAndExpr is called when exiting the filterAndExpr production.
	ExitFilterAndExpr(c *FilterAndExprContext)

	// ExitFilterNotExpr is called when exiting the filterNotExpr production.
	ExitFilterNotExpr(c *FilterNotExprContext)

	// ExitFilterAtom is called when exiting the filterAtom production.
	ExitFilterAtom(c *FilterAtomContext)

	// ExitFilterPredicate is called when exiting the filterPredicate production.
	ExitFilterPredicate(c *FilterPredicateContext)

	// ExitFilterFieldPredicate is called when exiting the filterFieldPredicate production.
	ExitFilterFieldPredicate(c *FilterFieldPredicateContext)

	// ExitFilterTagPredicate is called when exiting the filterTagPredicate production.
	ExitFilterTagPredicate(c *FilterTagPredicateContext)

	// ExitFilterTextPredicate is called when exiting the filterTextPredicate production.
	ExitFilterTextPredicate(c *FilterTextPredicateContext)

	// ExitFilterValue is called when exiting the filterValue production.
	ExitFilterValue(c *FilterValueContext)

	// ExitFilterValueWord is called when exiting the filterValueWord production.
	ExitFilterValueWord(c *FilterValueWordContext)

	// ExitOrOp is called when exiting the orOp production.
	ExitOrOp(c *OrOpContext)

	// ExitAndOp is called when exiting the andOp production.
	ExitAndOp(c *AndOpContext)

	// ExitNotOp is called when exiting the notOp production.
	ExitNotOp(c *NotOpContext)

	// ExitViewCommand is called when exiting the viewCommand production.
	ExitViewCommand(c *ViewCommandContext)

	// ExitViewTarget is called when exiting the viewTarget production.
	ExitViewTarget(c *ViewTargetContext)

	// ExitContextCommand is called when exiting the contextCommand production.
	ExitContextCommand(c *ContextCommandContext)

	// ExitContextArg is called when exiting the contextArg production.
	ExitContextArg(c *ContextArgContext)

	// ExitSetOp is called when exiting the setOp production.
	ExitSetOp(c *SetOpContext)

	// ExitAddFieldOp is called when exiting the addFieldOp production.
	ExitAddFieldOp(c *AddFieldOpContext)

	// ExitRemoveFieldOp is called when exiting the removeFieldOp production.
	ExitRemoveFieldOp(c *RemoveFieldOpContext)

	// ExitClearOp is called when exiting the clearOp production.
	ExitClearOp(c *ClearOpContext)

	// ExitTagOp is called when exiting the tagOp production.
	ExitTagOp(c *TagOpContext)

	// ExitOpValue is called when exiting the opValue production.
	ExitOpValue(c *OpValueContext)

	// ExitOpValueWord is called when exiting the opValueWord production.
	ExitOpValueWord(c *OpValueWordContext)

	// ExitIdent is called when exiting the ident production.
	ExitIdent(c *IdentContext)
}
