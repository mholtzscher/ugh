# TUI Implementation Plan (Bubble Tea + Huh)

## Overview

This document describes a staged plan to add a terminal UI to `ugh` using:

- `bubbletea` for the main interactive application
- `bubbles` for reusable widgets (list, text input, viewport, etc.)
- `huh` for structured forms (add/edit/config)

The TUI should reuse existing task/business logic from `internal/service` and preserve current CLI behavior.

## Goals

1. Launch the TUI from the root command (`ugh`) as the interactive entry point.
2. Make common workflows faster than raw CLI flags:
   - triage tasks
   - state transitions
   - add/edit
   - filtering/search
3. Keep behavior consistent with existing commands (`add`, `list`, `edit`, `done`, `undo`, `rm`, etc.).
4. Keep architecture testable and avoid duplicating domain logic in UI code.
5. Make bare `ugh` open the TUI in interactive terminals.
6. Support basic ANSI theming using the standard palette from https://ansicolor.com/.

## Non-Goals (Initial Rollout)

- Replacing existing CLI commands.
- Adding remote-first or daemon-only data paths.
- Building a fully mouse-driven UI.
- Implementing a full plugin-based theme system in the first release.

## Command Entry Behavior

Target behavior:

- `ugh` (no subcommand, interactive TTY): launch TUI.
- `ugh <subcommand> ...`: run existing CLI command behavior unchanged.
- `ugh --help` or `ugh help`: show help, do not launch TUI.
- non-interactive invocation (stdin/stdout not TTY): do not launch TUI.

Optional safeguards for script compatibility:

- add `--no-tui` global flag and/or `UGH_NO_TUI=1`
- do not auto-launch TUI when `--json` is provided

## Theming (ANSI Colors)

The TUI should support a small set of built-in themes based on the ANSI color palette documented at https://ansicolor.com/.

Base palette (ANSI 0-15):

- `0` black
- `1` red
- `2` green
- `3` yellow
- `4` blue
- `5` magenta
- `6` cyan
- `7` light gray
- `8` dark gray
- `9` light red
- `10` light green
- `11` light yellow
- `12` light blue
- `13` light magenta
- `14` light cyan
- `15` white

Initial theme presets:

- `ansi-default` (balanced contrast)
- `ansi-light` (for light terminal themes)
- `ansi-high-contrast` (max readability)

Color token model (theme maps these tokens to ANSI indices):

- `bg`, `surface`, `border`
- `text`, `textMuted`
- `accent`, `accentSoft`
- `success`, `warning`, `danger`, `info`
- `selectionBg`, `selectionFg`, `focusBorder`

Implementation notes:

- Use `lipgloss.Color("<index>")` tokens for ANSI compatibility.
- Keep visuals readable in 16-color terminals; do not require truecolor.
- Respect `NO_COLOR` and existing `--no-color` behavior.
- If color is disabled, render using plain text + structural indicators only.

Theme selection:

- Config value: `[ui] theme = "ansi-default"`
- Theme is config-managed (no CLI theme flag)

Example:

```toml
[ui]
theme = "ansi-default"
```

Precedence for theme choice:

1. config `[ui].theme`
2. default `ansi-default`

## Proposed Main Layout

```text
+--------------------------------------------------------------------------------------------------+
| ugh                                      [Tasks]  Calendar  System                 Sync: OK      |
+------------------------+-------------------------------------------+-----------------------------+
| NAV                    | TASKS                                     | DETAILS                     |
| > Inbox (12)           | / search: phone      [todo] [due-only]    | #42 Call mom                |
|   Now (5)              +----+---------+------------+----------+-----+ state: waiting             |
|   Waiting (3)          | ID | State   | Due        | Waiting  | Task| due: 2026-02-12            |
|   Later (9)            +----+---------+------------+----------+-----+ waiting_for: Alice         |
|   Done (34)            | 42 | waiting | 2026-02-12 | Alice    | Call| projects: family           |
|                        | 41 | now     | -          | -        | Pay | contexts: phone            |
| Projects               | 17 | inbox   | -          | -        | Buy | meta: channel=sms          |
|   work (7)             | 11 | inbox   | 2026-02-20 | -        | Book| notes: ...                |
|   family (4)           |                                           |                             |
| Contexts               |                                           | [e] edit  [x] done [D] rm  |
|   phone (3)            |                                           |                             |
|   home (6)             |                                           |                             |
+------------------------+-------------------------------------------+-----------------------------+
| j/k move  tab pane  enter open  a add  e edit  i/n/w/l state  x done  u undo  / search  ? help |
+--------------------------------------------------------------------------------------------------+
```

## UX Model

Primary screens:

1. **Tasks** (default)
   - left nav: built-in states + projects + contexts
   - middle list: tasks for current scope and filters
   - right detail: full task info and action hints
2. **Calendar**
   - due-date-focused view (group by day)
3. **System**
   - sync status/actions and daemon status/actions

Fallback for narrow terminals:

- collapse to one-column list + modal detail panel

## Keybindings (Initial)

Navigation:

- `j/k` or up/down: move selection
- `g/G`: top/bottom of list
- `tab`: cycle focused pane
- `enter`: open/expand selected task detail
- `/`: search
- `esc`: clear search / close modal

Task actions:

- `a`: add task
- `e`: edit selected task
- `x`: mark done
- `u`: undo (reopen)
- `i`: move to inbox
- `n`: move to now
- `w`: move to waiting
- `l`: move to later
- `D`: delete task (confirm)

Filtering:

