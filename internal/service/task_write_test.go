//nolint:testpackage // Tests validate unexported helper behavior.
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyMeta(t *testing.T) {
	t.Parallel()

	assert.Nil(t, copyMeta(nil), "copyMeta(nil) should return nil")

	original := map[string]string{"a": "1", "b": "2"}
	clone := copyMeta(original)
	assert.Equal(t, original, clone, "copyMeta(original) should return equal map")

	clone["a"] = "changed"
	assert.Equal(t, "1", original["a"], "copyMeta should not mutate original map")
}

func TestContainsString(t *testing.T) {
	t.Parallel()

	values := []string{"alpha", "beta", "gamma"}
	assert.True(t, containsString(values, "beta"), "containsString should return true for existing value")
	assert.False(t, containsString(values, "delta"), "containsString should return false for missing value")
}

func TestRemoveStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		toRemove []string
		want     []string
	}{
		{
			name:     "remove some values",
			input:    []string{"a", "b", "c", "d"},
			toRemove: []string{"b", "d"},
			want:     []string{"a", "c"},
		},
		{
			name:     "remove repeated target",
			input:    []string{"x", "y", "y", "z"},
			toRemove: []string{"y"},
			want:     []string{"x", "z"},
		},
		{
			name:     "nothing removed",
			input:    []string{"one", "two"},
			toRemove: nil,
			want:     []string{"one", "two"},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := removeStrings(tc.input, tc.toRemove)
			assert.Equal(t, tc.want, got, "removeStrings(%v, %v) mismatch", tc.input, tc.toRemove)
		})
	}
}
