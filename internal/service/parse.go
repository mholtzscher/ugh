package service

import (
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/domain"
	"github.com/mholtzscher/ugh/internal/store"
)

func parseMetaFlags(meta []string) (map[string]string, error) {
	result := map[string]string{}
	for _, m := range meta {
		k, v, ok := strings.Cut(m, domain.MetaSeparatorColon)
		if !ok {
			return nil, domain.InvalidMetaFormatError(m)
		}
		result[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return result, nil
}

func parseDay(value string) (*time.Time, error) {
	parsed, err := time.Parse(domain.DateLayoutYYYYMMDD, value)
	if err != nil {
		return nil, domain.InvalidDateFormatError(value)
	}
	utc := parsed.UTC()
	return &utc, nil
}

func normalizeState(value string) (store.State, error) {
	normalized, err := domain.NormalizeState(value)
	if err != nil {
		return "", err
	}
	return store.State(normalized), nil
}