- `t`: cycle `todo -> all -> done`
- `.`: toggle due-only
- `p`: filter by project
- `c`: filter by context

Global:

- `?`: help overlay
- `q`: quit / back

## Huh Form Usage

Use Huh for structured, multi-field flows:

- Add task
- Edit task
- Config setup flows (`db.sync_url`, `db.auth_token`, `db.sync_on_write`)

Prefer inline interactions for fast repetitive actions:

- state changes
- done/undo
- filtering and navigation

## Architecture Plan

Add a new package tree:

```text
internal/tui/
  app.go                # program bootstrap (bubbletea Program)
  model.go              # root model/state
  messages.go           # tea.Msg types
  keymap.go             # centralized key definitions
  theme.go              # ANSI theme tokens + selection
  styles.go             # lipgloss styles
  layout.go             # responsive layout calculation
  view_tasks.go         # tasks workspace rendering
  view_calendar.go      # calendar rendering
  view_system.go        # sync/daemon rendering
  actions.go            # async commands and dispatch
  forms_task.go         # Huh add/edit task forms
  forms_config.go       # Huh config forms
  filters.go            # list filter state + conversions
```

Add root command integration:

- wire root `Action` in `cmd/root.go` to launch TUI when invoked as bare `ugh` in interactive mode
- preserve current subcommand behavior for all existing commands

Add config integration:

- extend config model with `[ui]` section and `theme` field
- load theme from config during TUI startup

### Service integration

- Reuse `service.Service` interface from `internal/service/interface.go`.
- Route all task operations through service calls used by existing CLI commands.
- Keep state transition and validation logic in service/domain layers.

### Async execution model

- Every DB/network operation runs in `tea.Cmd` to avoid blocking render loop.
- Show in-progress indicators while commands run.
- Emit explicit success/error messages (status bar + optional toast).

### Error handling

- Non-fatal errors: status bar + keep current selection.
- Fatal init errors (open store/config): fail startup with clear message.
- Database lock errors: show actionable hint (daemon lock/retry guidance).

## Rollout Plan

### Phase 1: Skeleton + Read/Triage MVP

Scope:

- Add dependencies (`bubbletea`, `bubbles`, `lipgloss`, `huh`).
- Add root command launch path and initialize app model.
- Add ANSI theme token system and three built-in presets.
- Add config-backed theme selection (`[ui].theme`).
- Implement tasks workspace with:
  - nav sections (states/projects/contexts)
  - task list for selected scope
  - task detail pane
  - keyboard navigation + help footer
- Implement actions:
  - done/undo
  - state transitions (`inbox/now/waiting/later`)
  - delete with confirm

Acceptance:

- User can triage tasks without leaving TUI.
- Behavior matches existing CLI state semantics.
- TUI theme can be switched and remains readable with ANSI 16 colors.

### Phase 2: Huh Forms + Filtering

Scope:

- Add Huh add/edit task forms with all supported fields.
- Add search and quick filters (`todo/all/done`, `due-only`, project/context filters).
- Add optimistic refresh strategy after write operations.

Acceptance:

- User can create/edit tasks with full schema from TUI.
- Filter results match `ugh list` semantics.

### Phase 3: Calendar View

Scope:

- Add due-date-focused calendar screen backed by `DueOnly` list queries.
- Group tasks by date and support quick open/edit.

Acceptance:

- Calendar reflects `ugh calendar` behavior.

### Phase 4: System View (Sync + Daemon)

Scope:

- Sync actions/status (`sync`, `pull`, `push`, `status`).
- Daemon status + control actions (`start`, `stop`, `restart`, `install`, `uninstall`, `logs` entry guidance).

Acceptance:

- User can inspect and trigger sync/daemon operations in one place.

### Phase 5: Polish

Scope:

- Better empty/loading/error states.
- Narrow terminal layout polish.
- Multi-select bulk actions.
- Keymap discoverability improvements.

Acceptance:

- Stable and ergonomic daily-use experience.

## Testing Plan

1. **Unit tests**
   - model update transitions (key events -> state changes)
   - filter conversion logic
   - async action message handling
2. **Service contract tests**
   - ensure TUI action paths produce same results as CLI operations
3. **Golden/snapshot tests (selective)**
   - critical render states for tasks view and empty/error states
4. **Manual smoke checks**
   - add/edit/done/undo/delete
   - filter parity with `list`
   - sync/daemon actions where configured
   - verify readability and selection/focus cues across all built-in themes

## Risks and Mitigations

1. **UI freeze on I/O**
   - Mitigation: async `tea.Cmd` calls only; no direct DB/network in `Update`/`View`.
2. **Divergent behavior vs CLI**
   - Mitigation: keep all business logic in service/domain, not view layer.
3. **DB lock contention with daemon**
   - Mitigation: reuse existing retry behavior and show actionable lock messaging.
4. **Keybinding sprawl**
   - Mitigation: centralized keymap + context-aware help overlay.

## Implementation Checklist

- [ ] Wire root command default action to launch TUI in interactive mode
- [ ] Add `internal/tui` package skeleton
- [ ] Add ANSI theme tokens and built-in presets
- [ ] Implement tasks workspace rendering and navigation
- [ ] Implement async read/write actions via `service.Service`
- [ ] Implement Huh add/edit forms
- [ ] Implement search/filter controls
- [ ] Implement calendar view
- [ ] Implement system view (sync + daemon)
- [ ] Add tests for model/actions/filters
- [ ] Update README usage to document root-launched TUI (`ugh`)
