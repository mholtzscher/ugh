package registry

import (
	"testing"

	"github.com/urfave/cli/v3"
)

func TestBuildTree(t *testing.T) {
	r := New()
	if err := r.AddAll(
		Spec{ID: "parent", Build: func() *cli.Command { return &cli.Command{Name: "parent"} }},
		Spec{ID: "child.one", ParentID: "parent", Build: func() *cli.Command { return &cli.Command{Name: "one"} }},
		Spec{ID: "child.two", ParentID: "parent", Build: func() *cli.Command { return &cli.Command{Name: "two"} }},
		Spec{ID: "top", Build: func() *cli.Command { return &cli.Command{Name: "top"} }},
	); err != nil {
		t.Fatalf("add specs: %v", err)
	}

	commands, err := r.Build()
	if err != nil {
		t.Fatalf("build commands: %v", err)
	}

	if len(commands) != 2 {
		t.Fatalf("expected 2 root commands, got %d", len(commands))
	}
	if commands[0].Name != "parent" {
		t.Fatalf("expected first root command to be parent, got %q", commands[0].Name)
	}
	if commands[1].Name != "top" {
		t.Fatalf("expected second root command to be top, got %q", commands[1].Name)
	}

	if len(commands[0].Commands) != 2 {
		t.Fatalf("expected 2 child commands, got %d", len(commands[0].Commands))
	}
	if commands[0].Commands[0].Name != "one" {
		t.Fatalf("expected first child to be one, got %q", commands[0].Commands[0].Name)
	}
	if commands[0].Commands[1].Name != "two" {
		t.Fatalf("expected second child to be two, got %q", commands[0].Commands[1].Name)
	}
}

func TestDuplicateID(t *testing.T) {
	r := New()
	if err := r.Add(Spec{ID: "one", Build: func() *cli.Command { return &cli.Command{Name: "one"} }}); err != nil {
		t.Fatalf("first add failed: %v", err)
	}
	if err := r.Add(Spec{ID: "one", Build: func() *cli.Command { return &cli.Command{Name: "other"} }}); err == nil {
		t.Fatal("expected duplicate id error")
	}
}

func TestUnknownParent(t *testing.T) {
	r := New()
	if err := r.Add(Spec{ID: "child", ParentID: "missing", Build: func() *cli.Command { return &cli.Command{Name: "child"} }}); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if _, err := r.Build(); err == nil {
		t.Fatal("expected unknown parent error")
	}
}

func TestSiblingTokenConflict(t *testing.T) {
	r := New()
	if err := r.AddAll(
		Spec{ID: "first", Build: func() *cli.Command { return &cli.Command{Name: "first", Aliases: []string{"dup"}} }},
		Spec{ID: "second", Build: func() *cli.Command { return &cli.Command{Name: "dup"} }},
	); err != nil {
		t.Fatalf("add specs: %v", err)
	}
	if _, err := r.Build(); err == nil {
		t.Fatal("expected sibling token conflict error")
	}
}

func TestPreconfiguredChildrenConflict(t *testing.T) {
	r := New()
	if err := r.AddAll(
		Spec{ID: "parent", Build: func() *cli.Command {
			return &cli.Command{Name: "parent", Commands: []*cli.Command{{Name: "legacy"}}}
		}},
		Spec{ID: "child", ParentID: "parent", Build: func() *cli.Command { return &cli.Command{Name: "child"} }},
	); err != nil {
		t.Fatalf("add specs: %v", err)
	}
	if _, err := r.Build(); err == nil {
		t.Fatal("expected preconfigured children conflict error")
	}
}

func TestEmptyBuiltCommandName(t *testing.T) {
	r := New()
	if err := r.Add(Spec{ID: "bad", Build: func() *cli.Command { return &cli.Command{Name: " "} }}); err != nil {
		t.Fatalf("add spec: %v", err)
	}
	if _, err := r.Build(); err == nil {
		t.Fatal("expected empty name error")
	}
}

func TestCaseInsensitiveSiblingConflict(t *testing.T) {
	r := New()
	if err := r.AddAll(
		Spec{ID: "one", Build: func() *cli.Command { return &cli.Command{Name: "List"} }},
		Spec{ID: "two", Build: func() *cli.Command { return &cli.Command{Name: "list"} }},
	); err != nil {
		t.Fatalf("add specs: %v", err)
	}
	if _, err := r.Build(); err == nil {
		t.Fatal("expected case-insensitive sibling conflict error")
	}
}

func TestOwnTokenConflict(t *testing.T) {
	r := New()
	if err := r.Add(Spec{ID: "dup", Build: func() *cli.Command { return &cli.Command{Name: "list", Aliases: []string{"list"}} }}); err != nil {
		t.Fatalf("add spec: %v", err)
	}
	if _, err := r.Build(); err == nil {
		t.Fatal("expected own token conflict error")
	}
}
