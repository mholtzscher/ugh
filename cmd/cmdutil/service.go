package cmdutil

import (
	"context"
	"fmt"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

type ServiceFactory func(context.Context) (service.Service, error)

type ServiceAction func(context.Context, *cli.Command, service.Service) error

type SyncHook func(context.Context, service.Service) error

type SyncLabels struct {
	Before string
	After  string
}

func DefaultSyncLabels() SyncLabels {
	return SyncLabels{Before: "sync pull", After: "sync push"}
}

func WithService(factory ServiceFactory, action ServiceAction) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		svc, err := factory(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()
		return action(ctx, cmd, svc)
	}
}

func WithWriteSync(before SyncHook, after SyncHook, labels SyncLabels, action ServiceAction) ServiceAction {
	return func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
		if before != nil {
			if err := before(ctx, svc); err != nil {
				prefix := labels.Before
				if prefix == "" {
					prefix = "sync before write"
				}
				return fmt.Errorf("%s: %w", prefix, err)
			}
		}

		if err := action(ctx, cmd, svc); err != nil {
			return err
		}

		if after != nil {
			if err := after(ctx, svc); err != nil {
				prefix := labels.After
				if prefix == "" {
					prefix = "sync after write"
				}
				return fmt.Errorf("%s: %w", prefix, err)
			}
		}
		return nil
	}
}

func WithServiceAndWriteSync(factory ServiceFactory, before SyncHook, after SyncHook, action ServiceAction) cli.ActionFunc {
	return WithService(factory, WithWriteSync(before, after, DefaultSyncLabels(), action))
}
