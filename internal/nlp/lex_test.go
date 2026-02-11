package nlp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mholtzscher/ugh/internal/nlp"
)

func TestLexCompleteQuotedStringPrefersQuotedToken(t *testing.T) {
	t.Parallel()

	tokens, err := nlp.Lex(`add "email #hashtag"`)
	require.NoError(t, err, "lex error")

	names := tokenNames(tokens)
	assert.Contains(t, names, "Quoted", "quoted token should be emitted")
	assert.NotContains(
		t,
		names,
		"QuoteStart",
		"stateful quote-start token should not be used for complete quoted strings",
	)
	assert.NotContains(t, names, "ProjectTag", "hash inside quoted string must not be tokenized as project tag")
}

func TestLexUnterminatedQuoteStillTokenizes(t *testing.T) {
	t.Parallel()

	tokens, err := nlp.Lex(`add "email #hashtag`)
	require.NoError(t, err, "lex error")

	names := tokenNames(tokens)
	assert.Contains(t, names, "QuoteStart", "unterminated quoted input should emit QuoteStart")
	assert.Contains(t, names, "StringText", "unterminated quoted input should emit StringText")
	assert.NotContains(t, names, "ProjectTag", "hash inside open quote must not be tokenized as project tag")
}

func TestLexTagPrefixes(t *testing.T) {
	t.Parallel()

	tokens, err := nlp.Lex(`add task # @`)
	require.NoError(t, err, "lex error")

	names := tokenNames(tokens)
	assert.Contains(t, names, "ProjectTagPrefix", "standalone # should emit project tag prefix token")
	assert.Contains(t, names, "ContextTagPrefix", "standalone @ should emit context tag prefix token")
}

func tokenNames(tokens []nlp.LexToken) []string {
	names := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		names = append(names, tok.Name)
	}
	return names
}
