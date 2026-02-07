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
