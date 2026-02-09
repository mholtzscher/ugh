package nlp

type CreateCommand struct {
	Title string
	Ops   []Operation
}

type UpdateCommand struct {
	Target TargetRef
	Ops    []Operation
}

type FilterCommand struct {
	Expr FilterExpr
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
	Field Field
	Value string
}

func (SetOp) operation() {}

type AddOp struct {
	Field Field
	Value string
}

func (AddOp) operation() {}

type RemoveOp struct {
	Field Field
	Value string
}

func (RemoveOp) operation() {}

type ClearOp struct {
	Field Field
}

func (ClearOp) operation() {}

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
