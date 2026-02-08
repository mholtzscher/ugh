package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

type viewKind int

const (
	viewTasks viewKind = iota
)

const (
	defaultWidth      = 120
	defaultHeight     = 32
	statusLoadingText = "loading tasks..."
	statusRefreshText = "refreshing..."
	searchInputWidth  = 48
	searchMinWidth    = 1
	searchNarrowPad   = 6
	modalOuterPad     = 2
	modalMinWidth     = 52
	modalMinInput     = 1
	modalInputPad     = 8
	modalWidthNum     = 2
	modalWidthDen     = 3
	modalMinBodyWidth = 1
	keyCtrlC          = "ctrl+c"
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
	taskForm     taskFormState
}

type tabItem struct {
	label string
	state string
	count int64
}

func newModel(svc service.Service, opts Options) model {
	search := newSearchInput(searchInputWidth)
	tabs := defaultTabs()
	return model{
		svc:         svc,
		keys:        defaultKeyMap(),
		styles:      newStyles(SelectTheme(opts.ThemeName), opts.NoColor),
		viewportW:   defaultWidth,
		viewportH:   defaultHeight,
		layout:      calculateLayout(defaultWidth, defaultHeight),
		view:        viewTasks,
		filters:     defaultFiltersWithState(tabs[0].state),
		tabSelected: 0,
		tabs:        tabs,
		loading:     true,
		status:      statusLoadingText,
		searchInput: search,
		taskForm:    inactiveTaskForm(searchInputWidth),
	}
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
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportW = value.Width
		m.viewportH = value.Height
		m.layout = calculateLayout(value.Width, value.Height)
		m.searchInput.Width = searchWidthForLayout(m.layout)
		m.taskForm = m.taskForm.withWidth(searchWidthForLayout(m.layout))
		return m, clearScreenCmd()
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
	m.status = fmt.Sprintf("loaded %d tasks", len(m.tasks))
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
	if key.Matches(msg, m.keys.CycleCompletion) {
		m.filters = m.filters.cycleCompletion()
		m.status = fmt.Sprintf("completion filter: %s", m.filters.completionText())
		return m.refreshWithStatus(m.status)
	}
	if key.Matches(msg, m.keys.Search) {
		return m.startSearch()
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

	task := m.selectedTask()
	if task == nil {
		return m, nil
	}

	return m.handleTaskActionKeys(msg, task)
}

func (m model) startSearch() (tea.Model, tea.Cmd) {
	m.searchMode = true
	m.searchInput.SetValue(m.filters.search)
	m.searchInput.CursorEnd()
	m.status = "search: type query and press enter"
	cmd := m.searchInput.Focus()
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

func (m model) startTaskAddForm() (tea.Model, tea.Cmd) {
	m.taskForm = startAddTaskForm(searchWidthForLayout(m.layout))
	m.status = "add task form"
	cmd := m.taskForm.input.Focus()
	return m, cmd
}

func (m model) startTaskEditForm() (tea.Model, tea.Cmd) {
	task := m.selectedTask()
	if task == nil {
		m.status = "no task selected"
		return m, nil
	}
	m.taskForm = startEditTaskForm(task, searchWidthForLayout(m.layout))
	m.status = fmt.Sprintf("edit task #%d", task.ID)
	cmd := m.taskForm.input.Focus()
	return m, cmd
}

func (m model) handleTaskFormInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == keyCtrlC {
		return m, tea.Quit
	}

	if key.Matches(msg, m.keys.Esc) {
		m.taskForm = inactiveTaskForm(searchWidthForLayout(m.layout))
		m.status = "task form cancelled"
		return m, nil
	}
	if msg.String() == "shift+tab" {
		m.taskForm = m.taskForm.commitInput().previousField()
		cmd := m.taskForm.input.Focus()
		return m, cmd
	}
	if key.Matches(msg, m.keys.Select) {
		form := m.taskForm.commitInput()
		next, done := form.nextField()
		if !done {
			m.taskForm = next
			cmd := m.taskForm.input.Focus()
			return m, cmd
		}

		if strings.TrimSpace(form.values.title) == "" {
			m.status = "title is required"
			m.errText = "title is required"
			m.taskForm = form.withField(taskFormFieldTitle)
			cmd := m.taskForm.input.Focus()
			return m, cmd
		}

		m.errText = ""
		m.taskForm = inactiveTaskForm(searchWidthForLayout(m.layout))
		if form.mode == taskFormEdit {
			return m.startAction("updating task...", fullUpdateTaskCmd(m.svc, form.fullUpdateRequest()))
		}
		return m.startAction("creating task...", createTaskCmd(m.svc, form.createRequest()))
	}

	var cmd tea.Cmd
	m.taskForm.input, cmd = m.taskForm.input.Update(msg)
	return m, cmd
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
	m.selected = clampZeroToMax(m.selected+delta, len(m.tasks)-1)
	return m, nil
}

func (m model) jumpCurrentSelection(top bool) (tea.Model, tea.Cmd) {
	m.deleteTaskID = 0
	if len(m.tasks) == 0 {
		return m, nil
	}
	if top {
		m.selected = 0
	} else {
		m.selected = len(m.tasks) - 1
	}
	return m, nil
}

func (m model) refreshWithStatus(status string) (tea.Model, tea.Cmd) {
	m.loading = true
	m.deleteTaskID = 0
	m.status = status
	return m, refreshDataCmd(m.svc, m.filters)
}

func (m model) startAction(status string, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	m.loading = true
	m.deleteTaskID = 0
	m.status = status
	return m, cmd
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
	tabs := m.renderTabs()

	body := m.viewTasks()

	if m.taskForm.active() {
		body = m.renderTaskFormModal()
	}

	statusLine := m.renderStatusLine()
	footer := m.renderFooter()
	parts := []string{tabs}
	if m.searchMode {
		parts = append(parts, m.styles.panel.Render(m.searchInput.View()))
	}
	parts = append(parts, body, statusLine, footer)
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
	return m.styles.foot.Render(m.status + "  " + m.filters.statusText())
}

func (m model) renderFooter() string {
	keys := []string{
		m.styles.key.Render("j/k") + " move",
		m.styles.key.Render("[/]") + " tab",
		m.styles.key.Render("t") + " todo/all/done",
		m.styles.key.Render("/") + " search",
		m.styles.key.Render("p/c") + " project/context",
		m.styles.key.Render("a/e") + " add/edit",
		m.styles.key.Render("x/u") + " done/undo",
		m.styles.key.Render("D") + " delete",
		m.styles.key.Render("?") + " help",
		m.styles.key.Render("q") + " quit",
	}
	return m.styles.foot.Render(strings.Join(keys, "  "))
}

func (m model) renderHelp() string {
	text := strings.Join([]string{
		"TUI keybindings",
		"- j/k or up/down: move selection",
		"- [/]: switch state tab",
		"- t: cycle todo -> all -> done",
		"- /: edit search query (enter applies, esc cancels)",
		"- p/c: cycle project or context filter",
		"- a/e: add or edit task form",
		"- x/u: done or undo selected task",
		"- i/n/w/l: move selected task state",
		"- D: delete selected task (double press to confirm)",
		"- q: quit",
	}, "\n")
	return m.styles.help.Render(text)
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
