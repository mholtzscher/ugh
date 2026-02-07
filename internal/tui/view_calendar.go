package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
)

func (m model) viewCalendar() string {
	width := m.layout.listWidth
	if !m.layout.narrow {
		width = m.layout.navWidth + m.layout.listWidth + m.layout.detailWidth
	}

	return m.styles.panel.Width(width).Height(m.layout.bodyHeight).Render(m.renderCalendarList())
}

func (m model) renderCalendarList() string {
	if m.loading {
		return m.styles.muted.Render("Loading calendar...")
	}
	if len(m.calendarTasks) == 0 {
		return m.styles.muted.Render("No due tasks found.")
	}

	tasks := sortCalendarTasks(m.calendarTasks)
	selected := m.selectedCalendarTask()
	selectedID := int64(0)
	if selected != nil {
		selectedID = selected.ID
	}

	lines := []string{
		m.styles.title.Render("CALENDAR"),
		m.styles.muted.Render("Due-date view (enter opens selected task in Tasks)"),
	}

	lastDay := ""
	for _, task := range tasks {
		day := calendarDayLabel(task)
		if day != lastDay {
			if lastDay != "" {
				lines = append(lines, "")
			}
			lines = append(lines, m.styles.muted.Render(day))
			lastDay = day
		}

		line := fmt.Sprintf("  %3d  %-7s  %s", task.ID, task.State, task.Title)
		if task.ID == selectedID {
			line = m.styles.selected.Render("> " + strings.TrimPrefix(line, "  "))
		}
		lines = append(lines, line)
	}

	if selected != nil {
		lines = append(lines, "")
		lines = append(lines, m.styles.muted.Render("Selected"))
		lines = append(lines, fmt.Sprintf("#%d %s", selected.ID, selected.Title))
		lines = append(lines, fmt.Sprintf("due: %s", dueText(selected)))
		lines = append(lines, fmt.Sprintf("state: %s", selected.State))
	}

	return strings.Join(lines, "\n")
}

func sortCalendarTasks(tasks []*store.Task) []*store.Task {
	copied := append([]*store.Task(nil), tasks...)
	sort.SliceStable(copied, func(i, j int) bool {
		left := copied[i]
		right := copied[j]

		leftDate := calendarSortTime(left)
		rightDate := calendarSortTime(right)
		if !leftDate.Equal(rightDate) {
			return leftDate.Before(rightDate)
		}
		return left.ID < right.ID
	})
	return copied
}

func calendarSortTime(task *store.Task) time.Time {
	if task == nil || task.DueOn == nil {
		return time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	}
	return task.DueOn.UTC()
}

func calendarDayLabel(task *store.Task) string {
	if task == nil || task.DueOn == nil {
		return "No due date"
	}
	return task.DueOn.Format("2006-01-02 (Mon)")
}
