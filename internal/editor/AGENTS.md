# internal/editor AGENTS

Implements the `ugh edit` $EDITOR TOML round-trip with optional Taplo JSON schema hints and strict post-edit validation.

## Where To Look
- `internal/editor/editor.go` end-to-end flow: Task -> TOML -> temp file -> spawn editor -> parse -> validate.
- `internal/editor/task.schema.json` embedded JSON Schema; written into temp dir for editor/LSP support.
- `cmd/edit.go` command wiring; calls editor then `service.FullUpdateTask`.

## Conventions
- Editor selection: VISUAL then EDITOR, else first found of vim/vi/nano, else vi.
- Temp workspace: prefer temp dir under current working directory; fallback to system temp.
- Schema hinting (best-effort): write `ugh-task.schema.json` + `taplo.toml`, add `#:schema file://...` header.
- Parsing: decode TOML into TaskTOML with BurntSushi/toml.
- Validation/normalization: trim title (required), normalize state (default inbox), validate due_on (YYYY-MM-DD), trim waiting_for, clean/dedupe tags.
- Keep schema and struct in sync: TaskTOML tags/defaults/allowed values match `internal/editor/task.schema.json`.

## Anti-Patterns
- Adding editable fields without updating both TaskTOML and schema.
- Relying on schema for correctness (schema is advisory; validate gates writes).
- Making schema setup fatal (schema files are best-effort).
- Allowing partial updates here (editor path is full replace).
- Lossy conversions without explicit rules; due_on is date-only.

## Notes
- "No changes" is exact file-content match; whitespace-only edits count as changes.
- toml.Decode ignores # comments; header mainly for humans and schema tooling.
