package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"

	"github.com/mholtzscher/ugh/internal/service"
)

// Mode defines how the shell operates.
type Mode int

const (
	ModeInteractive Mode = iota
	ModeScriptFile
	ModeScriptStdin
)

// Options configures the shell behavior.
type Options struct {
	Mode        Mode
	InputFile   string
	HistoryPath string
}

// SessionState tracks the current shell session context.
type SessionState struct {
	SelectedTaskID *int64
	LastTaskIDs    []int64
	ContextProject string
	ContextContext string
	StartTime      time.Time
	CommandCount   int
}

// REPL manages the interactive shell session.
type REPL struct {
	service  service.Service
	options  Options
	state    *SessionState
	prompt   *Prompt
	executor *Executor
	display  *Display
	history  *History
}

// NewREPL creates a new REPL instance.
func NewREPL(svc service.Service, opts Options) *REPL {
	return &REPL{
		service: svc,
		options: opts,
		state: &SessionState{
			StartTime:    time.Now(),
			CommandCount: 0,
		},
		display: NewDisplay(),
		history: NewHistory(svc),
	}
}

// Run starts the REPL loop.
func (r *REPL) Run(ctx context.Context) error {
	r.executor = NewExecutor(r.service, r.state)

	switch r.options.Mode {
	case ModeInteractive:
		return r.runInteractive(ctx)
	case ModeScriptFile:
		return r.runScriptFile(ctx, r.options.InputFile)
	case ModeScriptStdin:
		return r.runScriptStdin(ctx)
	default:
		return r.runInteractive(ctx)
	}
}

func (r *REPL) runInteractive(ctx context.Context) error {
	prompt, err := NewPrompt(r.options.HistoryPath)
	if err != nil {
		return fmt.Errorf("initialize prompt: %w", err)
	}
	r.prompt = prompt
	defer r.prompt.Close()

	_, _ = fmt.Fprintln(os.Stdout, "ugh shell - Interactive NLP mode")
	_, _ = fmt.Fprintln(os.Stdout, "Type 'help' for available commands, 'quit' to exit")
	_, _ = fmt.Fprintln(os.Stdout, "")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, readErr := r.prompt.Readline()
		if readErr != nil {
			if errors.Is(readErr, readline.ErrInterrupt) {
				continue
			}
			if errors.Is(readErr, io.EOF) {
				_, _ = fmt.Fprintln(os.Stdout, "")
				return nil
			}
			return fmt.Errorf("read input: %w", readErr)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if procErr := r.processCommand(ctx, line); procErr != nil {
			if errors.Is(procErr, errQuit) {
				return nil
			}
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", procErr)
		}
	}
}

func (r *REPL) runScriptFile(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open script file: %w", err)
	}
	defer file.Close()

	return r.runScript(ctx, file)
}

func (r *REPL) runScriptStdin(ctx context.Context) error {
	return r.runScript(ctx, os.Stdin)
}

func (r *REPL) runScript(ctx context.Context, rdr io.Reader) error {
	scanner := NewScriptScanner(rdr)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if procErr := r.processCommand(ctx, line); procErr != nil {
			if errors.Is(procErr, errQuit) {
				return nil
			}
			return fmt.Errorf("line %d: %s: %w", scanner.LineNumber(), line, procErr)
		}
	}

	return scanner.Err()
}

var errQuit = errors.New("quit requested")

func (r *REPL) processCommand(ctx context.Context, input string) error {
	cmd := strings.ToLower(strings.TrimSpace(input))

	switch cmd {
	case "quit", "exit", "q":
		return errQuit
	case "help", "?":
		r.showHelp()
		return nil
	case "clear":
		r.display.Clear()
		return nil
	}

	r.state.CommandCount++

	result, err := r.executor.Execute(ctx, input)
	if err != nil {
		return err
	}

	if histErr := r.history.Record(ctx, input, true, result.Summary); histErr != nil {
		_ = histErr
	}

	r.display.ShowResult(result)
	return nil
}

func (r *REPL) showHelp() {
	_, _ = fmt.Fprintln(os.Stdout, "Available commands:")
	_, _ = fmt.Fprintln(os.Stdout, "")
	_, _ = fmt.Fprintln(os.Stdout, "  Navigation:")
	_, _ = fmt.Fprintln(os.Stdout, "    quit, exit, q    Exit the shell")
	_, _ = fmt.Fprintln(os.Stdout, "    help, ?          Show this help")
	_, _ = fmt.Fprintln(os.Stdout, "    clear            Clear the screen")
	_, _ = fmt.Fprintln(os.Stdout, "")
	_, _ = fmt.Fprintln(os.Stdout, "  Natural language patterns:")
	_, _ = fmt.Fprintln(os.Stdout, "    buy milk tomorrow #groceries @store")
	_, _ = fmt.Fprintln(os.Stdout, "    find tasks about report")
	_, _ = fmt.Fprintln(os.Stdout, "    mark selected as done")
	_, _ = fmt.Fprintln(os.Stdout, "    show all work tasks")
	_, _ = fmt.Fprintln(os.Stdout, "    delete the grocery task")
	_, _ = fmt.Fprintln(os.Stdout, "")
	_, _ = fmt.Fprintln(os.Stdout, "  Projects: #project-name")
	_, _ = fmt.Fprintln(os.Stdout, "  Contexts: @context-name")
	_, _ = fmt.Fprintln(os.Stdout, "  States: inbox, now, waiting, later, done")
}
