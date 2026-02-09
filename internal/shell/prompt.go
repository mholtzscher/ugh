package shell

import (
	"context"
	"fmt"

	"github.com/chzyer/readline"
	"github.com/pterm/pterm"

	"github.com/mholtzscher/ugh/internal/service"
)

// Prompt wraps readline functionality.
type Prompt struct {
	rl *readline.Instance
}

// NewPrompt creates a new interactive prompt with history loaded from SQLite.
func NewPrompt(svc service.Service, noColor bool) (*Prompt, error) {
	promptText := "ugh> "
	if !noColor {
		promptText = pterm.Cyan("➜ ") + pterm.Magenta("ugh> ")
	}

	cfg := &readline.Config{
		Prompt:          promptText,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	rl, err := readline.NewEx(cfg)
	if err != nil {
		return nil, fmt.Errorf("create readline: %w", err)
	}

	// Load recent history from SQLite into readline
	ctx := context.Background()
	const historyLoadLimit = 100
	history, err := svc.ListShellHistory(ctx, historyLoadLimit)
	if err == nil {
		// Add oldest first so newest ends up at the end (most recent)
		for i := len(history) - 1; i >= 0; i-- {
			_ = rl.SaveHistory(history[i].Command)
		}
	}

	return &Prompt{rl: rl}, nil
}

// Readline reads a single line of input.
func (p *Prompt) Readline() (string, error) {
	return p.rl.Readline()
}

// Close closes the prompt.
func (p *Prompt) Close() error {
	return p.rl.Close()
}

// History manages command history storage.
type History struct {
	svc service.Service
}

// NewHistory creates a new history manager.
func NewHistory(svc service.Service) *History {
	return &History{svc: svc}
}

// Record records a command in history.
func (h *History) Record(ctx context.Context, command string, success bool, summary string, intent string) error {
	_, err := h.svc.RecordShellHistory(ctx, command, success, summary, intent)
	return err
}
