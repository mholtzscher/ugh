package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

func orderedTaskStates() []store.State {
	return []store.State{store.StateInbox, store.StateNow, store.StateWaiting, store.StateLater, store.StateDone}
}

func loadTasksCmd(svc service.Service, filters listFilters) tea.Cmd {
	req := filters.toListTasksRequest()
	return func() tea.Msg {
		tasks, err := svc.ListTasks(context.Background(), req)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func loadTagCountsCmd(svc service.Service, filters listFilters) tea.Cmd {
	tagReq := filters.toListTagsRequest()
	return func() tea.Msg {
		projects, err := svc.ListProjects(context.Background(), tagReq)
		if err != nil {
			return tagCountsLoadedMsg{err: err}
		}
		contexts, err := svc.ListContexts(context.Background(), tagReq)
		if err != nil {
			return tagCountsLoadedMsg{err: err}
		}
		return tagCountsLoadedMsg{projects: projects, contexts: contexts}
	}
}

func loadStateCountsCmd(svc service.Service, filters listFilters) tea.Cmd {
	req := filters.toListTasksRequest()
	req.State = ""
	req.Project = ""
	req.Context = ""
	req.Search = ""

	return func() tea.Msg {
		states := orderedTaskStates()
		counts := make(map[store.State]int64, len(states))
		for _, state := range states {
			stateReq := req
			stateReq.State = string(state)
			tasks, err := svc.ListTasks(context.Background(), stateReq)
			if err != nil {
				return stateCountsLoadedMsg{err: err}
			}
			counts[state] = int64(len(tasks))
		}
		return stateCountsLoadedMsg{counts: counts}
	}
}

func refreshDataCmd(svc service.Service, filters listFilters) tea.Cmd {
	return tea.Batch(
		loadTasksCmd(svc, filters),
		loadTagCountsCmd(svc, filters),
		loadStateCountsCmd(svc, filters),
	)
}

func setDoneCmd(svc service.Service, id int64, done bool) tea.Cmd {
	return func() tea.Msg {
		count, err := svc.SetDone(context.Background(), []int64{id}, done)
		if err != nil {
			return actionResultMsg{err: err}
		}
		if count == 0 {
			return actionResultMsg{err: fmt.Errorf("task %d not found", id)}
		}
		if done {
			return actionResultMsg{status: fmt.Sprintf("task #%d marked done", id)}
		}
		return actionResultMsg{status: fmt.Sprintf("task #%d reopened", id)}
	}
}

func setStateCmd(svc service.Service, id int64, state string) tea.Cmd {
	return func() tea.Msg {
		_, err := svc.UpdateTask(context.Background(), service.UpdateTaskRequest{ID: id, State: &state})
		if err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: fmt.Sprintf("task #%d moved to %s", id, state)}
	}
}

func deleteTaskCmd(svc service.Service, id int64) tea.Cmd {
	return func() tea.Msg {
		count, err := svc.DeleteTasks(context.Background(), []int64{id})
		if err != nil {
			return actionResultMsg{err: err}
		}
		if count == 0 {
			return actionResultMsg{err: fmt.Errorf("task %d not found", id)}
		}
		return actionResultMsg{status: fmt.Sprintf("task #%d deleted", id)}
	}
}

func createTaskCmd(svc service.Service, req service.CreateTaskRequest) tea.Cmd {
	return func() tea.Msg {
		task, err := svc.CreateTask(context.Background(), req)
		if err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: fmt.Sprintf("created task #%d", task.ID)}
	}
}

func fullUpdateTaskCmd(svc service.Service, req service.FullUpdateTaskRequest) tea.Cmd {
	return func() tea.Msg {
		task, err := svc.FullUpdateTask(context.Background(), req)
		if err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: fmt.Sprintf("updated task #%d", task.ID)}
	}
}

func updateTaskCmd(svc service.Service, req service.UpdateTaskRequest) tea.Cmd {
	return func() tea.Msg {
		task, err := svc.UpdateTask(context.Background(), req)
		if err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: fmt.Sprintf("updated task #%d", task.ID)}
	}
}
