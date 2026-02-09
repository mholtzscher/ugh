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

type paneFocusMode int

const (
	paneFocusTasks paneFocusMode = iota
	paneFocusDetail
	paneFocusDetailEdit
)

func (m model) viewTasks() string {
	if m.layout.narrow {
		return m.viewTasksNarrow()
	}
	m = m.withPaneSizes(m.layout.listWidth, m.layout.bodyHeight, m.layout.detailWidth, m.layout.bodyHeight)

	listPanel, detailPanel := m.focusedPaneStyles()
	listStyle := listPanel.Width(m.layout.listWidth).Height(m.layout.bodyHeight)
	detailStyle := detailPanel.Width(m.layout.detailWidth).Height(m.layout.bodyHeight)

	list := listStyle.Render(m.renderTaskList())
	detail := detailStyle.Render(m.renderTaskDetail())

	return lipgloss.JoinHorizontal(lipgloss.Top, list, detail)
}

func (m model) viewTasksNarrow() string {
	if m.layout.bodyHeight < narrowSplitMinimum {
		m = m.withPaneSizes(m.layout.listWidth, m.layout.bodyHeight, m.layout.listWidth, m.layout.bodyHeight)
		listPanel, _ := m.focusedPaneStyles()
		listStyle := listPanel.
			Width(m.layout.listWidth).
			Height(m.layout.bodyHeight)
		return listStyle.Render(m.renderTaskList())
	}

	detailHeight := min(narrowDetailHeight, m.layout.bodyHeight/narrowHeightDivisor)
	listHeight := m.layout.bodyHeight - detailHeight

	listPanel, detailPanel := m.focusedPaneStyles()
	listStyle := listPanel.
		Width(m.layout.listWidth).
		Height(listHeight)
	detailStyle := detailPanel.
		Width(m.layout.listWidth).
		Height(detailHeight)
	m = m.withPaneSizes(m.layout.listWidth, listHeight, m.layout.listWidth, detailHeight)

	list := listStyle.Render(m.renderTaskList())
	detail := detailStyle.Render(m.renderTaskDetail())
	return lipgloss.JoinVertical(lipgloss.Left, list, detail)
}

func (m model) focusedPaneMode() paneFocusMode {
	if !m.taskForm.active() {
		return paneFocusTasks
	}
	if m.taskForm.editing {
		return paneFocusDetailEdit
	}
	return paneFocusDetail
}

func (m model) focusedPaneStyles() (lipgloss.Style, lipgloss.Style) {
	listPanel := m.styles.panel
	detailPanel := m.styles.panel

	switch m.focusedPaneMode() {
	case paneFocusTasks:
		listPanel = m.styles.panelFocus
	case paneFocusDetail:
		detailPanel = m.styles.panelFocus
	case paneFocusDetailEdit:
		detailPanel = m.styles.panelEdit
	}

	return listPanel, detailPanel
}

func (m model) withPaneSizes(listWidth int, listHeight int, detailWidth int, detailHeight int) model {
	listContentWidth := max(1, listWidth-panelPadW)
	listContentHeight := max(1, listHeight-panelPadH)
	showState := m.filters.state == ""

	m.taskTable.SetWidth(listContentWidth)
	m.taskTable.SetHeight(listContentHeight)
	setTaskTableData(
		&m.taskTable,
		taskTableColumns(showState, listContentWidth),
		taskTableRows(m.tasks, showState),
		clampZeroToMax(m.selected, len(m.tasks)-1),
	)

	detailContentWidth := max(1, detailWidth-panelPadW)
	detailContentHeight := max(1, detailHeight-panelPadH)
	m.detail.Width = detailContentWidth
	m.detail.Height = detailContentHeight

	return m
}

func (m model) renderTaskList() string {
	if m.loading && len(m.tasks) == 0 {
		return m.styles.muted.Render("Loading tasks...")
	}
	if len(m.tasks) == 0 {
		return m.styles.muted.Render("No tasks found for current filters.")
	}
	return m.taskTable.View()
}

func (m model) renderTaskDetail() string {
	m.detail.SetContent(m.renderTaskDetailContent())
	return m.detail.View()
}

func (m model) renderTaskDetailContent() string {
	if m.taskForm.active() {
		form := m.taskForm.withWidth(max(1, m.detail.Width))
		return form.render(m.styles)
	}

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
