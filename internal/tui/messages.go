package tui

import (
	"github.com/mholtzscher/ugh/internal/config"
	daemonservice "github.com/mholtzscher/ugh/internal/daemon/service"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

type tasksLoadedMsg struct {
	tasks []*store.Task
	err   error
}

type calendarTasksLoadedMsg struct {
	tasks []*store.Task
	err   error
}

type tagCountsLoadedMsg struct {
	projects []store.NameCount
	contexts []store.NameCount
	err      error
}

type stateCountsLoadedMsg struct {
	counts map[store.State]int64
	err    error
}

type syncStatusLoadedMsg struct {
	status *service.SyncStatus
	err    error
}

type daemonStatusLoadedMsg struct {
	managerName string
	status      *daemonservice.Status
	logPath     string
	err         error
}

type configSavedMsg struct {
	cfg  config.Config
	err  error
	path string
}

type actionResultMsg struct {
	status string
	err    error
}
