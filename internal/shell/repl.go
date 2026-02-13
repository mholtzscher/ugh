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
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"

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
	Mode      Mode
	InputFile string
	NoColor   bool
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
		display: NewDisplay(opts.NoColor),
		history: NewHistory(svc),
	}
}

// Run starts the REPL loop.
func (r *REPL) Run(ctx context.Context) error {
	r.executor = NewExecutor(r.service, r.state, r.options.NoColor)

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

//nolint:gocognit // REPL loop with pterm styling has higher complexity but is maintainable
func (r *REPL) runInteractive(ctx context.Context) error {
	prompt, err := NewPrompt(r.service, r.options.NoColor)
	if err != nil {
		return fmt.Errorf("initialize prompt: %w", err)
	}
	r.prompt = prompt
	defer r.prompt.Close()

	if r.options.NoColor {
		_, _ = fmt.Fprintln(os.Stdout, "ugh")
		_, _ = fmt.Fprintln(os.Stdout, "Type 'help' for available commands, 'quit' to exit")
	} else {
		bigText, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString("ugh")).Srender()
		pterm.Println(bigText)
		pterm.Info.Println("Type 'help' for available commands, 'quit' to exit")
	}
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
			if r.options.NoColor {
				_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", procErr)
			} else {
				pterm.Error.Println(procErr.Error())
			}
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

	if histErr := r.history.Record(ctx, input, true, result.Summary, result.Intent); histErr != nil {
		_ = histErr
	}

	r.display.ShowResult(result)
	return nil
}

func (r *REPL) showHelp() {
	if r.options.NoColor {
		r.showPlainHelp()
	} else {
		r.showColorHelp()
	}
}

//nolint:funlen // Plain help is intentionally explicit for readability in terminal output.
func (r *REPL) showPlainHelp() {
	_, _ = fmt.Fprintln(os.Stdout, "Available Commands")
	_, _ = fmt.Fprintln(os.Stdout)

	_, _ = fmt.Fprintln(os.Stdout, "Navigation:")
	_, _ = fmt.Fprintln(os.Stdout, "  quit, exit, q  Exit the shell    help, ?  Show this help")
	_, _ = fmt.Fprintln(os.Stdout, "  clear          Clear the screen")
	_, _ = fmt.Fprintln(os.Stdout, "")

	_, _ = fmt.Fprintln(os.Stdout, "Views:")
	_, _ = fmt.Fprintln(os.Stdout, "  view i/inbox      Inbox      view n/now        Now")
	_, _ = fmt.Fprintln(os.Stdout, "  view w/waiting    Waiting    view l/later      Later")
	_, _ = fmt.Fprintln(os.Stdout, "  view c/calendar   Tasks with due dates")
	_, _ = fmt.Fprintln(os.Stdout, "")

	_, _ = fmt.Fprintln(os.Stdout, "Examples:")
	_, _ = fmt.Fprintln(os.Stdout, "  add buy milk tomorrow #groceries @store")
	_, _ = fmt.Fprintln(os.Stdout, "  add task due:tomorrow state:inbox")
	_, _ = fmt.Fprintln(os.Stdout, "  set selected state:done")
	_, _ = fmt.Fprintln(os.Stdout, "  set 123 title:new title +project:work")
	_, _ = fmt.Fprintln(os.Stdout, "  find state:now")
	_, _ = fmt.Fprintln(os.Stdout, "  find state:now and project:work")
	_, _ = fmt.Fprintln(os.Stdout, "  find state:now or not state:done")
	_, _ = fmt.Fprintln(os.Stdout, "  find (state:now or state:waiting) and project:work")
	_, _ = fmt.Fprintln(os.Stdout, "  show 3")
	_, _ = fmt.Fprintln(os.Stdout, "  show #work")
	_, _ = fmt.Fprintln(os.Stdout, "  filter context:urgent")
	_, _ = fmt.Fprintln(os.Stdout)

	_, _ = fmt.Fprintln(os.Stdout, "Syntax:")
	_, _ = fmt.Fprintln(os.Stdout, "  add/create/new <title> [operations...]")
	_, _ = fmt.Fprintln(os.Stdout, "  set/edit/update <target> [operations...]")
	_, _ = fmt.Fprintln(os.Stdout, "  find/show/list/filter <expr>   (and/or/not with parentheses)")
	_, _ = fmt.Fprintln(os.Stdout)

	_, _ = fmt.Fprintln(os.Stdout, "Operations:")
	_, _ = fmt.Fprintln(os.Stdout, "  field:value       Set field (title, notes, due, waiting, state)")
	_, _ = fmt.Fprintln(os.Stdout, "  +field:value      Add to list (projects, contexts, meta)")
	_, _ = fmt.Fprintln(os.Stdout, "  -field:value      Remove from list")
	_, _ = fmt.Fprintln(os.Stdout, "  !field            Clear field")
	_, _ = fmt.Fprintln(os.Stdout, "  #project          Add project tag")
	_, _ = fmt.Fprintln(os.Stdout, "  @context          Add context tag")
	_, _ = fmt.Fprintln(os.Stdout)

	_, _ = fmt.Fprintln(os.Stdout, "Predicates:")
	_, _ = fmt.Fprintln(os.Stdout, "  state:inbox|now|waiting|later|done")
	_, _ = fmt.Fprintln(os.Stdout, "  due:today|tomorrow|YYYY-MM-DD")
	_, _ = fmt.Fprintln(os.Stdout, "  project:name, context:name, text:search")
	_, _ = fmt.Fprintln(os.Stdout, "  id:123 or just 123  Find by task ID")
	_, _ = fmt.Fprintln(os.Stdout, "  done visibility: hidden unless expression mentions state:done")
	_, _ = fmt.Fprintln(os.Stdout)

	_, _ = fmt.Fprintln(os.Stdout, "Targets:")
	_, _ = fmt.Fprintln(os.Stdout, "  selected          Currently selected task")
	_, _ = fmt.Fprintln(os.Stdout, "  #123              Task ID")
	_, _ = fmt.Fprintln(os.Stdout)

	_, _ = fmt.Fprintln(os.Stdout, "Context (sticky filters):")
	_, _ = fmt.Fprintln(os.Stdout, "  context            Show current context state")
	_, _ = fmt.Fprintln(os.Stdout, "  context #project   Set default project context")
	_, _ = fmt.Fprintln(os.Stdout, "  context @context   Set default context filter")
	_, _ = fmt.Fprintln(os.Stdout, "  context clear      Clear all context filters")
}

