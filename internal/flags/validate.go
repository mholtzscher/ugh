package flags

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

type StringRule func(*cli.Command, string) error

type StringSliceRule func(*cli.Command, []string) error

type BoolRule func(*cli.Command, bool) error

func StringAction(rules ...StringRule) func(context.Context, *cli.Command, string) error {
	return func(_ context.Context, cmd *cli.Command, value string) error {
		for _, rule := range rules {
			if rule == nil {
				continue
			}
			if err := rule(cmd, value); err != nil {
				return err
			}
		}
		return nil
	}
}

func StringSliceAction(rules ...StringSliceRule) func(context.Context, *cli.Command, []string) error {
	return func(_ context.Context, cmd *cli.Command, values []string) error {
		for _, rule := range rules {
			if rule == nil {
				continue
			}
			if err := rule(cmd, values); err != nil {
				return err
			}
		}
		return nil
	}
}

func BoolAction(rules ...BoolRule) func(context.Context, *cli.Command, bool) error {
	return func(_ context.Context, cmd *cli.Command, value bool) error {
		for _, rule := range rules {
			if rule == nil {
				continue
			}
			if err := rule(cmd, value); err != nil {
				return err
			}
		}
		return nil
	}
}

func OneOfCI(fieldName string, allowed ...string) StringRule {
	allowedSet := make(map[string]bool, len(allowed))
	normalizedAllowed := make([]string, 0, len(allowed))
	for _, option := range allowed {
		normalized := normalizeToken(option)
		if normalized == "" || allowedSet[normalized] {
			continue
		}
		allowedSet[normalized] = true
		normalizedAllowed = append(normalizedAllowed, normalized)
	}
	expected := strings.Join(normalizedAllowed, "|")

	return func(_ *cli.Command, value string) error {
		normalized := normalizeToken(value)
		if normalized == "" {
			return nil
		}
		if allowedSet[normalized] {
			return nil
		}
		return fmt.Errorf("invalid %s %q (expected %s)", fieldName, normalized, expected)
	}
}

func OneOfCaseInsensitiveRule(fieldName string, allowed ...string) StringRule {
	return OneOfCI(fieldName, allowed...)
}

func DateLayoutRule(fieldName string, layout string, expected string) StringRule {
	return func(_ *cli.Command, value string) error {
		value = strings.TrimSpace(value)
		if value == "" {
			return nil
		}
		if _, err := time.Parse(layout, value); err != nil {
			return fmt.Errorf("invalid %s format: %s (expected %s)", fieldName, value, expected)
		}
		return nil
	}
}

func EachContainsSeparatorRule(fieldName string, separator string, expected string) StringSliceRule {
	return func(_ *cli.Command, values []string) error {
		for _, value := range values {
			if _, _, ok := strings.Cut(value, separator); !ok {
				return fmt.Errorf("invalid %s format: %s (expected %s)", fieldName, value, expected)
			}
		}
		return nil
	}
}

func MutuallyExclusiveBoolFlagsRule(flagNames ...string) BoolRule {
	return func(cmd *cli.Command, value bool) error {
		if !value {
			return nil
		}
		active := make([]string, 0, len(flagNames))
		for _, flagName := range flagNames {
			if cmd.Bool(flagName) {
				active = append(active, flagName)
			}
		}
		if len(active) < 2 {
			return nil
		}
		left, right := active[0], active[1]
		if left > right {
			left, right = right, left
		}
		if len(flagNames) == 2 {
			return fmt.Errorf("cannot use both --%s and --%s", left, right)
		}
		return pairError("cannot combine", left, right)
	}
}

func BoolRequiresStringOneOfCaseInsensitiveRule(boolFlagName string, stringFlagName string, allowed ...string) BoolRule {
	allowedSet := make(map[string]bool, len(allowed))
	normalizedAllowed := make([]string, 0, len(allowed))
	for _, option := range allowed {
		normalized := normalizeToken(option)
		if normalized == "" || allowedSet[normalized] {
			continue
		}
		allowedSet[normalized] = true
		normalizedAllowed = append(normalizedAllowed, normalized)
	}

	return func(cmd *cli.Command, value bool) error {
		if !value {
			return nil
		}
		selected := normalizeToken(cmd.String(stringFlagName))
		if selected == "" || allowedSet[selected] {
			return nil
		}
		if len(normalizedAllowed) == 1 {
			return fmt.Errorf("cannot combine --%s with --%s other than %s", boolFlagName, stringFlagName, normalizedAllowed[0])
		}
		return fmt.Errorf("cannot combine --%s with --%s outside %s", boolFlagName, stringFlagName, strings.Join(normalizedAllowed, "|"))
	}
}

func pairError(prefix string, left string, right string) error {
	if left > right {
		left, right = right, left
	}
	return fmt.Errorf("%s --%s with --%s", prefix, left, right)
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
