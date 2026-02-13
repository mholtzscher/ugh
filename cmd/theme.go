package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pterm/pterm"

	"github.com/mholtzscher/ugh/internal/config"
)

type themePalette struct {
	primary        pterm.Color
	secondary      pterm.Color
	highlight      pterm.Color
	info           pterm.Color
	success        pterm.Color
	section        pterm.Color
	spinner        pterm.Color
	tableHeader    pterm.Color
	tableSeparator pterm.Color
	box            pterm.Color
	barLabel       pterm.Color
	bar            pterm.Color
}

type themeBuilder func() pterm.Theme

func themeRegistry() map[string]themeBuilder {
	return map[string]themeBuilder{
		config.DefaultUITheme: ansiDefaultTheme,
		"ocean": func() pterm.Theme {
			return themedVariant(themePalette{
				primary:        pterm.FgCyan,
				secondary:      pterm.FgBlue,
				highlight:      pterm.FgLightCyan,
				info:           pterm.FgCyan,
				success:        pterm.FgLightGreen,
				section:        pterm.FgLightCyan,
				spinner:        pterm.FgCyan,
				tableHeader:    pterm.FgLightCyan,
				tableSeparator: pterm.FgCyan,
				box:            pterm.FgCyan,
				barLabel:       pterm.FgLightCyan,
				bar:            pterm.FgCyan,
			})
		},
		"forest": func() pterm.Theme {
			return themedVariant(themePalette{
				primary:        pterm.FgGreen,
				secondary:      pterm.FgLightGreen,
				highlight:      pterm.FgYellow,
				info:           pterm.FgLightGreen,
				success:        pterm.FgGreen,
				section:        pterm.FgGreen,
				spinner:        pterm.FgGreen,
				tableHeader:    pterm.FgLightGreen,
				tableSeparator: pterm.FgGreen,
				box:            pterm.FgGreen,
				barLabel:       pterm.FgLightGreen,
				bar:            pterm.FgGreen,
			})
		},
		"ansi-high-contrast": func() pterm.Theme {
			return themedVariant(themePalette{
				primary:        pterm.FgLightBlue,
				secondary:      pterm.FgLightMagenta,
				highlight:      pterm.FgLightYellow,
				info:           pterm.FgLightCyan,
				success:        pterm.FgLightGreen,
				section:        pterm.FgLightYellow,
				spinner:        pterm.FgLightCyan,
				tableHeader:    pterm.FgLightWhite,
				tableSeparator: pterm.FgLightWhite,
				box:            pterm.FgLightWhite,
				barLabel:       pterm.FgLightWhite,
				bar:            pterm.FgLightCyan,
			})
		},
	}
}

func themeByName(name string) (pterm.Theme, bool) {
	builder, ok := themeRegistry()[name]
	if !ok {
		return pterm.Theme{}, false
	}
	return builder(), true
}

func themedVariant(p themePalette) pterm.Theme {
	theme := ansiDefaultTheme()
	theme.PrimaryStyle = pterm.Style{p.primary}
	theme.SecondaryStyle = pterm.Style{p.secondary}
	theme.HighlightStyle = pterm.Style{pterm.Bold, p.highlight}
	theme.InfoMessageStyle = pterm.Style{p.info}
	theme.InfoPrefixStyle = pterm.Style{pterm.FgBlack, pterm.BgCyan}
	theme.SuccessMessageStyle = pterm.Style{p.success}
	theme.SuccessPrefixStyle = pterm.Style{pterm.FgBlack, pterm.BgLightGreen}
	theme.WarningMessageStyle = pterm.Style{pterm.FgYellow}
	theme.WarningPrefixStyle = pterm.Style{pterm.FgBlack, pterm.BgYellow}
	theme.ErrorMessageStyle = pterm.Style{pterm.FgLightRed}
	theme.ErrorPrefixStyle = pterm.Style{pterm.FgBlack, pterm.BgLightRed}
	theme.SectionStyle = pterm.Style{pterm.Bold, p.section}
	theme.SpinnerStyle = pterm.Style{p.spinner}
	theme.TableHeaderStyle = pterm.Style{p.tableHeader}
	theme.TableSeparatorStyle = pterm.Style{p.tableSeparator}
	theme.BoxStyle = pterm.Style{p.box}
	theme.BarLabelStyle = pterm.Style{p.barLabel}
	theme.BarStyle = pterm.Style{p.bar}
	return theme
}

