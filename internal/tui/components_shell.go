package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

const taskFormFooterHint = "form: j/k move  enter or i edit  enter/ctrl+j next  shift+tab/ctrl+k previous  " +
	"ctrl+s save  esc stop/cancel"

func newHelpModel(styleSet styles, width int) help.Model {
	helpModel := help.New()
	helpModel.Width = width
	helpModel.Styles.ShortKey = styleSet.key
	helpModel.Styles.ShortDesc = styleSet.foot
	helpModel.Styles.ShortSeparator = styleSet.separator
	helpModel.Styles.FullKey = styleSet.key
	helpModel.Styles.FullDesc = styleSet.foot
	helpModel.Styles.FullSeparator = styleSet.separator
	helpModel.Styles.Ellipsis = styleSet.muted
	return helpModel
}

func newSpinnerModel(styleSet styles) spinner.Model {
	spinnerModel := spinner.New(spinner.WithSpinner(spinner.Dot))
	spinnerModel.Style = styleSet.key
	return spinnerModel
}

func (m model) renderTabs() string {
	renderTab := func(active bool, label string) string {
		if active {
			return m.styles.tabActive.Render(label)
		}
		return m.styles.tabPassive.Render(label)
	}

	parts := make([]string, 0, len(m.tabs))
	for i, tab := range m.tabs {
		label := fmt.Sprintf("%s (%d)", tab.label, tab.count)
		parts = append(parts, renderTab(i == m.tabSelected, label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func (m model) renderStatusLine() string {
	if m.errText != "" {
		return m.styles.errorText.Render(m.errText)
	}

	var parts []string
	if m.loading {
		parts = append(parts, m.spinner.View())
	}
	if strings.TrimSpace(m.status) != "" {
		parts = append(parts, m.status)
	}
	if filterText := strings.TrimSpace(m.filters.statusText()); filterText != "" {
		parts = append(parts, filterText)
	}
	if len(parts) == 0 {
		return ""
	}

	return m.styles.foot.Render(strings.Join(parts, "  "))
}

func (m model) renderFooter() string {
	if m.taskForm.active() {
		return m.styles.foot.Render(taskFormFooterHint)
	}
	helpModel := m.help
	helpModel.Width = m.viewportW
	return helpModel.ShortHelpView(m.keys.ShortHelp())
}

func (m model) renderHelp() string {
	helpModel := m.help
	helpModel.Width = m.viewportW
	return m.styles.help.Render(helpModel.FullHelpView(m.keys.FullHelp()))
}
