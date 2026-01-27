package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mholtzscher/ugh/internal/service"
)

func parseIDs(args []string) ([]int64, error) {
	ids := make([]int64, 0, len(args))
	for _, arg := range args {
		val, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid id %q", arg)
		}
		ids = append(ids, val)
	}
	if len(ids) == 0 {
		return nil, errors.New("at least one id is required")
	}
	return ids, nil
}

func newTaskService(ctx context.Context) (*service.TaskService, error) {
	st, err := openStore(ctx)
	if err != nil {
		return nil, err
	}
	return service.NewTaskService(st), nil
}

func autoSyncEnabled() bool {
	return loadedConfig.DB.SyncOnWrite && loadedConfig.DB.SyncURL != ""
}

func maybeSyncBeforeWrite(ctx context.Context, svc *service.TaskService) error {
	if !autoSyncEnabled() {
		return nil
	}
	return svc.Sync(ctx)
}

func maybeSyncAfterWrite(ctx context.Context, svc *service.TaskService) error {
	if !autoSyncEnabled() {
		return nil
	}
	return svc.Push(ctx)
}
