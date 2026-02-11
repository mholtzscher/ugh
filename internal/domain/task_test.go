package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mholtzscher/ugh/internal/domain"
)

func TestNormalizeState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty defaults to inbox", input: "", want: domain.TaskStateInbox},
		{name: "todo alias maps to inbox", input: "todo", want: domain.TaskStateInbox},
		{name: "trim and lowercase", input: "  NOW ", want: domain.TaskStateNow},
		{name: "valid state passes through", input: "later", want: domain.TaskStateLater},
		{name: "invalid state", input: "bogus", wantErr: true},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NormalizeState(tc.input)
			if tc.wantErr {
				require.Error(t, err, "NormalizeState(%q) should return error", tc.input)
				return
			}

			require.NoError(t, err, "NormalizeState(%q) error", tc.input)
			assert.Equal(t, tc.want, got, "NormalizeState(%q) mismatch", tc.input)
		})
	}
}
