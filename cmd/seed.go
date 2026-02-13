package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

const (
	defaultSeedValue     = uint64(1)
	defaultSeedTaskCount = 200
	defaultSeedChurn     = 5
	seedMutateMaxTry     = 8

	maxProjectsPerTask = 2
	maxContextsPerTask = 2
	maxMetaPerTask     = 3

	percentBase                      = 100
	titleVersionMax                  = 10
	addTagBiasPercent                = 65
	initialDuePercent                = 35
	waitingHasAssigneePercent        = 80
	initialDueDayRange               = 60
	initialDueDayOffset              = 30
	mutationDueDayRange              = 90
	mutationDueDayOffset             = 45
	mutationTitleIndex               = 0
	mutationNotesIndex               = 1
	mutationStateIndex               = 2
	mutationDueIndex                 = 3
	mutationWaitingIndex             = 4
	mutationProjectsIndex            = 5
	mutationContextsIndex            = 6
	mutationMetaIndex                = 7
	mutationCount                    = 9
	stateCutoffInbox                 = 35
	stateCutoffNow                   = 55
	stateCutoffLater                 = 75
	stateCutoffWaiting               = 90
	seedSalt                  uint64 = 0x9e3779b97f4a7c15
)

type seedResult struct {
	DBPath string `json:"dbPath"`
	Seed   uint64 `json:"seed"`
	Count  int    `json:"count"`
	Churn  int    `json:"churn"`
}

type seedDataGenerator struct {
	rng         *rand.Rand
	anchorDate  time.Time
	titleVerbs  []string
	titleNouns  []string
	titleExtras []string
	notes       []string
	projects    []string
	contexts    []string
	waiters     []string
	metaKeys    []string
	metaValues  []string
	states      []string
}

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var seedCmd = &cli.Command{
	Name:     "seed",
	Usage:    "Seed a temp SQLite db with realistic tasks",
	Category: "System",
	Hidden:   true,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:  flags.FlagSeed,
			Usage: "seed for deterministic content",
			Value: defaultSeedValue,
		},
		&cli.IntFlag{
			Name:  flags.FlagCount,
			Usage: "number of tasks to create",
			Value: defaultSeedTaskCount,
		},
		&cli.IntFlag{
			Name:  flags.FlagChurn,
			Usage: "post-create mutations per task",
			Value: defaultSeedChurn,
		},
		&cli.StringFlag{
			Name:  flags.FlagOut,
			Usage: "output sqlite path (defaults to unique temp file)",
		},
		&cli.BoolFlag{
			Name:  flags.FlagForce,
			Usage: "overwrite --out path if it exists",
		},
	},
	Action: runSeed,
}

func runSeed(ctx context.Context, cmd *cli.Command) error {
	count := cmd.Int(flags.FlagCount)
	if count < 0 {
		return errors.New("count must be >= 0")
	}

	churn := cmd.Int(flags.FlagChurn)
	if churn < 0 {
		return errors.New("churn must be >= 0")
	}

	seed := cmd.Uint64(flags.FlagSeed)
	outPath, err := seedDBPath(cmd.String(flags.FlagOut), seed, cmd.Bool(flags.FlagForce))
	if err != nil {
		return err
	}

	st, err := store.Open(ctx, store.Options{Path: outPath})
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	svc := service.NewTaskService(st)
	gen := newSeedDataGenerator(seed)

	for i := range count {
		task, createErr := svc.CreateTask(ctx, gen.nextCreateRequest(i))
		if createErr != nil {
			return createErr
		}

		for range churn {
			task, err = applyMutation(ctx, svc, gen, task)
			if err != nil {
				return err
			}
		}
	}

	return writeSeedResult(outputWriter(), seedResult{DBPath: outPath, Seed: seed, Count: count, Churn: churn})
}

