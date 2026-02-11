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

type ViewVerb string

type ContextVerb string

const (
	viewNameInbox    = "inbox"
	viewNameNow      = "now"
	viewNameWaiting  = "waiting"
	viewNameLater    = "later"
	viewNameCalendar = "calendar"
)

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

type ViewCommand struct {
	Verb   ViewVerb    `parser:"@@"`
	Target *ViewTarget `parser:"@@?"`
}

func (*ViewCommand) command() {}

type ContextCommand struct {
	Verb ContextVerb `parser:"@@"`
	Arg  *ContextArg `parser:"@@?"`
}

func (*ContextCommand) command() {}

type ViewTarget struct {
	Name string
}

type ContextArg struct {
	Clear   bool
	Project string
	Context string
}

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

type Predicate struct {
	Kind PredicateKind

	Text string
}

func (Predicate) filterExpr() {}

// HasProjectTag returns true if the command has an explicit project tag.
func (c *CreateCommand) HasProjectTag() bool {
	if c == nil {
		return false
	}
	for _, op := range c.Ops {
		if tag, ok := op.(TagOp); ok && tag.Kind == TagProject {
			return true
		}
	}
	return false
}

// HasContextTag returns true if the command has an explicit context tag.
func (c *CreateCommand) HasContextTag() bool {
	if c == nil {
		return false
	}
	for _, op := range c.Ops {
		if tag, ok := op.(TagOp); ok && tag.Kind == TagContext {
			return true
		}
	}
	return false
}

// InjectProject adds a project tag if one doesn't already exist.
func (c *CreateCommand) InjectProject(name string) {
	if c == nil || c.HasProjectTag() || name == "" {
		return
	}
	c.Ops = append(c.Ops, TagOp{Kind: TagProject, Value: name})
}

// InjectContext adds a context tag if one doesn't already exist.
func (c *CreateCommand) InjectContext(name string) {
	if c == nil || c.HasContextTag() || name == "" {
		return
	}
	c.Ops = append(c.Ops, TagOp{Kind: TagContext, Value: name})
}

// HasProjectPredicate returns true if the filter has a project predicate.
func (f *FilterCommand) HasProjectPredicate() bool {
	if f == nil || f.Expr == nil {
		return false
	}
	return hasPredicate(f.Expr, PredProject)
}

// HasContextPredicate returns true if the filter has a context predicate.
func (f *FilterCommand) HasContextPredicate() bool {
	if f == nil || f.Expr == nil {
		return false
	}
	return hasPredicate(f.Expr, PredContext)
}

// InjectProject adds a project predicate to the filter expression.
func (f *FilterCommand) InjectProject(name string) {
	if f == nil || f.HasProjectPredicate() || name == "" {
		return
	}
	pred := Predicate{Kind: PredProject, Text: name}
	if f.Expr == nil {
		f.Expr = pred
	} else {
		f.Expr = FilterBinary{Op: FilterAnd, Left: f.Expr, Right: pred}
	}
}

// InjectContext adds a context predicate to the filter expression.
func (f *FilterCommand) InjectContext(name string) {
	if f == nil || f.HasContextPredicate() || name == "" {
		return
	}
	pred := Predicate{Kind: PredContext, Text: name}
	if f.Expr == nil {
		f.Expr = pred
	} else {
		f.Expr = FilterBinary{Op: FilterAnd, Left: f.Expr, Right: pred}
	}
}

// HasProjectTag returns true if the update has an explicit project tag.
func (u *UpdateCommand) HasProjectTag() bool {
	if u == nil {
		return false
	}
	for _, op := range u.Ops {
		if tag, ok := op.(TagOp); ok && tag.Kind == TagProject {
			return true
		}
	}
	return false
}

// HasContextTag returns true if the update has an explicit context tag.
func (u *UpdateCommand) HasContextTag() bool {
	if u == nil {
		return false
	}
	for _, op := range u.Ops {
		if tag, ok := op.(TagOp); ok && tag.Kind == TagContext {
			return true
		}
	}
	return false
}

// InjectProject adds a project tag if one doesn't already exist.
func (u *UpdateCommand) InjectProject(name string) {
	if u == nil || u.HasProjectTag() || name == "" {
		return
	}
	u.Ops = append(u.Ops, TagOp{Kind: TagProject, Value: name})
}

// InjectContext adds a context tag if one doesn't already exist.
func (u *UpdateCommand) InjectContext(name string) {
	if u == nil || u.HasContextTag() || name == "" {
		return
	}
	u.Ops = append(u.Ops, TagOp{Kind: TagContext, Value: name})
}

// hasPredicate recursively checks if an expression contains a predicate of the given kind.
func hasPredicate(expr FilterExpr, kind PredicateKind) bool {
	switch typed := expr.(type) {
	case Predicate:
		return typed.Kind == kind
	case FilterBinary:
		return hasPredicate(typed.Left, kind) || hasPredicate(typed.Right, kind)
	case FilterNot:
		return hasPredicate(typed.Expr, kind)
	default:
		return false
	}
}
