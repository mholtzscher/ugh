package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mholtzscher/ugh/internal/service"

	"github.com/urfave/cli/v2"
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

// newService returns a Service implementation using direct database access.
func newService(ctx context.Context, c *cli.Context) (service.Service, error) {
	st, err := openStore(ctx, c)
	if err != nil {
		return nil, err
	}
	return service.NewTaskService(st), nil
}

func autoSyncEnabled() bool {
	return loadedConfig != nil && loadedConfig.DB.SyncOnWrite && loadedConfig.DB.SyncURL != ""
}

func maybeSyncBeforeWrite(ctx context.Context, svc service.Service) error {
	if !autoSyncEnabled() {
		return nil
	}
	return svc.Sync(ctx)
}

func maybeSyncAfterWrite(ctx context.Context, svc service.Service) error {
	if !autoSyncEnabled() {
		return nil
	}
	return svc.Push(ctx)
}
