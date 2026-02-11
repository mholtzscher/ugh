package compile

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/domain"
	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/service"
)

type Plan struct {
	Intent nlp.Intent

	Create *service.CreateTaskRequest
	Update *service.UpdateTaskRequest
	Filter *service.ListTasksRequest

	Target nlp.TargetRef
}

type BuildOptions struct {
	SelectedTaskID *int64
	Now            time.Time
}

const (
	splitNParts     = 2
	nextWeekDaySpan = 7
)

func Build(result nlp.ParseResult, opts BuildOptions) (Plan, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}

	switch cmd := result.Command.(type) {
	case *nlp.CreateCommand:
		req, err := buildCreateRequest(cmd, opts)
		if err != nil {
			return Plan{}, err
		}
		return Plan{Intent: nlp.IntentCreate, Create: &req}, nil
	case *nlp.UpdateCommand:
		req, target, err := buildUpdateRequest(cmd, opts)
		if err != nil {
			return Plan{}, err
		}
		return Plan{Intent: nlp.IntentUpdate, Update: &req, Target: target}, nil
	case *nlp.FilterCommand:
		req, err := buildFilterRequest(cmd, opts)
		if err != nil {
			return Plan{}, err
		}
		return Plan{Intent: nlp.IntentFilter, Filter: &req}, nil
	case *nlp.ViewCommand:
		return Plan{Intent: nlp.IntentView}, nil
	case *nlp.ContextCommand:
		return Plan{Intent: nlp.IntentContext}, nil
	default:
		return Plan{}, fmt.Errorf("unsupported parse command type %T", result.Command)
	}
}

func buildCreateRequest(cmd *nlp.CreateCommand, opts BuildOptions) (service.CreateTaskRequest, error) {
	req := service.CreateTaskRequest{Title: strings.TrimSpace(cmd.Title), State: domain.TaskStateInbox}

	for _, op := range cmd.Ops {
		switch typed := op.(type) {
		case nlp.SetOp:
			if err := applyCreateSet(&req, typed, opts.Now); err != nil {
				return service.CreateTaskRequest{}, err
			}
		case nlp.AddOp:
			if err := applyCreateAdd(&req, typed); err != nil {
				return service.CreateTaskRequest{}, err
			}
		case nlp.RemoveOp:
			return service.CreateTaskRequest{}, errors.New("remove operations are not supported during create")
		case nlp.ClearOp:
			if err := applyCreateClear(&req, typed); err != nil {
				return service.CreateTaskRequest{}, err
			}
		case nlp.TagOp:
			applyTag(&req, typed)
		default:
			return service.CreateTaskRequest{}, fmt.Errorf("unsupported create op type %T", op)
		}
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return service.CreateTaskRequest{}, errors.New("title is required")
	}

	return req, nil
}

//nolint:gocognit // update operation compilation is intentionally explicit by op type.
func buildUpdateRequest(cmd *nlp.UpdateCommand, opts BuildOptions) (service.UpdateTaskRequest, nlp.TargetRef, error) {
	resolvedTarget := nlp.TargetRef{Kind: nlp.TargetSelected}
	if cmd.Target != nil {
		resolvedTarget = *cmd.Target
	}
	if resolvedTarget.Kind == nlp.TargetSelected {
		if opts.SelectedTaskID == nil || *opts.SelectedTaskID <= 0 {
			return service.UpdateTaskRequest{}, nlp.TargetRef{}, errors.New("selected target requires SelectedTaskID")
		}
		resolvedTarget = nlp.TargetRef{Kind: nlp.TargetID, ID: *opts.SelectedTaskID}
	}

	if resolvedTarget.Kind != nlp.TargetID || resolvedTarget.ID <= 0 {
		return service.UpdateTaskRequest{}, nlp.TargetRef{}, errors.New("update target must resolve to a task id")
	}

	req := service.UpdateTaskRequest{
		ID:      resolvedTarget.ID,
		SetMeta: map[string]string{},
	}

	for _, op := range cmd.Ops {
		switch typed := op.(type) {
		case nlp.SetOp:
			if err := applyUpdateSet(&req, typed, opts.Now); err != nil {
				return service.UpdateTaskRequest{}, nlp.TargetRef{}, err
			}
		case nlp.AddOp:
			if err := applyUpdateAdd(&req, typed); err != nil {
				return service.UpdateTaskRequest{}, nlp.TargetRef{}, err
			}
		case nlp.RemoveOp:
			if err := applyUpdateRemove(&req, typed); err != nil {
				return service.UpdateTaskRequest{}, nlp.TargetRef{}, err
			}
		case nlp.ClearOp:
			if err := applyUpdateClear(&req, typed); err != nil {
				return service.UpdateTaskRequest{}, nlp.TargetRef{}, err
			}
		case nlp.TagOp:
			applyUpdateTag(&req, typed)
		default:
			return service.UpdateTaskRequest{}, nlp.TargetRef{}, fmt.Errorf("unsupported update op type %T", op)
		}
	}

	req.AddProjects = unique(req.AddProjects)
	req.AddContexts = unique(req.AddContexts)
	req.RemoveProjects = unique(req.RemoveProjects)
	req.RemoveContexts = unique(req.RemoveContexts)
	req.RemoveMetaKeys = unique(req.RemoveMetaKeys)

	return req, resolvedTarget, nil
}

