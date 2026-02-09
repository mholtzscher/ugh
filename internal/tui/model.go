package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	nlp "github.com/mholtzscher/ugh/internal/nlp"
	nlpcompile "github.com/mholtzscher/ugh/internal/nlp/compile"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

type viewKind int

const (
	viewTasks viewKind = iota
)

const (
	defaultWidth            = 120
	defaultHeight           = 32
	statusLoadingText       = "loading tasks..."
	statusRefreshText       = "refreshing..."
	statusCommandCancelled  = "command cancelled"
	statusCommandRejected   = "command rejected"
	searchInputWidth        = 48
	searchMinWidth          = 1
	searchNarrowPad         = 6
	keyCtrlC                = "ctrl+c"
	statusTitleRequired     = "title is required"
	statusTaskFormLastField = "last field selected - press ctrl+s to save"
)

type model struct {
	svc          service.Service
	keys         keyMap
	styles       styles
	viewportW    int
	viewportH    int
	layout       layoutSpec
	view         viewKind
	filters      listFilters
	tabSelected  int
	tabs         []tabItem
	projectTags  []store.NameCount
	contextTags  []store.NameCount
	tasks        []*store.Task
	selected     int
	loading      bool
	status       string
	errText      string
	showHelp     bool
	deleteTaskID int64
	searchInput  textinput.Model
	searchMode   bool
	commandInput textinput.Model
	commandMode  bool
	taskForm     taskFormState
	taskTable    table.Model
	detail       viewport.Model
	help         help.Model
	spinner      spinner.Model
}

type tabItem struct {
	label string
	state string
	count int64
}

func newModel(svc service.Service, opts Options) model {
	search := newSearchInput(searchInputWidth)
	command := newCommandInput(searchInputWidth)
	tabs := defaultTabs()
	layout := calculateLayout(defaultWidth, defaultHeight)
	styleSet := newStyles(SelectTheme(opts.ThemeName), opts.NoColor)
	taskTable := newTaskTable(styleSet, layout)
	helpModel := newHelpModel(styleSet, defaultWidth)
	spinnerModel := newSpinnerModel(styleSet)
	return model{
		svc:          svc,
		keys:         defaultKeyMap(),
		styles:       styleSet,
		viewportW:    defaultWidth,
		viewportH:    defaultHeight,
		layout:       layout,
		view:         viewTasks,
		filters:      defaultFiltersWithState(tabs[0].state),
		tabSelected:  0,
		tabs:         tabs,
		loading:      true,
		status:       statusLoadingText,
		searchInput:  search,
		commandInput: command,
		taskForm:     inactiveTaskForm(searchInputWidth),
		taskTable:    taskTable,
		detail:       viewport.New(searchInputWidth, defaultHeight),
		help:         helpModel,
		spinner:      spinnerModel,
	}
}

func newCommandInput(width int) textinput.Model {
	input := textinput.New()
	input.Prompt = ": "
	input.Placeholder = "create/update/filter command"
	input.CharLimit = 512
	input.Width = width
	return input
}

func newSearchInput(width int) textinput.Model {
	input := textinput.New()
	input.Prompt = "/ "
	input.Placeholder = "search tasks"
	input.CharLimit = 256
	input.Width = width
	return input
}

