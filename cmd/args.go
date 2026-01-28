package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"
)

func normalizeArgs(app *cli.App, args []string) []string {
	if len(args) <= 1 || app == nil {
		return args
	}

	globalFlags, globalValueFlags := buildGlobalFlagSpec(app.Flags)

	var globals []string
	var rest []string
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			rest = append(rest, args[i:]...)
			break
		}
		if isGlobalFlag(arg, globalFlags) {
			globals = append(globals, arg)
			if needsValue(arg, globalValueFlags) && i+1 < len(args) {
				globals = append(globals, args[i+1])
				i++
			}
			continue
		}
		rest = append(rest, arg)
	}

	cmdPath, remaining, chain := splitCommandPath(rest, app.Commands)
	valueFlags := collectValueFlags(chain)
	flags, positionals := reorderFlags(valueFlags, remaining)

	result := make([]string, 0, 1+len(globals)+len(cmdPath)+len(flags)+len(positionals))
	result = append(result, args[0])
	result = append(result, globals...)
	result = append(result, cmdPath...)
	result = append(result, flags...)
	result = append(result, positionals...)
	return result
}

func buildGlobalFlagSpec(flags []cli.Flag) (map[string]bool, map[string]bool) {
	flagSet := map[string]bool{}
	valueSet := map[string]bool{}
	for _, f := range flags {
		addFlagSpec(flagSet, valueSet, f)
	}
	return flagSet, valueSet
}

func addFlagSpec(flagSet map[string]bool, valueSet map[string]bool, f cli.Flag) {
	names := f.Names()
	if len(names) == 0 {
		return
	}
	takesValue := flagTakesValue(f)
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		prefix := "--"
		if len(name) == 1 {
			prefix = "-"
		}
		key := prefix + name
		flagSet[key] = true
		if takesValue {
			valueSet[key] = true
		}
	}
}

func addValueFlags(valueSet map[string]bool, f cli.Flag) {
	names := f.Names()
	if len(names) == 0 {
		return
	}
	if !flagTakesValue(f) {
		return
	}
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		prefix := "--"
		if len(name) == 1 {
			prefix = "-"
		}
		valueSet[prefix+name] = true
	}
}

func flagTakesValue(f cli.Flag) bool {
	if docFlag, ok := f.(cli.DocGenerationFlag); ok {
		return docFlag.TakesValue()
	}
	if _, ok := f.(*cli.BoolFlag); ok {
		return false
	}
	return true
}

func isGlobalFlag(arg string, globalFlags map[string]bool) bool {
	if strings.HasPrefix(arg, "--") {
		name := strings.SplitN(arg, "=", 2)[0]
		return globalFlags[name]
	}
	if strings.HasPrefix(arg, "-") && len(arg) == 2 {
		return globalFlags[arg]
	}
	return false
}

func needsValue(arg string, valueFlags map[string]bool) bool {
	if strings.HasPrefix(arg, "--") {
		if strings.Contains(arg, "=") {
			return false
		}
		name := strings.SplitN(arg, "=", 2)[0]
		return valueFlags[name]
	}
	if strings.HasPrefix(arg, "-") && len(arg) == 2 {
		return valueFlags[arg]
	}
	return false
}

func splitCommandPath(args []string, commands []*cli.Command) ([]string, []string, []*cli.Command) {
	if len(args) == 0 {
		return nil, nil, nil
	}
	var path []string
	var chain []*cli.Command
	remaining := args
	for len(remaining) > 0 {
		cmd := findCommand(remaining[0], commands)
		if cmd == nil {
			break
		}
		path = append(path, remaining[0])
		chain = append(chain, cmd)
		remaining = remaining[1:]
		commands = cmd.Subcommands
	}
	return path, remaining, chain
}

func findCommand(name string, commands []*cli.Command) *cli.Command {
	for _, cmd := range commands {
		for _, candidate := range cmd.Names() {
			if candidate == name {
				return cmd
			}
		}
	}
	return nil
}

func collectValueFlags(chain []*cli.Command) map[string]bool {
	valueFlags := map[string]bool{}
	for _, cmd := range chain {
		for _, f := range cmd.Flags {
			addValueFlags(valueFlags, f)
		}
	}
	return valueFlags
}

func reorderFlags(valueFlags map[string]bool, args []string) ([]string, []string) {
	var flags []string
	var positionals []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positionals = append(positionals, args[i:]...)
			break
		}
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			if needsValue(arg, valueFlags) && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i++
			}
			continue
		}
		positionals = append(positionals, arg)
	}
	return flags, positionals
}
