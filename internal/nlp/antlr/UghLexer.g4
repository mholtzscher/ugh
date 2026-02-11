// UghLexer.g4 — Lexer grammar for the ugh DSL.
//
// Key design: Command verbs (add, set, find, view, context, etc.) are
// recognized as distinct keyword tokens. This lets the parser grammar
// unambiguously dispatch to the correct command alternative.
lexer grammar UghLexer;

// ─── Quoted strings (highest priority) ──────────────────────────────────────
QUOTED      : '"' ( '\\' . | ~["\\] )* '"' ;

// ─── Hash-number IDs like #123  (must precede PROJECT_TAG) ──────────────────
HASH_NUMBER : '#' [0-9]+ ;

// ─── Tags ───────────────────────────────────────────────────────────────────
PROJECT_TAG : '#' [a-zA-Z_] [a-zA-Z0-9_-]* ;
CONTEXT_TAG : '@' [a-zA-Z_] [a-zA-Z0-9_-]* ;

// ─── Tag prefixes (bare # or @) for interactive completion ──────────────────
PROJECT_TAG_PREFIX : '#' ;
CONTEXT_TAG_PREFIX : '@' ;

// ─── Field setters ──────────────────────────────────────────────────────────
SET_FIELD
    : ( 'title' | 'notes' | 'due' | 'waiting' | 'waiting-for' | 'waiting_for'
      | 'state' | 'project' | 'projects' | 'context' | 'contexts'
      | 'meta' | 'id' | 'text'
      ) WS_INLINE* ':'
    ;

ADD_FIELD
    : '+' WS_INLINE* ( 'project' | 'projects' | 'context' | 'contexts' | 'meta' ) WS_INLINE* ':'
    ;

REMOVE_FIELD
    : '-' WS_INLINE* ( 'project' | 'projects' | 'context' | 'contexts' | 'meta' ) WS_INLINE* ':'
    ;

CLEAR_FIELD
    : '!' WS_INLINE* ( 'notes' | 'due' | 'waiting' | 'waiting-for' | 'waiting_for'
                      | 'projects' | 'contexts' | 'meta' )
    ;

// ─── Command verb keywords (case-insensitive) ──────────────────────────────
// These must precede IDENT to get priority.
KW_ADD      : [aA][dD][dD] ;
KW_CREATE   : [cC][rR][eE][aA][tT][eE] ;
KW_NEW      : [nN][eE][wW] ;
KW_SET      : [sS][eE][tT] ;
KW_EDIT     : [eE][dD][iI][tT] ;
KW_UPDATE   : [uU][pP][dD][aA][tT][eE] ;
KW_FIND     : [fF][iI][nN][dD] ;
KW_SHOW     : [sS][hH][oO][wW] ;
KW_LIST     : [lL][iI][sS][tT] ;
KW_FILTER   : [fF][iI][lL][tT][eE][rR] ;
KW_VIEW     : [vV][iI][eE][wW] ;
KW_CONTEXT  : [cC][oO][nN][tT][eE][xX][tT] ;

// ─── Boolean operator keywords ──────────────────────────────────────────────
KW_AND      : [aA][nN][dD] ;
KW_OR       : [oO][rR] ;
KW_NOT      : [nN][oO][tT] ;

// ─── Operators ──────────────────────────────────────────────────────────────
CLEAR_OP  : '!' ;
ADD_OP    : '+' ;
REMOVE_OP : '-' ;

// ─── Punctuation ────────────────────────────────────────────────────────────
COLON  : ':' ;
COMMA  : ',' ;
LPAREN : '(' ;
RPAREN : ')' ;

// ─── Logical operators (symbolic) ───────────────────────────────────────────
AND_OP : '&&' ;
OR_OP  : '||' ;

// ─── Identifiers / words ───────────────────────────────────────────────────
IDENT : [a-zA-Z0-9_-]+ ;

// ─── Whitespace (skip) ─────────────────────────────────────────────────────
WS : [ \t\r\n]+ -> skip ;

// Fragment used inside other rules (not a standalone token)
fragment WS_INLINE : [ \t] ;
