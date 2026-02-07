//nolint:testpackage // Tests verify unexported parsing helpers directly.
package service

import (
	"testing"

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
		{name: "trim and lowercase", input: "  NOW ", want: store.StateNow},
		{name: "invalid state", input: "bogus", wantErr: true},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeState(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("normalizeState(%q) error = nil, want error", tc.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("normalizeState(%q) error = %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("normalizeState(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseDay(t *testing.T) {
	t.Parallel()

	got, err := parseDay("2026-02-05")
	if err != nil {
		t.Fatalf("parseDay(valid) error = %v", err)
	}
	if got == nil {
		t.Fatal("parseDay(valid) returned nil date")
	}
	if got.Format("2006-01-02") != "2026-02-05" {
		t.Fatalf("parseDay(valid) date = %s, want 2026-02-05", got.Format("2006-01-02"))
	}
	if got.Location().String() != "UTC" {
		t.Fatalf("parseDay(valid) location = %s, want UTC", got.Location())
	}

	_, err = parseDay("02/05/2026")
	if err == nil {
		t.Fatal("parseDay(invalid) error = nil, want error")
	}
}

func TestParseMetaFlags(t *testing.T) {
	t.Parallel()

	meta, err := parseMetaFlags([]string{"a:1", " b : 2 "})
	if err != nil {
		t.Fatalf("parseMetaFlags(valid) error = %v", err)
	}
	if len(meta) != 2 {
		t.Fatalf("parseMetaFlags(valid) len = %d, want 2", len(meta))
	}
	if meta["a"] != "1" {
		t.Fatalf("parseMetaFlags(valid) a = %q, want %q", meta["a"], "1")
	}
	if meta["b"] != "2" {
		t.Fatalf("parseMetaFlags(valid) b = %q, want %q", meta["b"], "2")
	}

	_, err = parseMetaFlags([]string{"missing-separator"})
	if err == nil {
		t.Fatal("parseMetaFlags(invalid) error = nil, want error")
	}
}
