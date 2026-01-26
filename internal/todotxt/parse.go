package todotxt

import (
	"strings"
	"time"
)

func ParseLine(line string) Parsed {
	line = strings.TrimSpace(line)
	parsed := Parsed{Meta: map[string]string{}}
	if line == "" {
		return parsed
	}

	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return parsed
	}

	idx := 0
	if tokens[idx] == "x" {
		parsed.Done = true
		idx++
		if idx < len(tokens) && isDate(tokens[idx]) {
			parsed.CompletionDate = mustParseDate(tokens[idx])
			idx++
		}
		if idx < len(tokens) && isDate(tokens[idx]) {
			parsed.CreationDate = mustParseDate(tokens[idx])
			idx++
		}
	} else if isPriority(tokens[idx]) {
		parsed.Priority = tokens[idx][1:2]
		idx++
		if idx < len(tokens) && isDate(tokens[idx]) {
			parsed.CreationDate = mustParseDate(tokens[idx])
			idx++
		}
	} else if isDate(tokens[idx]) {
		parsed.CreationDate = mustParseDate(tokens[idx])
		idx++
	}

	descTokens := make([]string, 0, len(tokens)-idx)
	unknownTokens := make([]string, 0)
	for ; idx < len(tokens); idx++ {
		tok := tokens[idx]
		switch {
		case isProject(tok):
			parsed.Projects = append(parsed.Projects, tok[1:])
		case isContext(tok):
			parsed.Contexts = append(parsed.Contexts, tok[1:])
		case isMeta(tok):
			key, val := splitMeta(tok)
			if key != "" {
				parsed.Meta[key] = val
			} else {
				unknownTokens = append(unknownTokens, tok)
			}
		case looksSpecial(tok):
			unknownTokens = append(unknownTokens, tok)
		default:
			descTokens = append(descTokens, tok)
		}
	}

	parsed.Description = strings.Join(descTokens, " ")
	parsed.Unknown = unknownTokens
	return parsed
}

func isPriority(token string) bool {
	return len(token) == 3 && token[0] == '(' && token[2] == ')' && token[1] >= 'A' && token[1] <= 'Z'
}

func isDate(token string) bool {
	if len(token) != 10 {
		return false
	}
	_, err := time.Parse("2006-01-02", token)
	return err == nil
}

func mustParseDate(token string) *time.Time {
	val, err := time.Parse("2006-01-02", token)
	if err != nil {
		return nil
	}
	utc := val.UTC()
	return &utc
}

func isProject(token string) bool {
	return len(token) > 1 && token[0] == '+' && isTokenValue(token[1:])
}

func isContext(token string) bool {
	return len(token) > 1 && token[0] == '@' && isTokenValue(token[1:])
}

func isMeta(token string) bool {
	if strings.Contains(token, "://") {
		return false
	}
	key, val := splitMeta(token)
	return key != "" && val != ""
}

func splitMeta(token string) (string, string) {
	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	key := parts[0]
	val := parts[1]
	if key == "" || val == "" {
		return "", ""
	}
	if !isMetaKey(key) {
		return "", ""
	}
	return key, val
}

func isMetaKey(key string) bool {
	if key == "" {
		return false
	}
	if key[0] < 'A' || (key[0] > 'Z' && key[0] < 'a') || key[0] > 'z' {
		return false
	}
	for i := 1; i < len(key); i++ {
		ch := key[i]
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			continue
		}
		return false
	}
	return true
}

func isTokenValue(val string) bool {
	for i := 0; i < len(val); i++ {
		ch := val[i]
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			continue
		}
		return false
	}
	return val != ""
}

func looksSpecial(token string) bool {
	if token == "x" {
		return true
	}
	if strings.HasPrefix(token, "+") || strings.HasPrefix(token, "@") {
		return true
	}
	if strings.Contains(token, ":") {
		return true
	}
	if strings.HasPrefix(token, "(") {
		return true
	}
	return false
}
