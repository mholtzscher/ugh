package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context, path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("db path is required")
	}

	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve db path: %w", err)
	}

	db, err := sql.Open("sqlite", abspath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	store := &Store{db: db}
	if err := store.bootstrap(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) bootstrap(ctx context.Context) error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA foreign_keys=ON;",
	}
	for _, pragma := range pragmas {
		if _, err := s.db.ExecContext(ctx, pragma); err != nil {
			return fmt.Errorf("apply pragma: %w", err)
		}
	}

	if _, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_version (
  version INTEGER NOT NULL
);
`); err != nil {
		return fmt.Errorf("create schema_version: %w", err)
	}

	version, err := s.currentVersion(ctx)
	if err != nil {
		return err
	}
	if version == 0 {
		if err := s.applyV1(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) currentVersion(ctx context.Context) (int, error) {
	row := s.db.QueryRowContext(ctx, "SELECT version FROM schema_version LIMIT 1")
	var version int
	if err := row.Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("read schema version: %w", err)
	}
	return version, nil
}

func (s *Store) applyV1(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			done INTEGER NOT NULL DEFAULT 0,
			priority TEXT,
			completion_date TEXT,
			creation_date TEXT,
			description TEXT NOT NULL DEFAULT "",
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS task_projects (
			task_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS task_contexts (
			task_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS task_meta (
			task_id INTEGER NOT NULL,
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS task_unknown (
			task_id INTEGER NOT NULL,
			ordinal INTEGER NOT NULL,
			token TEXT NOT NULL,
			FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
		);`,
		`INSERT INTO schema_version(version) VALUES (1);`,
	}

	for _, stmt := range stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("apply migration: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration: %w", err)
	}
	return nil
}

func (s *Store) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	if task == nil {
		return nil, errors.New("task is required")
	}
	if task.Meta == nil {
		task.Meta = map[string]string{}
	}

	now := time.Now().UTC()
	createdAt := now.Unix()
	updatedAt := now.Unix()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin create: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(ctx, `
INSERT INTO tasks (done, priority, completion_date, creation_date, description, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
`, boolToInt(task.Done), nullString(task.Priority), formatDate(task.CompletionDate), formatDate(task.CreationDate), task.Description, createdAt, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("task id: %w", err)
	}

	if err := insertTokens(ctx, tx, id, task); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit task: %w", err)
	}

	return s.GetTask(ctx, id)
}

func (s *Store) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	if task == nil {
		return nil, errors.New("task is required")
	}
	if task.ID == 0 {
		return nil, errors.New("task id is required")
	}
	if task.Meta == nil {
		task.Meta = map[string]string{}
	}

	updatedAt := time.Now().UTC().Unix()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin update: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, `
UPDATE tasks
SET done = ?, priority = ?, completion_date = ?, creation_date = ?, description = ?, updated_at = ?
WHERE id = ?
`, boolToInt(task.Done), nullString(task.Priority), formatDate(task.CompletionDate), formatDate(task.CreationDate), task.Description, updatedAt, task.ID)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	if err := deleteTokens(ctx, tx, task.ID); err != nil {
		return nil, err
	}
	if err := insertTokens(ctx, tx, task.ID, task); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit update: %w", err)
	}

	return s.GetTask(ctx, task.ID)
}

