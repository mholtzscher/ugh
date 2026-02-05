//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate -f ../../sqlc.yaml

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	tursogo "turso.tech/database/tursogo"

	"github.com/mholtzscher/ugh/internal/store/sqlc"
	"github.com/pressly/goose/v3"
)

type Store struct {
	db      *sql.DB
	syncDb  *tursogo.TursoSyncDb
	queries *sqlc.Queries
}

type Options struct {
	Path        string
	SyncURL     string
	AuthToken   string
	BusyTimeout int // Milliseconds to wait for locks (default: 5000)
}

func Open(ctx context.Context, opts Options) (*Store, error) {
	if opts.Path == "" {
		return nil, errors.New("db path is required")
	}

	abspath, err := filepath.Abs(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("resolve db path: %w", err)
	}

	var db *sql.DB
	var syncDb *tursogo.TursoSyncDb

	if opts.SyncURL != "" {
		authToken := opts.AuthToken
		if authToken == "" {
			if envToken := os.Getenv("LIBSQL_AUTH_TOKEN"); envToken != "" {
				authToken = envToken
			}
		}
		if authToken == "" {
			return nil, errors.New("auth token required when sync_url is set (use db.auth_token in config or LIBSQL_AUTH_TOKEN env var)")
		}

		trueVal := true
		cfg := tursogo.TursoSyncDbConfig{
			Path:             abspath,
			RemoteUrl:        opts.SyncURL,
			AuthToken:        authToken,
			BootstrapIfEmpty: &trueVal,
		}
		sdb, err := tursogo.NewTursoSyncDb(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("create sync db: %w", err)
		}
		db, err = sdb.Connect(ctx)
		if err != nil {
			return nil, fmt.Errorf("connect sync db: %w", err)
		}
		syncDb = sdb
	} else {
		var err error
		db, err = sql.Open("turso", abspath)
		if err != nil {
			return nil, fmt.Errorf("open db: %w", err)
		}
	}

	// Set busy_timeout FIRST so subsequent pragmas wait if database is locked
	busyTimeout := opts.BusyTimeout
	if busyTimeout <= 0 {
		busyTimeout = 5000 // Default 5 seconds
	}
	if _, err := db.ExecContext(ctx, fmt.Sprintf("PRAGMA busy_timeout=%d;", busyTimeout)); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply pragma: %w", err)
	}
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply pragma: %w", err)
	}
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys=ON;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply pragma: %w", err)
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	goose.SetLogger(log.New(io.Discard, "", 0))
	if err := goose.Up(db, "migrations"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	store := &Store{db: db, syncDb: syncDb, queries: sqlc.New(db)}
	return store, nil
}

func (s *Store) Sync(ctx context.Context) error {
	if s.syncDb == nil {
		return errors.New("sync is not configured")
	}
	_, err := s.syncDb.Pull(ctx)
	return err
}

func (s *Store) Push(ctx context.Context) error {
	if s.syncDb == nil {
		return errors.New("sync is not configured")
	}
	return s.syncDb.Push(ctx)
}

func (s *Store) SyncStats(ctx context.Context) (*tursogo.TursoSyncDbStats, error) {
	if s.syncDb == nil {
		return nil, errors.New("sync is not configured")
	}
	stats, err := s.syncDb.Stats(ctx)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	if task == nil {
		return nil, errors.New("task is required")
	}
	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		return nil, errors.New("title is required")
	}
	if task.State == "" {
		task.State = StateInbox
	}
	if task.Meta == nil {
		task.Meta = map[string]string{}
	}

	now := time.Now().UTC()
	createdAt := now.Unix()
	updatedAt := now.Unix()

	completedAt := task.CompletedAt
	if task.State == StateDone {
		if completedAt == nil {
			completedAt = &now
		}
	} else {
		completedAt = nil
	}

	prevStateNull := sql.NullString{}
	if task.PrevState != nil && *task.PrevState != "" {
		prevStateNull = sql.NullString{String: string(*task.PrevState), Valid: true}
	}

	params := sqlc.InsertTaskParams{
		State:       string(task.State),
		PrevState:   prevStateNull,
		Title:       task.Title,
		Notes:       task.Notes,
		DueOn:       nullDate(task.DueOn),
		WaitingFor:  nullString(task.WaitingFor),
		CompletedAt: nullUnixTime(completedAt),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	row, err := s.queries.InsertTask(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("insert task: %w", err)
	}

	if err := s.insertDetails(ctx, row.ID, task); err != nil {
		return nil, err
	}

	return s.GetTask(ctx, row.ID)
}