func defaultTabs() []tabItem {
	return []tabItem{
		{label: "Inbox", state: "inbox"},
		{label: "Now", state: "now"},
		{label: "Waiting", state: "waiting"},
		{label: "Later", state: "later"},
		{label: "Done", state: "done"},
		{label: "All", state: ""},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		refreshDataCmd(m.svc, m.filters),
		tea.WindowSize(),
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportW = value.Width
		m.viewportH = value.Height
		m.layout = calculateLayout(value.Width, value.Height)
		m.help.Width = value.Width
		m.searchInput.Width = searchWidthForLayout(m.layout)
		m.commandInput.Width = searchWidthForLayout(m.layout)
		m.taskForm = m.taskForm.withWidth(searchWidthForLayout(m.layout))
		return m, clearScreenCmd()
	case spinner.TickMsg:
		if !m.loading {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(value)
		return m, cmd
	case tasksLoadedMsg:
		m = m.applyTasksLoaded(value)
		return m, nil
	case tagCountsLoadedMsg:
		m = m.applyTagCountsLoaded(value)
		return m, nil
	case stateCountsLoadedMsg:
		m = m.applyStateCountsLoaded(value)
		return m, nil
	case actionResultMsg:
		return m.applyActionResult(value)
	case tea.KeyMsg:
		return m.handleKey(value)
	default:
		return m, nil
	}
}

func (m model) applyTasksLoaded(msg tasksLoadedMsg) model {
	m.loading = false
	if msg.err != nil {
		m.errText = msg.err.Error()
		m.status = "failed to load tasks"
		return m
	}

	selectedID := selectedTaskID(m.tasks, m.selected)
	m.errText = ""
	m.tasks = msg.tasks
	m.selected = clampTaskSelection(m.tasks, selectedID)
	if selectedTaskID(m.tasks, m.selected) != selectedID {
		m.detail.SetYOffset(0)
	}
	showState := m.filters.state == ""
	setTaskTableData(
		&m.taskTable,
		taskTableColumns(showState, m.taskTable.Width()),
		taskTableRows(m.tasks, showState),
		clampZeroToMax(m.selected, len(m.tasks)-1),
	)
	m.status = ""
	return m
}

func (m model) applyTagCountsLoaded(msg tagCountsLoadedMsg) model {
	if msg.err != nil {
		m.errText = msg.err.Error()
		m.status = "failed to load tags"
		return m
	}
	m.projectTags = msg.projects
	m.contextTags = msg.contexts
	return m
}

func (m model) applyStateCountsLoaded(msg stateCountsLoadedMsg) model {
	if msg.err != nil {
		m.errText = msg.err.Error()
		m.status = "failed to load states"
		return m
	}
	m.tabs = updateTabCounts(defaultTabs(), msg.counts)
	return m
}

func updateTabCounts(tabs []tabItem, counts map[store.State]int64) []tabItem {
	result := make([]tabItem, len(tabs))
	for i, tab := range tabs {
		result[i] = tab
		switch tab.label {
		case "All":
			for _, c := range counts {
				result[i].count += c
			}
		case "Inbox":
			result[i].count = counts[store.StateInbox]
		case "Now":
			result[i].count = counts[store.StateNow]
		case "Waiting":
			result[i].count = counts[store.StateWaiting]
		case "Later":
			result[i].count = counts[store.StateLater]
		case "Done":
			result[i].count = counts[store.StateDone]
		}
	}
	return result
}

func (m model) applyActionResult(msg actionResultMsg) (tea.Model, tea.Cmd) {
	m.loading = false
	if msg.err != nil {
		m.errText = msg.err.Error()
		m.status = "action failed"
		return m, nil
	}
	m.errText = ""
	m.status = msg.status
	return m, refreshDataCmd(m.svc, m.filters)
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.taskForm.active() {
		return m.handleTaskFormInput(msg)
	}

	if m.searchMode {
		return m.handleSearchInput(msg)
	}

	if m.commandMode {
		return m.handleCommandInput(msg)
	}

	if modelValue, cmd, handled := m.handleGlobalKey(msg); handled {
		return modelValue, cmd
	}

	return m.handleTaskViewKey(msg)
}

func (m model) handleGlobalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	if key.Matches(msg, m.keys.Quit) {
		return m, tea.Quit, true
	}

	if key.Matches(msg, m.keys.Help) {
		m.showHelp = !m.showHelp
		return m, nil, true
	}
	if key.Matches(msg, m.keys.Refresh) {
		modelValue, cmd := m.refreshWithStatus(statusRefreshText)
		return modelValue, cmd, true
	}
	if key.Matches(msg, m.keys.Esc) {
		if m.filters.search != "" {
			m.filters.search = ""
			modelValue, cmd := m.refreshWithStatus("search cleared")
			return modelValue, cmd, true
		}
		m.deleteTaskID = 0
		m.showHelp = false
		return m, nil, true
	}

	return m, nil, false
}

