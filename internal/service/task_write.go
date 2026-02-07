package service

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
)

func (s *TaskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*store.Task, error) {
	meta, err := parseMetaFlags(req.Meta)
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	state, err := normalizeState(req.State)
	if err != nil {
		return nil, err
	}

	var dueOn *time.Time
	if strings.TrimSpace(req.DueOn) != "" {
		parsed, parseErr := parseDay(req.DueOn)
		if parseErr != nil {
			return nil, parseErr
		}
		dueOn = parsed
	}
	task := &store.Task{
		State:      state,
		Title:      req.Title,
		Notes:      req.Notes,
		DueOn:      dueOn,
		WaitingFor: strings.TrimSpace(req.WaitingFor),
		Projects:   req.Projects,
		Contexts:   req.Contexts,
		Meta:       meta,
	}

	return s.store.CreateTask(ctx, task)
}

func (s *TaskService) SetDone(ctx context.Context, ids []int64, done bool) (int64, error) {
	return s.store.SetDone(ctx, ids, done)
}

func (s *TaskService) DeleteTasks(ctx context.Context, ids []int64) (int64, error) {
	return s.store.DeleteTasks(ctx, ids)
}

//nolint:gocognit,nestif // UpdateTask applies many optional mutations in one place.
func (s *TaskService) UpdateTask(ctx context.Context, req UpdateTaskRequest) (*store.Task, error) {
	current, err := s.store.GetTask(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	updated := &store.Task{
		ID:          current.ID,
		State:       current.State,
		PrevState:   current.PrevState,
		Title:       current.Title,
		Notes:       current.Notes,
		DueOn:       current.DueOn,
		WaitingFor:  current.WaitingFor,
		CompletedAt: current.CompletedAt,
		Projects:    append([]string(nil), current.Projects...),
		Contexts:    append([]string(nil), current.Contexts...),
		Meta:        copyMeta(current.Meta),
	}

	if req.Title != nil {
		updated.Title = *req.Title
	}
	if req.Notes != nil {
		updated.Notes = *req.Notes
	}
	if req.State != nil {
		state, stateErr := normalizeState(*req.State)
		if stateErr != nil {
			return nil, stateErr
		}
		// Transition rules for done state.
		if state == store.StateDone && updated.State != store.StateDone {
			prev := updated.State
			updated.PrevState = &prev
			updated.CompletedAt = nil
		}
		if state != store.StateDone && updated.State == store.StateDone {
			updated.PrevState = nil
			updated.CompletedAt = nil
		}
		updated.State = state
	}
	// done is represented as state=done; completion toggles are handled by the done/undo commands.
	if req.ClearDueOn {
		updated.DueOn = nil
	} else if req.DueOn != nil {
		if strings.TrimSpace(*req.DueOn) == "" {
			updated.DueOn = nil
		} else {
			parsed, parseErr := parseDay(*req.DueOn)
			if parseErr != nil {
				return nil, parseErr
			}
			updated.DueOn = parsed
		}
	}
	if req.ClearWaitingFor {
		updated.WaitingFor = ""
	} else if req.WaitingFor != nil {
		updated.WaitingFor = strings.TrimSpace(*req.WaitingFor)
	}

	for _, p := range req.AddProjects {
		if !containsString(updated.Projects, p) {
			updated.Projects = append(updated.Projects, p)
		}
	}
	updated.Projects = removeStrings(updated.Projects, req.RemoveProjects)

	for _, c := range req.AddContexts {
		if !containsString(updated.Contexts, c) {
			updated.Contexts = append(updated.Contexts, c)
		}
	}
	updated.Contexts = removeStrings(updated.Contexts, req.RemoveContexts)

	if len(req.SetMeta) > 0 {
		if updated.Meta == nil {
			updated.Meta = map[string]string{}
		}
		maps.Copy(updated.Meta, req.SetMeta)
	}
	for _, k := range req.RemoveMetaKeys {
		delete(updated.Meta, k)
	}

	return s.store.UpdateTask(ctx, updated)
}

func (s *TaskService) FullUpdateTask(ctx context.Context, req FullUpdateTaskRequest) (*store.Task, error) {
	current, err := s.store.GetTask(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	state, err := normalizeState(req.State)
	if err != nil {
		return nil, err
	}
	var dueOn *time.Time
	if strings.TrimSpace(req.DueOn) != "" {
		parsed, parseErr := parseDay(req.DueOn)
		if parseErr != nil {
			return nil, parseErr
		}
		dueOn = parsed
	}
	updated := &store.Task{
		ID:          current.ID,
		State:       state,
		Title:       req.Title,
		Notes:       req.Notes,
		DueOn:       dueOn,
		WaitingFor:  strings.TrimSpace(req.WaitingFor),
		Projects:    req.Projects,
		Contexts:    req.Contexts,
		Meta:        req.Meta,
		CompletedAt: current.CompletedAt,
		PrevState:   current.PrevState,
	}
	if updated.State != store.StateDone {
		updated.CompletedAt = nil
	}
	if updated.State == store.StateDone && current.State != store.StateDone {
		prev := current.State
		updated.PrevState = &prev
	}
	if updated.State != store.StateDone {
		updated.PrevState = nil
	}

	return s.store.UpdateTask(ctx, updated)
}

func copyMeta(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	maps.Copy(result, m)
	return result
}

func containsString(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

func removeStrings(slice []string, toRemove []string) []string {
	if len(toRemove) == 0 {
		return slice
	}
	removeSet := make(map[string]bool, len(toRemove))
	for _, s := range toRemove {
		removeSet[s] = true
	}
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if !removeSet[s] {
			result = append(result, s)
		}
	}
	return result
}
