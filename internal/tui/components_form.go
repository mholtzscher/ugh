package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const (
	modalOuterPad     = 2
	modalPadSides     = 2
	modalMinWidth     = 52
	modalMinInput     = 1
	modalInputPad     = 8
	modalWidthNum     = 2
	modalWidthDen     = 3
	modalMinBodyWidth = 1
	modalHintText     = "enter: apply  esc: close"
)

func (m model) renderTaskFormModal() string {
	canvasWidth, canvasHeight := m.modalCanvasSize()

	availableWidth := max(modalMinBodyWidth, canvasWidth-(modalOuterPad*modalPadSides))
	preferredWidth := max(modalMinWidth, (canvasWidth*modalWidthNum)/modalWidthDen)
	modalWidth := min(availableWidth, preferredWidth)
	modalInputWidth := max(modalMinInput, modalWidth-modalInputPad)

	form := m.taskForm.withWidth(modalInputWidth)
	content := m.styles.panelFocus.Width(modalWidth).Render(
		form.render(m.styles) + "\n" + m.styles.muted.Render("esc: close"),
	)

	return lipgloss.Place(
		canvasWidth,
		canvasHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m model) renderInputModal(input textinput.Model, title string) string {
	canvasWidth, canvasHeight := m.modalCanvasSize()

	availableWidth := max(modalMinBodyWidth, canvasWidth-(modalOuterPad*modalPadSides))
	preferredWidth := max(modalMinWidth, (canvasWidth*modalWidthNum)/modalWidthDen)
	modalWidth := min(availableWidth, preferredWidth)
	inputWidth := max(modalMinInput, modalWidth-modalInputPad)

	input.Width = inputWidth

	var content string
	if title != "" {
		titleLine := m.styles.title.Render(title)
		content = m.styles.panelFocus.Width(modalWidth).Render(
			titleLine + "\n" + input.View() + "\n" + m.styles.muted.Render(modalHintText),
		)
	} else {
		content = m.styles.panelFocus.Width(modalWidth).Render(
			input.View() + "\n" + m.styles.muted.Render(modalHintText),
		)
	}

	return lipgloss.Place(
		canvasWidth,
		canvasHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m model) modalCanvasSize() (int, int) {
	canvasWidth := m.viewportW
	canvasHeight := m.viewportH

	if canvasWidth <= 0 {
		canvasWidth = m.layout.listWidth
		if !m.layout.narrow {
			canvasWidth = m.layout.listWidth + m.layout.detailWidth
		}
	}
	if canvasHeight <= 0 {
		canvasHeight = m.layout.bodyHeight
	}

	return max(modalMinBodyWidth, canvasWidth), max(1, canvasHeight)
}
