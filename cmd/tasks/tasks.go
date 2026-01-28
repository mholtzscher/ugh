package tasks

import (
	"context"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v2"
)

// Deps holds dependencies injected from the parent cmd package.
// This avoids circular imports between cmd and cmd/tasks.
type Deps struct {
	// WithTimeout returns a context with timeout.
	WithTimeout func(context.Context) (context.Context, context.CancelFunc)
	// NewService creates a new service instance.
	NewService func(*cli.Context) (service.Service, error)
	// ParseIDs parses ID arguments.
	ParseIDs func([]string) ([]int64, error)
	// OutputWriter returns the configured output writer.
	OutputWriter func(*cli.Context) output.Writer
	// MaybeSyncBeforeWrite runs sync before write operations if auto-sync is enabled.
	MaybeSyncBeforeWrite func(context.Context, service.Service) error
	// MaybeSyncAfterWrite runs sync after write operations if auto-sync is enabled.
	MaybeSyncAfterWrite func(context.Context, service.Service) error
	// FlagString gets a string flag value.
	FlagString func(*cli.Context, string) string
	// FlagBool gets a bool flag value.
	FlagBool func(*cli.Context, string) bool
	// OpenStore opens the store for direct access.
	OpenStore func(context.Context, *cli.Context) (*store.Store, error)
}

var deps Deps

// Init sets dependencies for task commands.
func Init(d Deps) {
	deps = d
}

// Commands returns the task commands without a parent group.
func Commands() []*cli.Command {
	commands := []*cli.Command{
		addCommand(),
		listCommand(),
		showCommand(),
		editCommand(),
		doneCommand(),
		undoCommand(),
		rmCommand(),
		importCommand(),
		exportCommand(),
		contextsCommand(),
		projectsCommand(),
		syncCommand(),
	}
	for _, cmd := range commands {
		cmd.Category = "Tasks"
	}
	return commands
}
