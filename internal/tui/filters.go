package tui

import (
	"strings"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

type listFilters struct {
	state   string
	project string
	context string
	search  string
}

func defaultFiltersWithState(state string) listFilters {
	return listFilters{state: state}
}

func (f listFilters) toListTasksRequest() service.ListTasksRequest {
	req := service.ListTasksRequest{
		State:   f.state,
		Project: f.project,
		Context: f.context,
		Search:  f.search,
	}

	switch f.state {
	case "":
		req.All = true
	case string(store.StateDone):
		req.DoneOnly = true
	default:
		req.TodoOnly = true
	}

	return req
}

func (f listFilters) toListTagsRequest() service.ListTagsRequest {
	var req service.ListTagsRequest

	switch f.state {
	case "":
		req.All = true
	case string(store.StateDone):
		req.DoneOnly = true
	default:
		req.TodoOnly = true
	}

	return req
}

func (f listFilters) statusText() string {
	parts := []string{}

	if f.project != "" {
		parts = append(parts, "project:"+f.project)
	}
	if f.context != "" {
		parts = append(parts, "context:"+f.context)
	}
	if f.search != "" {
		parts = append(parts, "search:"+f.search)
	}

	return strings.Join(parts, " ")
}
