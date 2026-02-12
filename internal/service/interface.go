package service

import (
	"context"

	"github.com/mholtzscher/ugh/internal/store"
)

// Service defines the interface for task operations.
// This is implemented by TaskService.
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
	SyncStatus(ctx context.Context) (*SyncStatus, error)
	Close() error

	// Shell history operations
	RecordShellHistory(
		ctx context.Context, command string, success bool, summary string, intent string,
	) (*store.ShellHistory, error)
	ListShellHistory(ctx context.Context, limit int64) ([]*store.ShellHistory, error)
	SearchShellHistory(
		ctx context.Context, search, intent string, success *bool, limit int64,
	) ([]*store.ShellHistory, error)
	UpdateShellHistory(ctx context.Context, id int64, success bool, summary string, intent string) error
	ClearShellHistory(ctx context.Context) error
	ListTaskEvents(ctx context.Context, taskID int64, limit int64) ([]*store.TaskEvent, error)
}

// Ensure TaskService implements Service.
var _ Service = (*TaskService)(nil)
