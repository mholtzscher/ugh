// Code generated from /home/michael/code/ugh/internal/nlp/antlr/UghParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // UghParser
import "github.com/antlr4-go/antlr/v4"

// BaseUghParserListener is a complete listener for a parse tree produced by UghParser.
type BaseUghParserListener struct{}

var _ UghParserListener = &BaseUghParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseUghParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseUghParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseUghParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseUghParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterRoot is called when production root is entered.
func (s *BaseUghParserListener) EnterRoot(ctx *RootContext) {}

// ExitRoot is called when production root is exited.
func (s *BaseUghParserListener) ExitRoot(ctx *RootContext) {}

// EnterCommand is called when production command is entered.
func (s *BaseUghParserListener) EnterCommand(ctx *CommandContext) {}

// ExitCommand is called when production command is exited.
func (s *BaseUghParserListener) ExitCommand(ctx *CommandContext) {}

// EnterCreateCommand is called when production createCommand is entered.
func (s *BaseUghParserListener) EnterCreateCommand(ctx *CreateCommandContext) {}

// ExitCreateCommand is called when production createCommand is exited.
func (s *BaseUghParserListener) ExitCreateCommand(ctx *CreateCommandContext) {}

// EnterCreateVerb is called when production createVerb is entered.
func (s *BaseUghParserListener) EnterCreateVerb(ctx *CreateVerbContext) {}

// ExitCreateVerb is called when production createVerb is exited.
func (s *BaseUghParserListener) ExitCreateVerb(ctx *CreateVerbContext) {}

// EnterCreatePart is called when production createPart is entered.
func (s *BaseUghParserListener) EnterCreatePart(ctx *CreatePartContext) {}

// ExitCreatePart is called when production createPart is exited.
func (s *BaseUghParserListener) ExitCreatePart(ctx *CreatePartContext) {}

// EnterCreateOp is called when production createOp is entered.
func (s *BaseUghParserListener) EnterCreateOp(ctx *CreateOpContext) {}

// ExitCreateOp is called when production createOp is exited.
func (s *BaseUghParserListener) ExitCreateOp(ctx *CreateOpContext) {}

// EnterCreateText is called when production createText is entered.
func (s *BaseUghParserListener) EnterCreateText(ctx *CreateTextContext) {}

// ExitCreateText is called when production createText is exited.
func (s *BaseUghParserListener) ExitCreateText(ctx *CreateTextContext) {}

// EnterUpdateCommand is called when production updateCommand is entered.
func (s *BaseUghParserListener) EnterUpdateCommand(ctx *UpdateCommandContext) {}

// ExitUpdateCommand is called when production updateCommand is exited.
func (s *BaseUghParserListener) ExitUpdateCommand(ctx *UpdateCommandContext) {}

// EnterUpdateVerb is called when production updateVerb is entered.
func (s *BaseUghParserListener) EnterUpdateVerb(ctx *UpdateVerbContext) {}

// ExitUpdateVerb is called when production updateVerb is exited.
func (s *BaseUghParserListener) ExitUpdateVerb(ctx *UpdateVerbContext) {}

// EnterTargetRef is called when production targetRef is entered.
func (s *BaseUghParserListener) EnterTargetRef(ctx *TargetRefContext) {}

// ExitTargetRef is called when production targetRef is exited.
func (s *BaseUghParserListener) ExitTargetRef(ctx *TargetRefContext) {}

// EnterOperation is called when production operation is entered.
func (s *BaseUghParserListener) EnterOperation(ctx *OperationContext) {}

// ExitOperation is called when production operation is exited.
func (s *BaseUghParserListener) ExitOperation(ctx *OperationContext) {}

// EnterFilterCommand is called when production filterCommand is entered.
func (s *BaseUghParserListener) EnterFilterCommand(ctx *FilterCommandContext) {}

// ExitFilterCommand is called when production filterCommand is exited.
func (s *BaseUghParserListener) ExitFilterCommand(ctx *FilterCommandContext) {}

// EnterFilterVerb is called when production filterVerb is entered.
func (s *BaseUghParserListener) EnterFilterVerb(ctx *FilterVerbContext) {}

