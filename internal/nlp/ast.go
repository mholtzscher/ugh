package nlp

type Node interface {
	NodeSpan() Span
}

type CommandAST interface {
	Node
	commandNode()
}

type CreateCommand struct {
	Title string
	Ops   []Operation
	Span  Span
}

func (c CreateCommand) commandNode()   {}
func (c CreateCommand) NodeSpan() Span { return c.Span }

type UpdateCommand struct {
	Target TargetRef
	Ops    []Operation
	Span   Span
}

func (c UpdateCommand) commandNode()   {}
func (c UpdateCommand) NodeSpan() Span { return c.Span }

type FilterCommand struct {
	Expr FilterExpr
	Span Span
}

func (c FilterCommand) commandNode()   {}
func (c FilterCommand) NodeSpan() Span { return c.Span }

type TargetKind int

const (
	TargetSelected TargetKind = iota
	TargetID
)

type TargetRef struct {
	Kind TargetKind
	ID   int64
	Span Span
}

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
	Node
	opNode()
}

type SetOp struct {
	Field Field
	Value Value
	Span  Span
}

func (o SetOp) opNode()        {}
func (o SetOp) NodeSpan() Span { return o.Span }

type AddOp struct {
	Field Field
	Value Value
	Span  Span
}

func (o AddOp) opNode()        {}
func (o AddOp) NodeSpan() Span { return o.Span }

type RemoveOp struct {
	Field Field
	Value Value
	Span  Span
}

func (o RemoveOp) opNode()        {}
func (o RemoveOp) NodeSpan() Span { return o.Span }

type ClearOp struct {
	Field Field
	Span  Span
}

func (o ClearOp) opNode()        {}
func (o ClearOp) NodeSpan() Span { return o.Span }

type TagKind int

const (
	TagProject TagKind = iota
	TagContext
)

type TagOp struct {
	Kind  TagKind
	Value string
	Span  Span
}

func (o TagOp) opNode()        {}
func (o TagOp) NodeSpan() Span { return o.Span }

type Value struct {
	Raw    string
	Quoted bool
	Span   Span
}

type FilterExpr interface {
	Node
	filterNode()
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
	Span  Span
}

func (e FilterBinary) filterNode()    {}
func (e FilterBinary) NodeSpan() Span { return e.Span }

type FilterNot struct {
	Expr FilterExpr
	Span Span
}

func (e FilterNot) filterNode()    {}
func (e FilterNot) NodeSpan() Span { return e.Span }

type PredicateKind int

const (
	PredState PredicateKind = iota
	PredDue
	PredProject
	PredContext
	PredText
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

	Span Span
}

func (e Predicate) filterNode()    {}
func (e Predicate) NodeSpan() Span { return e.Span }
