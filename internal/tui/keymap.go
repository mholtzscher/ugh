package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit            key.Binding
	Help            key.Binding
	NextView        key.Binding
	PrevView        key.Binding
	Refresh         key.Binding
	Up              key.Binding
	Down            key.Binding
	Top             key.Binding
	Bottom          key.Binding
	Select          key.Binding
	CycleCompletion key.Binding
	Search          key.Binding
	ProjectFilter   key.Binding
	ContextFilter   key.Binding
	Add             key.Binding
	Edit            key.Binding
	Done            key.Binding
	Undo            key.Binding
	Inbox           key.Binding
	Now             key.Binding
	Waiting         key.Binding
	Later           key.Binding
	Delete          key.Binding
	Esc             key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit:            key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:            key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		NextView:        key.NewBinding(key.WithKeys("]"), key.WithHelp("]", "next view")),
		PrevView:        key.NewBinding(key.WithKeys("["), key.WithHelp("[", "prev view")),
		Refresh:         key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		Up:              key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/up", "move up")),
		Down:            key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/down", "move down")),
		Top:             key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "top")),
		Bottom:          key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "bottom")),
		Select:          key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
		CycleCompletion: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "todo/all/done")),
		Search:          key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		ProjectFilter:   key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "project filter")),
		ContextFilter:   key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "context filter")),
		Add:             key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Edit:            key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Done:            key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "done")),
		Undo:            key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "undo")),
		Inbox:           key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "inbox")),
		Now:             key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "now")),
		Waiting:         key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "waiting")),
		Later:           key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "later")),
		Delete:          key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete")),
		Esc:             key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear")),
	}
}