func writeSeedResult(writer output.Writer, result seedResult) error {
	if writer.JSON {
		enc := json.NewEncoder(writer.Out)
		return enc.Encode(result)
	}

	return writer.WriteKeyValues([]output.KeyValue{
		{Key: "db", Value: result.DBPath},
		{Key: "seed", Value: strconv.FormatUint(result.Seed, 10)},
		{Key: "count", Value: strconv.Itoa(result.Count)},
		{Key: "churn", Value: strconv.Itoa(result.Churn)},
	})
}

func seedDBPath(out string, seed uint64, force bool) (string, error) {
	if out == "" {
		file, err := os.CreateTemp("", fmt.Sprintf("ugh-seed-%d-*.sqlite", seed))
		if err != nil {
			return "", fmt.Errorf("create temp db: %w", err)
		}
		name := file.Name()
		if closeErr := file.Close(); closeErr != nil {
			return "", fmt.Errorf("close temp db: %w", closeErr)
		}
		return name, nil
	}

	path, err := config.ResolveDBPath("", out)
	if err != nil {
		return "", fmt.Errorf("resolve --out path: %w", err)
	}

	if err = os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return "", fmt.Errorf("create db dir: %w", err)
	}

	_, err = os.Stat(path)
	if err == nil {
		if !force {
			return "", fmt.Errorf("output db exists: %s (use --%s)", path, flags.FlagForce)
		}
		if removeErr := os.Remove(path); removeErr != nil {
			return "", fmt.Errorf("remove existing db: %w", removeErr)
		}
		return path, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("check output db: %w", err)
	}

	return path, nil
}

func newSeedDataGenerator(seed uint64) *seedDataGenerator {
	//nolint:gosec // Deterministic pseudo-random data is intentional for local seeding.
	rng := rand.New(rand.NewPCG(seed, seed^seedSalt))

	return &seedDataGenerator{
		rng:         rng,
		anchorDate:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		titleVerbs:  []string{"Draft", "Review", "Prepare", "Plan", "Refine", "Ship", "Audit", "Sketch"},
		titleNouns:  []string{"quarterly report", "roadmap", "team update", "onboarding flow", "bug backlog"},
		titleExtras: []string{"for Q1", "for team", "for launch", "for customer", "for demo"},
		notes: []string{
			"collect references and outline scope",
			"sync with stakeholders before finalizing",
			"capture open questions and risks",
			"validate assumptions with real usage",
			"",
		},
		projects:   []string{"work", "ops", "home", "finance", "product", "research"},
		contexts:   []string{"desk", "meeting", "email", "phone", "home", "office"},
		waiters:    []string{"Alice", "Bob", "Priya", "Jordan", "Casey", "Dana"},
		metaKeys:   []string{"prio", "est", "ticket", "source"},
		metaValues: []string{"low", "medium", "high", "S", "M", "L", "ops", "support"},
		states:     []string{flags.TaskStateInbox, flags.TaskStateNow, flags.TaskStateLater, flags.TaskStateWaiting},
	}
}

func (g *seedDataGenerator) nextCreateRequest(index int) service.CreateTaskRequest {
	state := g.initialState()
	title := g.pick(g.titleVerbs) + " " + g.pick(g.titleNouns) + " " + g.pick(g.titleExtras)
	title += " #" + strconv.Itoa(index+1)

	request := service.CreateTaskRequest{
		Title:    title,
		Notes:    g.pick(g.notes),
		State:    state,
		Projects: g.sampleUnique(g.projects, g.rng.IntN(maxProjectsPerTask+1)),
		Contexts: g.sampleUnique(g.contexts, g.rng.IntN(maxContextsPerTask+1)),
		Meta:     g.sampleMeta(),
	}

	if g.rng.IntN(percentBase) < initialDuePercent {
		offset := g.rng.IntN(initialDueDayRange) - initialDueDayOffset
		request.DueOn = g.anchorDate.AddDate(0, 0, offset).Format(flags.DateLayoutYYYYMMDD)
	}

	if state == flags.TaskStateWaiting && g.rng.IntN(percentBase) < waitingHasAssigneePercent {
		request.WaitingFor = g.pick(g.waiters)
	}

	return request
}