func ansiDefaultTheme() pterm.Theme {
	return pterm.Theme{
		DefaultText:             pterm.Style{pterm.FgDefault, pterm.BgDefault},
		PrimaryStyle:            pterm.Style{pterm.FgLightCyan},
		SecondaryStyle:          pterm.Style{pterm.FgLightMagenta},
		HighlightStyle:          pterm.Style{pterm.Bold, pterm.FgYellow},
		InfoMessageStyle:        pterm.Style{pterm.FgLightCyan},
		InfoPrefixStyle:         pterm.Style{pterm.FgBlack, pterm.BgCyan},
		SuccessMessageStyle:     pterm.Style{pterm.FgGreen},
		SuccessPrefixStyle:      pterm.Style{pterm.FgBlack, pterm.BgGreen},
		WarningMessageStyle:     pterm.Style{pterm.FgYellow},
		WarningPrefixStyle:      pterm.Style{pterm.FgBlack, pterm.BgYellow},
		ErrorMessageStyle:       pterm.Style{pterm.FgLightRed},
		ErrorPrefixStyle:        pterm.Style{pterm.FgBlack, pterm.BgLightRed},
		FatalMessageStyle:       pterm.Style{pterm.FgLightRed},
		FatalPrefixStyle:        pterm.Style{pterm.FgBlack, pterm.BgLightRed},
		DescriptionMessageStyle: pterm.Style{pterm.FgDefault},
		DescriptionPrefixStyle:  pterm.Style{pterm.FgLightWhite, pterm.BgDarkGray},
		ScopeStyle:              pterm.Style{pterm.FgGray},
		ProgressbarBarStyle:     pterm.Style{pterm.FgCyan},
		ProgressbarTitleStyle:   pterm.Style{pterm.FgLightCyan},
		HeaderTextStyle:         pterm.Style{pterm.FgLightWhite, pterm.Bold},
		HeaderBackgroundStyle:   pterm.Style{pterm.BgGray},
		SpinnerStyle:            pterm.Style{pterm.FgLightCyan},
		SpinnerTextStyle:        pterm.Style{pterm.FgLightWhite},
		TimerStyle:              pterm.Style{pterm.FgGray},
		TableStyle:              pterm.Style{pterm.FgDefault},
		TableHeaderStyle:        pterm.Style{pterm.FgLightCyan},
		TableSeparatorStyle:     pterm.Style{pterm.FgGray},
		HeatmapStyle:            pterm.Style{pterm.FgDefault},
		HeatmapHeaderStyle:      pterm.Style{pterm.FgLightCyan},
		HeatmapSeparatorStyle:   pterm.Style{pterm.FgDefault},
		SectionStyle:            pterm.Style{pterm.Bold, pterm.FgYellow},
		BulletListTextStyle:     pterm.Style{pterm.FgDefault},
		BulletListBulletStyle:   pterm.Style{pterm.FgGray},
		TreeStyle:               pterm.Style{pterm.FgGray},
		TreeTextStyle:           pterm.Style{pterm.FgDefault},
		LetterStyle:             pterm.Style{pterm.FgDefault},
		DebugMessageStyle:       pterm.Style{pterm.FgGray},
		DebugPrefixStyle:        pterm.Style{pterm.FgBlack, pterm.BgGray},
		BoxStyle:                pterm.Style{pterm.FgDefault},
		BoxTextStyle:            pterm.Style{pterm.FgDefault},
		BarLabelStyle:           pterm.Style{pterm.FgLightCyan},
		BarStyle:                pterm.Style{pterm.FgCyan},
		Checkmark: pterm.Checkmark{
			Checked:   pterm.Green("✓"),
			Unchecked: pterm.Red("✗"),
		},
	}
}

func applyTheme(themeName string) error {
	name := strings.TrimSpace(themeName)
	if name == "" {
		name = config.DefaultUITheme
	}

	theme, ok := themeByName(name)
	if !ok {
		return fmt.Errorf("unknown ui.theme %q (available: %s)", name, strings.Join(availableThemeNames(), ", "))
	}

	pterm.ThemeDefault = theme //nolint:reassign // ThemeDefault is the official pterm theme switch mechanism.
	return nil
}

func availableThemeNames() []string {
	names := make([]string, 0, len(themeRegistry()))
	for name := range themeRegistry() {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
