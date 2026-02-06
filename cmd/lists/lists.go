package lists

import (
	"context"

	"github.com/mholtzscher/ugh/cmd/cmdutil"
	"github.com/mholtzscher/ugh/cmd/meta"
	"github.com/mholtzscher/ugh/cmd/registry"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"

	"github.com/urfave/cli/v3"
)

// Deps holds dependencies injected from the parent cmd package.
type Deps struct {
	NewService   func(context.Context) (service.Service, error)
	OutputWriter func() output.Writer
}

const (
	inboxID    registry.ID = "lists.inbox"
	nowID      registry.ID = "lists.now"
	waitingID  registry.ID = "lists.waiting"
	laterID    registry.ID = "lists.later"
	calendarID registry.ID = "lists.calendar"
)

// Register adds list command specs to the registry.
func Register(r *registry.Registry, d Deps) error {
	return r.AddAll(
		registry.Spec{ID: inboxID, Source: "cmd/lists", Build: func() *cli.Command { return newInboxCmd(d) }},
		registry.Spec{ID: nowID, Source: "cmd/lists", Build: func() *cli.Command { return newNowCmd(d) }},
		registry.Spec{ID: waitingID, Source: "cmd/lists", Build: func() *cli.Command { return newWaitingCmd(d) }},
		registry.Spec{ID: laterID, Source: "cmd/lists", Build: func() *cli.Command { return newLaterCmd(d) }},
		registry.Spec{ID: calendarID, Source: "cmd/lists", Build: func() *cli.Command { return newCalendarCmd(d) }},
	)
}

func newInboxCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:     "inbox",
		Aliases:  []string{"i"},
		Usage:    "List inbox tasks",
		Category: meta.ListsCategory.String(),
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
				TodoOnly: true,
				State:    flags.TaskStateInbox,
			})
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			return writer.WriteTasks(tasks)
		}),
	}
}

func newNowCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:     "now",
		Aliases:  []string{"n"},
		Usage:    "List tasks you can act on now",
		Category: meta.ListsCategory.String(),
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
				TodoOnly: true,
				State:    flags.TaskStateNow,
			})
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			return writer.WriteTasks(tasks)
		}),
	}
}

func newWaitingCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:     "waiting",
		Aliases:  []string{"w"},
		Usage:    "List waiting-for items",
		Category: meta.ListsCategory.String(),
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
				TodoOnly: true,
				State:    flags.TaskStateWaiting,
			})
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			return writer.WriteTasks(tasks)
		}),
	}
}

func newLaterCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:     "later",
		Aliases:  []string{"sd"},
		Usage:    "List tasks you are not doing now",
		Category: meta.ListsCategory.String(),
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
				TodoOnly: true,
				State:    flags.TaskStateLater,
			})
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			return writer.WriteTasks(tasks)
		}),
	}
}

func newCalendarCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:     "calendar",
		Aliases:  []string{"cal"},
		Usage:    "List items with due dates",
		Category: meta.ListsCategory.String(),
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
				TodoOnly: true,
				DueOnly:  true,
			})
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			return writer.WriteTasks(tasks)
		}),
	}
}
