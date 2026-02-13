package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pterm/pterm"

	"github.com/mholtzscher/ugh/internal/config"
)

type themeDefinition struct {
	Name    string
	Palette []string
}

// builtinDefaultTheme captures PTerm's built-in defaults so we can reliably
// switch back to them after applying a custom theme.
//
//nolint:gochecknoglobals // Captured once at init; used as a stable baseline.
var builtinDefaultTheme = pterm.ThemeDefault

// themeDefinitions contains all themes sourced from ansicolor.com.
//
// Add new themes here by pasting their hex palette.
//
//nolint:gochecknoglobals // Theme registry.
var themeDefinitions = []themeDefinition{
	{
		Name: "andromeda",
		Palette: []string{
			"#00e8c6",
			"#96e072",
			"#f3d56e",
			"#f39c12",
			"#ff007a",
			"#f92672",
			"#c74ded",
			"#00b8d4",
			"#fefefe",
			"#1a1d23",
			"#1e252c",
		},
	},
	{
		Name: "ayu-dark",
		Palette: []string{
			"#f07178",
			"#f29668",
			"#ff8f40",
			"#ffb454",
			"#e6c08a",
			"#d2a6ff",
			"#aad94c",
			"#95e6cb",
			"#59c2ff",
			"#39bae6",
		},
	},
	{
		Name: "ayu-light",
		Palette: []string{
			"#f07171",
			"#f2a191",
			"#fa8532",
			"#eba400",
			"#e59645",
			"#a37acc",
			"#86b300",
			"#4cbf99",
			"#22a4e6",
			"#55b4d4",
		},
	},
	{
		Name: "bearded-dark",
		Palette: []string{
			"#3398db",
			"#37ae6f",
			"#7e9e2d",
			"#d26d32",
			"#cc71bc",
			"#935cd1",
			"#c13838",
			"#de456b",
			"#24b5a8",
			"#c9a022",
			"#181a1f",
			"#a2abb6",
		},
	},
	{
		Name: "bearded-light",
		Palette: []string{
			"#0073d1",
			"#189433",
			"#5e8516",
			"#d06200",
			"#e022b4",
			"#8737e6",
			"#d03333",
			"#e8386a",
			"#009999",
			"#bb9600",
			"#f3f4f5",
			"#22a5c9",
		},
	},
	{
		Name: "catppuccin-frappe",
		Palette: []string{
			"#f2d5cf",
			"#eebebe",
			"#f4b8e4",
			"#ca9ee6",
			"#e78284",
			"#ea999c",
			"#ef9f76",
			"#e5c890",
			"#a6d189",
			"#81c8be",
			"#99d1db",
			"#85c1dc",
			"#8caaee",
			"#babbf1",
		},
	},
	{
		Name: "catppuccin-latte",
		Palette: []string{
			"#dc8a78",
			"#dd7878",
			"#ea76cb",
			"#8839ef",
			"#d20f39",
			"#e64553",
			"#fe640b",
			"#df8e1d",
			"#40a02b",
			"#179299",
			"#04a5e5",
			"#209fb5",
			"#1e66f5",
			"#7287fd",
		},
	},
	{
		Name: "catppuccin-macchiato",
		Palette: []string{
			"#f4dbd6",
			"#f0c6c6",
			"#f5bde6",
			"#c6a0f6",
			"#ed8796",
			"#ee99a0",
			"#f5a97f",
			"#eed49f",
			"#a6da95",
			"#8bd5ca",
			"#91d7e3",
			"#7dc4e4",
			"#8aadf4",
			"#b7bdf8",
		},
	},
	{
		Name: "catppuccin-mocha",
		Palette: []string{
			"#f5e0dc",
			"#f2cdcd",
			"#f5c2e7",
			"#cba6f7",
			"#f38ba8",
			"#eba0ac",
			"#fab387",
			"#f9e2af",
			"#a6e3a1",
			"#94e2d5",
			"#89dceb",
			"#74c7ec",
			"#89b4fa",
			"#b4befe",
		},
	},
	{
		Name: "dracula",
		Palette: []string{
			"#ff5555",
			"#ffb86c",
			"#f1fa8c",
			"#50fa7b",
			"#8be9fd",
			"#bd93f9",
			"#ff79c6",
			"#f8f8f2",
			"#282a36",
		},
	},
	{
		Name: "github-dark",
		Palette: []string{
			"#1f70eb",
			"#238636",
			"#fcb32c",
			"#ff5a02",
			"#b62324",
			"#8957e5",
			"#f0f6fd",
			"#000408",
		},
	},
	{
		Name: "gruvbox",
		Palette: []string{
			"#cc241d",
			"#98971a",
			"#d79921",
			"#458588",
			"#b16286",
			"#689d6a",
			"#d65d0e",
		},
	},
	{
		Name: "monokai-dark",
		Palette: []string{
			"#ff6188",
			"#fc9867",
			"#ffd966",
			"#a8dc76",
			"#78dce8",
			"#ab9ef2",
		},
	},
	{
		Name: "monokai-light",
		Palette: []string{
			"#e14774",
			"#e16031",
			"#cc7a0a",
			"#269d69",
			"#1c8ca8",
			"#7058be",
		},
	},
	{
		Name: "nord",
		Palette: []string{
			"#2e3440",
			"#3b4252",
			"#434c5e",
			"#4c566a",
			"#d8dee9",
			"#e5e9f0",
			"#eceff4",
			"#8fbcbb",
			"#88c0d0",
			"#81a1c1",
			"#5e81ac",
			"#bf616a",
			"#d08770",
			"#ebcb8b",
			"#a3be8c",
			"#b48ead",
		},
	},
	{
		Name: "one-dark",
		Palette: []string{
			"#e06c75",
			"#98c379",
			"#d19a66",
			"#61afef",
			"#c678dd",
			"#56b6c2",
			"#abb2bf",
			"#282c34",
		},
	},
	{
		Name: "one-light",
		Palette: []string{
			"#de3d35",
			"#3e953a",
			"#d2b67b",
			"#2f5af3",
			"#a00095",
			"#3e953a",
			"#2a2b33",
			"#f8f8f8",
		},
	},
	{
		Name: "palenight",
		Palette: []string{
			"#ff5572",
			"#a9c77d",
			"#ffcb6b",
			"#82aaff",
			"#c792ea",
			"#89ddff",
			"#676e95",
		},
	},
	{
		Name: "panda",
		Palette: []string{
			"#ff2c6d",
			"#19f9d8",
			"#ffb86c",
			"#45a9f9",
			"#ff75b5",
			"#b084eb",
			"#cccccc",
			"#292a2b",
		},
	},
	{
		Name: "solarized-dark",
		Palette: []string{
			"#b58900",
			"#cb4b16",
			"#dc322f",
			"#d33682",
			"#6c71c4",
			"#268bd2",
			"#2aa198",
			"#859900",
			"#002b36",
			"#073642",
			"#586e75",
			"#657b83",
			"#839496",
			"#93a1a1",
			"#eee8d5",
			"#fdf6e3",
		},
	},
	{
		Name: "solarized-light",
		Palette: []string{
			"#b58900",
			"#cb4b16",
			"#dc322f",
			"#d33682",
			"#6c71c4",
			"#268bd2",
			"#2aa198",
			"#859900",
			"#002b36",
			"#073642",
			"#586e75",
			"#657b83",
			"#839496",
			"#93a1a1",
			"#eee8d5",
			"#fdf6e3",
		},
	},
	{
		Name: "synthwave-84",
		Palette: []string{
			"#03edf9",
			"#fe7edb",
			"#f3e70f",
			"#ff8b3a",
			"#fe4450",
			"#72f1b8",
			"#272334",
		},
	},
	{
		Name: "tailwind",
		Palette: []string{
			"#f43f5e",
			"#ef4444",
			"#f97316",
			"#f59e0b",
			"#eab308",
			"#84cc16",
			"#22c55e",
			"#10b981",
			"#14b8a6",
			"#06b6d4",
			"#0ea5e9",
			"#3b82f6",
			"#6366f1",
			"#8b5cf6",
			"#a855f7",
			"#d946ef",
			"#ec4899",
			"#64748b",
			"#6b7280",
			"#71717a",
			"#737373",
			"#78716c",
		},
	},
	{
		Name: "tokyo-night",
		Palette: []string{
			"#f7768e",
			"#ff9e64",
			"#e0af68",
			"#9ece6a",
			"#73daca",
			"#b4f9f8",
			"#2ac3de",
			"#7dcfff",
			"#7aa2f7",
			"#bb9af7",
			"#a9b1d6",
			"#1a1b26",
		},
	},
	{
		Name: "tokyo-night-light",
		Palette: []string{
			"#8c4351",
			"#965027",
			"#8f5e15",
			"#385f0d",
			"#33635c",
			"#006c86",
			"#0f4b6e",
			"#2959aa",
			"#5a3e8e",
			"#343b58",
			"#e6e7ed",
		},
	},
}

