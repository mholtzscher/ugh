package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/service"
)

// newService returns a Service implementation using direct database access.
func (a *App) newService(ctx context.Context) (service.Service, error) {
	return a.runtime.NewService(ctx)
}

func (a *App) maybeSyncBeforeWrite(ctx context.Context, svc service.Service) error {
	return a.runtime.MaybeSyncBeforeWrite(ctx, svc)
}

func (a *App) maybeSyncAfterWrite(ctx context.Context, svc service.Service) error {
	return a.runtime.MaybeSyncAfterWrite(ctx, svc)
}