func (m model) handleTaskViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.NextView) {
		m.tabSelected = (m.tabSelected + 1) % len(m.tabs)
		m.filters.state = m.tabs[m.tabSelected].state
		return m.refreshWithStatus(fmt.Sprintf("filter: %s", m.tabs[m.tabSelected].label))
	}
	if key.Matches(msg, m.keys.PrevView) {
		if m.tabSelected == 0 {
			m.tabSelected = len(m.tabs) - 1
		} else {
			m.tabSelected--
		}
		m.filters.state = m.tabs[m.tabSelected].state
		return m.refreshWithStatus(fmt.Sprintf("filter: %s", m.tabs[m.tabSelected].label))
	}
	if key.Matches(msg, m.keys.Search) {
		return m.startSearch()
	}
	if key.Matches(msg, m.keys.Command) {
		return m.startCommand()
	}
	if key.Matches(msg, m.keys.Add) {
		return m.startTaskAddForm()
	}
	if key.Matches(msg, m.keys.Edit) {
		return m.startTaskEditForm()
	}
	if key.Matches(msg, m.keys.ProjectFilter) {
		return m.cycleProjectFilter()
	}
	if key.Matches(msg, m.keys.ContextFilter) {
		return m.cycleContextFilter()
	}

	if key.Matches(msg, m.keys.Up) {
		return m.moveCurrentSelection(-1)
	}
	if key.Matches(msg, m.keys.Down) {
		return m.moveCurrentSelection(1)
	}
	if key.Matches(msg, m.keys.Top) {
		return m.jumpCurrentSelection(true)
	}
	if key.Matches(msg, m.keys.Bottom) {
		return m.jumpCurrentSelection(false)
	}
	if key.Matches(msg, m.keys.DetailUp) {
		return m.scrollDetailUp()
	}
	if key.Matches(msg, m.keys.DetailDown) {
		return m.scrollDetailDown()
	}

	task := m.selectedTask()
	if task == nil {
		return m, nil
	}

	return m.handleTaskActionKeys(msg, task)
}

func (m model) startSearch() (tea.Model, tea.Cmd) {
	m.commandMode = false
	m.commandInput.Blur()
	m.searchMode = true
	m.searchInput.SetValue(m.filters.search)
	m.searchInput.CursorEnd()
	m.status = "search: type query and press enter"
	cmd := m.searchInput.Focus()
	return m, cmd
}

func (m model) startCommand() (tea.Model, tea.Cmd) {
	m.searchMode = false
	m.searchInput.Blur()
	m.commandMode = true
	m.commandInput.SetValue("")
	m.commandInput.CursorEnd()
	m.status = "command: enter natural language command"
	cmd := m.commandInput.Focus()
	return m, cmd
}

func (m model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == keyCtrlC {
		return m, tea.Quit
	}
	if key.Matches(msg, m.keys.Esc) {
		m.searchMode = false
		m.searchInput.Blur()
		m.status = "search cancelled"
		return m, nil
	}
	if key.Matches(msg, m.keys.Select) {
		m.searchMode = false
		m.searchInput.Blur()
		m.filters.search = strings.TrimSpace(m.searchInput.Value())
		if m.filters.search == "" {
			return m.refreshWithStatus("search cleared")
		}
		return m.refreshWithStatus("search set")
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m model) handleCommandInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == keyCtrlC {
		return m, tea.Quit
	}
	if key.Matches(msg, m.keys.Esc) {
		m.commandMode = false
		m.commandInput.Blur()
		m.status = statusCommandCancelled
		return m, nil
	}
	if key.Matches(msg, m.keys.Select) {
		return m.submitCommandInput()
	}

	var cmd tea.Cmd
	m.commandInput, cmd = m.commandInput.Update(msg)
	return m, cmd
}

func (m model) submitCommandInput() (tea.Model, tea.Cmd) {
	commandText := strings.TrimSpace(m.commandInput.Value())
	m.commandMode = false
	m.commandInput.Blur()

	if commandText == "" {
		m.status = statusCommandCancelled
		return m, nil
	}

	parsed, err := nlp.Parse(commandText, nlp.ParseOptions{Now: time.Now()})
	if err != nil {
		m.status = "command parse failed"
		m.errText = err.Error()
		return m, nil
	}

	selectedID := selectedTaskID(m.tasks, m.selected)
	var selectedIDPtr *int64
	if selectedID > 0 {
		selectedIDPtr = &selectedID
	}

	plan, err := nlpcompile.Build(parsed, nlpcompile.BuildOptions{SelectedTaskID: selectedIDPtr, Now: time.Now()})
	if err != nil {
		m.status = statusCommandRejected
		m.errText = err.Error()
		return m, nil
	}

	m.errText = ""
	switch plan.Intent {
	case nlp.IntentCreate:
		if plan.Create == nil {
			m.status = statusCommandRejected
			m.errText = "create plan missing request"
			return m, nil
		}
		return m.startAction("creating task...", createTaskCmd(m.svc, *plan.Create))
	case nlp.IntentUpdate:
		if plan.Update == nil {
			m.status = statusCommandRejected
			m.errText = "update plan missing request"
			return m, nil
		}
		return m.startAction("updating task...", updateTaskCmd(m.svc, *plan.Update))
	case nlp.IntentFilter:
		if plan.Filter == nil {
			m.status = statusCommandRejected
			m.errText = "filter plan missing request"
			return m, nil
		}
		modelValue, cmd := m.applyCommandFilter(*plan.Filter)
		return modelValue, cmd
	case nlp.IntentUnknown:
		m.status = statusCommandRejected
		m.errText = "unknown command intent"
		return m, nil
	default:
		m.status = statusCommandRejected
		m.errText = "unknown command intent"
		return m, nil
	}
}

