//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate -f ../../sqlc.yaml

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	tursogo "turso.tech/database/tursogo"

	"github.com/pressly/goose/v3"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/store/sqlc"
)

type Store struct {
	db      *sql.DB
	syncDB  *tursogo.TursoSyncDb
	queries *sqlc.Queries
}

type Options struct {
	Path        string
	SyncURL     string
	AuthToken   string
	BusyTimeout int // Milliseconds to wait for locks (default: 5000)
}

//nolint:gocognit,nestif,funlen // Store initialization handles sync, pragmas, and migrations in one flow.
func Open(ctx context.Context, opts Options) (*Store, error) {
	if opts.Path == "" {
		return nil, errors.New("db path is required")
	}

	abspath, err := filepath.Abs(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("resolve db path: %w", err)
	}

	var db *sql.DB
	var syncDB *tursogo.TursoSyncDb

	if opts.SyncURL != "" {
		authToken := opts.AuthToken
		if authToken == "" {
			if envToken := os.Getenv("LIBSQL_AUTH_TOKEN"); envToken != "" {
				authToken = envToken
			}
		}
		if authToken == "" {
			return nil, errors.New(
				"auth token required when sync_url is set (use db.auth_token in config or LIBSQL_AUTH_TOKEN env var)",
			)
		}

		trueVal := true
		cfg := tursogo.TursoSyncDbConfig{
			Path:             abspath,
			RemoteUrl:        opts.SyncURL,
			AuthToken:        authToken,
			BootstrapIfEmpty: &trueVal,
		}
		var sdb *tursogo.TursoSyncDb
		sdb, err = tursogo.NewTursoSyncDb(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("create sync db: %w", err)
		}
		db, err = sdb.Connect(ctx)
		if err != nil {
			return nil, fmt.Errorf("connect sync db: %w", err)
		}
		syncDB = sdb
	} else {
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
	_, err = db.ExecContext(ctx, fmt.Sprintf("PRAGMA busy_timeout=%d;", busyTimeout))
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply pragma: %w", err)
	}
	_, err = db.ExecContext(ctx, "PRAGMA journal_mode=WAL;")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply pragma: %w", err)
	}
	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys=ON;")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply pragma: %w", err)
	}

	goose.SetBaseFS(migrationsFS)
	err = goose.SetDialect("sqlite3")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	goose.SetLogger(goose.NopLogger())
	err = goose.Up(db, "migrations")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	store := &Store{db: db, syncDB: syncDB, queries: sqlc.New(db)}
	return store, nil
}

func (s *Store) Sync(ctx context.Context) error {
	if s.syncDB == nil {
		return errors.New("sync is not configured")
	}
	_, err := s.syncDB.Pull(ctx)
	return err
}

func (s *Store) Push(ctx context.Context) error {
	if s.syncDB == nil {
		return errors.New("sync is not configured")
	}
	return s.syncDB.Push(ctx)
}

