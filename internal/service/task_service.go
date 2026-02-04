package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
)

type TaskService struct {
	store *store.Store
}

func NewTaskService(store *store.Store) *TaskService {
	return &TaskService{
		store: store,
	}
}

type CreateTaskRequest struct {
	Title      string
	Notes      string
	Status     string
	Priority   string
	Projects   []string
	Contexts   []string
	Meta       []string
	DueOn      string
	DeferUntil string
	WaitingFor string
	Done       bool
}

func (s *TaskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*store.Task, error) {
	meta, err := parseMetaFlags(req.Meta)
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	status, err := normalizeStatus(req.Status)
	if err != nil {
		return nil, err
	}

	var dueOn *time.Time
	if strings.TrimSpace(req.DueOn) != "" {
		parsed, err := parseDay(req.DueOn)
		if err != nil {
			return nil, err
		}
		dueOn = parsed
	}
	var deferUntil *time.Time
	if strings.TrimSpace(req.DeferUntil) != "" {
		parsed, err := parseDay(req.DeferUntil)
		if err != nil {
			return nil, err
		}
		deferUntil = parsed
	}

	task := &store.Task{
		Done:       req.Done,
		Status:     status,
		Priority:   normalizePriority(req.Priority),
		Title:      req.Title,
		Notes:      req.Notes,
		DueOn:      dueOn,
		DeferUntil: deferUntil,
		WaitingFor: strings.TrimSpace(req.WaitingFor),
		Projects:   req.Projects,
		Contexts:   req.Contexts,
		Meta:       meta,
	}

	return s.store.CreateTask(ctx, task)
}

type ListTasksRequest struct {
	All             bool
	DoneOnly        bool
	TodoOnly        bool
	Status          string
	Project         string
	Context         string
	Priority        string
	Search          string
	DueOnly         bool
	DeferAfter      string
	DeferOnOrBefore string
}

func (s *TaskService) ListTasks(ctx context.Context, req ListTasksRequest) ([]*store.Task, error) {
	filters := store.Filters{
		All:             req.All,
		DoneOnly:        req.DoneOnly,
		TodoOnly:        req.TodoOnly,
		Status:          strings.TrimSpace(req.Status),
		Project:         req.Project,
		Context:         req.Context,
		Priority:        req.Priority,
		Search:          req.Search,
		DueSetOnly:      req.DueOnly,
		DeferAfter:      strings.TrimSpace(req.DeferAfter),
		DeferOnOrBefore: strings.TrimSpace(req.DeferOnOrBefore),
	}

	if !filters.All && !filters.DoneOnly && !filters.TodoOnly {
		filters.TodoOnly = true
	}

	return s.store.ListTasks(ctx, filters)
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*store.Task, error) {
	return s.store.GetTask(ctx, id)
}

func (s *TaskService) SetDone(ctx context.Context, ids []int64, done bool) (int64, error) {
	return s.store.SetDone(ctx, ids, done)
}

func (s *TaskService) DeleteTasks(ctx context.Context, ids []int64) (int64, error) {
	return s.store.DeleteTasks(ctx, ids)
}

func (s *TaskService) Close() error {
	return s.store.Close()
}

func (s *TaskService) Sync(ctx context.Context) error {
	return s.store.Sync(ctx)
}

func (s *TaskService) Push(ctx context.Context) error {
	return s.store.Push(ctx)
}

type ListTagsRequest struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
}

func (s *TaskService) ListProjects(ctx context.Context, req ListTagsRequest) ([]store.NameCount, error) {
	var status any
	if req.DoneOnly {
		status = int64(1)
	} else if req.TodoOnly {
		status = int64(0)
	}

	if !req.All && !req.DoneOnly && !req.TodoOnly {
		status = int64(0)
	}

	return s.store.ListProjectCounts(ctx, status)
}

func (s *TaskService) ListContexts(ctx context.Context, req ListTagsRequest) ([]store.NameCount, error) {
	var status any
	if req.DoneOnly {
		status = int64(1)
	} else if req.TodoOnly {
		status = int64(0)
	}

	if !req.All && !req.DoneOnly && !req.TodoOnly {
		status = int64(0)
	}

	return s.store.ListContextCounts(ctx, status)
}

type UpdateTaskRequest struct {
	ID              int64
	Title           *string
	Notes           *string
	Status          *string
	Priority        *string
	Done            *bool
	DueOn           *string
	DeferUntil      *string
	WaitingFor      *string
	AddProjects     []string
	AddContexts     []string
	SetMeta         map[string]string
	RemoveProjects  []string
	RemoveContexts  []string
	RemoveMetaKeys  []string
	RemovePriority  bool
	ClearDueOn      bool
	ClearDeferUntil bool
	ClearWaitingFor bool
}

