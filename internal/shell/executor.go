package shell

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/compile"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
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
	return e.executePlan(ctx, plan, parseResult.Intent)
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

	// Add context filters if set
	if e.state.ContextProject != "" && !strings.Contains(input, "#") {
		input = input + " #" + e.state.ContextProject
	}
	if e.state.ContextContext != "" && !strings.Contains(input, "@") {
		input = input + " @" + e.state.ContextContext
	}

	return input
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

func (e *Executor) executePlan(ctx context.Context, plan compile.Plan, intent nlp.Intent) (*ExecuteResult, error) {
	switch intent {
	case nlp.IntentCreate:
		return e.executeCreate(ctx, plan)
	case nlp.IntentUpdate:
		return e.executeUpdate(ctx, plan)
	case nlp.IntentFilter:
		return e.executeFilter(ctx, plan)
	case nlp.IntentUnknown:
		return nil, errors.New("unknown intent: could not determine command type")
	default:
		return nil, fmt.Errorf("unsupported intent: %v", intent)
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
		Message:   formatTaskList(tasks),
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

func formatTaskList(tasks []*store.Task) string {
	if len(tasks) == 0 {
		return "No tasks found"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d task(s):\n\n", len(tasks)))
	for _, task := range tasks {
		state := task.State
		if state == "" {
			state = "inbox"
		}
		sb.WriteString(fmt.Sprintf("  #%d %s [%s]\n", task.ID, task.Title, state))
	}
	return sb.String()
}
