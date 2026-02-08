package tui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/mholtzscher/ugh/internal/store"
)

const (
	panelPadW = 4
	panelPadH = 2

	tableIDWidth     = 4
	tableStateWidth  = 8
	tableDueWidth    = 10
	tableMinTitle    = 8
	tableGutterState = 6
	tableGutterBasic = 4
)

func newTaskTable(styleSet styles, layout layoutSpec) table.Model {
	tableStyles := table.DefaultStyles()
	tableStyles.Header = styleSet.muted.Bold(true)
	tableStyles.Cell = lipgloss.NewStyle()
	tableStyles.Selected = styleSet.selected

	taskTable := table.New(
		table.WithFocused(false),
		table.WithStyles(tableStyles),
		table.WithHeight(max(layoutMinBodyHeight, layout.bodyHeight-panelPadH)),
		table.WithWidth(max(layoutMinNarrowListWidth, layout.listWidth-panelPadW)),
	)
	taskTable.KeyMap = table.KeyMap{}
	taskTable.Help = help.New()

	return taskTable
}

func taskTableRows(tasks []*store.Task, showState bool) []table.Row {
	rows := make([]table.Row, 0, len(tasks))
	for _, task := range tasks {
		if task == nil {
			continue
		}
		if showState {
			rows = append(rows, table.Row{
				strconv.FormatInt(task.ID, 10),
				string(task.State),
				dueText(task),
				task.Title,
			})
			continue
		}

		rows = append(rows, table.Row{
			strconv.FormatInt(task.ID, 10),
			dueText(task),
			task.Title,
		})
	}
	return rows
}

func taskTableColumns(showState bool, totalWidth int) []table.Column {
	if showState {
		titleWidth := max(tableMinTitle, totalWidth-(tableIDWidth+tableStateWidth+tableDueWidth+tableGutterState))
		return []table.Column{
			{Title: "ID", Width: tableIDWidth},
			{Title: "State", Width: tableStateWidth},
			{Title: "Due", Width: tableDueWidth},
			{Title: "Task", Width: titleWidth},
		}
	}

	titleWidth := max(tableMinTitle, totalWidth-(tableIDWidth+tableDueWidth+tableGutterBasic))
	return []table.Column{
		{Title: "ID", Width: tableIDWidth},
		{Title: "Due", Width: tableDueWidth},
		{Title: "Task", Width: titleWidth},
	}
}

func setTaskTableData(taskTable *table.Model, columns []table.Column, rows []table.Row, cursor int) {
	if len(taskTable.Columns()) != len(columns) {
		taskTable.SetRows(nil)
	}

	taskTable.SetColumns(columns)
	taskTable.SetRows(rows)
	taskTable.SetCursor(cursor)
}