func (s *Store) GetTask(ctx context.Context, id int64) (*Task, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, done, priority, completion_date, creation_date, description, created_at, updated_at
FROM tasks WHERE id = ?
`, id)

	task, err := scanTask(row)
	if err != nil {
		return nil, err
	}
	if err := s.loadTokens(ctx, []*Task{task}); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Store) ListTasks(ctx context.Context, filters Filters) ([]*Task, error) {
	where := []string{}
	args := []any{}

	if filters.DoneOnly {
		where = append(where, "t.done = 1")
	} else if filters.TodoOnly {
		where = append(where, "t.done = 0")
	}
	if filters.Project != "" {
		where = append(where, "EXISTS (SELECT 1 FROM task_projects p WHERE p.task_id = t.id AND p.name = ?)")
		args = append(args, filters.Project)
	}
	if filters.Context != "" {
		where = append(where, "EXISTS (SELECT 1 FROM task_contexts c WHERE c.task_id = t.id AND c.name = ?)")
		args = append(args, filters.Context)
	}
	if filters.Priority != "" {
		where = append(where, "t.priority = ?")
		args = append(args, filters.Priority)
	}
	if filters.Search != "" {
		like := "%%%s%%"
		where = append(where, `(
  t.description LIKE ?
  OR EXISTS (SELECT 1 FROM task_projects p WHERE p.task_id = t.id AND p.name LIKE ?)
  OR EXISTS (SELECT 1 FROM task_contexts c WHERE c.task_id = t.id AND c.name LIKE ?)
  OR EXISTS (SELECT 1 FROM task_meta m WHERE m.task_id = t.id AND (m.key LIKE ? OR m.value LIKE ?))
)`)
		args = append(args, fmt.Sprintf(like, filters.Search), fmt.Sprintf(like, filters.Search), fmt.Sprintf(like, filters.Search), fmt.Sprintf(like, filters.Search), fmt.Sprintf(like, filters.Search))
	}

	query := `
SELECT t.id, t.done, t.priority, t.completion_date, t.creation_date, t.description, t.created_at, t.updated_at
FROM tasks t`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY CASE WHEN t.done = 1 THEN 1 ELSE 0 END, t.priority IS NULL, t.priority ASC, t.created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list tasks rows: %w", err)
	}

	if err := s.loadTokens(ctx, tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Store) SetDone(ctx context.Context, ids []int64, done bool) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	query := "UPDATE tasks SET done = ?, completion_date = ?, updated_at = ? WHERE id IN (" + placeholders(len(ids)) + ")"
	args := make([]any, 0, len(ids)+3)
	args = append(args, boolToInt(done))
	if done {
		today := time.Now().UTC().Format("2006-01-02")
		args = append(args, today)
	} else {
		args = append(args, nil)
	}
	args = append(args, time.Now().UTC().Unix())
	for _, id := range ids {
		args = append(args, id)
	}
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("update done: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update done count: %w", err)
	}
	return count, nil
}

func (s *Store) DeleteTasks(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	query := "DELETE FROM tasks WHERE id IN (" + placeholders(len(ids)) + ")"
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("delete tasks: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete count: %w", err)
	}
	return count, nil
}

func scanTask(scanner interface{ Scan(...any) error }) (*Task, error) {
	var task Task
	var done int
	var priority sql.NullString
	var completion sql.NullString
	var creation sql.NullString
	var description string
	var createdAt int64
	var updatedAt int64

	if err := scanner.Scan(&task.ID, &done, &priority, &completion, &creation, &description, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("scan task: %w", err)
	}

	task.Done = done == 1
	task.Priority = priority.String
	task.Description = description
	if creation.Valid {
		parsed, err := time.Parse("2006-01-02", creation.String)
		if err == nil {
			utc := parsed.UTC()
			task.CreationDate = &utc
		}
	}
	if completion.Valid {
		parsed, err := time.Parse("2006-01-02", completion.String)
		if err == nil {
			utc := parsed.UTC()
			task.CompletionDate = &utc
		}
	}
	task.CreatedAt = time.Unix(createdAt, 0).UTC()
	task.UpdatedAt = time.Unix(updatedAt, 0).UTC()

	return &task, nil
}

func (s *Store) loadTokens(ctx context.Context, tasks []*Task) error {
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
		task.Unknown = nil
		task.Meta = map[string]string{}
	}

	if err := loadSimpleTokens(ctx, s.db, "task_projects", "name", ids, byID, func(task *Task, value string) {
		task.Projects = append(task.Projects, value)
	}); err != nil {
		return err
	}

	if err := loadSimpleTokens(ctx, s.db, "task_contexts", "name", ids, byID, func(task *Task, value string) {
		task.Contexts = append(task.Contexts, value)
	}); err != nil {
		return err
	}

	if err := loadMeta(ctx, s.db, ids, byID); err != nil {
		return err
	}

	if err := loadUnknown(ctx, s.db, ids, byID); err != nil {
		return err
	}

	return nil
}

func loadSimpleTokens(ctx context.Context, db *sql.DB, table string, column string, ids []int64, byID map[int64]*Task, add func(*Task, string)) error {
	query := fmt.Sprintf("SELECT task_id, %s FROM %s WHERE task_id IN (%s) ORDER BY task_id", column, table, placeholders(len(ids)))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("load %s: %w", table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var value string
		if err := rows.Scan(&id, &value); err != nil {
			return fmt.Errorf("scan %s: %w", table, err)
		}
		task := byID[id]
		if task != nil {
			add(task, value)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows %s: %w", table, err)
	}
	return nil
}

func loadMeta(ctx context.Context, db *sql.DB, ids []int64, byID map[int64]*Task) error {
	query := "SELECT task_id, key, value FROM task_meta WHERE task_id IN (" + placeholders(len(ids)) + ") ORDER BY task_id"
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("load meta: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var key, value string
		if err := rows.Scan(&id, &key, &value); err != nil {
			return fmt.Errorf("scan meta: %w", err)
		}
		task := byID[id]
		if task != nil {
			task.Meta[key] = value
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows meta: %w", err)
	}
	return nil
}

func loadUnknown(ctx context.Context, db *sql.DB, ids []int64, byID map[int64]*Task) error {
	query := "SELECT task_id, ordinal, token FROM task_unknown WHERE task_id IN (" + placeholders(len(ids)) + ") ORDER BY task_id, ordinal"
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("load unknown: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var ordinal int
		var token string
		if err := rows.Scan(&id, &ordinal, &token); err != nil {
			return fmt.Errorf("scan unknown: %w", err)
		}
		task := byID[id]
		if task != nil {
			task.Unknown = append(task.Unknown, token)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows unknown: %w", err)
	}
	return nil
}

func insertTokens(ctx context.Context, tx *sql.Tx, id int64, task *Task) error {
	for _, project := range uniqueStrings(task.Projects) {
		if _, err := tx.ExecContext(ctx, "INSERT INTO task_projects (task_id, name) VALUES (?, ?)", id, project); err != nil {
			return fmt.Errorf("insert project: %w", err)
		}
	}
	for _, context := range uniqueStrings(task.Contexts) {
		if _, err := tx.ExecContext(ctx, "INSERT INTO task_contexts (task_id, name) VALUES (?, ?)", id, context); err != nil {
			return fmt.Errorf("insert context: %w", err)
		}
	}
	if task.Meta != nil {
		for key, value := range task.Meta {
			if _, err := tx.ExecContext(ctx, "INSERT INTO task_meta (task_id, key, value) VALUES (?, ?, ?)", id, key, value); err != nil {
				return fmt.Errorf("insert meta: %w", err)
			}
		}
	}
	for ordinal, token := range task.Unknown {
		if _, err := tx.ExecContext(ctx, "INSERT INTO task_unknown (task_id, ordinal, token) VALUES (?, ?, ?)", id, ordinal, token); err != nil {
			return fmt.Errorf("insert unknown: %w", err)
		}
	}
	return nil
}

func deleteTokens(ctx context.Context, tx *sql.Tx, id int64) error {
	tables := []string{"task_projects", "task_contexts", "task_meta", "task_unknown"}
	for _, table := range tables {
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+table+" WHERE task_id = ?", id); err != nil {
			return fmt.Errorf("delete tokens: %w", err)
		}
	}
	return nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func formatDate(val *time.Time) any {
	if val == nil {
		return nil
	}
	return val.Format("2006-01-02")
}

func placeholders(count int) string {
	parts := make([]string, count)
	for i := 0; i < count; i++ {
		parts[i] = "?"
	}
	return strings.Join(parts, ",")
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
