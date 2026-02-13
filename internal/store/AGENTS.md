# internal/store AGENTS

## OVERVIEW
SQLite/libSQL store: open DB, run migrations, keep append-only tasks consistent, and translate filters into safe SQL.

## WHERE TO LOOK
- `internal/store/store.go` open flow (pragmas -> goose Up), sync wiring, and all read/write paths
- `internal/store/migrations.go` embedded migration FS used by goose
- `internal/store/migrations/*.sql` schema history; append-only pivot and any destructive resets live here
- `db/queries/*.sql` source-of-truth SQL for sqlc (adds/changes start here)
- `internal/store/sqlc/` generated query layer (boundary: typed params/rows; do not hand-edit)
- `internal/store/filter_sql_builder.go` nlp.FilterExpr -> SQL WHERE clause (Squirrel + JSON1)
- `internal/store/filter_sql_builder_test.go` and `internal/store/store_filter_test.go` filter behavior/perf guardrails
- `internal/store/types.go` Go-side task shape; mirrors tasks_current columns (plus helpers)

## CONVENTIONS
- Open sequence: connect (local or sync) -> set PRAGMAs -> run goose migrations from embedded FS -> create `sqlc.New(db)`
- Append-only model: `tasks` is identity; `task_versions` is history; `tasks_current` is the read model
- Writes create a new `task_versions` row then update/replace the corresponding `tasks_current` row; never mutate old versions
- Deletes are logical (`deleted` flag in versions); keep `tasks_current` free of deleted rows
- JSON columns are stored as text; always keep them valid JSON arrays/objects (filters rely on json_valid/json_each)

## ANTI-PATTERNS
- Updating `tasks_current` without also writing a matching `task_versions` row (silently breaks history)
- Editing or depending on `internal/store/sqlc/*` directly; change `db/queries/*.sql` and regenerate instead
- Building SQL by string concat in filters; always parameterize and keep a strict predicate whitelist
- Using `tasks` alone for reads; most queries should target `tasks_current` (or explicit version scans)
- Adding migrations that rewrite history without a clear recovery story; one bad migration bricks existing DBs

## NOTES
- `internal/store/migrations/00007_append_only_tasks.sql` introduces `task_versions` + `tasks_current`; `00008` ensures `tasks_current` is a table (not a view)
- `internal/store/migrations/00009_reset_append_only_schema.sql` drops task tables; treat as data-loss and keep it out of normal upgrade paths
- Filter SQL assumes table alias `t` and uses SQLite JSON1 (`json_each`, `json_array_length`); keep column names stable or update builder + tests
- Sync uses `tursogo.TursoSyncDb`; auth token may come from `LIBSQL_AUTH_TOKEN`; `Sync()` pulls only (push is explicit)
- When debugging "missing" rows, check: latest version has `deleted=0`, `tasks_current.version_id` matches, and migrations ran on the target DB file
