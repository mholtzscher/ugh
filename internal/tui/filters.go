package tui

import (
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/service"
)

type completionFilter int

const (
	completionTodo completionFilter = iota
	completionAll
	completionDone
)

type listFilters struct {
	completion completionFilter
	state      string
	project    string
	context    string
	search     string
}

func defaultFiltersWithState(state string) listFilters {
	return listFilters{completion: completionTodo, state: state}
}

func (f listFilters) toListTasksRequest() service.ListTasksRequest {
	req := service.ListTasksRequest{
		State:   f.state,
		Project: f.project,
		Context: f.context,
		Search:  f.search,
	}

	switch f.completion {
	case completionAll:
		req.All = true
	case completionDone:
		req.DoneOnly = true
	case completionTodo:
		req.TodoOnly = true
	}

	return req
}

func (f listFilters) toListTagsRequest() service.ListTagsRequest {
	req := service.ListTagsRequest{}

	switch f.completion {
	case completionAll:
		req.All = true
	case completionDone:
		req.DoneOnly = true
	case completionTodo:
		req.TodoOnly = true
	}

	return req
}

func (f listFilters) completionText() string {
	switch f.completion {
	case completionAll:
		return "all"
	case completionDone:
		return "done"
	case completionTodo:
		return "todo"
	default:
		return "todo"
	}
}

func (f listFilters) cycleCompletion() listFilters {
	next := f
	switch f.completion {
	case completionTodo:
		next.completion = completionAll
	case completionAll:
		next.completion = completionDone
	case completionDone:
		next.completion = completionTodo
	default:
		next.completion = completionTodo
	}
	return next
}

func (f listFilters) statusText() string {
	parts := []string{fmt.Sprintf("mode:%s", f.completionText())}

	if f.state != "" {
		parts = append(parts, "state:"+f.state)
	}
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
