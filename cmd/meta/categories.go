package meta

type Category string

const (
	ListsCategory    Category = "Lists"
	TasksCategory    Category = "Tasks"
	ProjectsCategory Category = "Projects & Contexts"
	BackupCategory   Category = "Backup"
	SyncCategory     Category = "Sync"
	SystemCategory   Category = "System"
)

func (c Category) String() string {
	return string(c)
}
