package flags

import "github.com/mholtzscher/ugh/internal/domain"

const (
	FlagAll           = "all"
	FlagCompleted     = "completed"
	FlagConfigPath    = "config"
	FlagContext       = "context"
	FlagCounts        = "counts"
	FlagCreated       = "created"
	FlagDBPath        = "db"
	FlagDescription   = "description"
	FlagTitle         = "title"
	FlagDone          = "done"
	FlagEditor        = "editor"
	FlagJSON          = "json"
	FlagLines         = "lines"
	FlagMeta          = "meta"
	FlagNotes         = "notes"
	FlagNoColor       = "no-color"
	FlagNoTUI         = "no-tui"
	FlagNoFollow      = "no-follow"
	FlagNoDue         = "no-due"
	FlagNoWaitingFor  = "no-waiting-for"
	FlagProject       = "project"
	FlagRemoveContext = "remove-context"
	FlagRemoveMeta    = "remove-meta"
	FlagRemoveProject = "remove-project"
	FlagSearch        = "search"
	FlagState         = "state"
	FlagTodo          = "todo"
	FlagUndone        = "undone"
	FlagDueOn         = "due"
	FlagWaitingFor    = "waiting-for"
)

const (
	FieldState = "state"
	FieldDate  = "date"
	FieldMeta  = "meta"
)

const (
	TaskStateInbox   = domain.TaskStateInbox
	TaskStateNow     = domain.TaskStateNow
	TaskStateWaiting = domain.TaskStateWaiting
	TaskStateLater   = domain.TaskStateLater
	TaskStateDone    = domain.TaskStateDone

	TaskStatesUsage = domain.TaskStatesUsage
)

const (
	DateLayoutYYYYMMDD = domain.DateLayoutYYYYMMDD
	DateTextYYYYMMDD   = domain.DateTextYYYYMMDD

	MetaSeparatorColon = domain.MetaSeparatorColon
	MetaTextKeyValue   = domain.MetaTextKeyValue
)

func TaskStates() []string {
	return []string{TaskStateInbox, TaskStateNow, TaskStateWaiting, TaskStateLater, TaskStateDone}
}