func buildFilterRequest(cmd *nlp.FilterCommand, opts BuildOptions) (service.ListTasksRequest, error) {
	expr, err := NormalizeFilterExpr(cmd.Expr, opts)
	if err != nil {
		return service.ListTasksRequest{}, err
	}
	return service.ListTasksRequest{Filter: expr}, nil
}

func NormalizeFilterExpr(expr nlp.FilterExpr, opts BuildOptions) (nlp.FilterExpr, error) {
	if expr == nil {
		return expr, nil
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}
	return compileFilterExpr(expr, opts)
}

func compileFilterExpr(expr nlp.FilterExpr, opts BuildOptions) (nlp.FilterExpr, error) {
	switch typed := expr.(type) {
	case nlp.Predicate:
		return compilePredicate(typed, opts)
	case nlp.FilterBinary:
		left, err := compileFilterExpr(typed.Left, opts)
		if err != nil {
			return nil, err
		}
		right, err := compileFilterExpr(typed.Right, opts)
		if err != nil {
			return nil, err
		}
		return nlp.FilterBinary{Op: typed.Op, Left: left, Right: right}, nil
	case nlp.FilterNot:
		inner, err := compileFilterExpr(typed.Expr, opts)
		if err != nil {
			return nil, err
		}
		return nlp.FilterNot{Expr: inner}, nil
	default:
		return nil, fmt.Errorf("unsupported filter expression type %T", expr)
	}
}

func compilePredicate(pred nlp.Predicate, opts BuildOptions) (nlp.Predicate, error) {
	compiled := pred
	compiled.Text = strings.TrimSpace(pred.Text)

	switch pred.Kind {
	case nlp.PredState:
		state, err := domain.NormalizeState(compiled.Text)
		if err != nil {
			return nlp.Predicate{}, err
		}
		compiled.Text = state
	case nlp.PredDue:
		if compiled.Text == "" {
			return compiled, nil
		}
		dueDate, err := normalizeDate(compiled.Text, opts.Now)
		if err != nil {
			return nlp.Predicate{}, err
		}
		compiled.Text = dueDate
	case nlp.PredProject, nlp.PredContext, nlp.PredText:
		if compiled.Text == "" {
			return nlp.Predicate{}, errors.New("filter value cannot be empty")
		}
	case nlp.PredID:
		id, err := strconv.ParseInt(compiled.Text, 10, 64)
		if err != nil || id <= 0 {
			return nlp.Predicate{}, fmt.Errorf("invalid id filter %q", pred.Text)
		}
		compiled.Text = strconv.FormatInt(id, 10)
	default:
		return nlp.Predicate{}, fmt.Errorf("unsupported predicate kind %v", pred.Kind)
	}

	return compiled, nil
}

