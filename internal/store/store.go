//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/mholtzscher/ugh/internal/store/sqlc"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Store struct {
	db      *sql.DB
	queries *sqlc.Queries
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

	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply pragma: %w", err)
	}
	if _, err := db.ExecContext(ctx, "PRAGMA busy_timeout=5000;"); err != nil {
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
	if err := goose.Up(db, "migrations"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	store := &Store{db: db, queries: sqlc.New(db)}
	return store, nil
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
	if task.Meta == nil {
		task.Meta = map[string]string{}
	}

	now := time.Now().UTC()
	createdAt := now.Unix()
	updatedAt := now.Unix()

	params := sqlc.InsertTaskParams{
		Done:           boolToInt(task.Done),
		Priority:       nullString(task.Priority),
		CompletionDate: nullDate(task.CompletionDate),
		CreationDate:   nullDate(task.CreationDate),
		Description:    task.Description,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
	row, err := s.queries.InsertTask(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("insert task: %w", err)
	}

	if err := s.insertTokens(ctx, row.ID, task); err != nil {
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
	if task.Meta == nil {
		task.Meta = map[string]string{}
	}

	updatedAt := time.Now().UTC().Unix()

	params := sqlc.UpdateTaskParams{
		Done:           boolToInt(task.Done),
		Priority:       nullString(task.Priority),
		CompletionDate: nullDate(task.CompletionDate),
		CreationDate:   nullDate(task.CreationDate),
		Description:    task.Description,
		UpdatedAt:      updatedAt,
		ID:             task.ID,
	}
	if _, err := s.queries.UpdateTask(ctx, params); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	if err := s.deleteTokens(ctx, task.ID); err != nil {
		return nil, err
	}
	if err := s.insertTokens(ctx, task.ID, task); err != nil {
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
	if err := s.loadTokens(ctx, []*Task{task}); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Store) ListTasks(ctx context.Context, filters Filters) ([]*Task, error) {
	var status any
	if filters.DoneOnly {
		status = int64(1)
	} else if filters.TodoOnly {
		status = int64(0)
	}

	params := sqlc.ListTasksParams{
		Status:   status,
		Project:  nullAny(filters.Project),
		Context:  nullAny(filters.Context),
		Priority: nullAny(filters.Priority),
		Search:   nullAny(filters.Search),
	}
	rows, err := s.queries.ListTasks(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	tasks := make([]*Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, fromListRow(row))
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
	completion := sql.NullString{}
	if done {
		completion = sql.NullString{String: time.Now().UTC().Format("2006-01-02"), Valid: true}
	}
	params := sqlc.SetDoneParams{
		Done:           boolToInt(done),
		CompletionDate: completion,
		UpdatedAt:      time.Now().UTC().Unix(),
		Ids:            ids,
	}

	count, err := s.queries.SetDone(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("update done: %w", err)
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

	projects, err := s.queries.ListProjects(ctx, ids)
	if err != nil {
		return fmt.Errorf("load projects: %w", err)
	}
	for _, row := range projects {
		if task := byID[row.TaskID]; task != nil {
			task.Projects = append(task.Projects, row.Name)
		}
	}

	contexts, err := s.queries.ListContexts(ctx, ids)
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

	unknown, err := s.queries.ListUnknown(ctx, ids)
	if err != nil {
		return fmt.Errorf("load unknown: %w", err)
	}
	for _, row := range unknown {
		if task := byID[row.TaskID]; task != nil {
			task.Unknown = append(task.Unknown, row.Token)
		}
	}

	return nil
}

func (s *Store) insertTokens(ctx context.Context, id int64, task *Task) error {
	for _, project := range uniqueStrings(task.Projects) {
		if err := s.queries.InsertProject(ctx, sqlc.InsertProjectParams{TaskID: id, Name: project}); err != nil {
			return fmt.Errorf("insert project: %w", err)
		}
	}
	for _, context := range uniqueStrings(task.Contexts) {
		if err := s.queries.InsertContext(ctx, sqlc.InsertContextParams{TaskID: id, Name: context}); err != nil {
			return fmt.Errorf("insert context: %w", err)
		}
	}
	for key, value := range task.Meta {
		if err := s.queries.InsertMeta(ctx, sqlc.InsertMetaParams{TaskID: id, Key: key, Value: value}); err != nil {
			return fmt.Errorf("insert meta: %w", err)
		}
	}
	for ordinal, token := range task.Unknown {
		params := sqlc.InsertUnknownParams{TaskID: id, Ordinal: int64(ordinal), Token: token}
		if err := s.queries.InsertUnknown(ctx, params); err != nil {
			return fmt.Errorf("insert unknown: %w", err)
		}
	}
	return nil
}

func (s *Store) deleteTokens(ctx context.Context, id int64) error {
	if err := s.queries.DeleteTokens(ctx, id); err != nil {
		return fmt.Errorf("delete projects: %w", err)
	}
	if err := s.queries.DeleteContexts(ctx, id); err != nil {
		return fmt.Errorf("delete contexts: %w", err)
	}
	if err := s.queries.DeleteMeta(ctx, id); err != nil {
		return fmt.Errorf("delete meta: %w", err)
	}
	if err := s.queries.DeleteUnknown(ctx, id); err != nil {
		return fmt.Errorf("delete unknown: %w", err)
	}
	return nil
}

func fromGetRow(row sqlc.GetTaskRow) *Task {
	return &Task{
		ID:             row.ID,
		Done:           row.Done == 1,
		Priority:       row.Priority.String,
		CompletionDate: parseDate(row.CompletionDate),
		CreationDate:   parseDate(row.CreationDate),
		Description:    row.Description,
		CreatedAt:      time.Unix(row.CreatedAt, 0).UTC(),
		UpdatedAt:      time.Unix(row.UpdatedAt, 0).UTC(),
	}
}

func fromListRow(row sqlc.ListTasksRow) *Task {
	return &Task{
		ID:             row.ID,
		Done:           row.Done == 1,
		Priority:       row.Priority.String,
		CompletionDate: parseDate(row.CompletionDate),
		CreationDate:   parseDate(row.CreationDate),
		Description:    row.Description,
		CreatedAt:      time.Unix(row.CreatedAt, 0).UTC(),
		UpdatedAt:      time.Unix(row.UpdatedAt, 0).UTC(),
	}
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