func (s *Store) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	if task == nil {
		return nil, errors.New("task is required")
	}
	if task.ID == 0 {
		return nil, errors.New("task id is required")
	}
	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		return nil, errors.New("title is required")
	}
	if task.State == "" {
		task.State = StateInbox
	}
	if task.Meta == nil {
		task.Meta = map[string]string{}
	}

	now := time.Now().UTC()
	updatedAt := now.Unix()
	completedAt := task.CompletedAt
	if task.State == StateDone {
		if completedAt == nil {
			completedAt = &now
		}
	} else {
		completedAt = nil
	}

	prevStateNull := sql.NullString{}
	if task.PrevState != nil && *task.PrevState != "" {
		prevStateNull = sql.NullString{String: string(*task.PrevState), Valid: true}
	}

	params := sqlc.UpdateTaskParams{
		State:       string(task.State),
		PrevState:   prevStateNull,
		Title:       task.Title,
		Notes:       task.Notes,
		DueOn:       nullDate(task.DueOn),
		WaitingFor:  nullString(task.WaitingFor),
		CompletedAt: nullUnixTime(completedAt),
		UpdatedAt:   updatedAt,
		ID:          task.ID,
	}
	if _, err := s.queries.UpdateTask(ctx, params); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	if err := s.deleteDetails(ctx, task.ID); err != nil {
		return nil, err
	}
	if err := s.insertDetails(ctx, task.ID, task); err != nil {
		return nil, err
	}

	return s.GetTask(ctx, task.ID)
}

func (s *Store) GetTask(ctx context.Context, id int64) (*Task, error) {
	row, err := s.queries.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	task := fromGetRow(row)
	if err := s.loadDetails(ctx, []*Task{task}); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Store) ListTasks(ctx context.Context, filters Filters) ([]*Task, error) {
	state := strings.TrimSpace(filters.State)
	if filters.DoneOnly {
		state = string(StateDone)
	}

	excludeDone := int64(0)
	if filters.TodoOnly {
		excludeDone = 1
	}

	// Build search filter as sql.NullString
	searchNull := sql.NullString{}
	if filters.Search != "" {
		searchNull = sql.NullString{String: filters.Search, Valid: true}
	}

	dueSet := int64(0)
	if filters.DueSetOnly {
		dueSet = 1
	}

	// For (? IS NULL OR field = ?) pattern, pass the same value twice
	// When first param is nil, second doesn't matter (OR short-circuits)
	params := sqlc.ListTasksParams{
		Column1: excludeDone,
		Column2: nullAny(state),
		State:   state,

		Column4: nullAny(filters.Project),
		Name:    filters.Project,
		Column6: nullAny(filters.Context),
		Name_2:  filters.Context,

		Column8:  nullAny(filters.Search),
		Column9:  searchNull,
		Column10: searchNull,
		Column11: searchNull,
		Column12: searchNull,
		Column13: searchNull,
		Column14: searchNull,

		Column15: dueSet,
	}
	rows, err := s.queries.ListTasks(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	tasks := make([]*Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, fromListRow(row))
	}
	if err := s.loadDetails(ctx, tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Store) SetDone(ctx context.Context, ids []int64, done bool) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	updatedAt := time.Now().UTC().Unix()
	if done {
		params := sqlc.CompleteTasksParams{
			CompletedAt: sql.NullInt64{Int64: time.Now().UTC().Unix(), Valid: true},
			UpdatedAt:   updatedAt,
			Ids:         ids,
		}
		count, err := s.queries.CompleteTasks(ctx, params)
		if err != nil {
			return 0, fmt.Errorf("complete tasks: %w", err)
		}
		return count, nil
	}

	params := sqlc.ReopenTasksParams{UpdatedAt: updatedAt, Ids: ids}
	count, err := s.queries.ReopenTasks(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("reopen tasks: %w", err)
	}
	return count, nil
}

func (s *Store) DeleteTasks(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := s.queries.DeleteTasks(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("delete tasks: %w", err)
	}
	return count, nil
}

func (s *Store) loadDetails(ctx context.Context, tasks []*Task) error {
	if len(tasks) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(tasks))
	byID := make(map[int64]*Task, len(tasks))
	for _, task := range tasks {
		ids = append(ids, task.ID)
		byID[task.ID] = task
		task.Projects = nil
		task.Contexts = nil
		task.Meta = map[string]string{}
	}

	projects, err := s.queries.ListProjectsForTasks(ctx, ids)
	if err != nil {
		return fmt.Errorf("load projects: %w", err)
	}
	for _, row := range projects {
		if task := byID[row.TaskID]; task != nil {
			task.Projects = append(task.Projects, row.Name)
		}
	}

	contexts, err := s.queries.ListContextsForTasks(ctx, ids)
	if err != nil {
		return fmt.Errorf("load contexts: %w", err)
	}
	for _, row := range contexts {
		if task := byID[row.TaskID]; task != nil {
			task.Contexts = append(task.Contexts, row.Name)
		}
	}

	meta, err := s.queries.ListMeta(ctx, ids)
	if err != nil {
		return fmt.Errorf("load meta: %w", err)
	}
	for _, row := range meta {
		if task := byID[row.TaskID]; task != nil {
			task.Meta[row.Key] = row.Value
		}
	}

	return nil
}

