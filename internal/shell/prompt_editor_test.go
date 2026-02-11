//nolint:testpackage // tests cover unexported prompt helpers.
package shell

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellCompleterProjectTag(t *testing.T) {
	t.Parallel()

	completer := &shellCompleter{
		listProjects: func(context.Context) ([]string, error) {
			return []string{"work", "home"}, nil
		},
		listContexts: func(context.Context) ([]string, error) {
			return []string{"phone"}, nil
		},
	}

	suffixes, offset := completer.Do([]rune("add task #wo"), len([]rune("add task #wo")))
	require.Equal(t, 3, offset, "offset should match typed # fragment")
	assert.Contains(t, completionStrings(suffixes), "rk", "#wo should complete to #work")
}

func TestShellCompleterNoSuggestionsInsideOpenQuote(t *testing.T) {
	t.Parallel()

	completer := &shellCompleter{
		listProjects: func(context.Context) ([]string, error) { return []string{"work"}, nil },
		listContexts: func(context.Context) ([]string, error) { return []string{"phone"}, nil },
	}

	suffixes, offset := completer.Do([]rune(`add "email #wo`), len([]rune(`add "email #wo`)))
	assert.Empty(t, suffixes, "should not autocomplete inside an open quote")
	assert.Zero(t, offset, "offset should be zero when no completions are returned")
}

func TestShellCompleterStateValue(t *testing.T) {
	t.Parallel()

	completer := &shellCompleter{}
	suffixes, offset := completer.Do([]rune("find state:n"), len([]rune("find state:n")))

	require.Equal(t, len([]rune("state:n")), offset, "offset should match field fragment")
	assert.Contains(t, completionStrings(suffixes), "ow", "state:n should complete to state:now")
}

func TestShellCompleterContextCommandProject(t *testing.T) {
	t.Parallel()

	completer := &shellCompleter{
		listProjects: func(context.Context) ([]string, error) { return []string{"work"}, nil },
	}

	suffixes, offset := completer.Do([]rune("context #w"), len([]rune("context #w")))
	require.Equal(t, 2, offset, "offset should match #w fragment")
	assert.Contains(t, completionStrings(suffixes), "ork", "context #w should complete to #work")
}

func TestShellPainterHighlightsTokens(t *testing.T) {
	t.Parallel()

	painter := newShellPainter()
	line := `add "milk" #work @phone`
	painted := string(painter.Paint([]rune(line), len([]rune(line))))

	assert.Contains(t, painted, ansiYellow+`"milk"`+ansiReset, "quoted string should be highlighted")
	assert.Contains(t, painted, ansiBlue+"#work"+ansiReset, "project tag should be highlighted")
	assert.Contains(t, painted, ansiGreen+"@phone"+ansiReset, "context tag should be highlighted")
}

func completionStrings(suffixes [][]rune) []string {
	out := make([]string, 0, len(suffixes))
	for _, suffix := range suffixes {
		out = append(out, string(suffix))
	}
	return out
}
