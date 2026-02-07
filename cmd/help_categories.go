package cmd

import (
	"io"
	"sort"
	"strings"

	"github.com/urfave/cli/v3"
)

// customRootHelpTemplate is based on urfave/cli's RootCommandHelpTemplate, but
// uses orderedCategories to control category ordering.
const customRootHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.FullName}} {{if .VisibleFlags}}[global options]{{end}}{{if .VisibleCommands}} [command [command options]]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Arguments}} [arguments...]{{end}}{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}
{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range orderedCategories .VisibleCategories}}{{if .Name}}

   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{template "visibleCommandTemplate" .}}{{end}}{{end}}{{end}}{{if .VisibleFlagCategories}}

GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{template "copyrightTemplate" .}}{{end}}
`

func helpCategoryOrder() []string {
	return []string{
		"Lists",
		"Tasks",
		"Projects & Contexts",
		"Sync",
		"System",
	}
}

func orderedCategories(categories []cli.CommandCategory) []cli.CommandCategory {
	if len(categories) == 0 {
		return categories
	}

	byName := make(map[string]cli.CommandCategory, len(categories))
	remaining := make([]cli.CommandCategory, 0, len(categories))
	var unnamed cli.CommandCategory

	for _, cat := range categories {
		if cat == nil {
			continue
		}
		name := cat.Name()
		if strings.TrimSpace(name) == "" {
			unnamed = cat
			continue
		}
		byName[name] = cat
		remaining = append(remaining, cat)
	}

	result := make([]cli.CommandCategory, 0, len(categories))
	if unnamed != nil {
		result = append(result, unnamed)
	}

	used := map[string]bool{}
	for _, name := range helpCategoryOrder() {
		if cat, ok := byName[name]; ok {
			result = append(result, cat)
			used[name] = true
		}
	}

	other := make([]cli.CommandCategory, 0)
	for _, cat := range remaining {
		if cat == nil {
			continue
		}
		name := cat.Name()
		if used[name] {
			continue
		}
		other = append(other, cat)
	}
	sort.Slice(other, func(i, j int) bool {
		return strings.ToLower(other[i].Name()) < strings.ToLower(other[j].Name())
	})

	result = append(result, other...)
	return result
}

//nolint:gochecknoinits,reassign // urfave/cli requires global template hook customization during package initialization.
func init() {
	// Inject custom template funcs so our custom root template can reorder categories.
	old := cli.HelpPrinterCustom
	cli.HelpPrinterCustom = func(w io.Writer, templ string, data any, customFuncs map[string]any) {
		if customFuncs == nil {
			customFuncs = map[string]any{}
		}
		customFuncs["orderedCategories"] = orderedCategories
		old(w, templ, data, customFuncs)
	}

	// Use the custom root template.
	rootCmd.CustomRootCommandHelpTemplate = customRootHelpTemplate
}
