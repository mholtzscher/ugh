# Datetime Display Config

**Type**: Feature
**Effort**: S
**Status**: Implemented

## Problem

Users see timestamps in UTC (e.g., `2024-01-15T14:30:00Z`) requiring mental conversion. Need local timezone display with user's preferred format.

## Scope

- **In scope**: Human/tty output (tables, task details, history, due dates)
- **Out of scope**: JSON output, CLI config commands

## Config

New `[display]` section in `config.toml`:

```toml
[display]
datetime_format = "2006-01-02 15:04"  # Go reference time format
timezone = "local"                     # "local" or IANA like "America/New_York"
```

### Format Examples

| Config | Output |
|--------|--------|
| `"2006-01-02 15:04"` (default) | `2024-01-15 14:30` |
| `"2006-01-02 3:04 PM"` | `2024-01-15 2:30 PM` |
| `"02 Jan 2006 15:04"` | `15 Jan 2024 14:30` |
| `"Jan 2, 2006 3:04 PM"` | `Jan 15, 2024 2:30 PM` |
| `"2006-01-02"` | `2024-01-15` (date-only) |

**Default**: `"2006-01-02 15:04"` (ISO date, 24h time)

## Behavior

| Field | Format Source |
|-------|---------------|
| `CreatedAt`, `UpdatedAt`, `CompletedAt`, `DueOn` | `datetime_format` config |
| JSON output | UTC RFC3339 (unchanged) |

**Due dates**: Since `DueOn` has no time component, format will show `00:00` (midnight) if format includes time. User can set `"2006-01-02"` for date-only.

## Implementation

### D1: Config Struct (S)

**File**: `internal/config/config.go`

```go
type Config struct {
    Version int
    DB      DB
    Daemon  Daemon
    Display Display
}

type Display struct {
    DatetimeFormat string `toml:"datetime_format"` // Go time format, default "2006-01-02 15:04"
    Timezone       string `toml:"timezone"`        // "local" or IANA, default "local"
}
```

**Apply defaults** in `applyDefaults()`:
- `DatetimeFormat = "2006-01-02 15:04"` if empty
- `Timezone = "local"` if empty

### D2: Time Formatter (S)

**File**: `internal/output/time_format.go`

```go
type TimeFormatter struct {
    location *time.Location
    layout   string
}

func NewTimeFormatter(cfg config.Display) *TimeFormatter

func (f *TimeFormatter) Format(t time.Time) string
```

**Implementation**:
1. Load location: `"local"` → `time.Local`, IANA → `time.LoadLocation()`
2. On invalid IANA → fallback to UTC, log warning
3. Format: `t.In(location).Format(f.layout)`

### D3: Output Integration (S)

**File**: `internal/output/output.go`

- Add `TimeFormatter` to `Output` struct
- Replace `formatDateTime()`, `formatDateTimePtr()`, `formatDate()` with formatter calls
- Remove old formatting functions

**File**: `internal/output/human.go`

- Add `TimeFormatter` to `Human` struct
- Replace all hardcoded date/time formatting with formatter:
  - `dateTimeOrDash()` / `dateTimeFromTimeOrDash()`
  - `formatTaskDueDate()`, `formatDetailDate()`
  - History table (line 334)
  - Version diff header (line 351)
- Remove old formatting functions

### D4: Wire Config (S)

**File**: `cmd/root.go`

- Create `TimeFormatter` from config
- Pass to output constructors

## Files Changed

| File | Change |
|------|--------|
| `internal/config/config.go` | Add `Display` struct |
| `internal/output/time_format.go` | NEW |
| `internal/output/output.go` | Use formatter |
| `internal/output/human.go` | Use formatter |
| `cmd/root.go` | Wire config to output |

## Testing

1. **Unit**: `TimeFormatter` with various format strings and timezones
2. **Integration**: Human output matches configured format
3. **Edge cases**: Invalid IANA, invalid format string, unset TZ

## Risks

| Risk | Mitigation |
|------|------------|
| Invalid format string | No validation; Go produces weird output, user's responsibility |
| Invalid IANA timezone | Fallback to UTC, log warning |
| DueOn shows midnight | User can set date-only format |

## Acceptance Criteria

- [ ] Default format shows `2024-01-15 14:30`
- [ ] `"2006-01-02 3:04 PM"` shows `2024-01-15 2:30 PM`
- [ ] `"2006-01-02"` shows date-only (no time)
- [ ] `timezone = "local"` uses system TZ
- [ ] `timezone = "America/New_York"` converts to that TZ
- [ ] JSON output unchanged (UTC RFC3339)
- [ ] Invalid IANA falls back to UTC with warning