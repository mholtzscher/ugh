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

type UpdateTaskRequest struct {
	ID   int64
	Text string
}

func (s *TaskService) UpdateTaskText(ctx context.Context, req UpdateTaskRequest) (*store.Task, error) {
	current, err := s.store.GetTask(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	parsed := todotxt.ParseLine(req.Text)
	if parsed.CreationDate == nil {
		if current.CreationDate != nil {
			parsed.CreationDate = current.CreationDate
		} else if parsed.Done && parsed.CompletionDate != nil {
			parsed.CreationDate = parsed.CompletionDate
		} else {
			parsed.CreationDate = nowDate()
		}
	}

	updated := &store.Task{
		ID:             current.ID,
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

	return s.store.UpdateTask(ctx, updated)
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