// GetTheme looks up a theme by name.
//
// Special theme: "default" resets to PTerm's built-in theme.
func GetTheme(name string) (pterm.Theme, bool) {
	normalized := normalizeThemeName(name)
	if normalized == "" {
		normalized = normalizeThemeName(config.DefaultUITheme)
	}

	if normalized == "default" {
		return builtinDefaultTheme, true
	}

	for _, def := range themeDefinitions {
		if normalizeThemeName(def.Name) == normalized {
			return themeFromPalette(def.Palette), true
		}
	}

	return pterm.Theme{}, false
}

func normalizeThemeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}

func themeFromPalette(hexPalette []string) pterm.Theme {
	// We start from the built-in theme to inherit defaults for fields we don't
	// explicitly override.
	theme := builtinDefaultTheme

	const (
		idxPrimary        = 0
		idxSuccess        = 1
		idxHighlight      = 2
		idxSecondary      = 3
		idxError          = 4
		idxInfo           = 5
		idxSection        = 6
		idxSpinner        = 7
		idxTableHeader    = 8
		idxTableSeparator = 9
		idxBox            = 10
		idxBarLabel       = 11
	)

	hexAt := func(index int) string {
		if len(hexPalette) == 0 {
			return ""
		}
		return hexPalette[index%len(hexPalette)]
	}

	// Palette mapping: uses the first 12 colors.
	primaryHex := hexAt(idxPrimary)
	successHex := hexAt(idxSuccess)
	highlightHex := hexAt(idxHighlight)
	secondaryHex := hexAt(idxSecondary)
	errorHex := hexAt(idxError)
	infoHex := hexAt(idxInfo)
	sectionHex := hexAt(idxSection)
	spinnerHex := hexAt(idxSpinner)
	tableHeaderHex := hexAt(idxTableHeader)
	tableSeparatorHex := hexAt(idxTableSeparator)
	boxHex := hexAt(idxBox)
	barLabelHex := hexAt(idxBarLabel)
	barHex := primaryHex // Reuse primary to keep mapping to 12 colors.

	theme.PrimaryStyle = hexToFgStyle(primaryHex, theme.PrimaryStyle)
	theme.SecondaryStyle = hexToFgStyle(secondaryHex, theme.SecondaryStyle)
	theme.HighlightStyle = withBold(hexToFgStyle(highlightHex, theme.HighlightStyle))

	theme.InfoMessageStyle = hexToFgStyle(infoHex, theme.InfoMessageStyle)
	theme.SuccessMessageStyle = hexToFgStyle(successHex, theme.SuccessMessageStyle)
	theme.WarningMessageStyle = hexToFgStyle(highlightHex, theme.WarningMessageStyle)
	theme.ErrorMessageStyle = hexToFgStyle(errorHex, theme.ErrorMessageStyle)
	theme.FatalMessageStyle = theme.ErrorMessageStyle

	// Prefix styles: force explicit background colors.
	black := hexToFgStyle("#000000", pterm.Style{pterm.FgBlack})
	theme.InfoPrefixStyle = joinStyles(black, hexToBgStyle(infoHex, theme.InfoPrefixStyle))
	theme.SuccessPrefixStyle = joinStyles(black, hexToBgStyle(successHex, theme.SuccessPrefixStyle))
	theme.WarningPrefixStyle = joinStyles(black, hexToBgStyle(highlightHex, theme.WarningPrefixStyle))
	theme.ErrorPrefixStyle = joinStyles(black, hexToBgStyle(errorHex, theme.ErrorPrefixStyle))
	theme.FatalPrefixStyle = theme.ErrorPrefixStyle

	theme.SectionStyle = withBold(hexToFgStyle(sectionHex, theme.SectionStyle))
	theme.SpinnerStyle = hexToFgStyle(spinnerHex, theme.SpinnerStyle)

	theme.TableHeaderStyle = hexToFgStyle(tableHeaderHex, theme.TableHeaderStyle)
	theme.TableSeparatorStyle = hexToFgStyle(tableSeparatorHex, theme.TableSeparatorStyle)
	theme.BoxStyle = hexToFgStyle(boxHex, theme.BoxStyle)

	theme.ProgressbarTitleStyle = hexToFgStyle(barLabelHex, theme.ProgressbarTitleStyle)
	theme.ProgressbarBarStyle = hexToFgStyle(barHex, theme.ProgressbarBarStyle)
	theme.BarLabelStyle = hexToFgStyle(barLabelHex, theme.BarLabelStyle)
	theme.BarStyle = hexToFgStyle(barHex, theme.BarStyle)

	return theme
}

