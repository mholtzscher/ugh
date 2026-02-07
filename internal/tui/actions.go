package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mholtzscher/ugh/internal/config"
	daemonservice "github.com/mholtzscher/ugh/internal/daemon/service"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

func orderedTaskStates() []store.State {
	return []store.State{store.StateInbox, store.StateNow, store.StateWaiting, store.StateLater, store.StateDone}
}

func loadTasksCmd(svc service.Service, filters listFilters) tea.Cmd {
	req := filters.toListTasksRequest()
	return func() tea.Msg {
		tasks, err := svc.ListTasks(context.Background(), req)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func loadCalendarTasksCmd(svc service.Service, filters listFilters) tea.Cmd {
	req := filters.toListTasksRequest()
	req.DueOnly = true

	return func() tea.Msg {
		tasks, err := svc.ListTasks(context.Background(), req)
		return calendarTasksLoadedMsg{tasks: tasks, err: err}
	}
}

func loadTagCountsCmd(svc service.Service, filters listFilters) tea.Cmd {
	tagReq := filters.toListTagsRequest()
	return func() tea.Msg {
		projects, err := svc.ListProjects(context.Background(), tagReq)
		if err != nil {
			return tagCountsLoadedMsg{err: err}
		}
		contexts, err := svc.ListContexts(context.Background(), tagReq)
		if err != nil {
			return tagCountsLoadedMsg{err: err}
		}
		return tagCountsLoadedMsg{projects: projects, contexts: contexts}
	}
}

func loadStateCountsCmd(svc service.Service, filters listFilters) tea.Cmd {
	req := filters.toListTasksRequest()
	req.State = ""
	req.Project = ""
	req.Context = ""
	req.Search = ""

	return func() tea.Msg {
		states := orderedTaskStates()
		counts := make(map[store.State]int64, len(states))
		for _, state := range states {
			stateReq := req
			stateReq.State = string(state)
			tasks, err := svc.ListTasks(context.Background(), stateReq)
			if err != nil {
				return stateCountsLoadedMsg{err: err}
			}
			counts[state] = int64(len(tasks))
		}
		return stateCountsLoadedMsg{counts: counts}
	}
}

func refreshDataCmd(svc service.Service, filters listFilters) tea.Cmd {
	return tea.Batch(
		loadTasksCmd(svc, filters),
		loadCalendarTasksCmd(svc, filters),
		loadTagCountsCmd(svc, filters),
		loadStateCountsCmd(svc, filters),
		loadSyncStatusCmd(svc),
		loadDaemonStatusCmd(),
	)
}

func loadSyncStatusCmd(svc service.Service) tea.Cmd {
	return func() tea.Msg {
		status, err := svc.SyncStatus(context.Background())
		return syncStatusLoadedMsg{status: status, err: err}
	}
}

func loadDaemonStatusCmd() tea.Cmd {
	return func() tea.Msg {
		mgr, err := daemonservice.Detect()
		if err != nil {
			return daemonStatusLoadedMsg{err: err}
		}

		status, err := mgr.Status()
		if err != nil {
			return daemonStatusLoadedMsg{managerName: mgr.Name(), err: err}
		}

		return daemonStatusLoadedMsg{
			managerName: mgr.Name(),
			status:      &status,
			logPath:     mgr.LogPath(),
		}
	}
}

func setDoneCmd(svc service.Service, id int64, done bool) tea.Cmd {
	return func() tea.Msg {
		count, err := svc.SetDone(context.Background(), []int64{id}, done)
		if err != nil {
			return actionResultMsg{err: err}
		}
		if count == 0 {
			return actionResultMsg{err: fmt.Errorf("task %d not found", id)}
		}
		if done {
			return actionResultMsg{status: fmt.Sprintf("task #%d marked done", id)}
		}
		return actionResultMsg{status: fmt.Sprintf("task #%d reopened", id)}
	}
}

func setStateCmd(svc service.Service, id int64, state string) tea.Cmd {
	return func() tea.Msg {
		_, err := svc.UpdateTask(context.Background(), service.UpdateTaskRequest{ID: id, State: &state})
		if err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: fmt.Sprintf("task #%d moved to %s", id, state)}
	}
}

func deleteTaskCmd(svc service.Service, id int64) tea.Cmd {
	return func() tea.Msg {
		count, err := svc.DeleteTasks(context.Background(), []int64{id})
		if err != nil {
			return actionResultMsg{err: err}
		}
		if count == 0 {
			return actionResultMsg{err: fmt.Errorf("task %d not found", id)}
		}
		return actionResultMsg{status: fmt.Sprintf("task #%d deleted", id)}
	}
}

func createTaskCmd(svc service.Service, req service.CreateTaskRequest) tea.Cmd {
	return func() tea.Msg {
		task, err := svc.CreateTask(context.Background(), req)
		if err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: fmt.Sprintf("created task #%d", task.ID)}
	}
}

