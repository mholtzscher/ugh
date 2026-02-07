package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mholtzscher/ugh/internal/config"
	daemonservice "github.com/mholtzscher/ugh/internal/daemon/service"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

type viewKind int

const (
	viewTasks viewKind = iota
	viewCalendar
	viewSystem
)

const (
	defaultWidth      = 120
	defaultHeight     = 32
	viewCount         = 3
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
	svc              service.Service
	keys             keyMap
	styles           styles
	theme            Theme
	noColor          bool
	configPath       string
	appConfig        config.Config
	viewportW        int
	viewportH        int
	layout           layoutSpec
	view             viewKind
	focus            paneFocus
	filters          listFilters
	navItems         []navItem
	navSelected      int
	projectTags      []store.NameCount
	contextTags      []store.NameCount
	stateCounts      map[store.State]int64
	tasks            []*store.Task
	selected         int
	calendarTasks    []*store.Task
	calendarSelected int
	syncStatus       *service.SyncStatus
	syncStatusErr    string
	daemonManager    string
	daemonStatus     *daemonservice.Status
	daemonStatusErr  string
	daemonLogPath    string
	loading          bool
	status           string
	errText          string
	showHelp         bool
	deleteTaskID     int64
	searchInput      textinput.Model
	searchMode       bool
	taskForm         taskFormState
	configForm       configFormState
}