func withBold(style pterm.Style) pterm.Style {
	return append(pterm.Style{pterm.Bold}, style...)
}

func joinStyles(a, b pterm.Style) pterm.Style {
	out := make(pterm.Style, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	return out
}

// hexToFgStyle returns an SGR style for an exact foreground color: "38;2;R;G;B".
func hexToFgStyle(hex string, fallback pterm.Style) pterm.Style {
	const (
		sgrForegroundTrueColor = 38
		sgrTrueColorMode       = 2
	)

	r, g, b, ok := parseHexColor(hex)
	if !ok {
		return fallback
	}
	return pterm.Style{
		pterm.Color(sgrForegroundTrueColor),
		pterm.Color(sgrTrueColorMode),
		pterm.Color(r),
		pterm.Color(g),
		pterm.Color(b),
	}
}

// hexToBgStyle returns an SGR style for an exact background color: "48;2;R;G;B".
func hexToBgStyle(hex string, fallback pterm.Style) pterm.Style {
	const (
		sgrBackgroundTrueColor = 48
		sgrTrueColorMode       = 2
	)

	r, g, b, ok := parseHexColor(hex)
	if !ok {
		return fallback
	}
	return pterm.Style{
		pterm.Color(sgrBackgroundTrueColor),
		pterm.Color(sgrTrueColorMode),
		pterm.Color(r),
		pterm.Color(g),
		pterm.Color(b),
	}
}

func parseHexColor(hex string) (uint8, uint8, uint8, bool) {
	const hexLength = 6

	value := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(value) != hexLength {
		return 0, 0, 0, false
	}

	rv, err := strconv.ParseUint(value[0:2], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}
	gv, err := strconv.ParseUint(value[2:4], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}
	bv, err := strconv.ParseUint(value[4:6], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}

	return uint8(rv), uint8(gv), uint8(bv), true
}

func applyTheme(themeName string) error {
	name := strings.TrimSpace(themeName)
	if name == "" {
		name = config.DefaultUITheme
	}

	theme, ok := GetTheme(name)
	if !ok {
		return fmt.Errorf("unknown ui.theme %q (available: %s)", name, strings.Join(AvailableThemeNames(), ", "))
	}

	pterm.ThemeDefault = theme //nolint:reassign // ThemeDefault is the official pterm theme switch mechanism.
	return nil
}

func AvailableThemeNames() []string {
	names := make([]string, 0, len(themeDefinitions)+1)
	names = append(names, "default")
	for _, def := range themeDefinitions {
		names = append(names, def.Name)
	}
	sort.Strings(names)
	return names
}
