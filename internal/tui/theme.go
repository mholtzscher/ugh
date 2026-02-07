package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type colorTokens struct {
	bg          lipgloss.Color
	surface     lipgloss.Color
	border      lipgloss.Color
	text        lipgloss.Color
	textMuted   lipgloss.Color
	accent      lipgloss.Color
	accentSoft  lipgloss.Color
	success     lipgloss.Color
	warning     lipgloss.Color
	danger      lipgloss.Color
	info        lipgloss.Color
	selectionBg lipgloss.Color
	selectionFg lipgloss.Color
	focusBorder lipgloss.Color
}

type Theme struct {
	Name   string
	Tokens colorTokens
}

const (
	themeANSIDefault      = "ansi-default"
	themeANSILight        = "ansi-light"
	themeANSIHighContrast = "ansi-high-contrast"
)

func SelectTheme(name string) Theme {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", themeANSIDefault:
		return ansiDefaultTheme()
	case themeANSILight:
		return ansiLightTheme()
	case themeANSIHighContrast:
		return ansiHighContrastTheme()
	default:
		return ansiDefaultTheme()
	}
}

func ansiDefaultTheme() Theme {
	return Theme{
		Name: themeANSIDefault,
		Tokens: colorTokens{
			bg:          lipgloss.Color("0"),
			surface:     lipgloss.Color("0"),
			border:      lipgloss.Color("12"),
			text:        lipgloss.Color("15"),
			textMuted:   lipgloss.Color("7"),
			accent:      lipgloss.Color("14"),
			accentSoft:  lipgloss.Color("12"),
			success:     lipgloss.Color("10"),
			warning:     lipgloss.Color("11"),
			danger:      lipgloss.Color("9"),
			info:        lipgloss.Color("14"),
			selectionBg: lipgloss.Color("12"),
			selectionFg: lipgloss.Color("0"),
			focusBorder: lipgloss.Color("14"),
		},
	}
}

func ansiLightTheme() Theme {
	return Theme{
		Name: themeANSILight,
		Tokens: colorTokens{
			bg:          lipgloss.Color("15"),
			surface:     lipgloss.Color("7"),
			border:      lipgloss.Color("8"),
			text:        lipgloss.Color("0"),
			textMuted:   lipgloss.Color("8"),
			accent:      lipgloss.Color("4"),
			accentSoft:  lipgloss.Color("6"),
			success:     lipgloss.Color("2"),
			warning:     lipgloss.Color("3"),
			danger:      lipgloss.Color("1"),
			info:        lipgloss.Color("4"),
			selectionBg: lipgloss.Color("12"),
			selectionFg: lipgloss.Color("15"),
			focusBorder: lipgloss.Color("4"),
		},
	}
}

func ansiHighContrastTheme() Theme {
	return Theme{
		Name: themeANSIHighContrast,
		Tokens: colorTokens{
			bg:          lipgloss.Color("0"),
			surface:     lipgloss.Color("0"),
			border:      lipgloss.Color("15"),
			text:        lipgloss.Color("15"),
			textMuted:   lipgloss.Color("7"),
			accent:      lipgloss.Color("14"),
			accentSoft:  lipgloss.Color("12"),
			success:     lipgloss.Color("10"),
			warning:     lipgloss.Color("11"),
			danger:      lipgloss.Color("9"),
			info:        lipgloss.Color("14"),
			selectionBg: lipgloss.Color("15"),
			selectionFg: lipgloss.Color("0"),
			focusBorder: lipgloss.Color("15"),
		},
	}
}