func (g *seedDataGenerator) initialState() string {
	roll := g.rng.IntN(percentBase)
	switch {
	case roll < stateCutoffInbox:
		return flags.TaskStateInbox
	case roll < stateCutoffNow:
		return flags.TaskStateNow
	case roll < stateCutoffLater:
		return flags.TaskStateLater
	case roll < stateCutoffWaiting:
		return flags.TaskStateWaiting
	default:
		return flags.TaskStateDone
	}
}

func (g *seedDataGenerator) sampleUnique(pool []string, count int) []string {
	if count <= 0 || len(pool) == 0 {
		return nil
	}
	if count > len(pool) {
		count = len(pool)
	}

	indexes := g.rng.Perm(len(pool))[:count]
	out := make([]string, 0, count)
	for _, index := range indexes {
		out = append(out, pool[index])
	}
	return out
}

func (g *seedDataGenerator) sampleMeta() []string {
	count := g.rng.IntN(maxMetaPerTask + 1)
	if count == 0 {
		return nil
	}

	keys := g.sampleUnique(g.metaKeys, count)
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, key+":"+g.pick(g.metaValues))
	}
	return out
}

func applyMutation(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, error) {
	if task == nil {
		return nil, errors.New("task is required")
	}

	for range seedMutateMaxTry {
		next, changed, err := tryMutation(ctx, svc, g, task)
		if err != nil {
			return nil, err
		}
		if changed {
			return next, nil
		}
	}

	fallback := task.Title + " update"
	return svc.UpdateTask(
		ctx,
		service.UpdateTaskRequest{ID: task.ID, Title: &fallback},
	)
}

func tryMutation(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	switch g.rng.IntN(mutationCount) {
	case mutationTitleIndex:
		return mutateTitle(ctx, svc, g, task)
	case mutationNotesIndex:
		return mutateNotes(ctx, svc, g, task)
	case mutationStateIndex:
		return mutateState(ctx, svc, g, task)
	case mutationDueIndex:
		return mutateDue(ctx, svc, g, task)
	case mutationWaitingIndex:
		return mutateWaiting(ctx, svc, g, task)
	case mutationProjectsIndex:
		return mutateProjects(ctx, svc, g, task)
	case mutationContextsIndex:
		return mutateContexts(ctx, svc, g, task)
	case mutationMetaIndex:
		return mutateMeta(ctx, svc, g, task)
	default:
		return toggleDone(ctx, svc, task)
	}
}

func mutateTitle(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	nextTitle := task.Title + " v" + strconv.Itoa(g.rng.IntN(titleVersionMax)+1)
	if nextTitle == task.Title {
		return task, false, nil
	}

	updated, err := svc.UpdateTask(
		ctx,
		service.UpdateTaskRequest{ID: task.ID, Title: &nextTitle},
	)
	return updated, true, err
}

func mutateNotes(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	nextNotes := g.pick(g.notes)
	if nextNotes == task.Notes {
		return task, false, nil
	}

	updated, err := svc.UpdateTask(
		ctx,
		service.UpdateTaskRequest{ID: task.ID, Notes: &nextNotes},
	)
	return updated, true, err
}

func mutateState(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	if task.State == store.StateDone {
		return toggleDone(ctx, svc, task)
	}

	nextState := g.pick(g.states)
	if string(task.State) == nextState {
		return task, false, nil
	}

	updated, err := svc.UpdateTask(
		ctx,
		service.UpdateTaskRequest{ID: task.ID, State: &nextState},
	)
	return updated, true, err
}

func mutateDue(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	if task.DueOn != nil {
		updated, err := svc.UpdateTask(
			ctx,
			service.UpdateTaskRequest{ID: task.ID, ClearDueOn: true},
		)
		return updated, true, err
	}

	offset := g.rng.IntN(mutationDueDayRange) - mutationDueDayOffset
	due := g.anchorDate.AddDate(0, 0, offset).Format(flags.DateLayoutYYYYMMDD)
	updated, err := svc.UpdateTask(
		ctx,
		service.UpdateTaskRequest{ID: task.ID, DueOn: &due},
	)
	return updated, true, err
}

