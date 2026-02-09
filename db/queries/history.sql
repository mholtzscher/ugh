-- name: InsertShellHistory :one
INSERT INTO shell_history (
  timestamp,
  command,
  success,
  result_summary,
  intent
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING
  id,
  timestamp,
  command,
  success,
  result_summary,
  intent;

-- name: ListShellHistory :many
SELECT
  id,
  timestamp,
  command,
  success,
  result_summary,
  intent
FROM shell_history
ORDER BY timestamp DESC
LIMIT ?;

-- name: SearchShellHistory :many
SELECT
  id,
  timestamp,
  command,
  success,
  result_summary,
  intent
FROM shell_history
WHERE (
  ? IS NULL OR command LIKE '%' || ? || '%'
)
  AND (? IS NULL OR intent = ?)
  AND (? IS NULL OR success = ?)
ORDER BY timestamp DESC
LIMIT ?;

-- name: ClearShellHistory :exec
DELETE FROM shell_history;
