//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate -f ../../sqlc.yaml

package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sort"
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

	_, err := s.queries.InsertTaskIdentity(ctx, createdAt)
	if err != nil {
		return nil, fmt.Errorf("insert task identity: %w", err)
	}
	var identityID int64
	err = s.db.QueryRowContext(ctx, "SELECT id FROM tasks ORDER BY id DESC LIMIT 1").Scan(&identityID)
	if err != nil {
		return nil, fmt.Errorf("read task identity id: %w", err)
	}

	projectsJSON, contextsJSON, metaJSON, err := encodeTaskDetails(task)
	if err != nil {
		return nil, fmt.Errorf("encode task details: %w", err)
	}

	versionID, err := s.queries.InsertTaskVersion(ctx, sqlc.InsertTaskVersionParams{
		TaskID:       identityID,
		State:        string(task.State),
		PrevState:    prevStateNull,
		Title:        task.Title,
		Notes:        task.Notes,
		DueOn:        nullDate(task.DueOn),
		WaitingFor:   nullString(task.WaitingFor),
		CompletedAt:  nullUnixTime(completedAt),
		UpdatedAt:    updatedAt,
		Deleted:      0,
		ProjectsJson: projectsJSON,
		ContextsJson: contextsJSON,
		MetaJson:     metaJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("insert task version: %w", err)
	}
	err = s.queries.UpsertTaskCurrent(ctx, sqlc.UpsertTaskCurrentParams{
		ID:           identityID,
		State:        string(task.State),
		PrevState:    prevStateNull,
		Title:        task.Title,
		Notes:        task.Notes,
		DueOn:        nullDate(task.DueOn),
		WaitingFor:   nullString(task.WaitingFor),
		CompletedAt:  nullUnixTime(completedAt),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		ProjectsJson: projectsJSON,
		ContextsJson: contextsJSON,
		MetaJson:     metaJSON,
		VersionID:    versionID,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert current task: %w", err)
	}
	return s.GetTask(ctx, identityID)
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
	current, err := s.GetTask(ctx, task.ID)
	if err != nil {
		return nil, err
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

	params := sqlc.InsertTaskVersionParams{
		TaskID:      task.ID,
		State:       string(task.State),
		PrevState:   prevStateNull,
		Title:       task.Title,
		Notes:       task.Notes,
		DueOn:       nullDate(task.DueOn),
		WaitingFor:  nullString(task.WaitingFor),
		CompletedAt: nullUnixTime(completedAt),
		UpdatedAt:   updatedAt,
		Deleted:     0,
	}

	projectsJSON, contextsJSON, metaJSON, err := encodeTaskDetails(task)
	if err != nil {
		return nil, fmt.Errorf("encode task details: %w", err)
	}
	params.ProjectsJson = projectsJSON
	params.ContextsJson = contextsJSON
	params.MetaJson = metaJSON

	versionID, err := s.queries.InsertTaskVersion(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("insert task version: %w", err)
	}
	if upsertErr := s.queries.UpsertTaskCurrent(ctx, sqlc.UpsertTaskCurrentParams{
		ID:           task.ID,
		State:        string(task.State),
		PrevState:    prevStateNull,
		Title:        task.Title,
		Notes:        task.Notes,
		DueOn:        nullDate(task.DueOn),
		WaitingFor:   nullString(task.WaitingFor),
		CompletedAt:  nullUnixTime(completedAt),
		CreatedAt:    current.CreatedAt.UTC().Unix(),
		UpdatedAt:    updatedAt,
		ProjectsJson: projectsJSON,
		ContextsJson: contextsJSON,
		MetaJson:     metaJSON,
		VersionID:    versionID,
	}); upsertErr != nil {
		return nil, fmt.Errorf("upsert current task: %w", upsertErr)
	}

	return s.GetTask(ctx, task.ID)
}

func (s *Store) GetTask(ctx context.Context, id int64) (*Task, error) {
	row, err := s.queries.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	return fromGetRow(row)
}

func (s *Store) ListTaskVersions(ctx context.Context, taskID int64, limit int64) ([]*TaskVersion, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.queries.ListTaskVersions(ctx, sqlc.ListTaskVersionsParams{TaskID: taskID, Limit: limit})
	if err != nil {
		return nil, fmt.Errorf("list task versions: %w", err)
	}
	versions := make([]*TaskVersion, 0, len(rows))
	for _, row := range rows {
		version, convErr := fromVersionRow(row)
		if convErr != nil {
			return nil, convErr
		}
		versions = append(versions, version)
	}
	return versions, nil
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
		"t.projects_json",
		"t.contexts_json",
		"t.meta_json",
	).
		From("tasks_current t")

	if opts.Recent {
		queryBuilder = queryBuilder.OrderBy("t.updated_at DESC", "t.version_id DESC")
	} else {
		queryBuilder = queryBuilder.OrderBy(
			"CASE WHEN t.state = 'done' THEN 1 ELSE 0 END",
			"CASE WHEN t.due_on IS NULL OR t.due_on = '' THEN 1 ELSE 0 END",
			"t.due_on ASC",
			"t.updated_at DESC",
			"t.version_id DESC",
		)
	}

	if opts.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(opts.Limit))
	}

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
			&row.ProjectsJSON,
			&row.ContextsJSON,
			&row.MetaJSON,
		); scanErr != nil {
			return nil, fmt.Errorf("scan task row: %w", scanErr)
		}
		task, convErr := fromListRow(row)
		if convErr != nil {
			return nil, convErr
		}
		tasks = append(tasks, task)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("iterate task rows: %w", rowsErr)
	}

	return tasks, nil
}

