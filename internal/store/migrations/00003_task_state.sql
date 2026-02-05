-- +goose Up

ALTER TABLE tasks ADD COLUMN state TEXT NOT NULL DEFAULT 'inbox';
ALTER TABLE tasks ADD COLUMN prev_state TEXT;

-- Migrate legacy fields (done/status) into the simplified one-dimensional state.
-- Mapping:
--   inbox   -> inbox
--   next    -> now
--   waiting -> waiting
--   someday -> later
--   done=1  -> done (and prev_state remembers the prior active state)
UPDATE tasks
SET
  state = CASE
    WHEN done = 1 THEN 'done'
    WHEN status = 'next' THEN 'now'
    WHEN status = 'someday' THEN 'later'
    ELSE status
  END,
  prev_state = CASE
    WHEN done = 1 THEN CASE
      WHEN status = 'next' THEN 'now'
      WHEN status = 'someday' THEN 'later'
      ELSE status
    END
    ELSE NULL
  END;

-- Ensure completed_at exists for done tasks.
UPDATE tasks
SET completed_at = COALESCE(completed_at, updated_at)
WHERE state = 'done' AND (completed_at IS NULL OR completed_at = 0);

CREATE INDEX IF NOT EXISTS idx_tasks_state ON tasks(state);

-- +goose Down

-- Best-effort downgrade: keep columns but restore legacy semantics.
UPDATE tasks
SET
  done = CASE WHEN state = 'done' THEN 1 ELSE 0 END,
  status = CASE
    WHEN state = 'now' THEN 'next'
    WHEN state = 'later' THEN 'someday'
    WHEN state = 'done' THEN CASE
      WHEN prev_state = 'now' THEN 'next'
      WHEN prev_state = 'later' THEN 'someday'
      WHEN prev_state IS NULL OR prev_state = '' THEN 'inbox'
      ELSE prev_state
    END
    ELSE state
  END;