func (m model) startTaskAddForm() (tea.Model, tea.Cmd) {
	m.taskForm = startAddTaskForm(searchWidthForLayout(m.layout))
	m.status = "add task form"
	return m, nil
}

func (m model) startTaskEditForm() (tea.Model, tea.Cmd) {
	task := m.selectedTask()
	if task == nil {
		m.status = "no task selected"
		return m, nil
	}
	m.taskForm = startEditTaskForm(task, searchWidthForLayout(m.layout))
	m.status = fmt.Sprintf("edit task #%d", task.ID)
	return m, nil
}

func (m model) handleTaskFormInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == keyCtrlC {
		return m, tea.Quit
	}
	if key.Matches(msg, m.keys.Save) {
		return m.submitTaskForm()
	}
	if key.Matches(msg, m.keys.Esc) {
		return m.handleTaskFormEsc()
	}
	if m.taskForm.editing {
		return m.handleTaskFormEditKey(msg)
	}
	return m.handleTaskFormNavKey(msg)
}

func (m model) handleTaskFormEsc() (tea.Model, tea.Cmd) {
	if m.taskForm.editing {
		m.taskForm = m.taskForm.withField(m.taskForm.field)
		m.status = "field edit cancelled"
		return m, nil
	}
	m.taskForm = inactiveTaskForm(searchWidthForLayout(m.layout))
	m.status = "task form cancelled"
	return m, nil
}

func (m model) handleTaskFormNavKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Up) || msg.String() == "k" || msg.String() == "shift+tab" {
		m.taskForm = m.taskForm.previousField()
		return m, nil
	}
	if key.Matches(msg, m.keys.Down) || msg.String() == "j" || msg.String() == "tab" {
		next, _ := m.taskForm.nextField()
		m.taskForm = next
		return m, nil
	}
	if key.Matches(msg, m.keys.Select) || msg.String() == "i" {
		var cmd tea.Cmd
		m.taskForm, cmd = m.taskForm.startEditing()
		return m, cmd
	}
	return m, nil
}

func (m model) handleTaskFormEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "shift+tab" || msg.String() == "ctrl+k" {
		return m.moveTaskFormField(false)
	}
	if msg.String() == "tab" || msg.String() == "ctrl+j" {
		return m.moveTaskFormField(true)
	}
	if key.Matches(msg, m.keys.Select) && m.taskForm.field != taskFormFieldNotes {
		return m.moveTaskFormField(true)
	}

	var cmd tea.Cmd
	m.taskForm, cmd = m.taskForm.update(msg)
	return m, cmd
}

func (m model) moveTaskFormField(forward bool) (tea.Model, tea.Cmd) {
	form := m.taskForm.commitInput()
	if !forward {
		previous := form.previousField()
		var cmd tea.Cmd
		m.taskForm, cmd = previous.startEditing()
		return m, cmd
	}

	next, done := form.nextField()
	if done {
		m.taskForm = form.stopEditing()
		m.status = statusTaskFormLastField
		return m, nil
	}

	var cmd tea.Cmd
	m.taskForm, cmd = next.startEditing()
	return m, cmd
}

func (m model) submitTaskForm() (tea.Model, tea.Cmd) {
	form := m.taskForm.commitInput().stopEditing()
	if strings.TrimSpace(form.values.title) == "" {
		m.status = statusTitleRequired
		m.errText = statusTitleRequired
		m.taskForm = form.withField(taskFormFieldTitle)
		return m, nil
	}

	m.errText = ""
	m.taskForm = inactiveTaskForm(searchWidthForLayout(m.layout))
	if form.mode == taskFormEdit {
		return m.startAction("updating task...", fullUpdateTaskCmd(m.svc, form.fullUpdateRequest()))
	}
	return m.startAction("creating task...", createTaskCmd(m.svc, form.createRequest()))
}