func applyCreateSet(req *service.CreateTaskRequest, op nlp.SetOp, now time.Time) error {
	value := strings.TrimSpace(string(op.Value))
	switch op.Field {
	case nlp.FieldTitle:
		req.Title = value
	case nlp.FieldNotes:
		req.Notes = value
	case nlp.FieldDue:
		due, err := normalizeDate(value, now)
		if err != nil {
			return err
		}
		req.DueOn = due
	case nlp.FieldWaiting:
		req.WaitingFor = value
	case nlp.FieldState:
		state, err := domain.NormalizeState(value)
		if err != nil {
			return err
		}
		req.State = state
	case nlp.FieldProjects:
		req.Projects = parseList(value)
	case nlp.FieldContexts:
		req.Contexts = parseList(value)
	case nlp.FieldMeta:
		meta := parseList(value)
		req.Meta = unique(append([]string(nil), meta...))
	default:
		return fmt.Errorf("unsupported create set field %v", op.Field)
	}
	return nil
}

func applyCreateAdd(req *service.CreateTaskRequest, op nlp.AddOp) error {
	value := strings.TrimSpace(string(op.Value))
	switch op.Field {
	case nlp.FieldTitle, nlp.FieldNotes, nlp.FieldDue, nlp.FieldWaiting, nlp.FieldState:
		return errors.New("+ supports projects/contexts/meta only")
	case nlp.FieldProjects:
		req.Projects = unique(append(req.Projects, parseList(value)...))
	case nlp.FieldContexts:
		req.Contexts = unique(append(req.Contexts, parseList(value)...))
	case nlp.FieldMeta:
		req.Meta = unique(append(req.Meta, parseList(value)...))
	default:
		return fmt.Errorf("unsupported add field %v", op.Field)
	}
	return nil
}

func applyCreateClear(req *service.CreateTaskRequest, op nlp.ClearOp) error {
	switch op.Field {
	case nlp.FieldTitle, nlp.FieldState:
		return fmt.Errorf("cannot clear field %v in create request", op.Field)
	case nlp.FieldNotes:
		req.Notes = ""
	case nlp.FieldDue:
		req.DueOn = ""
	case nlp.FieldWaiting:
		req.WaitingFor = ""
	case nlp.FieldProjects:
		req.Projects = nil
	case nlp.FieldContexts:
		req.Contexts = nil
	case nlp.FieldMeta:
		req.Meta = nil
	default:
		return fmt.Errorf("cannot clear field %v in create request", op.Field)
	}
	return nil
}

func applyTag(req *service.CreateTaskRequest, op nlp.TagOp) {
	if op.Kind == nlp.TagProject {
		req.Projects = unique(append(req.Projects, strings.TrimSpace(op.Value)))
		return
	}
	req.Contexts = unique(append(req.Contexts, strings.TrimSpace(op.Value)))
}

func applyUpdateSet(req *service.UpdateTaskRequest, op nlp.SetOp, now time.Time) error {
	value := strings.TrimSpace(string(op.Value))
	switch op.Field {
	case nlp.FieldTitle:
		req.Title = ptr(value)
	case nlp.FieldNotes:
		req.Notes = ptr(value)
	case nlp.FieldDue:
		due, err := normalizeDate(value, now)
		if err != nil {
			return err
		}
		req.DueOn = ptr(due)
		req.ClearDueOn = false
	case nlp.FieldWaiting:
		req.WaitingFor = ptr(value)
		req.ClearWaitingFor = false
	case nlp.FieldState:
		state, err := domain.NormalizeState(value)
		if err != nil {
			return err
		}
		req.State = ptr(state)
	case nlp.FieldMeta:
		k, v, err := parseMetaValue(value)
		if err != nil {
			return err
		}
		req.SetMeta[k] = v
	case nlp.FieldProjects, nlp.FieldContexts:
		return fmt.Errorf("set %q is not supported; use + or - operations", op.Field)
	default:
		return fmt.Errorf("unsupported set field %v", op.Field)
	}
	return nil
}

