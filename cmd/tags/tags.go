package tags

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
	projectsID registry.ID = "tags.projects"
	contextsID registry.ID = "tags.contexts"
)

// Register adds projects/contexts command specs to the registry.
func Register(r *registry.Registry, d Deps) error {
	return r.AddAll(
		registry.Spec{ID: projectsID, Source: "cmd/tags", Build: func() *cli.Command { return newProjectsCmd(d) }},
		registry.Spec{ID: contextsID, Source: "cmd/tags", Build: func() *cli.Command { return newContextsCmd(d) }},
	)
}

func newProjectsCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:     "projects",
		Aliases:  []string{"proj"},
		Usage:    "List projects",
		Category: meta.ProjectsCategory.String(),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    flags.FlagAll,
				Aliases: []string{"a"},
				Usage:   "include completed tasks",
				Action: flags.BoolAction(
					flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
				),
			},
			&cli.BoolFlag{
				Name:    flags.FlagDone,
				Aliases: []string{"x"},
				Usage:   "only completed tasks",
				Action: flags.BoolAction(
					flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
				),
			},
			&cli.BoolFlag{
				Name:    flags.FlagTodo,
				Aliases: []string{"t"},
				Usage:   "only pending tasks",
				Action: flags.BoolAction(
					flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
				),
			},
			&cli.BoolFlag{
				Name:  flags.FlagCounts,
				Usage: "include counts",
			},
		},
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			tags, err := svc.ListProjects(ctx, service.ListTagsRequest{
				All:      cmd.Bool(flags.FlagAll),
				DoneOnly: cmd.Bool(flags.FlagDone),
				TodoOnly: cmd.Bool(flags.FlagTodo),
			})
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			if cmd.Bool(flags.FlagCounts) {
				return writer.WriteTagsWithCounts(tags)
			}
			return writer.WriteTags(tags)
		}),
	}
}

func newContextsCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:     "contexts",
		Aliases:  []string{"ctx"},
		Usage:    "List contexts",
		Category: meta.ProjectsCategory.String(),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    flags.FlagAll,
				Aliases: []string{"a"},
				Usage:   "include completed tasks",
				Action: flags.BoolAction(
					flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
				),
			},
			&cli.BoolFlag{
				Name:    flags.FlagDone,
				Aliases: []string{"x"},
				Usage:   "only completed tasks",
				Action: flags.BoolAction(
					flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
				),
			},
			&cli.BoolFlag{
				Name:    flags.FlagTodo,
				Aliases: []string{"t"},
				Usage:   "only pending tasks",
				Action: flags.BoolAction(
					flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
				),
			},
			&cli.BoolFlag{
				Name:  flags.FlagCounts,
				Usage: "include counts",
			},
		},
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			tags, err := svc.ListContexts(ctx, service.ListTagsRequest{
				All:      cmd.Bool(flags.FlagAll),
				DoneOnly: cmd.Bool(flags.FlagDone),
				TodoOnly: cmd.Bool(flags.FlagTodo),
			})
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			if cmd.Bool(flags.FlagCounts) {
				return writer.WriteTagsWithCounts(tags)
			}
			return writer.WriteTags(tags)
		}),
	}
}
