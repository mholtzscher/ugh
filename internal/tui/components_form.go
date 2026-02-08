package tui

import "github.com/charmbracelet/lipgloss"

const (
	modalOuterPad     = 2
	modalMinWidth     = 52
	modalMinInput     = 1
	modalInputPad     = 8
	modalWidthNum     = 2
	modalWidthDen     = 3
	modalMinBodyWidth = 1
)

func (m model) renderTaskFormModal() string {
	bodyWidth := m.layout.listWidth
	if !m.layout.narrow {
		bodyWidth = m.layout.listWidth + m.layout.detailWidth
	}
	bodyWidth = max(modalMinBodyWidth, bodyWidth)

	availableWidth := max(modalMinBodyWidth, bodyWidth-modalOuterPad)
	preferredWidth := max(modalMinWidth, (bodyWidth*modalWidthNum)/modalWidthDen)
	modalWidth := min(availableWidth, preferredWidth)
	modalInputWidth := max(modalMinInput, modalWidth-modalInputPad)

	form := m.taskForm.withWidth(modalInputWidth)
	content := m.styles.panelFocus.Width(modalWidth).Render(form.render(m.styles))

	return lipgloss.Place(
		bodyWidth,
		m.layout.bodyHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