func (m model) cycleProjectFilter() (tea.Model, tea.Cmd) {
	next, ok := nextTagFilterValue(m.filters.project, m.projectTags)
	if !ok {
		m.status = "no projects available"
		return m, nil
	}
	m.filters.project = next
	if next == "" {
		return m.refreshWithStatus("project filter cleared")
	}
	return m.refreshWithStatus("project filter: " + next)
}

func (m model) cycleContextFilter() (tea.Model, tea.Cmd) {
	next, ok := nextTagFilterValue(m.filters.context, m.contextTags)
	if !ok {
		m.status = "no contexts available"
		return m, nil
	}
	m.filters.context = next
	if next == "" {
		return m.refreshWithStatus("context filter cleared")
	}
	return m.refreshWithStatus("context filter: " + next)
}

func (m model) handleTaskActionKeys(msg tea.KeyMsg, task *store.Task) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Done) {
		return m.startAction(fmt.Sprintf("marking task #%d done...", task.ID), setDoneCmd(m.svc, task.ID, true))
	}
	if key.Matches(msg, m.keys.Undo) {
		return m.startAction(fmt.Sprintf("reopening task #%d...", task.ID), setDoneCmd(m.svc, task.ID, false))
	}
	if key.Matches(msg, m.keys.Inbox) {
		return m.startAction(
			fmt.Sprintf("moving task #%d to inbox...", task.ID),
			setStateCmd(m.svc, task.ID, string(store.StateInbox)),
		)
	}
	if key.Matches(msg, m.keys.Now) {
		return m.startAction(
			fmt.Sprintf("moving task #%d to now...", task.ID),
			setStateCmd(m.svc, task.ID, string(store.StateNow)),
		)
	}
	if key.Matches(msg, m.keys.Waiting) {
		return m.startAction(
			fmt.Sprintf("moving task #%d to waiting...", task.ID),
			setStateCmd(m.svc, task.ID, string(store.StateWaiting)),
		)
	}
	if key.Matches(msg, m.keys.Later) {
		return m.startAction(
			fmt.Sprintf("moving task #%d to later...", task.ID),
			setStateCmd(m.svc, task.ID, string(store.StateLater)),
		)
	}
	if key.Matches(msg, m.keys.Delete) {
		if m.deleteTaskID == task.ID {
			return m.startAction(fmt.Sprintf("deleting task #%d...", task.ID), deleteTaskCmd(m.svc, task.ID))
		}
		m.deleteTaskID = task.ID
		m.status = fmt.Sprintf("press D again to delete task #%d", task.ID)
		return m, nil
	}

	return m, nil
}

func (m model) moveCurrentSelection(delta int) (tea.Model, tea.Cmd) {
	m.deleteTaskID = 0
	next := clampZeroToMax(m.selected+delta, len(m.tasks)-1)
	if next != m.selected {
		m.detail.SetYOffset(0)
	}
	m.selected = next
	m.taskTable.SetCursor(m.selected)
	return m, nil
}

func (m model) jumpCurrentSelection(top bool) (tea.Model, tea.Cmd) {
	m.deleteTaskID = 0
	if len(m.tasks) == 0 {
		return m, nil
	}
	previous := m.selected
	if top {
		m.selected = 0
	} else {
		m.selected = len(m.tasks) - 1
	}
	if m.selected != previous {
		m.detail.SetYOffset(0)
	}
	m.taskTable.SetCursor(m.selected)
	return m, nil
}

func (m model) scrollDetailUp() (tea.Model, tea.Cmd) {
	if m.selectedTask() == nil {
		return m, nil
	}
	m.detail.HalfPageUp()
	return m, nil
}

func (m model) scrollDetailDown() (tea.Model, tea.Cmd) {
	if m.selectedTask() == nil {
		return m, nil
	}
	m.detail.HalfPageDown()
	return m, nil
}

func (m model) refreshWithStatus(status string) (tea.Model, tea.Cmd) {
	m.loading = true
	m.deleteTaskID = 0
	m.status = status
	return m, tea.Batch(refreshDataCmd(m.svc, m.filters), m.spinner.Tick)
}