type listTaskRow struct {
	ID           int64
	State        string
	PrevState    sql.NullString
	Title        string
	Notes        string
	DueOn        sql.NullString
	WaitingFor   sql.NullString
	CompletedAt  sql.NullInt64
	CreatedAt    int64
	UpdatedAt    int64
	ProjectsJSON string
	ContextsJSON string
	MetaJSON     string
}

//nolint:gocognit,nestif // Done/undo snapshot logic is centralized for consistency.
func (s *Store) SetDone(ctx context.Context, ids []int64, done bool) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	now := time.Now().UTC()
	updatedAt := now.Unix()
	completedAt := sql.NullInt64{Int64: now.Unix(), Valid: true}

	var changed int64
	for _, id := range ids {
		task, err := s.GetTask(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return 0, err
		}

		next := *task
		next.Projects = append([]string(nil), task.Projects...)
		next.Contexts = append([]string(nil), task.Contexts...)
		next.Meta = copyStringMap(task.Meta)

		if done {
			if task.State == StateDone {
				continue
			}
			prev := task.State
			next.PrevState = &prev
			next.State = StateDone
			next.CompletedAt = &now
		} else {
			if task.State != StateDone {
				continue
			}
			if task.PrevState != nil {
				next.State = *task.PrevState
			} else {
				next.State = StateInbox
			}
			next.PrevState = nil
			next.CompletedAt = nil
		}

		projectsJSON, contextsJSON, metaJSON, encErr := encodeTaskDetails(&next)
		if encErr != nil {
			return 0, fmt.Errorf("encode task details: %w", encErr)
		}

		prevStateNull := sql.NullString{}
		if next.PrevState != nil && *next.PrevState != "" {
			prevStateNull = sql.NullString{String: string(*next.PrevState), Valid: true}
		}

		completedValue := sql.NullInt64{}
		if next.CompletedAt != nil {
			completedValue = completedAt
		}

		versionID, insertErr := s.queries.InsertTaskVersion(ctx, sqlc.InsertTaskVersionParams{
			TaskID:       task.ID,
			State:        string(next.State),
			PrevState:    prevStateNull,
			Title:        next.Title,
			Notes:        next.Notes,
			DueOn:        nullDate(next.DueOn),
			WaitingFor:   nullString(next.WaitingFor),
			CompletedAt:  completedValue,
			UpdatedAt:    updatedAt,
			Deleted:      0,
			ProjectsJson: projectsJSON,
			ContextsJson: contextsJSON,
			MetaJson:     metaJSON,
		})
		if insertErr != nil {
			return 0, fmt.Errorf("insert task version: %w", insertErr)
		}
		err = s.queries.UpsertTaskCurrent(ctx, sqlc.UpsertTaskCurrentParams{
			ID:           task.ID,
			State:        string(next.State),
			PrevState:    prevStateNull,
			Title:        next.Title,
			Notes:        next.Notes,
			DueOn:        nullDate(next.DueOn),
			WaitingFor:   nullString(next.WaitingFor),
			CompletedAt:  completedValue,
			CreatedAt:    task.CreatedAt.UTC().Unix(),
			UpdatedAt:    updatedAt,
			ProjectsJson: projectsJSON,
			ContextsJson: contextsJSON,
			MetaJson:     metaJSON,
			VersionID:    versionID,
		})
		if err != nil {
			return 0, fmt.Errorf("upsert current task: %w", err)
		}
		changed++
	}

	return changed, nil
}

