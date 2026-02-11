// UghParser.g4 — Parser grammar for the ugh DSL.
//
// Verb keywords are distinct lexer tokens, so the parser can unambiguously
// dispatch to the correct command type.
parser grammar UghParser;

options { tokenVocab = UghLexer; }

// ─── Entry point ────────────────────────────────────────────────────────────
root
    : command EOF
    ;

command
    : createCommand
    | updateCommand
    | filterCommand
    | viewCommand
    | contextCommand
    ;

// ─── CREATE ─────────────────────────────────────────────────────────────────
createCommand
    : createVerb createPart*
    ;

createVerb
    : KW_ADD | KW_CREATE | KW_NEW
    ;

createPart
    : createOp
    | createText
    ;

createOp
    : setOp
    | addFieldOp
    | removeFieldOp
    | clearOp
    | tagOp
    ;

// Words that form the title text (any ident-like token that isn't an operation).
createText
    : ident
    | QUOTED
    | HASH_NUMBER
    | COMMA
    ;

// ─── UPDATE ─────────────────────────────────────────────────────────────────
updateCommand
    : updateVerb targetRef? operation*
    ;

updateVerb
    : KW_SET | KW_EDIT | KW_UPDATE
    ;

targetRef
    : HASH_NUMBER
    | IDENT
    ;

operation
    : setOp
    | addFieldOp
    | removeFieldOp
    | clearOp
    | tagOp
    ;

// ─── FILTER ─────────────────────────────────────────────────────────────────
filterCommand
    : filterVerb filterOrExpr
    ;

filterVerb
    : KW_FIND | KW_SHOW | KW_LIST | KW_FILTER
    ;

// OR has lowest precedence
filterOrExpr
    : filterAndExpr ( orOp filterAndExpr )*
    ;

// AND has higher precedence than OR
filterAndExpr
    : filterNotExpr ( andOp filterNotExpr )*
    ;

// NOT (unary prefix)
filterNotExpr
    : notOp filterAtom
    | filterAtom
    ;

filterAtom
    : LPAREN filterOrExpr RPAREN
    | filterPredicate
    ;

filterPredicate
    : filterFieldPredicate
    | filterTagPredicate
    | filterTextPredicate
    ;

filterFieldPredicate
    : SET_FIELD filterValue
    ;

filterTagPredicate
    : PROJECT_TAG
    | CONTEXT_TAG
    ;

filterTextPredicate
    : filterValue
    ;

filterValue
    : QUOTED
    | filterValueWord+
    ;

filterValueWord
    : ident
    | HASH_NUMBER
    | COLON
    | COMMA
    ;

// ─── Boolean operators ──────────────────────────────────────────────────────
orOp  : OR_OP  | KW_OR  ;
andOp : AND_OP | KW_AND ;
notOp : CLEAR_OP | KW_NOT ;

// ─── VIEW ───────────────────────────────────────────────────────────────────
viewCommand
    : KW_VIEW viewTarget?
    ;

viewTarget
    : ident
    ;

// ─── CONTEXT ────────────────────────────────────────────────────────────────
contextCommand
    : KW_CONTEXT contextArg?
    ;

contextArg
    : PROJECT_TAG
    | CONTEXT_TAG
    | ident       // "clear" or other ident
    ;

// ─── Shared operation rules ─────────────────────────────────────────────────
setOp
    : SET_FIELD opValue
    ;

addFieldOp
    : ADD_FIELD opValue
    ;

removeFieldOp
    : REMOVE_FIELD opValue
    ;

clearOp
    : CLEAR_FIELD
    | CLEAR_OP ident   // "! fieldname" form
    ;

tagOp
    : PROJECT_TAG
    | CONTEXT_TAG
    ;

opValue
    : QUOTED
    | opValueWord+
    ;

opValueWord
    : ident
    | HASH_NUMBER
    | COLON
    | COMMA
    ;

// ─── Ident catch-all ────────────────────────────────────────────────────────
// Since keywords are separate tokens, we need a rule that accepts both
// IDENT and keyword tokens in positions where any word is valid (title text,
// op values, view targets, etc.).
ident
    : IDENT
    | KW_ADD | KW_CREATE | KW_NEW
    | KW_SET | KW_EDIT | KW_UPDATE
    | KW_FIND | KW_SHOW | KW_LIST | KW_FILTER
    | KW_VIEW | KW_CONTEXT
    | KW_AND | KW_OR | KW_NOT
    ;