func (m model) startAction(status string, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	m.loading = true
	m.deleteTaskID = 0
	m.status = status
	return m, tea.Batch(cmd, m.spinner.Tick)
}

func (m model) selectedTask() *store.Task {
	if len(m.tasks) == 0 {
		return nil
	}
	if m.selected < 0 || m.selected >= len(m.tasks) {
		return nil
	}
	return m.tasks[m.selected]
}

func (m model) View() string {
	if m.searchMode {
		return m.renderModalWithBackdrop(m.renderInputModal(m.searchInput, "Search"))
	}

	if m.commandMode {
		return m.renderModalWithBackdrop(m.renderInputModal(m.commandInput, "Command"))
	}

	tabs := m.renderTabs()
	statusLine := m.renderStatusLine()

	renderModel := m
	if statusLine == "" {
		renderModel.layout.bodyHeight++
	}

	body := renderModel.viewTasks()

	footer := m.renderFooter()
	parts := []string{tabs}
	parts = append(parts, body)
	if statusLine != "" {
		parts = append(parts, statusLine)
	}
	parts = append(parts, footer)
	if m.showHelp {
		parts = append(parts, m.renderHelp())
	}

	content := strings.Join(parts, "\n")
	if m.viewportW > 0 && m.viewportH > 0 {
		content = lipgloss.Place(
			m.viewportW,
			m.viewportH,
			lipgloss.Left,
			lipgloss.Top,
			content,
		)
	}

	return m.styles.app.Render(content)
}

func (m model) renderModalWithBackdrop(modalContent string) string {
	canvasWidth, canvasHeight := m.modalCanvasSize()
	content := lipgloss.Place(
		canvasWidth,
		canvasHeight,
		lipgloss.Left,
		lipgloss.Top,
		modalContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceBackground(m.styles.modalShade.GetBackground()),
	)
	return m.styles.app.Render(content)
}

func selectedTaskID(tasks []*store.Task, selected int) int64 {
	if selected < 0 || selected >= len(tasks) || tasks[selected] == nil {
		return 0
	}
	return tasks[selected].ID
}

func clampTaskSelection(tasks []*store.Task, preferredID int64) int {
	if len(tasks) == 0 {
		return 0
	}
	if preferredID == 0 {
		return 0
	}
	for i, task := range tasks {
		if task != nil && task.ID == preferredID {
			return i
		}
	}
	return 0
}

func clampZeroToMax(value int, maximum int) int {
	if maximum < 0 {
		return 0
	}
	if value < 0 {
		return 0
	}
	if value > maximum {
		return maximum
	}
	return value
}

func searchWidthForLayout(layout layoutSpec) int {
	if layout.narrow {
		return max(searchMinWidth, layout.listWidth-searchNarrowPad)
	}
	return searchInputWidth
}

func (m model) applyCommandFilter(req service.ListTasksRequest) (tea.Model, tea.Cmd) {
	m.filters.project = req.Project
	m.filters.context = req.Context
	m.filters.search = req.Search
	m.filters.dueOnly = req.DueOnly

	switch {
	case req.State != "":
		m.filters.state = req.State
	case req.All:
		m.filters.state = ""
	case req.DoneOnly:
		m.filters.state = string(store.StateDone)
	}

	m.tabSelected = tabIndexForState(m.tabs, m.filters.state, m.tabSelected)
	if m.filters.state == "" {
		return m.refreshWithStatus("filter command applied")
	}
	return m.refreshWithStatus(fmt.Sprintf("filter: %s", m.tabs[m.tabSelected].label))
}

func tabIndexForState(tabs []tabItem, state string, fallback int) int {
	for i, tab := range tabs {
		if tab.state == state {
			return i
		}
	}
	if fallback >= 0 && fallback < len(tabs) {
		return fallback
	}
	return 0
}

func clearScreenCmd() tea.Cmd {
	return tea.ClearScreen
}

func nextTagFilterValue(current string, tags []store.NameCount) (string, bool) {
	if len(tags) == 0 {
		return "", false
	}

	if current == "" {
		return tags[0].Name, true
	}

	for idx, tag := range tags {
		if tag.Name != current {
			continue
		}
		next := idx + 1
		if next >= len(tags) {
			return "", true
		}
		return tags[next].Name, true
	}

	return tags[0].Name, true
}