// ExitFilterVerb is called when production filterVerb is exited.
func (s *BaseUghParserListener) ExitFilterVerb(ctx *FilterVerbContext) {}

// EnterFilterOrExpr is called when production filterOrExpr is entered.
func (s *BaseUghParserListener) EnterFilterOrExpr(ctx *FilterOrExprContext) {}

// ExitFilterOrExpr is called when production filterOrExpr is exited.
func (s *BaseUghParserListener) ExitFilterOrExpr(ctx *FilterOrExprContext) {}

// EnterFilterAndExpr is called when production filterAndExpr is entered.
func (s *BaseUghParserListener) EnterFilterAndExpr(ctx *FilterAndExprContext) {}

// ExitFilterAndExpr is called when production filterAndExpr is exited.
func (s *BaseUghParserListener) ExitFilterAndExpr(ctx *FilterAndExprContext) {}

// EnterFilterNotExpr is called when production filterNotExpr is entered.
func (s *BaseUghParserListener) EnterFilterNotExpr(ctx *FilterNotExprContext) {}

// ExitFilterNotExpr is called when production filterNotExpr is exited.
func (s *BaseUghParserListener) ExitFilterNotExpr(ctx *FilterNotExprContext) {}

// EnterFilterAtom is called when production filterAtom is entered.
func (s *BaseUghParserListener) EnterFilterAtom(ctx *FilterAtomContext) {}

// ExitFilterAtom is called when production filterAtom is exited.
func (s *BaseUghParserListener) ExitFilterAtom(ctx *FilterAtomContext) {}

// EnterFilterPredicate is called when production filterPredicate is entered.
func (s *BaseUghParserListener) EnterFilterPredicate(ctx *FilterPredicateContext) {}

// ExitFilterPredicate is called when production filterPredicate is exited.
func (s *BaseUghParserListener) ExitFilterPredicate(ctx *FilterPredicateContext) {}

// EnterFilterFieldPredicate is called when production filterFieldPredicate is entered.
func (s *BaseUghParserListener) EnterFilterFieldPredicate(ctx *FilterFieldPredicateContext) {}

// ExitFilterFieldPredicate is called when production filterFieldPredicate is exited.
func (s *BaseUghParserListener) ExitFilterFieldPredicate(ctx *FilterFieldPredicateContext) {}

// EnterFilterTagPredicate is called when production filterTagPredicate is entered.
func (s *BaseUghParserListener) EnterFilterTagPredicate(ctx *FilterTagPredicateContext) {}

// ExitFilterTagPredicate is called when production filterTagPredicate is exited.
func (s *BaseUghParserListener) ExitFilterTagPredicate(ctx *FilterTagPredicateContext) {}

// EnterFilterTextPredicate is called when production filterTextPredicate is entered.
func (s *BaseUghParserListener) EnterFilterTextPredicate(ctx *FilterTextPredicateContext) {}

// ExitFilterTextPredicate is called when production filterTextPredicate is exited.
func (s *BaseUghParserListener) ExitFilterTextPredicate(ctx *FilterTextPredicateContext) {}

// EnterFilterValue is called when production filterValue is entered.
func (s *BaseUghParserListener) EnterFilterValue(ctx *FilterValueContext) {}

// ExitFilterValue is called when production filterValue is exited.
func (s *BaseUghParserListener) ExitFilterValue(ctx *FilterValueContext) {}

// EnterFilterValueWord is called when production filterValueWord is entered.
func (s *BaseUghParserListener) EnterFilterValueWord(ctx *FilterValueWordContext) {}

// ExitFilterValueWord is called when production filterValueWord is exited.
func (s *BaseUghParserListener) ExitFilterValueWord(ctx *FilterValueWordContext) {}

// EnterOrOp is called when production orOp is entered.
func (s *BaseUghParserListener) EnterOrOp(ctx *OrOpContext) {}

// ExitOrOp is called when production orOp is exited.
func (s *BaseUghParserListener) ExitOrOp(ctx *OrOpContext) {}

// EnterAndOp is called when production andOp is entered.
func (s *BaseUghParserListener) EnterAndOp(ctx *AndOpContext) {}

