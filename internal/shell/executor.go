package shell

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/compile"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

const (
	viewNameInbox    = "inbox"
	viewNameNow      = "now"
	viewNameWaiting  = "waiting"
	viewNameLater    = "later"
	viewNameCalendar = "calendar"
)

// Executor bridges NLP parsing to service execution.
type Executor struct {
	svc    service.Service
	state  *SessionState
	parser nlp.Parser
}

// NewExecutor creates a new executor.
func NewExecutor(svc service.Service, state *SessionState) *Executor {
	return &Executor{
		svc:    svc,
		state:  state,
		parser: nlp.NewParser(),
	}
}

// Execute parses and executes a natural language command.
func (e *Executor) Execute(ctx context.Context, input string) (*ExecuteResult, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("empty command")
	}

	if badRune, runePos, ok := firstDisallowedControlRune(input); ok {
		return nil, fmt.Errorf(
			"command contains non-printable control character %s at rune position %d; paste plain text and try again",
			formatControlRune(badRune),
			runePos,
		)
	}

	// Pre-process: resolve pronouns and context
	input = e.preprocessInput(input)

	// Parse the natural language input
	parseOpts := nlp.ParseOptions{
		Mode: nlp.ModeAuto,
		Now:  time.Now(),
	}

	parseResult, err := e.parser.Parse(ctx, input, parseOpts)
	if err != nil {
		return nil, err
	}

	// Check for parse diagnostics
	if len(parseResult.Diagnostics) > 0 {
		for _, diag := range parseResult.Diagnostics {
			if diag.Severity == nlp.SeverityError {
				return nil, nlp.NewDiagnosticError(errors.New(diag.Message), []nlp.Diagnostic{diag})
			}
		}
	}

	// Inject sticky context into the AST (post-parse, pre-compile)
	e.injectContext(parseResult)

	// Compile the parse result to an execution plan
	buildOpts := compile.BuildOptions{
		SelectedTaskID: e.state.SelectedTaskID,
		Now:            time.Now(),
	}

	plan, err := compile.Build(parseResult, buildOpts)
	if err != nil {
		return nil, fmt.Errorf("compile: %w", err)
	}

	// Execute the plan
	return e.executePlan(ctx, plan, parseResult)
}

func (e *Executor) preprocessInput(input string) string {
	// Replace pronouns with actual IDs (only if IDs exist)
	if lastID := e.getLastTaskID(); lastID != 0 {
		input = strings.ReplaceAll(input, " it ", fmt.Sprintf(" %d ", lastID))
		input = strings.ReplaceAll(input, " this ", fmt.Sprintf(" %d ", lastID))
		input = strings.ReplaceAll(input, " last ", fmt.Sprintf(" %d ", lastID))
	}
	if secondToLastID := e.getSecondToLastTaskID(); secondToLastID != 0 {
		input = strings.ReplaceAll(input, " that ", fmt.Sprintf(" %d ", secondToLastID))
	}
	if selectedID := e.getSelectedTaskID(); selectedID != 0 {
		input = strings.ReplaceAll(input, " selected ", fmt.Sprintf(" %d ", selectedID))
	}

	return input
}

func (e *Executor) injectContext(parseResult nlp.ParseResult) {
	switch cmd := parseResult.Command.(type) {
	case *nlp.CreateCommand:
		cmd.InjectProject(e.state.ContextProject)
		cmd.InjectContext(e.state.ContextContext)
	case *nlp.UpdateCommand:
		cmd.InjectProject(e.state.ContextProject)
		cmd.InjectContext(e.state.ContextContext)
	case *nlp.FilterCommand:
		cmd.InjectProject(e.state.ContextProject)
		cmd.InjectContext(e.state.ContextContext)
	}
}

func (e *Executor) getLastTaskID() int64 {
	if len(e.state.LastTaskIDs) > 0 {
		return e.state.LastTaskIDs[len(e.state.LastTaskIDs)-1]
	}
	return 0
}

func (e *Executor) getSecondToLastTaskID() int64 {
	if len(e.state.LastTaskIDs) > 1 {
		return e.state.LastTaskIDs[len(e.state.LastTaskIDs)-2]
	}
	return e.getLastTaskID()
}

func (e *Executor) getSelectedTaskID() int64 {
	if e.state.SelectedTaskID != nil {
		return *e.state.SelectedTaskID
	}
	return 0
}