func mutateWaiting(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	if task.WaitingFor != "" {
		updated, err := svc.UpdateTask(
			ctx,
			service.UpdateTaskRequest{ID: task.ID, ClearWaitingFor: true},
		)
		return updated, true, err
	}

	waiter := g.pick(g.waiters)
	updated, err := svc.UpdateTask(
		ctx,
		service.UpdateTaskRequest{ID: task.ID, WaitingFor: &waiter},
	)
	return updated, true, err
}

func mutateProjects(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	return mutateMembership(
		ctx,
		svc,
		g,
		task,
		task.Projects,
		g.projects,
		func(value string) service.UpdateTaskRequest {
			return service.UpdateTaskRequest{ID: task.ID, AddProjects: []string{value}}
		},
		func(value string) service.UpdateTaskRequest {
			return service.UpdateTaskRequest{ID: task.ID, RemoveProjects: []string{value}}
		},
	)
}

func mutateContexts(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	return mutateMembership(
		ctx,
		svc,
		g,
		task,
		task.Contexts,
		g.contexts,
		func(value string) service.UpdateTaskRequest {
			return service.UpdateTaskRequest{ID: task.ID, AddContexts: []string{value}}
		},
		func(value string) service.UpdateTaskRequest {
			return service.UpdateTaskRequest{ID: task.ID, RemoveContexts: []string{value}}
		},
	)
}

func mutateMeta(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
) (*store.Task, bool, error) {
	key := g.pick(g.metaKeys)
	value := g.pick(g.metaValues)

	if currentValue, ok := task.Meta[key]; !ok || currentValue != value {
		updated, err := svc.UpdateTask(
			ctx,
			service.UpdateTaskRequest{ID: task.ID, SetMeta: map[string]string{key: value}},
		)
		return updated, true, err
	}

	updated, err := svc.UpdateTask(
		ctx,
		service.UpdateTaskRequest{ID: task.ID, RemoveMetaKeys: []string{key}},
	)
	return updated, true, err
}

func mutateMembership(
	ctx context.Context,
	svc service.Service,
	g *seedDataGenerator,
	task *store.Task,
	current []string,
	pool []string,
	addRequest func(string) service.UpdateTaskRequest,
	removeRequest func(string) service.UpdateTaskRequest,
) (*store.Task, bool, error) {
	if len(current) < len(pool) && g.rng.IntN(percentBase) < addTagBiasPercent {
		value := g.pickMissing(pool, current)
		if value == "" {
			return task, false, nil
		}

		updated, err := svc.UpdateTask(ctx, addRequest(value))
		return updated, true, err
	}

	if len(current) == 0 {
		return task, false, nil
	}

	value := g.pick(current)
	updated, err := svc.UpdateTask(ctx, removeRequest(value))
	return updated, true, err
}

func toggleDone(
	ctx context.Context,
	svc service.Service,
	task *store.Task,
) (*store.Task, bool, error) {
	markDone := task.State != store.StateDone
	changed, err := svc.SetDone(ctx, []int64{task.ID}, markDone)
	if err != nil {
		return nil, false, err
	}
	if changed == 0 {
		return task, false, nil
	}

	updated, err := svc.GetTask(ctx, task.ID)
	if err != nil {
		return nil, false, err
	}
	return updated, true, nil
}

func (g *seedDataGenerator) pick(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[g.rng.IntN(len(values))]
}

func (g *seedDataGenerator) pickMissing(pool []string, existing []string) string {
	if len(pool) == 0 {
		return ""
	}

	existingSet := make(map[string]bool, len(existing))
	for _, value := range existing {
		existingSet[value] = true
	}

	candidates := make([]string, 0, len(pool))
	for _, value := range pool {
		if !existingSet[value] {
			candidates = append(candidates, value)
		}
	}

	if len(candidates) == 0 {
		return ""
	}

	return g.pick(candidates)
}
