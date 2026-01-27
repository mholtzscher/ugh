package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
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
	Line      string
	Priority  string
	Projects  []string
	Contexts  []string
	Meta      []string
	Done      bool
	Created   string
	Completed string
}

func (s *TaskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*store.Task, error) {
	parsed := todotxt.ParseLine(req.Line)

	if req.Priority != "" {
		parsed.Priority = normalizePriority(req.Priority)
	}
	if len(req.Projects) > 0 {
		parsed.Projects = append(parsed.Projects, req.Projects...)
	}
	if len(req.Contexts) > 0 {
		parsed.Contexts = append(parsed.Contexts, req.Contexts...)
	}
	if req.Done {
		parsed.Done = true
		if parsed.CompletionDate == nil {
			parsed.CompletionDate = nowDate()
		}
	}
	if req.Created != "" {
		date, err := parseDate(req.Created)
		if err != nil {
			return nil, err
		}
		parsed.CreationDate = date
	}
	if req.Completed != "" {
		date, err := parseDate(req.Completed)
		if err != nil {
			return nil, err
		}
		parsed.CompletionDate = date
	}
	if parsed.CreationDate == nil {
		if parsed.Done && parsed.CompletionDate != nil {
			parsed.CreationDate = parsed.CompletionDate
		} else {
			parsed.CreationDate = nowDate()
		}
	}

	meta, err := parseMetaFlags(req.Meta)
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}
	if len(meta) > 0 {
		if parsed.Meta == nil {
			parsed.Meta = map[string]string{}
		}
		for key, value := range meta {
			parsed.Meta[key] = value
		}
	}

	task := &store.Task{
		Done:           parsed.Done,
		Priority:       parsed.Priority,
		CompletionDate: parsed.CompletionDate,
		CreationDate:   parsed.CreationDate,
		Description:    parsed.Description,
		Projects:       parsed.Projects,
		Contexts:       parsed.Contexts,
		Meta:           parsed.Meta,
		Unknown:        parsed.Unknown,
	}

	return s.store.CreateTask(ctx, task)
}

type ListTasksRequest struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
	Project  string
	Context  string
	Priority string
	Search   string
}

func (s *TaskService) ListTasks(ctx context.Context, req ListTasksRequest) ([]*store.Task, error) {
	filters := store.Filters{
		All:      req.All,
		DoneOnly: req.DoneOnly,
		TodoOnly: req.TodoOnly,
		Project:  req.Project,
		Context:  req.Context,
		Priority: req.Priority,
		Search:   req.Search,
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
	ID             int64
	Description    *string
	Priority       *string
	Done           *bool
	AddProjects    []string
	AddContexts    []string
	SetMeta        map[string]string
	RemoveProjects []string
	RemoveContexts []string
	RemoveMetaKeys []string
	RemovePriority bool
}

func (s *TaskService) UpdateTask(ctx context.Context, req UpdateTaskRequest) (*store.Task, error) {
	current, err := s.store.GetTask(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	updated := &store.Task{
		ID:             current.ID,
		Done:           current.Done,
		Priority:       current.Priority,
		CompletionDate: current.CompletionDate,
		CreationDate:   current.CreationDate,
		Description:    current.Description,
		Projects:       append([]string(nil), current.Projects...),
		Contexts:       append([]string(nil), current.Contexts...),
		Meta:           copyMeta(current.Meta),
		Unknown:        append([]string(nil), current.Unknown...),
	}

	if req.Description != nil {
		updated.Description = *req.Description
	}
	if req.RemovePriority {
		updated.Priority = ""
	} else if req.Priority != nil {
		updated.Priority = normalizePriority(*req.Priority)
	}
	if req.Done != nil {
		updated.Done = *req.Done
		if *req.Done && updated.CompletionDate == nil {
			updated.CompletionDate = nowDate()
		} else if !*req.Done {
			updated.CompletionDate = nil
		}
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
	ID          int64
	Description string
	Priority    string
	Done        bool
	Projects    []string
	Contexts    []string
	Meta        map[string]string
}

func (s *TaskService) FullUpdateTask(ctx context.Context, req FullUpdateTaskRequest) (*store.Task, error) {
	current, err := s.store.GetTask(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	updated := &store.Task{
		ID:           current.ID,
		Done:         req.Done,
		Priority:     normalizePriority(req.Priority),
		CreationDate: current.CreationDate,
		Description:  req.Description,
		Projects:     req.Projects,
		Contexts:     req.Contexts,
		Meta:         req.Meta,
		Unknown:      current.Unknown,
	}

	if req.Done && !current.Done {
		updated.CompletionDate = nowDate()
	} else if req.Done && current.Done {
		updated.CompletionDate = current.CompletionDate
	} else {
		updated.CompletionDate = nil
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

func parseDate(value string) (*time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", value)
	}
	utc := parsed.UTC()
	return &utc, nil
}

func nowDate() *time.Time {
	now := time.Now().UTC()
	return &now
}
