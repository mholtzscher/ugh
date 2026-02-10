package service

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
	State    string
	Project  string
	Context  string
	Search   string
	DueOnly  bool
	DueOn    string // Date in YYYY-MM-DD format for exact due date matching
	ID       int64  // Specific task ID to fetch (0 means no ID filter)
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
