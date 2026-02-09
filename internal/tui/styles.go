package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	app        lipgloss.Style
	header     lipgloss.Style
	panel      lipgloss.Style
	panelFocus lipgloss.Style
	panelEdit  lipgloss.Style
	modalShade lipgloss.Style
	title      lipgloss.Style
	muted      lipgloss.Style
	selected   lipgloss.Style
	errorText  lipgloss.Style
	warning    lipgloss.Style
	success    lipgloss.Style
	foot       lipgloss.Style
	help       lipgloss.Style
	key        lipgloss.Style
	separator  lipgloss.Style
	tabActive  lipgloss.Style
	tabPassive lipgloss.Style
}

func newStyles(theme Theme, noColor bool) styles {
	base := lipgloss.NewStyle()
	if noColor {
		panel := base.Border(lipgloss.NormalBorder()).Padding(0, 1)
		return styles{
			app:        base,
			header:     base.Bold(true),
			panel:      panel,
			panelFocus: panel,
			panelEdit:  panel,
			modalShade: base,
			title:      base.Bold(true),
			muted:      base,
			selected:   base.Bold(true).Underline(true),
			errorText:  base.Bold(true),
			warning:    base.Bold(true),
			success:    base.Bold(true),
			foot:       base,
			help:       panel,
			key:        base.Bold(true),
			separator:  base,
			tabActive:  base.Bold(true),
			tabPassive: base,
		}
	}

	panel := base.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Tokens.border).
		Foreground(theme.Tokens.text).
		Padding(0, 1)

	return styles{
		app:        base.Foreground(theme.Tokens.text),
		header:     base.Bold(true).Foreground(theme.Tokens.accent),
		panel:      panel,
		panelFocus: panel.BorderForeground(theme.Tokens.focusBorder),
		panelEdit:  panel.BorderForeground(theme.Tokens.success),
		modalShade: base.Background(theme.Tokens.surface),
		title:      base.Bold(true).Foreground(theme.Tokens.accentSoft),
		muted:      base.Foreground(theme.Tokens.textMuted),
		selected:   base.Foreground(theme.Tokens.selectionFg).Background(theme.Tokens.selectionBg).Bold(true),
		errorText:  base.Foreground(theme.Tokens.danger).Bold(true),
		warning:    base.Foreground(theme.Tokens.warning).Bold(true),
		success:    base.Foreground(theme.Tokens.success).Bold(true),
		foot:       base.Foreground(theme.Tokens.textMuted),
		help:       panel.BorderForeground(theme.Tokens.focusBorder),
		key:        base.Foreground(theme.Tokens.accent).Bold(true),
		separator:  base.Foreground(theme.Tokens.border),
		tabActive: base.
			Foreground(theme.Tokens.selectionFg).
			Background(theme.Tokens.selectionBg).
			Bold(true).
			Padding(0, 1),
		tabPassive: base.Foreground(theme.Tokens.textMuted).Padding(0, 1),
	}
}
