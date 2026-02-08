package tui

import "github.com/mholtzscher/ugh/internal/store"

type tasksLoadedMsg struct {
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

type actionResultMsg struct {
	status string
	err    error
}
