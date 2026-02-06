package tasks

import (
	"context"

	"github.com/mholtzscher/ugh/cmd/registry"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"

	"github.com/urfave/cli/v3"
)

// Deps holds dependencies injected from the parent cmd package.
type Deps struct {
	NewService           func(context.Context) (service.Service, error)
	MaybeSyncBeforeWrite func(context.Context, service.Service) error
	MaybeSyncAfterWrite  func(context.Context, service.Service) error
	OutputWriter         func() output.Writer
}

const (
	addID  registry.ID = "tasks.add"
	listID registry.ID = "tasks.list"
	showID registry.ID = "tasks.show"
	editID registry.ID = "tasks.edit"
	doneID registry.ID = "tasks.done"
	undoID registry.ID = "tasks.undo"
	rmID   registry.ID = "tasks.rm"
)

// Register adds task command specs to the registry.
func Register(r *registry.Registry, d Deps) error {
	return r.AddAll(
		registry.Spec{ID: addID, Source: "cmd/tasks", Build: func() *cli.Command { return newAddCmd(d) }},
		registry.Spec{ID: listID, Source: "cmd/tasks", Build: func() *cli.Command { return newListCmd(d) }},
		registry.Spec{ID: showID, Source: "cmd/tasks", Build: func() *cli.Command { return newShowCmd(d) }},
		registry.Spec{ID: editID, Source: "cmd/tasks", Build: func() *cli.Command { return newEditCmd(d) }},
		registry.Spec{ID: doneID, Source: "cmd/tasks", Build: func() *cli.Command { return newDoneCmd(d) }},
		registry.Spec{ID: undoID, Source: "cmd/tasks", Build: func() *cli.Command { return newUndoCmd(d) }},
		registry.Spec{ID: rmID, Source: "cmd/tasks", Build: func() *cli.Command { return newRmCmd(d) }},
	)
}