func (s *Store) insertDetails(ctx context.Context, taskID int64, task *Task) error {
	now := time.Now().UTC().Unix()

	for _, project := range uniqueStrings(cleanNames(task.Projects)) {
		projectID, err := s.queries.EnsureProject(ctx, sqlc.EnsureProjectParams{Name: project, CreatedAt: now, UpdatedAt: now})
		if err != nil {
			return fmt.Errorf("ensure project %q: %w", project, err)
		}
		if err := s.queries.InsertTaskProjectLink(ctx, sqlc.InsertTaskProjectLinkParams{TaskID: taskID, ProjectID: projectID}); err != nil {
			return fmt.Errorf("insert project link: %w", err)
		}
	}

	for _, context := range uniqueStrings(cleanNames(task.Contexts)) {
		contextID, err := s.queries.EnsureContext(ctx, sqlc.EnsureContextParams{Name: context, CreatedAt: now, UpdatedAt: now})
		if err != nil {
			return fmt.Errorf("ensure context %q: %w", context, err)
		}
		if err := s.queries.InsertTaskContextLink(ctx, sqlc.InsertTaskContextLinkParams{TaskID: taskID, ContextID: contextID}); err != nil {
			return fmt.Errorf("insert context link: %w", err)
		}
	}

	for key, value := range task.Meta {
		if err := s.queries.InsertMeta(ctx, sqlc.InsertMetaParams{TaskID: taskID, Key: key, Value: value}); err != nil {
			return fmt.Errorf("insert meta: %w", err)
		}
	}

	return nil
}

func (s *Store) ListProjectCounts(ctx context.Context, onlyDone bool, excludeDone bool) ([]NameCount, error) {
	params := sqlc.ListProjectCountsParams{
		Column1: boolToInt(onlyDone),
		Column2: boolToInt(excludeDone),
	}
	rows, err := s.queries.ListProjectCounts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list project counts: %w", err)
	}
	result := make([]NameCount, 0, len(rows))
	for _, row := range rows {
		result = append(result, NameCount{Name: row.Name, Count: row.Count})
	}
	return result, nil
}

func (s *Store) ListContextCounts(ctx context.Context, onlyDone bool, excludeDone bool) ([]NameCount, error) {
	params := sqlc.ListContextCountsParams{
		Column1: boolToInt(onlyDone),
		Column2: boolToInt(excludeDone),
	}
	rows, err := s.queries.ListContextCounts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list context counts: %w", err)
	}
	result := make([]NameCount, 0, len(rows))
	for _, row := range rows {
		result = append(result, NameCount{Name: row.Name, Count: row.Count})
	}
	return result, nil
}

func (s *Store) deleteDetails(ctx context.Context, taskID int64) error {
	if err := s.queries.DeleteTaskProjectLinks(ctx, taskID); err != nil {
		return fmt.Errorf("delete project links: %w", err)
	}
	if err := s.queries.DeleteTaskContextLinks(ctx, taskID); err != nil {
		return fmt.Errorf("delete context links: %w", err)
	}
	if err := s.queries.DeleteMeta(ctx, taskID); err != nil {
		return fmt.Errorf("delete meta: %w", err)
	}
	return nil
}

func fromGetRow(row sqlc.GetTaskRow) *Task {
	return &Task{
		ID:          row.ID,
		State:       State(row.State),
		PrevState:   parseStatePtr(row.PrevState),
		Title:       row.Title,
		Notes:       row.Notes,
		DueOn:       parseDate(row.DueOn),
		WaitingFor:  row.WaitingFor.String,
		CompletedAt: parseUnixTime(row.CompletedAt),
		CreatedAt:   time.Unix(row.CreatedAt, 0).UTC(),
		UpdatedAt:   time.Unix(row.UpdatedAt, 0).UTC(),
	}
}

func fromListRow(row sqlc.ListTasksRow) *Task {
	return &Task{
		ID:          row.ID,
		State:       State(row.State),
		PrevState:   parseStatePtr(row.PrevState),
		Title:       row.Title,
		Notes:       row.Notes,
		DueOn:       parseDate(row.DueOn),
		WaitingFor:  row.WaitingFor.String,
		CompletedAt: parseUnixTime(row.CompletedAt),
		CreatedAt:   time.Unix(row.CreatedAt, 0).UTC(),
		UpdatedAt:   time.Unix(row.UpdatedAt, 0).UTC(),
	}
}

func parseStatePtr(val sql.NullString) *State {
	if !val.Valid || val.String == "" {
		return nil
	}
	s := State(val.String)
	return &s
}

func parseDate(val sql.NullString) *time.Time {
	if !val.Valid || val.String == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", val.String)
	if err != nil {
		return nil
	}
	utc := parsed.UTC()
	return &utc
}

func parseUnixTime(val sql.NullInt64) *time.Time {
	if !val.Valid {
		return nil
	}
	t := time.Unix(val.Int64, 0).UTC()
	return &t
}

func boolToInt(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullDate(value *time.Time) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: value.Format("2006-01-02"), Valid: true}
}

func nullUnixTime(value *time.Time) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: value.UTC().Unix(), Valid: true}
}

func nullAny(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func cleanNames(values []string) []string {
	if len(values) == 0 {
		return values
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.TrimPrefix(value, "+")
		value = strings.TrimPrefix(value, "@")
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}
