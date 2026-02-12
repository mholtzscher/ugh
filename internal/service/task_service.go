package service

import (
	"context"

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

func (s *TaskService) Close() error {
	return s.store.Close()
}

func (s *TaskService) Sync(ctx context.Context) error {
	return s.store.Sync(ctx)
}

func (s *TaskService) Push(ctx context.Context) error {
	return s.store.Push(ctx)
}

func (s *TaskService) SyncStatus(ctx context.Context) (*SyncStatus, error) {
	stats, err := s.store.SyncStats(ctx)
	if err != nil {
		return nil, err
	}
	return &SyncStatus{
		LastPullUnixTime: stats.LastPullUnixTime,
		LastPushUnixTime: stats.LastPushUnixTime,
		PendingChanges:   stats.CdcOperations,
		NetworkSentBytes: stats.NetworkSentBytes,
		NetworkRecvBytes: stats.NetworkReceivedBytes,
		Revision:         stats.Revision,
	}, nil
}

func (s *TaskService) RecordShellHistory(
	ctx context.Context, command string, success bool, summary string, intent string,
) (*store.ShellHistory, error) {
	return s.store.RecordShellHistory(ctx, command, success, summary, intent)
}

func (s *TaskService) ListShellHistory(ctx context.Context, limit int64) ([]*store.ShellHistory, error) {
	return s.store.ListShellHistory(ctx, limit)
}

func (s *TaskService) SearchShellHistory(
	ctx context.Context, search, intent string, success *bool, limit int64,
) ([]*store.ShellHistory, error) {
	return s.store.SearchShellHistory(ctx, search, intent, success, limit)
}

func (s *TaskService) ClearShellHistory(ctx context.Context) error {
	return s.store.ClearShellHistory(ctx)
}
