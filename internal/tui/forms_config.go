package tui

import (
	"errors"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/mholtzscher/ugh/internal/config"
)

type configFormField int

const (
	configFormFieldSyncURL configFormField = iota
	configFormFieldAuthToken
	configFormFieldSyncOnWrite
	configFormFieldTheme
)

const configFormFieldCount = 4

const (
	maskTokenMinLength = 8
	maskTokenPrefixLen = 4
	maskTokenSuffixLen = 4
	maskTokenFallback  = "********"
)

type configFormValues struct {
	syncURL     string
	authToken   string
	syncOnWrite string
	theme       string
}

type configFormState struct {
	active bool
	field  configFormField
	values configFormValues
	input  textinput.Model
}

func inactiveConfigForm(width int) configFormState {
	input := textinput.New()
	input.Prompt = "> "
	input.CharLimit = 1024
	input.Width = width
	return configFormState{active: false, input: input}
}

func startConfigForm(cfg config.Config, width int) configFormState {
	form := inactiveConfigForm(width)
	form.active = true
	form.values = configFormValues{
		syncURL:     strings.TrimSpace(cfg.DB.SyncURL),
		authToken:   strings.TrimSpace(cfg.DB.AuthToken),
		syncOnWrite: strconv.FormatBool(cfg.DB.SyncOnWrite),
		theme:       strings.TrimSpace(cfg.UI.Theme),
	}
	if form.values.theme == "" {
		form.values.theme = config.DefaultUITheme
	}
	return form.withField(configFormFieldSyncURL)
}

func (f configFormState) isActive() bool {
	return f.active
}

func (f configFormState) withWidth(width int) configFormState {
	f.input.Width = width
	return f
}

func (f configFormState) withField(field configFormField) configFormState {
	f.field = field
	f.input.Prompt = configFormFieldLabel(field) + ": "
	f.input.Placeholder = configFormFieldPlaceholder(field)
	f.input.SetValue(f.valueForField(field))
	f.input.CursorEnd()
	return f
}

func (f configFormState) valueForField(field configFormField) string {
	switch field {
	case configFormFieldSyncURL:
		return f.values.syncURL
	case configFormFieldAuthToken:
		return f.values.authToken
	case configFormFieldSyncOnWrite:
		return f.values.syncOnWrite
	case configFormFieldTheme:
		return f.values.theme
	default:
		return ""
	}
}

func (f configFormState) commitInput() configFormState {
	value := strings.TrimSpace(f.input.Value())
	switch f.field {
	case configFormFieldSyncURL:
		f.values.syncURL = value
	case configFormFieldAuthToken:
		f.values.authToken = value
	case configFormFieldSyncOnWrite:
		f.values.syncOnWrite = value
	case configFormFieldTheme:
		f.values.theme = value
	}
	return f
}

func (f configFormState) nextField() (configFormState, bool) {
	if f.field >= configFormFieldCount-1 {
		return f, true
	}
	return f.withField(f.field + 1), false
}

func (f configFormState) previousField() configFormState {
	if f.field <= 0 {
		return f.withField(configFormFieldSyncURL)
	}
	return f.withField(f.field - 1)
}

func (f configFormState) render(styles styles) string {
	lines := []string{styles.title.Render("Edit Config")}
	for field := range configFormFieldCount {
		item := configFormField(field)
		value := f.valueForField(item)
		if value == "" {
			value = "-"
		}
		if item == configFormFieldAuthToken {
			value = maskToken(value)
		}
		line := configFormFieldLabel(item) + ": " + value
		if item == f.field {
			line = styles.selected.Render(line)
		}
		lines = append(lines, line)
	}
	lines = append(lines, "", f.input.View())
	lines = append(lines, styles.muted.Render("enter: next/save  shift+tab: previous  esc: cancel"))
	return strings.Join(lines, "\n")
}

func (f configFormState) applyToConfig(cfg config.Config) (config.Config, error) {
	updated := cfg
	updated.DB.SyncURL = strings.TrimSpace(f.values.syncURL)
	updated.DB.AuthToken = strings.TrimSpace(f.values.authToken)

	syncOnWriteRaw := strings.TrimSpace(f.values.syncOnWrite)
	if syncOnWriteRaw == "" {
		syncOnWriteRaw = "false"
	}
	parsed, err := strconv.ParseBool(syncOnWriteRaw)
	if err != nil {
		return cfg, errors.New("sync_on_write must be true or false")
	}
	updated.DB.SyncOnWrite = parsed

	theme := strings.TrimSpace(f.values.theme)
	if theme == "" {
		theme = config.DefaultUITheme
	}
	updated.UI.Theme = theme

	if updated.Version == 0 {
		updated.Version = config.DefaultVersion
	}

	return updated, nil
}

func configFormFieldLabel(field configFormField) string {
	switch field {
	case configFormFieldSyncURL:
		return "db.sync_url"
	case configFormFieldAuthToken:
		return "db.auth_token"
	case configFormFieldSyncOnWrite:
		return "db.sync_on_write"
	case configFormFieldTheme:
		return "ui.theme"
	default:
		return "field"
	}
}

func configFormFieldPlaceholder(field configFormField) string {
	switch field {
	case configFormFieldSyncURL:
		return "libsql://..."
	case configFormFieldAuthToken:
		return "optional"
	case configFormFieldSyncOnWrite:
		return "true or false"
	case configFormFieldTheme:
		return "ansi-default | ansi-light | ansi-high-contrast"
	default:
		return ""
	}
}

func maskToken(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	if len(value) <= maskTokenMinLength {
		return maskTokenFallback
	}
	return value[:maskTokenPrefixLen] + "..." + value[len(value)-maskTokenSuffixLen:]
}
