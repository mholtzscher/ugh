package todotxt

import (
	"sort"
	"strings"
)

func Format(parsed Parsed) string {
	parts := make([]string, 0)
	if parsed.Done {
		parts = append(parts, "x")
	}
	if parsed.Priority != "" {
		parts = append(parts, "("+parsed.Priority+")")
	}
	if parsed.Done {
		if parsed.CompletionDate != nil {
			parts = append(parts, parsed.CompletionDate.Format("2006-01-02"))
		}
		if parsed.CreationDate != nil {
			parts = append(parts, parsed.CreationDate.Format("2006-01-02"))
		}
	} else if parsed.CreationDate != nil {
		parts = append(parts, parsed.CreationDate.Format("2006-01-02"))
	}

	if parsed.Description != "" {
		parts = append(parts, strings.Fields(parsed.Description)...)
	}

	projects := append([]string(nil), parsed.Projects...)
	sort.Strings(projects)
	for _, project := range projects {
		parts = append(parts, "+"+project)
	}

	contexts := append([]string(nil), parsed.Contexts...)
	sort.Strings(contexts)
	for _, context := range contexts {
		parts = append(parts, "@"+context)
	}

	if len(parsed.Meta) > 0 {
		keys := make([]string, 0, len(parsed.Meta))
		for key := range parsed.Meta {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			parts = append(parts, key+":"+parsed.Meta[key])
		}
	}

	if len(parsed.Unknown) > 0 {
		parts = append(parts, parsed.Unknown...)
	}

	return strings.Join(parts, " ")
}
