package service

import (
	"reflect"
	"testing"
)

func TestCopyMeta(t *testing.T) {
	t.Parallel()

	if got := copyMeta(nil); got != nil {
		t.Fatalf("copyMeta(nil) = %#v, want nil", got)
	}

	original := map[string]string{"a": "1", "b": "2"}
	clone := copyMeta(original)
	if !reflect.DeepEqual(clone, original) {
		t.Fatalf("copyMeta(original) = %#v, want %#v", clone, original)
	}

	clone["a"] = "changed"
	if original["a"] != "1" {
		t.Fatalf("copyMeta should not mutate original map, got %q", original["a"])
	}
}

func TestContainsString(t *testing.T) {
	t.Parallel()

	values := []string{"alpha", "beta", "gamma"}
	if !containsString(values, "beta") {
		t.Fatal("containsString should return true for existing value")
	}
	if containsString(values, "delta") {
		t.Fatal("containsString should return false for missing value")
	}
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
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("removeStrings(%v, %v) = %v, want %v", tc.input, tc.toRemove, got, tc.want)
			}
		})
	}
}