func (r *REPL) showColorHelp() {
	primary := pterm.ThemeDefault.PrimaryStyle.Sprint
	secondary := pterm.ThemeDefault.SecondaryStyle.Sprint
	info := pterm.ThemeDefault.InfoMessageStyle.Sprint
	success := pterm.ThemeDefault.SuccessMessageStyle.Sprint
	warning := pterm.ThemeDefault.WarningMessageStyle.Sprint
	text := pterm.ThemeDefault.DefaultText.Sprint

	pterm.DefaultSection.Println("Available Commands")

	// Navigation panel
	pterm.DefaultBox.WithTitle(primary("Navigation")).WithRightPadding(1).WithLeftPadding(1).Println(
		primary("quit, exit, q") + "    Exit the shell\n" +
			primary("help, ?") + "          Show this help\n" +
			primary("clear") + "            Clear the screen")

	// Views panel
	pterm.DefaultBox.WithTitle(success("Views")).WithRightPadding(1).WithLeftPadding(1).Println(
		success("view") + "               Show available views\n" +
			success("view i/inbox") + "     Inbox tasks\n" +
			success("view n/now") + "       Now tasks\n" +
			success("view w/waiting") + "   Waiting tasks\n" +
			success("view l/later") + "     Later tasks\n" +
			success("view c/calendar") + "  Tasks with due dates")

	// Examples panel
	pterm.DefaultBox.WithTitle(success("Examples")).WithRightPadding(1).WithLeftPadding(1).Println(
		success("add buy milk tomorrow #groceries @store") + "\n" +
			success("add task due:tomorrow state:inbox") + "\n" +
			success("set selected state:done") + "\n" +
			success("set 123 title:new title +project:work") + "\n" +
			success("find state:now") + "\n" +
			success("find state:now and project:work") + "\n" +
			success("find state:now or not state:done") + "\n" +
			success("find (state:now or state:waiting) and project:work") + "\n" +
			success("show 3") + "\n" +
			success("show #work") + "\n" +
			success("filter context:urgent"))

	// Syntax panel - colors match the explanatory boxes
	syntaxContent := warning("add/create/new") + " " +
		text("<title>") + " " +
		secondary("[operations...]") + "\n" +
		warning("set/edit/update") + " " +
		text("<target>") + " " +
		secondary("[operations...]") + "\n" +
		warning("find/show/list/filter") + " " +
		info("<expr>") + " " +
		warning("(and/or/not, parentheses)")
	pterm.DefaultBox.WithTitle(warning("Syntax")).
		WithRightPadding(1).
		WithLeftPadding(1).
		Println(syntaxContent)

	// Operations panel
	pterm.DefaultBox.WithTitle(secondary("Operations")).WithRightPadding(1).WithLeftPadding(1).Println(
		secondary("field:value") + "       Set field (title, notes, due, waiting, state)\n" +
			secondary("+field:value") + "      Add to list (projects, contexts, meta)\n" +
			secondary("-field:value") + "      Remove from list\n" +
			secondary("!field") + "            Clear field\n" +
			secondary("#project") + "          Add project tag\n" +
			secondary("@context") + "          Add context tag")

	// Predicates panel
	pterm.DefaultBox.WithTitle(info("Predicates")).WithRightPadding(1).WithLeftPadding(1).Println(
		info("state:inbox|now|waiting|later|done") + "\n" +
			info("due:today|tomorrow|YYYY-MM-DD") + "\n" +
			info("project:name, context:name, text:search") + "\n" +
			info("id:123 or just 123") + "  Find by task ID\n" +
			info("done visibility") + "        Hidden unless expression mentions state:done")

	// Targets panel
	pterm.DefaultBox.WithTitle(text("Targets")).WithRightPadding(1).WithLeftPadding(1).Println(
		text("selected") + "          Currently selected task\n" +
			text("#123") + "              Task ID")

	// Context panel
	pterm.DefaultBox.WithTitle(primary("Context (sticky filters)")).WithRightPadding(1).WithLeftPadding(1).Println(
		primary("context") + "            Show current context state\n" +
			primary("context #project") + "   Set default project context\n" +
			primary("context @context") + "   Set default context filter\n" +
			primary("context clear") + "      Clear all context filters")
}
