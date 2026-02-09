package shell

import (
	"context"
	"fmt"

	"github.com/chzyer/readline"

	"github.com/mholtzscher/ugh/internal/service"
)

// Prompt wraps readline functionality.
type Prompt struct {
	rl *readline.Instance
}

// NewPrompt creates a new interactive prompt.
func NewPrompt(historyPath string) (*Prompt, error) {
	cfg := &readline.Config{
		Prompt:          "ugh> ",
		HistoryFile:     historyPath,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	rl, err := readline.NewEx(cfg)
	if err != nil {
		return nil, fmt.Errorf("create readline: %w", err)
	}

	return &Prompt{rl: rl}, nil
}

// Readline reads a single line of input.
func (p *Prompt) Readline() (string, error) {
	return p.rl.Readline()
}

// Close closes the prompt and saves history.
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
func (h *History) Record(_ context.Context, command string, success bool, summary string) error {
	// Placeholder - actual implementation would use store
	_ = command
	_ = success
	_ = summary
	return nil
}