func (s *Store) DeleteTasks(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	updatedAt := time.Now().UTC().Unix()
	var deleted int64
	for _, id := range ids {
		task, err := s.GetTask(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return 0, err
		}
		projectsJSON, contextsJSON, metaJSON, encErr := encodeTaskDetails(task)
		if encErr != nil {
			return 0, fmt.Errorf("encode task details: %w", encErr)
		}

		_, insertErr := s.queries.InsertTaskVersion(ctx, sqlc.InsertTaskVersionParams{
			TaskID:       task.ID,
			State:        string(task.State),
			PrevState:    nullState(task.PrevState),
			Title:        task.Title,
			Notes:        task.Notes,
			DueOn:        nullDate(task.DueOn),
			WaitingFor:   nullString(task.WaitingFor),
			CompletedAt:  nullUnixTime(task.CompletedAt),
			UpdatedAt:    updatedAt,
			Deleted:      1,
			ProjectsJson: projectsJSON,
			ContextsJson: contextsJSON,
			MetaJson:     metaJSON,
		})
		if insertErr != nil {
			return 0, fmt.Errorf("insert tombstone version: %w", insertErr)
		}
		if err = s.queries.DeleteTaskCurrent(ctx, task.ID); err != nil {
			return 0, fmt.Errorf("delete current task: %w", err)
		}
		deleted++
	}

	return deleted, nil
}

func encodeTaskDetails(task *Task) (string, string, string, error) {
	projects := uniqueStrings(cleanNames(task.Projects))
	contexts := uniqueStrings(cleanNames(task.Contexts))

	sort.Strings(projects)
	sort.Strings(contexts)

	meta := task.Meta
	if meta == nil {
		meta = map[string]string{}
	}

	projectsJSON, err := json.Marshal(projects)
	if err != nil {
		return "", "", "", err
	}
	contextsJSON, err := json.Marshal(contexts)
	if err != nil {
		return "", "", "", err
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return "", "", "", err
	}

	return string(projectsJSON), string(contextsJSON), string(metaJSON), nil
}