func (e *Executor) executePlan(
	ctx context.Context,
	plan compile.Plan,
	parseResult nlp.ParseResult,
) (*ExecuteResult, error) {
	switch parseResult.Intent {
	case nlp.IntentCreate:
		return e.executeCreate(ctx, plan)
	case nlp.IntentUpdate:
		return e.executeUpdate(ctx, plan)
	case nlp.IntentFilter:
		return e.executeFilter(ctx, plan, parseResult)
	case nlp.IntentView:
		return e.executeView(ctx, parseResult)
	case nlp.IntentContext:
		return e.executeContext(parseResult)
	case nlp.IntentLog:
		return e.executeLog(ctx, plan)
	case nlp.IntentUnknown:
		return nil, errors.New("unknown intent: could not determine command type")
	default:
		return nil, fmt.Errorf("unsupported intent: %v", parseResult.Intent)
	}
}

func (e *Executor) executeView(ctx context.Context, parseResult nlp.ParseResult) (*ExecuteResult, error) {
	cmd, ok := parseResult.Command.(*nlp.ViewCommand)
	if !ok || cmd == nil {
		return nil, errors.New("invalid view command")
	}

	if cmd.Target == nil {
		return e.showViewHelp(), nil
	}

	filterQuery, err := viewFilterQuery(cmd.Target.Name)
	if err != nil {
		return nil, err
	}

	return e.Execute(ctx, filterQuery)
}

func (e *Executor) executeContext(parseResult nlp.ParseResult) (*ExecuteResult, error) {
	cmd, ok := parseResult.Command.(*nlp.ContextCommand)
	if !ok || cmd == nil {
		return nil, errors.New("invalid context command")
	}

	if cmd.Arg == nil {
		return e.showContext(), nil
	}

	if cmd.Arg.Clear {
		e.state.ContextProject = ""
		e.state.ContextContext = ""
		return &ExecuteResult{
			Intent:    "context",
			Message:   "Context filters cleared",
			Level:     ResultLevelInfo,
			Summary:   "cleared context",
			Timestamp: time.Now(),
		}, nil
	}

	if cmd.Arg.Project != "" {
		e.state.ContextProject = cmd.Arg.Project
		return &ExecuteResult{
			Intent:    "context",
			Message:   fmt.Sprintf("Set project context to #%s", cmd.Arg.Project),
			Level:     ResultLevelInfo,
			Summary:   fmt.Sprintf("context project #%s", cmd.Arg.Project),
			Timestamp: time.Now(),
		}, nil
	}

	if cmd.Arg.Context != "" {
		e.state.ContextContext = cmd.Arg.Context
		return &ExecuteResult{
			Intent:    "context",
			Message:   fmt.Sprintf("Set context filter to @%s", cmd.Arg.Context),
			Level:     ResultLevelInfo,
			Summary:   fmt.Sprintf("context @%s", cmd.Arg.Context),
			Timestamp: time.Now(),
		}, nil
	}

	return nil, errors.New("invalid context command argument")
}

const defaultTaskLogLimit = 20

func (e *Executor) executeLog(ctx context.Context, plan compile.Plan) (*ExecuteResult, error) {
	taskID := plan.Target.ID
	if taskID <= 0 {
		return nil, errors.New("log command requires a task id")
	}

	versions, err := e.svc.ListTaskVersions(ctx, taskID, defaultTaskLogLimit)
	if err != nil {
		return nil, fmt.Errorf("list task versions: %w", err)
	}
	if len(versions) == 0 {
		return nil, errors.New("task not found")
	}

	e.state.LastTaskIDs = []int64{taskID}

	return &ExecuteResult{
		Intent:    "show log",
		Versions:  versions,
		Level:     ResultLevelInfo,
		TaskIDs:   []int64{taskID},
		Summary:   fmt.Sprintf("showed %d versions", len(versions)),
		Timestamp: time.Now(),
	}, nil
}

func viewFilterQuery(viewName string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(viewName)) {
	case viewNameInbox:
		return "find state:inbox", nil
	case viewNameNow:
		return "find state:now", nil
	case viewNameWaiting:
		return "find state:waiting", nil
	case viewNameLater:
		return "find state:later", nil
	case viewNameCalendar:
		return "find due:*", nil
	default:
		return "", fmt.Errorf("unknown view: %s", viewName)
	}
}

func (e *Executor) showContext() *ExecuteResult {
	status := output.ContextStatus{
		SelectedID: e.state.SelectedTaskID,
		LastIDs:    e.state.LastTaskIDs,
		Project:    e.state.ContextProject,
		Context:    e.state.ContextContext,
	}

	return &ExecuteResult{
		Intent:    "context",
		Context:   &status,
		Level:     ResultLevelInfo,
		Summary:   "showing context",
		Timestamp: time.Now(),
	}
}