func newModel(svc service.Service, opts Options) model {
	theme := SelectTheme(opts.ThemeName)
	stateCounts := defaultStateCounts()
	search := newSearchInput(searchInputWidth)
	appCfg := opts.Config
	if appCfg.Version == 0 {
		appCfg.Version = config.DefaultVersion
	}
	if strings.TrimSpace(appCfg.UI.Theme) == "" {
		appCfg.UI.Theme = config.DefaultUITheme
	}
	return model{
		svc:         svc,
		keys:        defaultKeyMap(),
		styles:      newStyles(theme, opts.NoColor),
		theme:       theme,
		noColor:     opts.NoColor,
		configPath:  opts.ConfigPath,
		appConfig:   appCfg,
		viewportW:   defaultWidth,
		viewportH:   defaultHeight,
		layout:      calculateLayout(defaultWidth, defaultHeight),
		view:        viewTasks,
		focus:       focusList,
		filters:     defaultFilters(),
		stateCounts: stateCounts,
		navItems:    buildNavItems(stateCounts, nil, nil),
		loading:     true,
		status:      statusLoadingText,
		searchInput: search,
		taskForm:    inactiveTaskForm(searchInputWidth),
		configForm:  inactiveConfigForm(searchInputWidth),
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
		m.configForm = m.configForm.withWidth(searchWidthForLayout(m.layout))
		return m, clearScreenCmd()
	case tasksLoadedMsg:
		m = m.applyTasksLoaded(value)
		return m, nil
	case calendarTasksLoadedMsg:
		m = m.applyCalendarTasksLoaded(value)
		return m, nil
	case syncStatusLoadedMsg:
		m = m.applySyncStatusLoaded(value)
		return m, nil
	case daemonStatusLoadedMsg:
		m = m.applyDaemonStatusLoaded(value)
		return m, nil
	case tagCountsLoadedMsg:
		m = m.applyTagCountsLoaded(value)
		return m, nil
	case stateCountsLoadedMsg:
		m = m.applyStateCountsLoaded(value)
		return m, nil
	case actionResultMsg:
		return m.applyActionResult(value)
	case configSavedMsg:
		return m.applyConfigSaved(value)
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

func (m model) applyCalendarTasksLoaded(msg calendarTasksLoadedMsg) model {
	if msg.err != nil {
		m.errText = msg.err.Error()
		m.status = "failed to load calendar tasks"
		return m
	}

	selectedID := selectedTaskID(m.calendarTasks, m.calendarSelected)
	m.calendarTasks = msg.tasks
	m.calendarSelected = clampTaskSelection(m.calendarTasks, selectedID)
	return m
}

func (m model) applySyncStatusLoaded(msg syncStatusLoadedMsg) model {
	if msg.err != nil {
		m.syncStatus = nil
		m.syncStatusErr = msg.err.Error()
		return m
	}

	m.syncStatusErr = ""
	m.syncStatus = msg.status
	return m
}

func (m model) applyDaemonStatusLoaded(msg daemonStatusLoadedMsg) model {
	if msg.err != nil {
		m.daemonStatus = nil
		m.daemonStatusErr = msg.err.Error()
		m.daemonLogPath = ""
		m.daemonManager = ""
		return m
	}

	m.daemonStatusErr = ""
	m.daemonManager = msg.managerName
	m.daemonStatus = msg.status
	m.daemonLogPath = msg.logPath
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
	m = m.rebuildNavItems()
	return m
}

func (m model) applyStateCountsLoaded(msg stateCountsLoadedMsg) model {
	if msg.err != nil {
		m.errText = msg.err.Error()
		m.status = "failed to load states"
		return m
	}
	m.stateCounts = msg.counts
	m = m.rebuildNavItems()
	return m
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

func (m model) applyConfigSaved(msg configSavedMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.errText = msg.err.Error()
		m.status = "failed to save config"
		return m, nil
	}

	m.errText = ""
	m.appConfig = msg.cfg
	m.configPath = msg.path
	if strings.TrimSpace(msg.cfg.UI.Theme) != "" {
		m.theme = SelectTheme(msg.cfg.UI.Theme)
		m.styles = newStyles(m.theme, m.noColor)
	}
	m.status = "config saved"
	return m, refreshDataCmd(m.svc, m.filters)
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.taskForm.active() {
		return m.handleTaskFormInput(msg)
	}
	if m.configForm.isActive() {
		return m.handleConfigFormInput(msg)
	}

	if m.searchMode {
		return m.handleSearchInput(msg)
	}

	if modelValue, cmd, handled := m.handleGlobalKey(msg); handled {
		return modelValue, cmd
	}

	if m.view == viewCalendar {
		if modelValue, cmd, handled := m.handleTaskScopeKey(msg); handled {
			return modelValue, cmd
		}
		return m.handleCalendarViewKey(msg)
	}

	if m.view == viewSystem {
		return m.handleSystemViewKey(msg)
	}

	if m.view != viewTasks {
		return m, nil
	}

	if modelValue, cmd, handled := m.handleTaskScopeKey(msg); handled {
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
	if key.Matches(msg, m.keys.NextView) {
		m.view = (m.view + 1) % viewCount
		return m, nil, true
	}
	if key.Matches(msg, m.keys.PrevView) {
		m.view = prevView(m.view)
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

func (m model) handleTaskScopeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	if key.Matches(msg, m.keys.NextPane) {
		m.focus = nextFocus(m.focus)
		return m, nil, true
	}
	if key.Matches(msg, m.keys.PrevPane) {
		m.focus = prevFocus(m.focus)
		return m, nil, true
	}
	if key.Matches(msg, m.keys.CycleCompletion) {
		m.filters = m.filters.cycleCompletion()
		m.status = fmt.Sprintf("completion filter: %s", m.filters.completionText())
		modelValue, cmd := m.refreshWithStatus(m.status)
		return modelValue, cmd, true
	}
	if key.Matches(msg, m.keys.ToggleDueOnly) {
		m.filters.dueOnly = !m.filters.dueOnly
		if m.filters.dueOnly {
			m.status = "due-only enabled"
		} else {
			m.status = "due-only disabled"
		}
		modelValue, cmd := m.refreshWithStatus(m.status)
		return modelValue, cmd, true
	}
	if key.Matches(msg, m.keys.Search) {
		modelValue, cmd := m.startSearch()
		return modelValue, cmd, true
	}
	if key.Matches(msg, m.keys.Add) {
		modelValue, cmd := m.startTaskAddForm()
		return modelValue, cmd, true
	}
	if key.Matches(msg, m.keys.Edit) {
		modelValue, cmd := m.startTaskEditForm()
		return modelValue, cmd, true
	}
	if key.Matches(msg, m.keys.ProjectFilter) {
		modelValue, cmd := m.cycleProjectFilter()
		return modelValue, cmd, true
	}
	if key.Matches(msg, m.keys.ContextFilter) {
		modelValue, cmd := m.cycleContextFilter()
		return modelValue, cmd, true
	}

	return m, nil, false
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
	task := m.selectedTaskForCurrentView()
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
	if key.Matches(msg, m.keys.PrevPane) {
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

func (m model) startConfigForm() (tea.Model, tea.Cmd) {
	m.configForm = startConfigForm(m.appConfig, searchWidthForLayout(m.layout))
	m.status = "config form"
	cmd := m.configForm.input.Focus()
	return m, cmd
}

func (m model) handleConfigFormInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == keyCtrlC {
		return m, tea.Quit
	}

	if key.Matches(msg, m.keys.Esc) {
		m.configForm = inactiveConfigForm(searchWidthForLayout(m.layout))
		m.status = "config form cancelled"
		return m, nil
	}
	if key.Matches(msg, m.keys.PrevPane) {
		m.configForm = m.configForm.commitInput().previousField()
		cmd := m.configForm.input.Focus()
		return m, cmd
	}
	if key.Matches(msg, m.keys.Select) {
		return m.handleConfigFormSubmit()
	}

	var cmd tea.Cmd
	m.configForm.input, cmd = m.configForm.input.Update(msg)
	return m, cmd
}

func (m model) handleConfigFormSubmit() (tea.Model, tea.Cmd) {
	form := m.configForm.commitInput()
	next, done := form.nextField()
	if !done {
		m.configForm = next
		cmd := m.configForm.input.Focus()
		return m, cmd
	}

	updated, err := form.applyToConfig(m.appConfig)
	if err != nil {
		m.errText = err.Error()
		m.status = "invalid config value"
		m.configForm = form
		cmd := m.configForm.input.Focus()
		return m, cmd
	}

	cfgPath, pathErr := m.resolveConfigPath()
	if pathErr != nil {
		m.errText = pathErr.Error()
		m.status = "unable to resolve config path"
		m.configForm = form
		cmd := m.configForm.input.Focus()
		return m, cmd
	}

	m.errText = ""
	m.configForm = inactiveConfigForm(searchWidthForLayout(m.layout))
	m.status = "saving config..."
	return m, saveConfigCmd(cfgPath, updated)
}

func (m model) resolveConfigPath() (string, error) {
	if strings.TrimSpace(m.configPath) != "" {
		return strings.TrimSpace(m.configPath), nil
	}
	return config.DefaultPath()
}

func (m model) daemonInstallConfigPath() string {
	cfgPath, err := m.resolveConfigPath()
	if err != nil {
		return ""
	}
	if _, statErr := os.Stat(cfgPath); statErr != nil {
		return ""
	}
	return cfgPath
}

func (m model) daemonLogsHint() string {
	if strings.TrimSpace(m.daemonLogPath) != "" {
		return "daemon log path: " + strings.TrimSpace(m.daemonLogPath)
	}
	if m.daemonManager == "systemd" {
		return "daemon logs: run `ugh daemon logs` (journalctl)"
	}
	return "daemon logs: run `ugh daemon logs`"
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

func (m model) handleCalendarViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Up) {
		m.deleteTaskID = 0
		m.calendarSelected = clampZeroToMax(m.calendarSelected-1, len(m.calendarTasks)-1)
		return m, nil
	}
	if key.Matches(msg, m.keys.Down) {
		m.deleteTaskID = 0
		m.calendarSelected = clampZeroToMax(m.calendarSelected+1, len(m.calendarTasks)-1)
		return m, nil
	}
	if key.Matches(msg, m.keys.Top) {
		m.deleteTaskID = 0
		if len(m.calendarTasks) > 0 {
			m.calendarSelected = 0
		}
		return m, nil
	}
	if key.Matches(msg, m.keys.Bottom) {
		m.deleteTaskID = 0
		if len(m.calendarTasks) > 0 {
			m.calendarSelected = len(m.calendarTasks) - 1
		}
		return m, nil
	}

	task := m.selectedCalendarTask()
	if task == nil {
		return m, nil
	}

	if key.Matches(msg, m.keys.Select) {
		m.view = viewTasks
		m.focus = focusList
		m.selected = clampTaskSelection(m.tasks, task.ID)
		m.status = fmt.Sprintf("opened task #%d in Tasks", task.ID)
		return m, nil
	}

	return m.handleTaskActionKeys(msg, task)
}

func (m model) handleSystemViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Config) {
		return m.startConfigForm()
	}
	if key.Matches(msg, m.keys.DaemonInstall) {
		return m.startAction("installing daemon service...", daemonInstallCmd(m.daemonInstallConfigPath()))
	}
	if key.Matches(msg, m.keys.DaemonUninstall) {
		return m.startAction("uninstalling daemon service...", daemonUninstallCmd())
	}
	if key.Matches(msg, m.keys.DaemonStart) {
		return m.startAction("starting daemon service...", daemonStartCmd())
	}
	if key.Matches(msg, m.keys.DaemonStop) {
		return m.startAction("stopping daemon service...", daemonStopCmd())
	}
	if key.Matches(msg, m.keys.DaemonRestart) {
		return m.startAction("restarting daemon service...", daemonRestartCmd())
	}
	if key.Matches(msg, m.keys.DaemonLogs) {
		m.status = m.daemonLogsHint()
		return m, nil
	}
	if key.Matches(msg, m.keys.SyncAll) {
		return m.startAction("syncing with remote...", syncAllCmd(m.svc))
	}
	if key.Matches(msg, m.keys.SyncPull) {
		return m.startAction("pulling changes from remote...", syncPullCmd(m.svc))
	}
	if key.Matches(msg, m.keys.SyncPush) {
		return m.startAction("pushing changes to remote...", syncPushCmd(m.svc))
	}

	return m, nil
}

func (m model) handleTaskViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	if key.Matches(msg, m.keys.Select) {
		return m.handleSelectKey()
	}

	if m.focus == focusNav {
		return m, nil
	}

	task := m.selectedTask()
	if task == nil {
		return m, nil
	}

	return m.handleTaskActionKeys(msg, task)
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
	if m.focus == focusNav {
		m.navSelected = clampZeroToMax(m.navSelected+delta, len(m.navItems)-1)
		return m, nil
	}

	m.selected = clampZeroToMax(m.selected+delta, len(m.tasks)-1)
	return m, nil
}

func (m model) jumpCurrentSelection(top bool) (tea.Model, tea.Cmd) {
	m.deleteTaskID = 0
	if m.focus == focusNav {
		if len(m.navItems) == 0 {
			return m, nil
		}
		if top {
			m.navSelected = 0
		} else {
			m.navSelected = len(m.navItems) - 1
		}
		return m, nil
	}

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

func (m model) handleSelectKey() (tea.Model, tea.Cmd) {
	switch m.focus {
	case focusNav:
		return m.applySelectedScope()
	case focusList:
		m.focus = focusDetail
		return m, nil
	case focusDetail:
		m.focus = focusList
		return m, nil
	default:
		return m, nil
	}
}

func (m model) applySelectedScope() (tea.Model, tea.Cmd) {
	if len(m.navItems) == 0 {
		return m, nil
	}

	item := m.navItems[m.navSelected]
	m.filters.state = ""
	m.filters.project = ""
	m.filters.context = ""

	switch item.kind {
	case navState:
		m.filters.state = item.value
	case navProject:
		m.filters.project = item.value
	case navContext:
		m.filters.context = item.value
	}

	m.status = fmt.Sprintf("scope set: %s", selectedScopeText(m.filters))
	return m.refreshWithStatus(m.status)
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

func (m model) selectedCalendarTask() *store.Task {
	if len(m.calendarTasks) == 0 {
		return nil
	}
	if m.calendarSelected < 0 || m.calendarSelected >= len(m.calendarTasks) {
		return nil
	}
	return m.calendarTasks[m.calendarSelected]
}

func (m model) selectedTaskForCurrentView() *store.Task {
	if m.view == viewCalendar {
		return m.selectedCalendarTask()
	}
	return m.selectedTask()
}

func (m model) rebuildNavItems() model {
	current := m.currentNavSelectionKey()
	m.navItems = buildNavItems(m.stateCounts, m.projectTags, m.contextTags)
	if len(m.navItems) == 0 {
		m.navSelected = 0
		return m
	}
	idx := findNavItemIndex(m.navItems, current.kind, current.value)
	if idx >= 0 {
		m.navSelected = idx
		return m
	}
	m.navSelected = clampZeroToMax(m.navSelected, len(m.navItems)-1)
	return m
}

type navSelectionKey struct {
	kind  navItemKind
	value string
}

func (m model) currentNavSelectionKey() navSelectionKey {
	if len(m.navItems) == 0 {
		return navSelectionKey{}
	}
	if m.navSelected < 0 || m.navSelected >= len(m.navItems) {
		return navSelectionKey{}
	}
	item := m.navItems[m.navSelected]
	return navSelectionKey{kind: item.kind, value: item.value}
}

func (m model) View() string {
	tab := m.renderTabs()

	var body string
	switch m.view {
	case viewTasks:
		body = m.viewTasks()
	case viewCalendar:
		body = m.viewCalendar()
	case viewSystem:
		body = m.viewSystem()
	}

	if m.configForm.isActive() {
		body = m.renderConfigFormModal()
	}
	if m.taskForm.active() {
		body = m.renderTaskFormModal()
	}

	statusLine := m.renderStatusLine()
	footer := m.renderFooter()
	parts := []string{tab}
	if m.searchMode {
		parts = append(parts, m.panelStyle(focusList).Render(m.searchInput.View()))
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
		bodyWidth = m.layout.navWidth + m.layout.listWidth + m.layout.detailWidth
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

func (m model) renderConfigFormModal() string {
	bodyWidth := m.layout.listWidth
	if !m.layout.narrow {
		bodyWidth = m.layout.navWidth + m.layout.listWidth + m.layout.detailWidth
	}
	bodyWidth = max(modalMinBodyWidth, bodyWidth)

	availableWidth := max(modalMinBodyWidth, bodyWidth-modalOuterPad)
	preferredWidth := max(modalMinWidth, (bodyWidth*modalWidthNum)/modalWidthDen)
	modalWidth := min(availableWidth, preferredWidth)
	modalInputWidth := max(modalMinInput, modalWidth-modalInputPad)

	form := m.configForm.withWidth(modalInputWidth)
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

	parts := []string{
		renderTab(m.view == viewTasks, "Tasks"),
		renderTab(m.view == viewCalendar, "Calendar"),
		renderTab(m.view == viewSystem, "System"),
		m.styles.muted.Render("theme: " + m.theme.Name),
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
	if m.view == viewSystem {
		keys := []string{
			m.styles.key.Render("s") + " sync",
			m.styles.key.Render("l") + " pull",
			m.styles.key.Render("p") + " push",
			m.styles.key.Render("o") + " config",
			m.styles.key.Render("I/U") + " install/uninstall",
			m.styles.key.Render("S/T/R") + " start/stop/restart",
			m.styles.key.Render("L") + " logs hint",
			m.styles.key.Render("r") + " refresh",
			m.styles.key.Render("[/]") + " view",
			m.styles.key.Render("?") + " help",
			m.styles.key.Render("q") + " quit",
		}
		return m.styles.foot.Render(strings.Join(keys, "  "))
	}

	keys := []string{
		m.styles.key.Render("j/k") + " move",
		m.styles.key.Render("tab") + " pane",
		m.styles.key.Render("enter") + " select",
		m.styles.key.Render("[/]") + " view",
		m.styles.key.Render("t") + " todo/all/done",
		m.styles.key.Render(".") + " due-only",
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
	if m.view == viewSystem {
		text := strings.Join([]string{
			"System view keybindings",
			"- s: sync (pull + push)",
			"- l: pull from remote",
			"- p: push to remote",
			"- o: edit sync/theme config",
			"- I/U: install or uninstall daemon service",
			"- S/T/R: start, stop, or restart daemon service",
			"- L: show daemon logs hint/path",
			"- r: refresh status",
			"- [ / ]: switch views",
			"- q: quit",
		}, "\n")
		return m.styles.help.Render(text)
	}

	text := strings.Join([]string{
		"TUI keybindings",
		"- tab / shift+tab: cycle focused pane",
		"- [ / ]: switch Tasks, Calendar, System",
		"- enter on NAV: apply selected scope",
		"- t: cycle todo -> all -> done",
		"- .: toggle due-only",
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

func defaultStateCounts() map[store.State]int64 {
	return map[store.State]int64{
		store.StateInbox:   0,
		store.StateNow:     0,
		store.StateWaiting: 0,
		store.StateLater:   0,
		store.StateDone:    0,
	}
}

func findNavItemIndex(items []navItem, kind navItemKind, value string) int {
	for i, item := range items {
		if item.kind == kind && item.value == value {
			return i
		}
	}
	return -1
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

func prevView(current viewKind) viewKind {
	if current == 0 {
		return viewCount - 1
	}
	return current - 1
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
