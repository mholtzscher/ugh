package service

import (
	"context"

	"github.com/mholtzscher/ugh/internal/store"
)

// Service defines the interface for task operations.
// This is implemented by both TaskService (direct DB) and APIService (HTTP client).
type Service interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) (*store.Task, error)
	ListTasks(ctx context.Context, req ListTasksRequest) ([]*store.Task, error)
	GetTask(ctx context.Context, id int64) (*store.Task, error)
	UpdateTask(ctx context.Context, req UpdateTaskRequest) (*store.Task, error)
	FullUpdateTask(ctx context.Context, req FullUpdateTaskRequest) (*store.Task, error)
	SetDone(ctx context.Context, ids []int64, done bool) (int64, error)
	DeleteTasks(ctx context.Context, ids []int64) (int64, error)
	ListProjects(ctx context.Context, req ListTagsRequest) ([]store.NameCount, error)
	ListContexts(ctx context.Context, req ListTagsRequest) ([]store.NameCount, error)
	Sync(ctx context.Context) error
	Push(ctx context.Context) error
	Close() error
}

// Ensure TaskService implements Service
var _ Service = (*TaskService)(nil)