func fullUpdateTaskCmd(svc service.Service, req service.FullUpdateTaskRequest) tea.Cmd {
	return func() tea.Msg {
		task, err := svc.FullUpdateTask(context.Background(), req)
		if err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: fmt.Sprintf("updated task #%d", task.ID)}
	}
}

func syncPullCmd(svc service.Service) tea.Cmd {
	return func() tea.Msg {
		if err := svc.Sync(context.Background()); err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: "pulled changes from remote"}
	}
}

func syncPushCmd(svc service.Service) tea.Cmd {
	return func() tea.Msg {
		if err := svc.Push(context.Background()); err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: "pushed changes to remote"}
	}
}

func syncAllCmd(svc service.Service) tea.Cmd {
	return func() tea.Msg {
		if err := svc.Sync(context.Background()); err != nil {
			return actionResultMsg{err: err}
		}
		if err := svc.Push(context.Background()); err != nil {
			return actionResultMsg{err: err}
		}
		return actionResultMsg{status: "synced with remote"}
	}
}

func saveConfigCmd(path string, cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		err := config.Save(path, cfg)
		return configSavedMsg{cfg: cfg, err: err, path: path}
	}
}

func daemonInstallCmd(configPath string) tea.Cmd {
	return func() tea.Msg {
		mgr, err := daemonservice.Detect()
		if err != nil {
			return actionResultMsg{err: fmt.Errorf("detect service manager: %w", err)}
		}

		binaryPath, err := executablePath()
		if err != nil {
			return actionResultMsg{err: fmt.Errorf("get binary path: %w", err)}
		}

		err = mgr.Install(daemonservice.InstallConfig{BinaryPath: binaryPath, ConfigPath: configPath})
		if err != nil {
			return actionResultMsg{err: mapDaemonActionError("install", err)}
		}

		status, _ := mgr.Status()
		return actionResultMsg{status: fmt.Sprintf("daemon installed at %s", status.ServicePath)}
	}
}

func daemonUninstallCmd() tea.Cmd {
	return func() tea.Msg {
		mgr, err := daemonservice.Detect()
		if err != nil {
			return actionResultMsg{err: fmt.Errorf("detect service manager: %w", err)}
		}

		status, _ := mgr.Status()
		err = mgr.Uninstall()
		if err != nil {
			return actionResultMsg{err: mapDaemonActionError("uninstall", err)}
		}
		return actionResultMsg{status: fmt.Sprintf("daemon uninstalled from %s", status.ServicePath)}
	}
}

func daemonStartCmd() tea.Cmd {
	return func() tea.Msg {
		mgr, err := daemonservice.Detect()
		if err != nil {
			return actionResultMsg{err: fmt.Errorf("detect service manager: %w", err)}
		}

		err = mgr.Start()
		if err != nil {
			return actionResultMsg{err: mapDaemonActionError("start", err)}
		}
		return actionResultMsg{status: "daemon started"}
	}
}

func daemonStopCmd() tea.Cmd {
	return func() tea.Msg {
		mgr, err := daemonservice.Detect()
		if err != nil {
			return actionResultMsg{err: fmt.Errorf("detect service manager: %w", err)}
		}

		err = mgr.Stop()
		if err != nil {
			return actionResultMsg{err: mapDaemonActionError("stop", err)}
		}
		return actionResultMsg{status: "daemon stopped"}
	}
}

func daemonRestartCmd() tea.Cmd {
	return func() tea.Msg {
		mgr, err := daemonservice.Detect()
		if err != nil {
			return actionResultMsg{err: fmt.Errorf("detect service manager: %w", err)}
		}

		_ = mgr.Stop()
		err = mgr.Start()
		if err != nil {
			return actionResultMsg{err: mapDaemonActionError("restart", err)}
		}
		return actionResultMsg{status: "daemon restarted"}
	}
}

func mapDaemonActionError(action string, err error) error {
	switch {
	case errors.Is(err, daemonservice.ErrNotInstalled):
		return errors.New("service not installed - run install first")
	case errors.Is(err, daemonservice.ErrAlreadyInstalled):
		return errors.New("service already installed")
	case errors.Is(err, daemonservice.ErrNotRunning):
		return errors.New("daemon is not running")
	default:
		return fmt.Errorf("%s daemon service: %w", action, err)
	}
}

func executablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return filepath.Abs(exe)
	}
	return resolved, nil
}
