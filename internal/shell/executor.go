package shell

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/compile"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

const (
	contextNoneValue = "none"

	viewNameInbox    = "inbox"
	viewNameNow      = "now"
	viewNameWaiting  = "waiting"
	viewNameLater    = "later"
	viewNameCalendar = "calendar"
)

// Executor bridges NLP parsing to service execution.
type Executor struct {
	svc     service.Service
	state   *SessionState
	parser  nlp.Parser
	noColor bool
}

// NewExecutor creates a new executor.
func NewExecutor(svc service.Service, state *SessionState, noColor bool) *Executor {
	return &Executor{
		svc:     svc,
		state:   state,
		parser:  nlp.NewParser(),
		noColor: noColor,
	}
}

// Execute parses and executes a natural language command.
func (e *Executor) Execute(ctx context.Context, input string) (*ExecuteResult, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("empty command")
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
		return nil, fmt.Errorf("parse: %w", err)
	}

	// Check for parse diagnostics
	if len(parseResult.Diagnostics) > 0 {
		for _, diag := range parseResult.Diagnostics {
			if diag.Severity == nlp.SeverityError {
				return nil, fmt.Errorf("parse error: %s", diag.Message)
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
		return e.executeFilter(ctx, plan)
	case nlp.IntentView:
		return e.executeView(ctx, parseResult)
	case nlp.IntentContext:
		return e.executeContext(parseResult)
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
			Summary:   "cleared context",
			Timestamp: time.Now(),
		}, nil
	}

	if cmd.Arg.Project != "" {
		e.state.ContextProject = cmd.Arg.Project
		return &ExecuteResult{
			Intent:    "context",
			Message:   fmt.Sprintf("Set project context to #%s", cmd.Arg.Project),
			Summary:   fmt.Sprintf("context project #%s", cmd.Arg.Project),
			Timestamp: time.Now(),
		}, nil
	}

	if cmd.Arg.Context != "" {
		e.state.ContextContext = cmd.Arg.Context
		return &ExecuteResult{
			Intent:    "context",
			Message:   fmt.Sprintf("Set context filter to @%s", cmd.Arg.Context),
			Summary:   fmt.Sprintf("context @%s", cmd.Arg.Context),
			Timestamp: time.Now(),
		}, nil
	}

	return nil, errors.New("invalid context command argument")
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
		return "find due:today", nil
	default:
		return "", fmt.Errorf("unknown view: %s", viewName)
	}
}

func (e *Executor) showContext() *ExecuteResult {
	selected := contextNoneValue
	if e.state.SelectedTaskID != nil {
		selected = fmt.Sprintf("#%d", *e.state.SelectedTaskID)
	}

	last := formatTaskIDs(e.state.LastTaskIDs)
	project := formatContextValue(e.state.ContextProject, "#")
	ctx := formatContextValue(e.state.ContextContext, "@")

	data := pterm.TableData{
		{"Selected", selected},
		{"Last", last},
		{"Project", project},
		{"Context", ctx},
	}

	var msg string
	if !e.noColor {
		table, _ := pterm.DefaultTable.WithData(data).Srender()
		msg = pterm.Yellow("Current Context:\n") + table
	} else {
		var b strings.Builder
		b.WriteString("Current Context:\n")
		for _, row := range data {
			_, _ = fmt.Fprintf(&b, "  %s: %s\n", row[0], row[1])
		}
		msg = b.String()
	}

	return &ExecuteResult{
		Intent:    "context",
		Message:   msg,
		Summary:   "showing context",
		Timestamp: time.Now(),
	}
}

func formatTaskIDs(ids []int64) string {
	if len(ids) == 0 {
		return contextNoneValue
	}
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = fmt.Sprintf("#%d", id)
	}
	return strings.Join(strs, ", ")
}

func formatContextValue(value, prefix string) string {
	if value == "" {
		return contextNoneValue
	}
	return prefix + value
}

func (e *Executor) showViewHelp() *ExecuteResult {
	var msg string
	if !e.noColor {
		msg = pterm.Green("Available Views:\n") +
			"  " + pterm.LightGreen("i, inbox") + "     Inbox tasks\n" +
			"  " + pterm.LightGreen("n, now") + "       Now tasks\n" +
			"  " + pterm.LightGreen("w, waiting") + "   Waiting tasks\n" +
			"  " + pterm.LightGreen("l, later") + "     Later tasks\n" +
			"  " + pterm.LightGreen("c, calendar") + "  Tasks due today\n" +
			"\nUsage: view <name> (e.g., view i or view inbox)"
	} else {
		var b strings.Builder
		b.WriteString("Available Views:\n")
		b.WriteString("  i, inbox     Inbox tasks\n")
		b.WriteString("  n, now       Now tasks\n")
		b.WriteString("  w, waiting   Waiting tasks\n")
		b.WriteString("  l, later     Later tasks\n")
		b.WriteString("  c, calendar  Tasks due today\n")
		b.WriteString("\nUsage: view <name> (e.g., view i or view inbox)")
		msg = b.String()
	}

	return &ExecuteResult{
		Intent:    "view",
		Message:   msg,
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
		Summary:   fmt.Sprintf("updated task #%d", task.ID),
		Timestamp: time.Now(),
	}, nil
}

func (e *Executor) executeFilter(ctx context.Context, plan compile.Plan) (*ExecuteResult, error) {
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

	return &ExecuteResult{
		Intent:    "filter",
		Message:   formatTaskList(tasks, e.noColor),
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

func formatTaskList(tasks []*store.Task, noColor bool) string {
	if len(tasks) == 0 {
		return "No tasks found"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d task(s):\n\n", len(tasks)))
	for _, task := range tasks {
		sb.WriteString(formatTaskLine(task, noColor) + "\n")
	}
	return sb.String()
}

func formatTaskLine(task *store.Task, noColor bool) string {
	state := task.State
	if state == "" {
		state = "inbox"
	}

	idStr := formatID(task.ID, noColor)
	stateStr := formatState(string(state), noColor)
	tags := formatTags(task.Projects, task.Contexts, noColor)
	dueStr := formatDueDate(task.DueOn, noColor)

	line := fmt.Sprintf("  %s %s %s", idStr, task.Title, stateStr)
	if tags != "" {
		line += " " + tags
	}
	if dueStr != "" {
		line += " " + dueStr
	}
	return line
}

func formatID(id int64, noColor bool) string {
	if noColor {
		return "#" + strconv.FormatInt(id, 10)
	}
	return pterm.Cyan("#" + strconv.FormatInt(id, 10))
}

func formatState(state string, noColor bool) string {
	if noColor {
		return "[" + state + "]"
	}
	return pterm.Magenta("[" + state + "]")
}

func formatTags(projects, contexts []string, noColor bool) string {
	var tags []string
	for _, p := range projects {
		if noColor {
			tags = append(tags, "#"+p)
		} else {
			tags = append(tags, pterm.Blue("#"+p))
		}
	}
	for _, c := range contexts {
		if noColor {
			tags = append(tags, "@"+c)
		} else {
			tags = append(tags, pterm.Green("@"+c))
		}
	}
	return strings.Join(tags, " ")
}

func formatDueDate(dueOn *time.Time, noColor bool) string {
	if dueOn == nil {
		return ""
	}
	if noColor {
		return dueOn.Format("2006-01-02")
	}
	return pterm.Yellow(dueOn.Format("2006-01-02"))
}