// ExitAndOp is called when production andOp is exited.
func (s *BaseUghParserListener) ExitAndOp(ctx *AndOpContext) {}

// EnterNotOp is called when production notOp is entered.
func (s *BaseUghParserListener) EnterNotOp(ctx *NotOpContext) {}

// ExitNotOp is called when production notOp is exited.
func (s *BaseUghParserListener) ExitNotOp(ctx *NotOpContext) {}

// EnterViewCommand is called when production viewCommand is entered.
func (s *BaseUghParserListener) EnterViewCommand(ctx *ViewCommandContext) {}

// ExitViewCommand is called when production viewCommand is exited.
func (s *BaseUghParserListener) ExitViewCommand(ctx *ViewCommandContext) {}

// EnterViewTarget is called when production viewTarget is entered.
func (s *BaseUghParserListener) EnterViewTarget(ctx *ViewTargetContext) {}

// ExitViewTarget is called when production viewTarget is exited.
func (s *BaseUghParserListener) ExitViewTarget(ctx *ViewTargetContext) {}

// EnterContextCommand is called when production contextCommand is entered.
func (s *BaseUghParserListener) EnterContextCommand(ctx *ContextCommandContext) {}

// ExitContextCommand is called when production contextCommand is exited.
func (s *BaseUghParserListener) ExitContextCommand(ctx *ContextCommandContext) {}

// EnterContextArg is called when production contextArg is entered.
func (s *BaseUghParserListener) EnterContextArg(ctx *ContextArgContext) {}

// ExitContextArg is called when production contextArg is exited.
func (s *BaseUghParserListener) ExitContextArg(ctx *ContextArgContext) {}

// EnterSetOp is called when production setOp is entered.
func (s *BaseUghParserListener) EnterSetOp(ctx *SetOpContext) {}

// ExitSetOp is called when production setOp is exited.
func (s *BaseUghParserListener) ExitSetOp(ctx *SetOpContext) {}

// EnterAddFieldOp is called when production addFieldOp is entered.
func (s *BaseUghParserListener) EnterAddFieldOp(ctx *AddFieldOpContext) {}

// ExitAddFieldOp is called when production addFieldOp is exited.
func (s *BaseUghParserListener) ExitAddFieldOp(ctx *AddFieldOpContext) {}

// EnterRemoveFieldOp is called when production removeFieldOp is entered.
func (s *BaseUghParserListener) EnterRemoveFieldOp(ctx *RemoveFieldOpContext) {}

// ExitRemoveFieldOp is called when production removeFieldOp is exited.
func (s *BaseUghParserListener) ExitRemoveFieldOp(ctx *RemoveFieldOpContext) {}

// EnterClearOp is called when production clearOp is entered.
func (s *BaseUghParserListener) EnterClearOp(ctx *ClearOpContext) {}

// ExitClearOp is called when production clearOp is exited.
func (s *BaseUghParserListener) ExitClearOp(ctx *ClearOpContext) {}

// EnterTagOp is called when production tagOp is entered.
func (s *BaseUghParserListener) EnterTagOp(ctx *TagOpContext) {}

// ExitTagOp is called when production tagOp is exited.
func (s *BaseUghParserListener) ExitTagOp(ctx *TagOpContext) {}

// EnterOpValue is called when production opValue is entered.
func (s *BaseUghParserListener) EnterOpValue(ctx *OpValueContext) {}

// ExitOpValue is called when production opValue is exited.
func (s *BaseUghParserListener) ExitOpValue(ctx *OpValueContext) {}

// EnterOpValueWord is called when production opValueWord is entered.
func (s *BaseUghParserListener) EnterOpValueWord(ctx *OpValueWordContext) {}

// ExitOpValueWord is called when production opValueWord is exited.
func (s *BaseUghParserListener) ExitOpValueWord(ctx *OpValueWordContext) {}

// EnterIdent is called when production ident is entered.
func (s *BaseUghParserListener) EnterIdent(ctx *IdentContext) {}

// ExitIdent is called when production ident is exited.
func (s *BaseUghParserListener) ExitIdent(ctx *IdentContext) {}
