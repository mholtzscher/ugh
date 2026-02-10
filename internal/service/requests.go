package service

import "github.com/mholtzscher/ugh/internal/nlp"

type CreateTaskRequest struct {
	Title      string
	Notes      string
	State      string
	Projects   []string
	Contexts   []string
	Meta       []string
	DueOn      string
	WaitingFor string
}

type ListTasksRequest struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
	Filter   nlp.FilterExpr
}

type ListTagsRequest struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
}

type UpdateTaskRequest struct {
	ID              int64
	Title           *string
	Notes           *string
	State           *string
	DueOn           *string
	WaitingFor      *string
	AddProjects     []string
	AddContexts     []string
	SetMeta         map[string]string
	RemoveProjects  []string
	RemoveContexts  []string
	RemoveMetaKeys  []string
	ClearDueOn      bool
	ClearWaitingFor bool
}

type FullUpdateTaskRequest struct {
	ID         int64
	Title      string
	Notes      string
	State      string
	Projects   []string
	Contexts   []string
	Meta       map[string]string
	DueOn      string
	WaitingFor string
}

type SyncStatus struct {
	LastPullUnixTime int64
	LastPushUnixTime int64
	PendingChanges   int64
	NetworkSentBytes int64
	NetworkRecvBytes int64
	Revision         string
}
