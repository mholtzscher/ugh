package tui

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/service"
)

type Options struct {
	ThemeName  string
	NoColor    bool
	ConfigPath string
	Config     config.Config
}

func Run(ctx context.Context, svc service.Service, opts Options) error {
	if svc == nil {
		return errors.New("tui service is nil")
	}

	_ = ctx
	program := tea.NewProgram(newModel(svc, opts), tea.WithAltScreen())
	_, err := program.Run()
	if err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	return nil
}
