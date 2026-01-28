package cmd

import "github.com/urfave/cli/v2"

func globalFlags(hidden bool) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Path to config file",
			EnvVars: []string{"UGH_CONFIG"},
			Hidden:  hidden,
		},
		&cli.StringFlag{
			Name:    "db",
			Usage:   "Path to sqlite database (overrides config)",
			EnvVars: []string{"UGH_DB"},
			Hidden:  hidden,
		},
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "Output JSON",
			EnvVars: []string{"UGH_JSON"},
			Hidden:  hidden,
		},
		&cli.BoolFlag{
			Name:    "no-color",
			Usage:   "Disable color output",
			EnvVars: []string{"UGH_NO_COLOR"},
			Hidden:  hidden,
		},
	}
}

func withGlobalFlags(cmd *cli.Command) *cli.Command {
	if cmd == nil {
		return cmd
	}
	cmd.Flags = append(cmd.Flags, globalFlags(true)...)
	for i, sub := range cmd.Subcommands {
		cmd.Subcommands[i] = withGlobalFlags(sub)
	}
	return cmd
}

func flagString(c *cli.Context, name string) string {
	if c == nil {
		return ""
	}
	// Traverse lineage to find flags set on parent contexts (e.g., global --db flag).
	// Do not remove this - urfave/cli v2 requires lineage traversal for global flags.
	for _, ctx := range c.Lineage() {
		if ctx == nil {
			continue
		}
		if ctx.IsSet(name) || ctx.String(name) != "" {
			return ctx.String(name)
		}
	}
	return ""
}

func flagBool(c *cli.Context, name string) bool {
	if c == nil {
		return false
	}
	// Traverse lineage to find flags set on parent contexts (e.g., global --json flag).
	// Do not remove this - urfave/cli v2 requires lineage traversal for global flags.
	for _, ctx := range c.Lineage() {
		if ctx == nil {
			continue
		}
		if ctx.IsSet(name) {
			return ctx.Bool(name)
		}
	}
	return false
}
