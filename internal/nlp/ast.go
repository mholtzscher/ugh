package nlp

// Root is the top-level grammar entrypoint for the DSL.
type Root struct {
	Cmd Command `parser:"@@"`
}

type Command interface {
	command()
}

type CreateVerb string

type UpdateVerb string

type FilterVerb string

type CreateCommand struct {
	Verb  CreateVerb   `parser:"@@"`
	Parts []CreatePart `parser:"@@*"`

	Title string
	Ops   []Operation
}

func (*CreateCommand) command() {}

type CreatePart interface {
	createPart()
}

type CreateOpPart struct {
	Op CreateOp `parser:"@@"`
}

func (*CreateOpPart) createPart() {}

type CreateText struct {
	Text OpValue `parser:"@(Ident | Quoted | HashNumber | Comma)"`
}

func (*CreateText) createPart() {}

type CreateOp interface {
	createOp()
}

type UpdateCommand struct {
	Verb   UpdateVerb  `parser:"@@"`
	Target *TargetRef  `parser:"@@?"`
	Ops    []Operation `parser:"@@*"`
}

func (*UpdateCommand) command() {}

type FilterCommand struct {
	Verb  FilterVerb     `parser:"@@"`
	Chain *FilterOrChain `parser:"@@"`

	Expr FilterExpr
}

func (*FilterCommand) command() {}

// FilterOrChain and related types build a parse tree that is converted into
// FilterExpr after parsing.
type FilterOrChain struct {
	Left  *FilterAndChain `parser:"@@"`
	Op    *OrOperator     `parser:"@@?"`
	Right *FilterOrChain  `parser:"@@?"`
}

type OrOperator struct{}

type FilterAndChain struct {
	Left  *FilterNotExpr  `parser:"@@"`
	Op    *AndOperator    `parser:"@@?"`
	Right *FilterAndChain `parser:"@@?"`
}

type AndOperator struct{}

type FilterNotExpr struct {
	Not  *NotOperator `parser:"@@?"`
	Atom *FilterAtom  `parser:"@@"`
}

type NotOperator struct{}

type FilterAtom struct {
	Paren *FilterOrChain   `parser:"'(' @@ ')'"`
	Pred  *FilterPredicate `parser:"| @@"`
}

type FilterPredicate struct {
	Field *FilterFieldPredicate `parser:"@@"`
	Tag   *FilterTagPredicate   `parser:"| @@"`
	Text  *FilterTextPredicate  `parser:"| @@"`
}

type FilterFieldPredicate struct {
	Field string      `parser:"@SetField"`
	Value FilterValue `parser:"@@"`
}

type FilterTagPredicate struct {
	Project string `parser:"@ProjectTag"`
	Context string `parser:"| @ContextTag"`
}

type FilterTextPredicate struct {
	Value FilterValue `parser:"@@"`
}

type TargetKind int

const (
	TargetSelected TargetKind = iota
	TargetID
)

type TargetRef struct {
	Kind TargetKind
	ID   int64
}

// OpValue is a string-like value that can span multiple tokens and is
// reconstructed with minimal normalization.
type OpValue string

// FilterValue is like OpValue, but parsing stops at boolean operators and
// delimiters used in filter expressions.
type FilterValue string

type Field int

const (
	FieldTitle Field = iota
	FieldNotes
	FieldDue
	FieldWaiting
	FieldState
	FieldProjects
	FieldContexts
	FieldMeta
)

type Operation interface {
	operation()
}

type SetOp struct {
	Field Field   `parser:"@SetField"`
	Value OpValue `parser:"@(Quoted | (Ident | HashNumber | Colon | Comma)+)"`
}

func (SetOp) operation() {}
func (SetOp) createOp()  {}

type AddOp struct {
	Field Field   `parser:"@AddField"`
	Value OpValue `parser:"@(Quoted | (Ident | HashNumber | Colon | Comma)+)"`
}

func (AddOp) operation() {}
func (AddOp) createOp()  {}

type RemoveOp struct {
	Field Field   `parser:"@RemoveField"`
	Value OpValue `parser:"@(Quoted | (Ident | HashNumber | Colon | Comma)+)"`
}

func (RemoveOp) operation() {}
func (RemoveOp) createOp()  {}

type ClearOp struct {
	Field Field `parser:"@(ClearField | ClearOp Ident)"`
}

func (ClearOp) operation() {}
func (ClearOp) createOp()  {}

type TagKind int

const (
	TagProject TagKind = iota
	TagContext
)

type TagOp struct {
	Kind  TagKind
	Value string
}

func (TagOp) operation() {}
func (TagOp) createOp()  {}

// DueShorthandOp is a create-only shorthand for setting due date.
// It is parsed from tokens like "today", "tomorrow", "next-week".
type DueShorthandOp struct {
	Value string
}

func (DueShorthandOp) operation() {}
func (DueShorthandOp) createOp()  {}

type FilterExpr interface {
	filterExpr()
}

type FilterBoolOp int

const (
	FilterAnd FilterBoolOp = iota
	FilterOr
)

type FilterBinary struct {
	Op    FilterBoolOp
	Left  FilterExpr
	Right FilterExpr
}

func (FilterBinary) filterExpr() {}

type FilterNot struct {
	Expr FilterExpr
}

func (FilterNot) filterExpr() {}

type PredicateKind int

const (
	PredState PredicateKind = iota
	PredDue
	PredProject
	PredContext
	PredText
	PredID
)

type DateCmp int

const (
	DateEq DateCmp = iota
	DateGT
	DateGTE
	DateLT
	DateLTE
)

type DateValueKind int

const (
	DateAbsolute DateValueKind = iota
	DateToday
	DateTomorrow
	DateNextWeek
)

type Predicate struct {
	Kind PredicateKind

	Text string

	DateCmp  DateCmp
	DateKind DateValueKind
	DateText string
}

func (Predicate) filterExpr() {}
