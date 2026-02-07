package tui

import (
	"fmt"

	"github.com/mholtzscher/ugh/internal/store"
)

type paneFocus int

const (
	focusNav paneFocus = iota
	focusList
	focusDetail
)

const focusPaneCount = 3

type navItemKind int

const (
	navState navItemKind = iota
	navProject
	navContext
)

type navItem struct {
	kind  navItemKind
	label string
	value string
	count int64
}

func buildNavItems(
	stateCounts map[store.State]int64,
	projects []store.NameCount,
	contexts []store.NameCount,
) []navItem {
	items := []navItem{
		{kind: navState, label: "All", value: "", count: totalStateCount(stateCounts)},
		{kind: navState, label: "Inbox", value: string(store.StateInbox), count: stateCounts[store.StateInbox]},
		{kind: navState, label: "Now", value: string(store.StateNow), count: stateCounts[store.StateNow]},
		{kind: navState, label: "Waiting", value: string(store.StateWaiting), count: stateCounts[store.StateWaiting]},
		{kind: navState, label: "Later", value: string(store.StateLater), count: stateCounts[store.StateLater]},
		{kind: navState, label: "Done", value: string(store.StateDone), count: stateCounts[store.StateDone]},
	}

	for _, project := range projects {
		items = append(items, navItem{kind: navProject, label: project.Name, value: project.Name, count: project.Count})
	}

	for _, context := range contexts {
		items = append(items, navItem{kind: navContext, label: context.Name, value: context.Name, count: context.Count})
	}

	return items
}

func totalStateCount(counts map[store.State]int64) int64 {
	return counts[store.StateInbox] +
		counts[store.StateNow] +
		counts[store.StateWaiting] +
		counts[store.StateLater] +
		counts[store.StateDone]
}

func selectedScopeText(filters listFilters) string {
	if filters.project != "" {
		return fmt.Sprintf("project:%s", filters.project)
	}
	if filters.context != "" {
		return fmt.Sprintf("context:%s", filters.context)
	}
	if filters.state != "" {
		return fmt.Sprintf("state:%s", filters.state)
	}
	return "scope:all"
}

func nextFocus(current paneFocus) paneFocus {
	return (current + 1) % focusPaneCount
}

func prevFocus(current paneFocus) paneFocus {
	if current == 0 {
		return focusPaneCount - 1
	}
	return current - 1
}
