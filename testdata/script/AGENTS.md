# testdata/script AGENTS

Testscript files exercise the CLI end-to-end via shell-like scripts and golden assertions.

## Where To Look
- `testdata/script/*.txt` scenarios (one file = one story).
- Inline fixtures via `-- name --` blocks (golden output, config TOML, etc.).
- Existing patterns: `basic.txt`, `errors.txt`, `show_edit.txt`, `repl_filters*.txt`, `sync*.txt`.

## Conventions
- Keep each script small and linear: setup -> action -> asserts.
- Prefer per-test isolation: write DB/config under $WORK (e.g. `--db $WORK/db.sqlite`, `--config $WORK/config.toml`).
- Use `exec` for commands; use `! exec` when failure is expected.
- Assert with `stdout <regex>` / `stderr <regex>` (regex, not literals).
- For full output, prefer `cmp stdout want-*.txt` over many `stdout` lines.
- Store expected multi-line output in a `-- want-*.txt --` block near the assertion.
- Use `exists <path>` for filesystem side-effects.

## Anti-Patterns
- Rely on user machine state (real $HOME, global config, existing DB files).
- Assert unstable text (timestamps, durations, platform-specific paths, unordered lists).
- Overuse `exec sh -c` to create files; prefer `-- file --` blocks.
- Over-broad regex like `stdout '.*'`.
- Repeat long command lines; keep flag ordering consistent.

## Notes
- $WORK is a fresh temp directory per script; use it for all writes.
- Runner sets hermetic HOME/XDG and `TZ=UTC`.
- stdout/stderr patterns are regex; escape `[` `]` `(` `)` `.` `?` `+` `*` when matching JSON.
- Negative checks: `! stdout <regex>` ensures absence.
