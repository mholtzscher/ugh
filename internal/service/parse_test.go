//nolint:testpackage // Tests verify unexported parsing helpers directly.
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mholtzscher/ugh/internal/store"
)

func TestNormalizeState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    store.State
		wantErr bool
	}{
		{name: "empty defaults to inbox", input: "", want: store.StateInbox},
		{name: "todo alias maps to inbox", input: "todo", want: store.StateInbox},
		{name: "trim and lowercase", input: "  NOW ", want: store.StateNow},
		{name: "invalid state", input: "bogus", wantErr: true},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeState(tc.input)
			if tc.wantErr {
				require.Error(t, err, "normalizeState(%q) should return error", tc.input)
				return
			}

			require.NoError(t, err, "normalizeState(%q) error", tc.input)
			assert.Equal(t, tc.want, got, "normalizeState(%q) mismatch", tc.input)
		})
	}
}

func TestParseDay(t *testing.T) {
	t.Parallel()

	got, err := parseDay("2026-02-05")
	require.NoError(t, err, "parseDay(valid) error")
	require.NotNil(t, got, "parseDay(valid) returned nil date")
	assert.Equal(t, "2026-02-05", got.Format("2006-01-02"), "parseDay(valid) date mismatch")
	assert.Equal(t, "UTC", got.Location().String(), "parseDay(valid) location mismatch")

	_, err = parseDay("02/05/2026")
	require.Error(t, err, "parseDay(invalid) should return error")
}

func TestParseMetaFlags(t *testing.T) {
	t.Parallel()

	meta, err := parseMetaFlags([]string{"a:1", " b : 2 "})
	require.NoError(t, err, "parseMetaFlags(valid) error")
	require.Len(t, meta, 2, "parseMetaFlags(valid) len mismatch")
	assert.Equal(t, "1", meta["a"], "parseMetaFlags(valid) a value mismatch")
	assert.Equal(t, "2", meta["b"], "parseMetaFlags(valid) b value mismatch")

	_, err = parseMetaFlags([]string{"missing-separator"})
	require.Error(t, err, "parseMetaFlags(invalid) should return error")
}