func (s *Store) SyncStats(ctx context.Context) (*tursogo.TursoSyncDbStats, error) {
	if s.syncDB == nil {
		return nil, errors.New("sync is not configured")
	}
	stats, err := s.syncDB.Stats(ctx)
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

	err = s.insertDetails(ctx, row.ID, task)
	if err != nil {
		return nil, err
	}

	createdTask, err := s.GetTask(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	err = s.recordTaskEvent(
		ctx,
		createdTask.ID,
		TaskEventKindCreate,
		fmt.Sprintf("created task #%d", createdTask.ID),
		buildCreateChanges(createdTask),
	)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (s *Store) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	if task == nil {
		return nil, errors.New("task is required")
	}
	if task.ID == 0 {
		return nil, errors.New("task id is required")
	}
	beforeTask, err := s.GetTask(ctx, task.ID)
	if err != nil {
		return nil, err
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
	if _, err = s.queries.UpdateTask(ctx, params); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	err = s.deleteDetails(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	err = s.insertDetails(ctx, task.ID, task)
	if err != nil {
		return nil, err
	}

	updatedTask, err := s.GetTask(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	changes := buildTaskDiff(beforeTask, updatedTask)
	if len(changes) > 0 {
		err = s.recordTaskEvent(
			ctx,
			updatedTask.ID,
			TaskEventKindUpdate,
			fmt.Sprintf("updated task #%d", updatedTask.ID),
			changes,
		)
		if err != nil {
			return nil, err
		}
	}

	return updatedTask, nil
}

func (s *Store) GetTask(ctx context.Context, id int64) (*Task, error) {
	row, err := s.queries.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	task := fromGetRow(row)
	err = s.loadDetails(ctx, []*Task{task})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Store) ListTasksByExpr(
	ctx context.Context,
	expr nlp.FilterExpr,
	opts ListTasksByExprOptions,
) ([]*Task, error) {
	conditions := make([]sq.Sqlizer, 0)
	if opts.OnlyDone {
		conditions = append(conditions, sq.Expr("t.state = 'done'"))
	} else if opts.ExcludeDone {
		conditions = append(conditions, sq.Expr("t.state != 'done'"))
	}

	if expr != nil {
		builder := &filterSQLBuilder{}
		exprClause, exprArgs, err := builder.Build(expr)
		if err != nil {
			return nil, fmt.Errorf("build filter SQL: %w", err)
		}
		conditions = append(conditions, sq.Expr(exprClause, exprArgs...))
	}

	queryBuilder := sq.Select(
		"t.id",
		"t.state",
		"t.prev_state",
		"CAST(t.title AS TEXT) AS title",
		"CAST(t.notes AS TEXT) AS notes",
		"t.due_on",
		"t.waiting_for",
		"t.completed_at",
		"t.created_at",
		"t.updated_at",
	).
		From("tasks t").
		OrderBy(
			"CASE WHEN t.state = 'done' THEN 1 ELSE 0 END",
			"CASE WHEN t.due_on IS NULL OR t.due_on = '' THEN 1 ELSE 0 END",
			"t.due_on ASC",
			"t.updated_at DESC",
		)

	for _, condition := range conditions {
		queryBuilder = queryBuilder.Where(condition)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list tasks query: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks by expression: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var row listTaskRow
		if scanErr := rows.Scan(
			&row.ID,
			&row.State,
			&row.PrevState,
			&row.Title,
			&row.Notes,
			&row.DueOn,
			&row.WaitingFor,
			&row.CompletedAt,
			&row.CreatedAt,
			&row.UpdatedAt,
		); scanErr != nil {
			return nil, fmt.Errorf("scan task row: %w", scanErr)
		}
		tasks = append(tasks, fromListRow(row))
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("iterate task rows: %w", rowsErr)
	}

	err = s.loadDetails(ctx, tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

type listTaskRow struct {
	ID          int64
	State       string
	PrevState   sql.NullString
	Title       string
	Notes       string
	DueOn       sql.NullString
	WaitingFor  sql.NullString
	CompletedAt sql.NullInt64
	CreatedAt   int64
	UpdatedAt   int64
}

func (s *Store) SetDone(ctx context.Context, ids []int64, done bool) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	beforeTasks, err := s.getExistingTasksByID(ctx, ids)
	if err != nil {
		return 0, err
	}

	updatedAt := time.Now().UTC().Unix()
	var count int64
	if done {
		params := sqlc.CompleteTasksParams{
			CompletedAt: sql.NullInt64{Int64: time.Now().UTC().Unix(), Valid: true},
			UpdatedAt:   updatedAt,
			Ids:         ids,
		}
		count, err = s.queries.CompleteTasks(ctx, params)
		if err != nil {
			return 0, fmt.Errorf("complete tasks: %w", err)
		}
	} else {
		params := sqlc.ReopenTasksParams{UpdatedAt: updatedAt, Ids: ids}
		count, err = s.queries.ReopenTasks(ctx, params)
		if err != nil {
			return 0, fmt.Errorf("reopen tasks: %w", err)
		}
	}

	afterTasks, err := s.getExistingTasksByID(ctx, ids)
	if err != nil {
		return 0, err
	}

	kind := TaskEventKindUndo
	verb := "reopened"
	if done {
		kind = TaskEventKindDone
		verb = "marked done"
	}

	for _, id := range ids {
		beforeTask := beforeTasks[id]
		afterTask := afterTasks[id]
		if beforeTask == nil || afterTask == nil {
			continue
		}
		changes := buildTaskDiff(beforeTask, afterTask)
		if len(changes) == 0 {
			continue
		}
		err = s.recordTaskEvent(
			ctx,
			id,
			kind,
			fmt.Sprintf("%s task #%d", verb, id),
			changes,
		)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (s *Store) DeleteTasks(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	beforeTasks, err := s.getExistingTasksByID(ctx, ids)
	if err != nil {
		return 0, err
	}

	count, err := s.queries.DeleteTasks(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("delete tasks: %w", err)
	}

	for _, id := range ids {
		beforeTask := beforeTasks[id]
		if beforeTask == nil {
			continue
		}
		err = s.recordTaskEvent(
			ctx,
			id,
			TaskEventKindDelete,
			fmt.Sprintf("deleted task #%d", id),
			buildDeleteChanges(beforeTask),
		)
		if err != nil {
			return 0, err
		}
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
		projectID, err := s.queries.EnsureProject(
			ctx,
			sqlc.EnsureProjectParams{Name: project, CreatedAt: now, UpdatedAt: now},
		)
		if err != nil {
			return fmt.Errorf("ensure project %q: %w", project, err)
		}
		err = s.queries.InsertTaskProjectLink(
			ctx,
			sqlc.InsertTaskProjectLinkParams{TaskID: taskID, ProjectID: projectID},
		)
		if err != nil {
			return fmt.Errorf("insert project link: %w", err)
		}
	}

	for _, context := range uniqueStrings(cleanNames(task.Contexts)) {
		contextID, err := s.queries.EnsureContext(
			ctx,
			sqlc.EnsureContextParams{Name: context, CreatedAt: now, UpdatedAt: now},
		)
		if err != nil {
			return fmt.Errorf("ensure context %q: %w", context, err)
		}
		err = s.queries.InsertTaskContextLink(
			ctx,
			sqlc.InsertTaskContextLinkParams{TaskID: taskID, ContextID: contextID},
		)
		if err != nil {
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

// RecordShellHistory records a command in shell history.
func (s *Store) RecordShellHistory(
	ctx context.Context, command string, success bool, summary string, intent string,
) (*ShellHistory, error) {
	params := sqlc.InsertShellHistoryParams{
		Timestamp:     time.Now().UTC().Unix(),
		Command:       command,
		Success:       success,
		ResultSummary: nullString(summary),
		Intent:        nullString(intent),
	}
	row, err := s.queries.InsertShellHistory(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("insert shell history: %w", err)
	}
	return &ShellHistory{
		ID:            row.ID,
		Timestamp:     row.Timestamp,
		Command:       row.Command,
		Success:       row.Success,
		ResultSummary: row.ResultSummary.String,
		Intent:        row.Intent.String,
	}, nil
}

// ListShellHistory returns recent shell history entries.
func (s *Store) ListShellHistory(ctx context.Context, limit int64) ([]*ShellHistory, error) {
	rows, err := s.queries.ListShellHistory(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list shell history: %w", err)
	}
	result := make([]*ShellHistory, 0, len(rows))
	for _, row := range rows {
		result = append(result, &ShellHistory{
			ID:            row.ID,
			Timestamp:     row.Timestamp,
			Command:       row.Command,
			Success:       row.Success,
			ResultSummary: row.ResultSummary.String,
			Intent:        row.Intent.String,
		})
	}
	return result, nil
}

// SearchShellHistory searches shell history with filters.
func (s *Store) SearchShellHistory(
	ctx context.Context, search, intent string, success *bool, limit int64,
) ([]*ShellHistory, error) {
	var searchNull sql.NullString
	if search != "" {
		searchNull = sql.NullString{String: search, Valid: true}
	}
	var intentNull sql.NullString
	if intent != "" {
		intentNull = sql.NullString{String: intent, Valid: true}
	}
	var successBool bool
	var successAny any
	if success != nil {
		successBool = *success
		successAny = *success
	}
	params := sqlc.SearchShellHistoryParams{
		Column1: searchNull,
		Column2: searchNull,
		Column3: intentNull,
		Intent:  intentNull,
		Column5: successAny,
		Success: successBool,
		Limit:   limit,
	}
	rows, err := s.queries.SearchShellHistory(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search shell history: %w", err)
	}
	result := make([]*ShellHistory, 0, len(rows))
	for _, row := range rows {
		result = append(result, &ShellHistory{
			ID:            row.ID,
			Timestamp:     row.Timestamp,
			Command:       row.Command,
			Success:       row.Success,
			ResultSummary: row.ResultSummary.String,
			Intent:        row.Intent.String,
		})
	}
	return result, nil
}

// UpdateShellHistory updates an existing shell history entry.
func (s *Store) UpdateShellHistory(
	ctx context.Context, id int64, success bool, summary string, intent string,
) error {
	params := sqlc.UpdateShellHistoryParams{
		Success:       success,
		ResultSummary: nullString(summary),
		Intent:        nullString(intent),
		ID:            id,
	}
	if err := s.queries.UpdateShellHistory(ctx, params); err != nil {
		return fmt.Errorf("update shell history: %w", err)
	}
	return nil
}

// ClearShellHistory clears all shell history.
func (s *Store) ClearShellHistory(ctx context.Context) error {
	if err := s.queries.ClearShellHistory(ctx); err != nil {
		return fmt.Errorf("clear shell history: %w", err)
	}
	return nil
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

func fromListRow(row listTaskRow) *Task {
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
