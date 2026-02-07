package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/mholtzscher/ugh/internal/store"
)

const (
	narrowDetailHeight  = 8
	narrowSplitMinimum  = 3
	narrowHeightDivisor = 2
)

func (m model) viewTasks() string {
	if m.layout.narrow {
		return m.viewTasksNarrow()
	}

	navStyle := m.panelStyle(focusNav).Width(m.layout.navWidth).Height(m.layout.bodyHeight)
	listStyle := m.panelStyle(focusList).Width(m.layout.listWidth).Height(m.layout.bodyHeight)
	detailStyle := m.panelStyle(focusDetail).Width(m.layout.detailWidth).Height(m.layout.bodyHeight)

	nav := navStyle.Render(m.renderNav())
	list := listStyle.Render(m.renderTaskList())
	detail := detailStyle.Render(m.renderTaskDetail())

	return lipgloss.JoinHorizontal(lipgloss.Top, nav, list, detail)
}

func (m model) viewTasksNarrow() string {
	if m.layout.bodyHeight < narrowSplitMinimum {
		listStyle := m.panelStyle(focusList).
			Width(m.layout.listWidth).
			Height(m.layout.bodyHeight)
		return listStyle.Render(m.renderTaskList())
	}

	detailHeight := min(narrowDetailHeight, m.layout.bodyHeight/narrowHeightDivisor)
	listHeight := m.layout.bodyHeight - detailHeight

	listStyle := m.panelStyle(focusList).
		Width(m.layout.listWidth).
		Height(listHeight)
	detailStyle := m.panelStyle(focusDetail).
		Width(m.layout.listWidth).
		Height(detailHeight)

	list := listStyle.Render(m.renderTaskList())
	detail := detailStyle.Render(m.renderTaskDetail())
	return lipgloss.JoinVertical(lipgloss.Left, list, detail)
}

func (m model) panelStyle(focus paneFocus) lipgloss.Style {
	if m.focus == focus {
		return m.styles.panelFocus
	}
	return m.styles.panel
}

func (m model) renderNav() string {
	if len(m.navItems) == 0 {
		return m.styles.muted.Render("No navigation items yet.")
	}

	lines := []string{m.styles.title.Render("NAV")}
	kind := navState
	lines = append(lines, m.styles.muted.Render("States"))

	for idx, item := range m.navItems {
		if item.kind != kind {
			kind = item.kind
			lines = append(lines, "", m.styles.muted.Render(navSectionLabel(kind)))
		}
		line := fmt.Sprintf("  %-16s (%d)", item.label, item.count)
		if idx == m.navSelected {
			line = "> " + strings.TrimPrefix(line, "  ")
			line = m.styles.selected.Render(line)
		}
		lines = append(lines, line)
	}

	if m.focus == focusNav {
		lines = append(lines, "", m.styles.muted.Render("enter: apply selected scope"))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderTaskList() string {
	if m.loading {
		return m.styles.muted.Render("Loading tasks...")
	}
	if len(m.tasks) == 0 {
		return m.styles.muted.Render("No tasks found for current filters.")
	}

	lines := []string{m.styles.title.Render("TASKS")}
	for i, task := range m.tasks {
		line := fmt.Sprintf("%3d  %-7s  %-10s  %s", task.ID, task.State, dueText(task), task.Title)
		if i == m.selected {
			line = m.styles.selected.Render(line)
		}
		lines = append(lines, line)
	}

	if m.focus == focusList {
		lines = append(lines, "", m.styles.muted.Render("enter: focus details"))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderTaskDetail() string {
	task := m.selectedTask()
	if task == nil {
		return m.styles.muted.Render("Select a task to view details.")
	}

	lines := []string{
		m.styles.title.Render("DETAILS"),
		fmt.Sprintf("#%d %s", task.ID, task.Title),
		fmt.Sprintf("state: %s", task.State),
		fmt.Sprintf("due: %s", dueText(task)),
		fmt.Sprintf("waiting: %s", emptyDash(task.WaitingFor)),
		fmt.Sprintf("projects: %s", joinedOrDash(task.Projects)),
		fmt.Sprintf("contexts: %s", joinedOrDash(task.Contexts)),
	}

	if task.Notes != "" {
		lines = append(lines, "", "notes:", task.Notes)
	}
	if len(task.Meta) > 0 {
		lines = append(lines, fmt.Sprintf("meta keys: %d", len(task.Meta)))
	}

	if m.deleteTaskID == task.ID {
		lines = append(lines, "", m.styles.warning.Render("Press D again to confirm delete."))
	}
	if m.focus == focusDetail {
		lines = append(lines, "", m.styles.muted.Render("enter: back to task list"))
	}

	return strings.Join(lines, "\n")
}

func navSectionLabel(kind navItemKind) string {
	switch kind {
	case navState:
		return "States"
	case navProject:
		return "Projects"
	case navContext:
		return "Contexts"
	default:
		return ""
	}
}

func dueText(task *store.Task) string {
	if task == nil || task.DueOn == nil {
		return "-"
	}
	return task.DueOn.Format("2006-01-02")
}

func emptyDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func joinedOrDash(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ",")
}
