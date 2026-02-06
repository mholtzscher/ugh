package cmd

import (
	"context"
	"fmt"
	"os"

	backupcmd "github.com/mholtzscher/ugh/cmd/backup"
	configcmd "github.com/mholtzscher/ugh/cmd/config"
	daemoncmd "github.com/mholtzscher/ugh/cmd/daemon"
	listscmd "github.com/mholtzscher/ugh/cmd/lists"
	"github.com/mholtzscher/ugh/cmd/registry"
	synccmd "github.com/mholtzscher/ugh/cmd/sync"
	tagscmd "github.com/mholtzscher/ugh/cmd/tags"
	taskscmd "github.com/mholtzscher/ugh/cmd/tasks"
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	appruntime "github.com/mholtzscher/ugh/internal/runtime"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v3"
)

// Version is set at build time.
var Version = "0.1.1" // x-release-please-version

type App struct {
	runtime *appruntime.Runtime
}

func NewApp() *App {
	return &App{runtime: appruntime.New()}
}

func Execute() {
	app := NewApp()
	rootCmd, err := app.RootCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func (a *App) RootCommand() (*cli.Command, error) {
	rootCmd := &cli.Command{
		Name:        "ugh",
		Usage:       "ugh is a task CLI",
		Description: "ugh is a task CLI with SQLite storage.",
		Version:     Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flags.FlagConfigPath,
				Usage: "path to config file",
			},
			&cli.StringFlag{
				Name:    flags.FlagDBPath,
				Aliases: []string{"d"},
				Usage:   "path to sqlite database (overrides config)",
			},
			&cli.BoolFlag{
				Name:    flags.FlagJSON,
				Aliases: []string{"j"},
				Usage:   "output json",
			},
			&cli.BoolFlag{
				Name:  flags.FlagNoColor,
				Usage: "disable color output",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			return ctx, a.loadConfig(cmd)
		},
	}

	cmdRegistry := registry.New()
	if err := a.registerCommands(cmdRegistry); err != nil {
		return nil, err
	}
	commands, err := cmdRegistry.Build()
	if err != nil {
		return nil, err
	}
	rootCmd.Commands = commands

	applyRootHelpCustomization(rootCmd)
	return rootCmd, nil
}

func (a *App) registerCommands(cmdRegistry *registry.Registry) error {
	if err := taskscmd.Register(cmdRegistry, taskscmd.Deps{
		NewService:           a.newService,
		MaybeSyncBeforeWrite: a.maybeSyncBeforeWrite,
		MaybeSyncAfterWrite:  a.maybeSyncAfterWrite,
		OutputWriter:         a.outputWriter,
	}); err != nil {
		return fmt.Errorf("register task commands: %w", err)
	}

	if err := listscmd.Register(cmdRegistry, listscmd.Deps{
		NewService:   a.newService,
		OutputWriter: a.outputWriter,
	}); err != nil {
		return fmt.Errorf("register list commands: %w", err)
	}

	if err := backupcmd.Register(cmdRegistry, backupcmd.Deps{
		NewService:           a.newService,
		MaybeSyncBeforeWrite: a.maybeSyncBeforeWrite,
		MaybeSyncAfterWrite:  a.maybeSyncAfterWrite,
		OutputWriter:         a.outputWriter,
	}); err != nil {
		return fmt.Errorf("register backup commands: %w", err)
	}

	if err := tagscmd.Register(cmdRegistry, tagscmd.Deps{
		NewService:   a.newService,
		OutputWriter: a.outputWriter,
	}); err != nil {
		return fmt.Errorf("register projects/contexts commands: %w", err)
	}

	if err := synccmd.Register(cmdRegistry, synccmd.Deps{
		OpenStore:    a.openStore,
		OutputWriter: a.outputWriter,
	}); err != nil {
		return fmt.Errorf("register sync commands: %w", err)
	}

	if err := configcmd.Register(cmdRegistry, configcmd.Deps{
		Config:             a.getConfig,
		SetConfig:          a.setConfig,
		SetConfigWasLoaded: a.setConfigWasLoaded,
		OutputWriter:       a.outputWriter,
		ConfigPath:         a.getConfigPath,
	}); err != nil {
		return fmt.Errorf("register config commands: %w", err)
	}

	if err := daemoncmd.Register(cmdRegistry, daemoncmd.Deps{
		Config:       a.getConfig,
		OutputWriter: a.outputWriter,
	}); err != nil {
		return fmt.Errorf("register daemon commands: %w", err)
	}

	return nil
}

func (a *App) getConfig() *config.Config {
	return a.runtime.Config()
}

func (a *App) setConfig(cfg *config.Config) {
	a.runtime.SetConfig(cfg)
}

func (a *App) setConfigWasLoaded(wasLoaded bool) {
	a.runtime.SetConfigWasLoaded(wasLoaded)
}

func (a *App) getConfigPath() string {
	return a.runtime.ConfigPath()
}

func (a *App) loadConfig(cmd *cli.Command) error {
	if cmd != nil {
		a.runtime.SetGlobalOptions(
			cmd.String(flags.FlagConfigPath),
			cmd.String(flags.FlagDBPath),
			cmd.Bool(flags.FlagJSON),
			cmd.Bool(flags.FlagNoColor),
		)
	}
	return a.runtime.LoadConfig()
}

func (a *App) openStore(ctx context.Context) (*store.Store, error) {
	return a.runtime.OpenStore(ctx)
}

func (a *App) outputWriter() output.Writer {
	return a.runtime.OutputWriter()
}