func applyUpdateAdd(req *service.UpdateTaskRequest, op nlp.AddOp) error {
	value := strings.TrimSpace(string(op.Value))
	switch op.Field {
	case nlp.FieldTitle, nlp.FieldNotes, nlp.FieldDue, nlp.FieldWaiting, nlp.FieldState:
		return fmt.Errorf("unsupported add field %v", op.Field)
	case nlp.FieldProjects:
		req.AddProjects = append(req.AddProjects, parseList(value)...)
	case nlp.FieldContexts:
		req.AddContexts = append(req.AddContexts, parseList(value)...)
	case nlp.FieldMeta:
		k, v, err := parseMetaValue(value)
		if err != nil {
			return err
		}
		req.SetMeta[k] = v
	default:
		return fmt.Errorf("unsupported add field %v", op.Field)
	}
	return nil
}

func applyUpdateRemove(req *service.UpdateTaskRequest, op nlp.RemoveOp) error {
	value := strings.TrimSpace(string(op.Value))
	switch op.Field {
	case nlp.FieldTitle, nlp.FieldNotes, nlp.FieldDue, nlp.FieldWaiting, nlp.FieldState:
		return fmt.Errorf("unsupported remove field %v", op.Field)
	case nlp.FieldProjects:
		req.RemoveProjects = append(req.RemoveProjects, parseList(value)...)
	case nlp.FieldContexts:
		req.RemoveContexts = append(req.RemoveContexts, parseList(value)...)
	case nlp.FieldMeta:
		for _, part := range parseList(value) {
			key := part
			if strings.Contains(part, ":") {
				key = strings.TrimSpace(strings.SplitN(part, ":", splitNParts)[0])
			}
			if key != "" {
				req.RemoveMetaKeys = append(req.RemoveMetaKeys, key)
			}
		}
	default:
		return fmt.Errorf("unsupported remove field %v", op.Field)
	}
	return nil
}

func applyUpdateClear(req *service.UpdateTaskRequest, op nlp.ClearOp) error {
	switch op.Field {
	case nlp.FieldTitle, nlp.FieldState:
		return fmt.Errorf("unsupported clear field %v", op.Field)
	case nlp.FieldDue:
		req.ClearDueOn = true
		req.DueOn = nil
	case nlp.FieldWaiting:
		req.ClearWaitingFor = true
		req.WaitingFor = nil
	case nlp.FieldNotes:
		req.Notes = ptr("")
	case nlp.FieldProjects, nlp.FieldContexts, nlp.FieldMeta:
		return fmt.Errorf("clear %v is not supported in patch updates", op.Field)
	default:
		return fmt.Errorf("unsupported clear field %v", op.Field)
	}
	return nil
}

func applyUpdateTag(req *service.UpdateTaskRequest, op nlp.TagOp) {
	if op.Kind == nlp.TagProject {
		req.AddProjects = append(req.AddProjects, strings.TrimSpace(op.Value))
		return
	}
	req.AddContexts = append(req.AddContexts, strings.TrimSpace(op.Value))
}

func normalizeDate(value string, now time.Time) (string, error) {
	lower := strings.ToLower(strings.TrimSpace(value))
	day := now
	switch lower {
	case "today":
		return day.Format(domain.DateLayoutYYYYMMDD), nil
	case "tomorrow":
		return day.AddDate(0, 0, 1).Format(domain.DateLayoutYYYYMMDD), nil
	case "next-week":
		return day.AddDate(0, 0, nextWeekDaySpan).Format(domain.DateLayoutYYYYMMDD), nil
	default:
		if _, err := time.Parse(domain.DateLayoutYYYYMMDD, lower); err != nil {
			return "", domain.InvalidDateFormatError(value)
		}
		return lower, nil
	}
}

func parseList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return unique(out)
}

func parseMetaValue(value string) (string, string, error) {
	k, v, ok := strings.Cut(value, domain.MetaSeparatorColon)
	if !ok {
		return "", "", domain.InvalidMetaFormatError(value)
	}
	k = strings.TrimSpace(k)
	v = strings.TrimSpace(v)
	if k == "" {
		return "", "", domain.InvalidMetaFormatError(value)
	}
	return k, v, nil
}

func unique(values []string) []string {
	if len(values) == 0 {
		return values
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || slices.Contains(result, value) {
			continue
		}
		result = append(result, value)
	}
	return result
}

func ptr(value string) *string {
	return &value
}