func firstDisallowedControlRune(input string) (rune, int, bool) {
	position := 1
	for _, r := range input {
		if isDisallowedControlRune(r) {
			return r, position, true
		}
		position++
	}

	return 0, 0, false
}

func isDisallowedControlRune(r rune) bool {
	if r == '\n' || r == '\r' || r == '\t' {
		return false
	}

	return unicode.IsControl(r)
}

func formatControlRune(r rune) string {
	if r >= 0 && r <= 0xFF {
		return fmt.Sprintf("U+%04X (\\x%02X)", r, r)
	}

	return fmt.Sprintf("U+%04X", r)
}

func (e *Executor) showViewHelp() *ExecuteResult {
	help := output.ViewHelp{
		Entries: []output.ViewHelpEntry{
			{Label: "i, inbox", Description: "Inbox tasks"},
			{Label: "n, now", Description: "Now tasks"},
			{Label: "w, waiting", Description: "Waiting tasks"},
			{Label: "l, later", Description: "Later tasks"},
			{Label: "c, calendar", Description: "Tasks with due dates"},
		},
		Usage: "view <name> (e.g., view i or view inbox)",
	}

	return &ExecuteResult{
		Intent:    "view",
		ViewHelp:  &help,
		Level:     ResultLevelInfo,
		Summary:   "showing view help",
		Timestamp: time.Now(),
	}
}

func (e *Executor) executeCreate(ctx context.Context, plan compile.Plan) (*ExecuteResult, error) {
	if plan.Create == nil {
		return nil, errors.New("no create request compiled")
	}

	task, err := e.svc.CreateTask(ctx, *plan.Create)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	e.state.LastTaskIDs = []int64{task.ID}

	return &ExecuteResult{
		Intent:    "create",
		Message:   formatTaskCreated(task),
		TaskIDs:   []int64{task.ID},
		Level:     ResultLevelSuccess,
		Summary:   fmt.Sprintf("created task #%d", task.ID),
		Timestamp: time.Now(),
	}, nil
}

func (e *Executor) executeUpdate(ctx context.Context, plan compile.Plan) (*ExecuteResult, error) {
	if plan.Update == nil {
		return nil, errors.New("no update request compiled")
	}

	task, err := e.svc.UpdateTask(ctx, *plan.Update)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	e.state.LastTaskIDs = []int64{task.ID}
	if plan.Target.Kind == nlp.TargetSelected {
		e.state.SelectedTaskID = &task.ID
	}

	return &ExecuteResult{
		Intent:    "update",
		Message:   formatTaskUpdated(task),
		TaskIDs:   []int64{task.ID},
		Level:     ResultLevelSuccess,
		Summary:   fmt.Sprintf("updated task #%d", task.ID),
		Timestamp: time.Now(),
	}, nil
}

func (e *Executor) executeFilter(
	ctx context.Context,
	plan compile.Plan,
	parseResult nlp.ParseResult,
) (*ExecuteResult, error) {
	if plan.Filter == nil {
		return nil, errors.New("no filter request compiled")
	}

	tasks, err := e.svc.ListTasks(ctx, *plan.Filter)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	// Update state with results
	taskIDs := make([]int64, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}
	e.state.LastTaskIDs = taskIDs

	showVerb := false
	if cmd, ok := parseResult.Command.(*nlp.FilterCommand); ok {
		showVerb = string(cmd.Verb) == "show"
	}

	if showVerb && len(tasks) == 1 {
		return &ExecuteResult{
			Intent:    "show",
			Task:      tasks[0],
			Level:     ResultLevelInfo,
			TaskIDs:   taskIDs,
			Summary:   fmt.Sprintf("showed task #%d", tasks[0].ID),
			Timestamp: time.Now(),
		}, nil
	}

	return &ExecuteResult{
		Intent:    "filter",
		Tasks:     tasks,
		Level:     ResultLevelInfo,
		TaskIDs:   taskIDs,
		Summary:   fmt.Sprintf("found %d tasks", len(tasks)),
		Timestamp: time.Now(),
	}, nil
}

func formatTaskCreated(task *store.Task) string {
	return fmt.Sprintf("Created task #%d: %s", task.ID, task.Title)
}

func formatTaskUpdated(task *store.Task) string {
	return fmt.Sprintf("Updated task #%d: %s", task.ID, task.Title)
}
