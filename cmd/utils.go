package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseIDs(args []string) ([]int64, error) {
	ids := make([]int64, 0, len(args))
	for _, arg := range args {
		val, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid id %q", arg)
		}
		ids = append(ids, val)
	}
	if len(ids) == 0 {
		return nil, errors.New("at least one id is required")
	}
	return ids, nil
}

func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, fmt.Errorf("invalid date %q (expected YYYY-MM-DD)", value)
	}
	utc := parsed.UTC()
	return &utc, nil
}

func parseMetaFlags(values []string) (map[string]string, error) {
	if len(values) == 0 {
		return nil, nil
	}
	meta := make(map[string]string)
	for _, item := range values {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		sep := strings.IndexAny(item, ":=")
		if sep <= 0 || sep >= len(item)-1 {
			return nil, fmt.Errorf("invalid meta %q (use key:value)", item)
		}
		key := item[:sep]
		value := item[sep+1:]
		meta[key] = value
	}
	return meta, nil
}

func normalizePriority(value string) string {
	value = strings.TrimSpace(value)
	if len(value) == 3 && value[0] == '(' && value[2] == ')' {
		value = value[1:2]
	}
	if value == "" {
		return ""
	}
	value = strings.ToUpper(value)
	if value[0] < 'A' || value[0] > 'Z' {
		return ""
	}
	return value[:1]
}