func (s *Store) ListProjectCounts(ctx context.Context, onlyDone bool, excludeDone bool) ([]NameCount, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT p.value AS name, COUNT(t.id) AS count
FROM tasks_current t
JOIN json_each(t.projects_json) p
WHERE (? = 0 OR t.state = 'done')
  AND (? = 0 OR t.state != 'done')
GROUP BY p.value
ORDER BY p.value ASC;`,
		boolToInt(onlyDone),
		boolToInt(excludeDone),
	)
	if err != nil {
		return nil, fmt.Errorf("list project counts: %w", err)
	}
	defer rows.Close()

	result := make([]NameCount, 0)
	for rows.Next() {
		var row NameCount
		if scanErr := rows.Scan(&row.Name, &row.Count); scanErr != nil {
			return nil, fmt.Errorf("scan project counts: %w", scanErr)
		}
		result = append(result, row)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("iterate project counts: %w", rowsErr)
	}

	return result, nil
}

func (s *Store) ListContextCounts(ctx context.Context, onlyDone bool, excludeDone bool) ([]NameCount, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT c.value AS name, COUNT(t.id) AS count
FROM tasks_current t
JOIN json_each(t.contexts_json) c
WHERE (? = 0 OR t.state = 'done')
  AND (? = 0 OR t.state != 'done')
GROUP BY c.value
ORDER BY c.value ASC;`,
		boolToInt(onlyDone),
		boolToInt(excludeDone),
	)
	if err != nil {
		return nil, fmt.Errorf("list context counts: %w", err)
	}
	defer rows.Close()

	result := make([]NameCount, 0)
	for rows.Next() {
		var row NameCount
		if scanErr := rows.Scan(&row.Name, &row.Count); scanErr != nil {
			return nil, fmt.Errorf("scan context counts: %w", scanErr)
		}
		result = append(result, row)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("iterate context counts: %w", rowsErr)
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

// ClearShellHistory clears all shell history.
func (s *Store) ClearShellHistory(ctx context.Context) error {
	if err := s.queries.ClearShellHistory(ctx); err != nil {
		return fmt.Errorf("clear shell history: %w", err)
	}
	return nil
}

func fromGetRow(row sqlc.GetTaskRow) (*Task, error) {
	projects, contexts, meta, err := decodeTaskDetails(row.ProjectsJson, row.ContextsJson, row.MetaJson)
	if err != nil {
		return nil, fmt.Errorf("decode task details: %w", err)
	}

	return &Task{
		ID:          row.ID,
		State:       State(row.State),
		PrevState:   parseStatePtr(row.PrevState),
		Title:       row.Title,
		Notes:       row.Notes,
		DueOn:       parseDate(row.DueOn),
		WaitingFor:  row.WaitingFor.String,
		CompletedAt: parseUnixTime(row.CompletedAt),
		Projects:    projects,
		Contexts:    contexts,
		Meta:        meta,
		CreatedAt:   time.Unix(row.CreatedAt, 0).UTC(),
		UpdatedAt:   time.Unix(row.UpdatedAt, 0).UTC(),
	}, nil
}

func fromListRow(row listTaskRow) (*Task, error) {
	projects, contexts, meta, err := decodeTaskDetails(row.ProjectsJSON, row.ContextsJSON, row.MetaJSON)
	if err != nil {
		return nil, fmt.Errorf("decode task details: %w", err)
	}

	return &Task{
		ID:          row.ID,
		State:       State(row.State),
		PrevState:   parseStatePtr(row.PrevState),
		Title:       row.Title,
		Notes:       row.Notes,
		DueOn:       parseDate(row.DueOn),
		WaitingFor:  row.WaitingFor.String,
		CompletedAt: parseUnixTime(row.CompletedAt),
		Projects:    projects,
		Contexts:    contexts,
		Meta:        meta,
		CreatedAt:   time.Unix(row.CreatedAt, 0).UTC(),
		UpdatedAt:   time.Unix(row.UpdatedAt, 0).UTC(),
	}, nil
}

func fromVersionRow(row sqlc.TaskVersion) (*TaskVersion, error) {
	projects, contexts, meta, err := decodeTaskDetails(row.ProjectsJson, row.ContextsJson, row.MetaJson)
	if err != nil {
		return nil, fmt.Errorf("decode task version details: %w", err)
	}

	return &TaskVersion{
		VersionID:   row.VersionID,
		TaskID:      row.TaskID,
		State:       State(row.State),
		PrevState:   parseStatePtr(row.PrevState),
		Title:       row.Title,
		Notes:       row.Notes,
		DueOn:       parseDate(row.DueOn),
		WaitingFor:  row.WaitingFor.String,
		CompletedAt: parseUnixTime(row.CompletedAt),
		UpdatedAt:   time.Unix(row.UpdatedAt, 0).UTC(),
		Deleted:     row.Deleted != 0,
		Projects:    projects,
		Contexts:    contexts,
		Meta:        meta,
	}, nil
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

func nullState(value *State) sql.NullString {
	if value == nil || *value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*value), Valid: true}
}

func decodeTaskDetails(projectsJSON, contextsJSON, metaJSON string) ([]string, []string, map[string]string, error) {
	projects := []string{}
	if strings.TrimSpace(projectsJSON) != "" {
		if err := json.Unmarshal([]byte(projectsJSON), &projects); err != nil {
			return nil, nil, nil, err
		}
	}

	contexts := []string{}
	if strings.TrimSpace(contextsJSON) != "" {
		if err := json.Unmarshal([]byte(contextsJSON), &contexts); err != nil {
			return nil, nil, nil, err
		}
	}

	meta := map[string]string{}
	if strings.TrimSpace(metaJSON) != "" {
		if err := json.Unmarshal([]byte(metaJSON), &meta); err != nil {
			return nil, nil, nil, err
		}
	}

	return projects, contexts, meta, nil
}

func copyStringMap(input map[string]string) map[string]string {
	if input == nil {
		return map[string]string{}
	}
	return maps.Clone(input)
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
