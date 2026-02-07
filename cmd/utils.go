package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"

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

func commandArgs(cmd *cli.Command) []string {
	if cmd == nil {
		return nil
	}
	argsLen := cmd.Args().Len()
	args := make([]string, argsLen)
	for i := range argsLen {
		args[i] = cmd.Args().Get(i)
	}
	return args
}

// newService returns a Service implementation using direct database access.
func newService(ctx context.Context) (service.Service, error) {
	st, err := openStore(ctx)
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
