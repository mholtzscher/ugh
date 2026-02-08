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

	listStyle := m.styles.panel.Width(m.layout.listWidth).Height(m.layout.bodyHeight)
	detailStyle := m.styles.panel.Width(m.layout.detailWidth).Height(m.layout.bodyHeight)

	list := listStyle.Render(m.renderTaskList())
	detail := detailStyle.Render(m.renderTaskDetail())

	return lipgloss.JoinHorizontal(lipgloss.Top, list, detail)
}

func (m model) viewTasksNarrow() string {
	if m.layout.bodyHeight < narrowSplitMinimum {
		listStyle := m.styles.panel.
			Width(m.layout.listWidth).
			Height(m.layout.bodyHeight)
		return listStyle.Render(m.renderTaskList())
	}

	detailHeight := min(narrowDetailHeight, m.layout.bodyHeight/narrowHeightDivisor)
	listHeight := m.layout.bodyHeight - detailHeight

	listStyle := m.styles.panel.
		Width(m.layout.listWidth).
		Height(listHeight)
	detailStyle := m.styles.panel.
		Width(m.layout.listWidth).
		Height(detailHeight)

	list := listStyle.Render(m.renderTaskList())
	detail := detailStyle.Render(m.renderTaskDetail())
	return lipgloss.JoinVertical(lipgloss.Left, list, detail)
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

	return strings.Join(lines, "\n")
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
