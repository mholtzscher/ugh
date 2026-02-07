//nolint:testpackage // Tests validate internal navigation helpers directly.
package tui

import (
	"testing"

	"github.com/mholtzscher/ugh/internal/store"
)

func TestBuildNavItemsIncludesAllState(t *testing.T) {
	counts := map[store.State]int64{
		store.StateInbox:   1,
		store.StateNow:     2,
		store.StateWaiting: 0,
		store.StateLater:   3,
		store.StateDone:    4,
	}

	items := buildNavItems(counts, nil, nil)
	if len(items) < 1 {
		t.Fatal("expected navigation items")
	}
	if items[0].label != "All" {
		t.Fatalf("expected first item to be All, got %q", items[0].label)
	}
	if items[0].count != 10 {
		t.Fatalf("expected All count 10, got %d", items[0].count)
	}
}

func TestSelectedScopeTextPrecedence(t *testing.T) {
	f := listFilters{state: "now", project: "work", context: "home"}
	got := selectedScopeText(f)
	if got != "project:work" {
		t.Fatalf("expected project scope precedence, got %q", got)
	}
}

func TestNextTagFilterValue(t *testing.T) {
	tags := []store.NameCount{{Name: "a"}, {Name: "b"}}

	next, ok := nextTagFilterValue("", tags)
	if !ok || next != "a" {
		t.Fatalf("expected first tag, got %q (ok=%v)", next, ok)
	}

	next, ok = nextTagFilterValue("a", tags)
	if !ok || next != "b" {
		t.Fatalf("expected second tag, got %q (ok=%v)", next, ok)
	}

	next, ok = nextTagFilterValue("b", tags)
	if !ok || next != "" {
		t.Fatalf("expected clear cycle, got %q (ok=%v)", next, ok)
	}
}