func (s *TaskService) UpdateTask(ctx context.Context, req UpdateTaskRequest) (*store.Task, error) {
	current, err := s.store.GetTask(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	updated := &store.Task{
		ID:          current.ID,
		Done:        current.Done,
		Status:      current.Status,
		Priority:    current.Priority,
		Title:       current.Title,
		Notes:       current.Notes,
		DueOn:       current.DueOn,
		DeferUntil:  current.DeferUntil,
		WaitingFor:  current.WaitingFor,
		CompletedAt: current.CompletedAt,
		Projects:    append([]string(nil), current.Projects...),
		Contexts:    append([]string(nil), current.Contexts...),
		Meta:        copyMeta(current.Meta),
	}

	if req.Title != nil {
		updated.Title = *req.Title
	}
	if req.Notes != nil {
		updated.Notes = *req.Notes
	}
	if req.Status != nil {
		status, err := normalizeStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		updated.Status = status
	}
	if req.RemovePriority {
		updated.Priority = ""
	} else if req.Priority != nil {
		updated.Priority = normalizePriority(*req.Priority)
	}
	if req.Done != nil {
		updated.Done = *req.Done
		if !*req.Done {
			updated.CompletedAt = nil
		}
	}
	if req.ClearDueOn {
		updated.DueOn = nil
	} else if req.DueOn != nil {
		if strings.TrimSpace(*req.DueOn) == "" {
			updated.DueOn = nil
		} else {
			parsed, err := parseDay(*req.DueOn)
			if err != nil {
				return nil, err
			}
			updated.DueOn = parsed
		}
	}
	if req.ClearDeferUntil {
		updated.DeferUntil = nil
	} else if req.DeferUntil != nil {
		if strings.TrimSpace(*req.DeferUntil) == "" {
			updated.DeferUntil = nil
		} else {
			parsed, err := parseDay(*req.DeferUntil)
			if err != nil {
				return nil, err
			}
			updated.DeferUntil = parsed
		}
	}
	if req.ClearWaitingFor {
		updated.WaitingFor = ""
	} else if req.WaitingFor != nil {
		updated.WaitingFor = strings.TrimSpace(*req.WaitingFor)
	}

	for _, p := range req.AddProjects {
		if !containsString(updated.Projects, p) {
			updated.Projects = append(updated.Projects, p)
		}
	}
	updated.Projects = removeStrings(updated.Projects, req.RemoveProjects)

	for _, c := range req.AddContexts {
		if !containsString(updated.Contexts, c) {
			updated.Contexts = append(updated.Contexts, c)
		}
	}
	updated.Contexts = removeStrings(updated.Contexts, req.RemoveContexts)

	if len(req.SetMeta) > 0 {
		if updated.Meta == nil {
			updated.Meta = map[string]string{}
		}
		for k, v := range req.SetMeta {
			updated.Meta[k] = v
		}
	}
	for _, k := range req.RemoveMetaKeys {
		delete(updated.Meta, k)
	}

	return s.store.UpdateTask(ctx, updated)
}

type FullUpdateTaskRequest struct {
	ID         int64
	Title      string
	Notes      string
	Status     string
	Priority   string
	Done       bool
	Projects   []string
	Contexts   []string
	Meta       map[string]string
	DueOn      string
	DeferUntil string
	WaitingFor string
}

func (s *TaskService) FullUpdateTask(ctx context.Context, req FullUpdateTaskRequest) (*store.Task, error) {
	current, err := s.store.GetTask(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	status, err := normalizeStatus(req.Status)
	if err != nil {
		return nil, err
	}
	var dueOn *time.Time
	if strings.TrimSpace(req.DueOn) != "" {
		parsed, err := parseDay(req.DueOn)
		if err != nil {
			return nil, err
		}
		dueOn = parsed
	}
	var deferUntil *time.Time
	if strings.TrimSpace(req.DeferUntil) != "" {
		parsed, err := parseDay(req.DeferUntil)
		if err != nil {
			return nil, err
		}
		deferUntil = parsed
	}

	updated := &store.Task{
		ID:          current.ID,
		Done:        req.Done,
		Status:      status,
		Priority:    normalizePriority(req.Priority),
		Title:       req.Title,
		Notes:       req.Notes,
		DueOn:       dueOn,
		DeferUntil:  deferUntil,
		WaitingFor:  strings.TrimSpace(req.WaitingFor),
		Projects:    req.Projects,
		Contexts:    req.Contexts,
		Meta:        req.Meta,
		CompletedAt: current.CompletedAt,
	}
	if !req.Done {
		updated.CompletedAt = nil
	}

	return s.store.UpdateTask(ctx, updated)
}

func copyMeta(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeStrings(slice []string, toRemove []string) []string {
	if len(toRemove) == 0 {
		return slice
	}
	removeSet := make(map[string]bool, len(toRemove))
	for _, s := range toRemove {
		removeSet[s] = true
	}
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if !removeSet[s] {
			result = append(result, s)
		}
	}
	return result
}

func normalizePriority(p string) string {
	return strings.ToUpper(strings.TrimSpace(p))
}

func parseMetaFlags(meta []string) (map[string]string, error) {
	result := map[string]string{}
	for _, m := range meta {
		k, v, ok := strings.Cut(m, ":")
		if !ok {
			return nil, fmt.Errorf("invalid meta format: %s (expected key:value)", m)
		}
		result[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return result, nil
}

func parseDay(value string) (*time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", value)
	}
	utc := parsed.UTC()
	return &utc, nil
}

func normalizeStatus(value string) (store.Status, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return store.StatusInbox, nil
	}
	switch value {
	case string(store.StatusInbox), string(store.StatusNext), string(store.StatusWaiting), string(store.StatusSomeday):
		return store.Status(value), nil
	default:
		return "", fmt.Errorf("invalid status %q (expected inbox|next|waiting|someday)", value)
	}
}
