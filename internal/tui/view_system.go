package tui

import (
	"fmt"
	"strings"
	"time"
)

func (m model) viewSystem() string {
	width := m.layout.listWidth
	if !m.layout.narrow {
		width = m.layout.navWidth + m.layout.listWidth + m.layout.detailWidth
	}

	return m.styles.panel.Width(width).Height(m.layout.bodyHeight).Render(m.renderSystem())
}

func (m model) renderSystem() string {
	lines := []string{
		m.styles.title.Render("SYSTEM"),
		"",
		m.styles.title.Render("Sync"),
		m.styles.muted.Render("s: sync   l: pull   p: push   o: config   r: refresh"),
	}

	switch {
	case m.syncStatusErr != "":
		lines = append(lines, m.styles.warning.Render("sync status unavailable: "+m.syncStatusErr))
	case m.syncStatus != nil:
		lines = append(lines,
			fmt.Sprintf("last pull: %s", formatUnixTime(m.syncStatus.LastPullUnixTime)),
			fmt.Sprintf("last push: %s", formatUnixTime(m.syncStatus.LastPushUnixTime)),
			fmt.Sprintf("pending changes: %d", m.syncStatus.PendingChanges),
			fmt.Sprintf("network sent: %d bytes", m.syncStatus.NetworkSentBytes),
			fmt.Sprintf("network recv: %d bytes", m.syncStatus.NetworkRecvBytes),
			fmt.Sprintf("revision: %s", emptyDash(m.syncStatus.Revision)),
		)
	default:
		lines = append(lines, m.styles.muted.Render("sync status loading..."))
	}

	lines = append(lines,
		"",
		m.styles.title.Render("Config"),
		fmt.Sprintf("db.sync_url: %s", emptyDash(m.appConfig.DB.SyncURL)),
		fmt.Sprintf("db.auth_token: %s", maskToken(m.appConfig.DB.AuthToken)),
		fmt.Sprintf("db.sync_on_write: %t", m.appConfig.DB.SyncOnWrite),
		fmt.Sprintf("ui.theme: %s", emptyDash(m.appConfig.UI.Theme)),
		m.styles.muted.Render("press o to edit config values"),
	)

	lines = append(lines,
		"",
		m.styles.title.Render("Daemon"),
		m.styles.muted.Render("I/U install/uninstall   S/T/R start/stop/restart   L logs hint"),
	)

	switch {
	case m.daemonStatusErr != "":
		lines = append(lines, m.styles.warning.Render("daemon status unavailable: "+m.daemonStatusErr))
	case m.daemonStatus != nil:
		lines = append(lines,
			fmt.Sprintf("manager: %s", emptyDash(m.daemonManager)),
			fmt.Sprintf("installed: %t", m.daemonStatus.Installed),
			fmt.Sprintf("running: %t", m.daemonStatus.Running),
			fmt.Sprintf("pid: %d", m.daemonStatus.PID),
			fmt.Sprintf("service path: %s", emptyDash(m.daemonStatus.ServicePath)),
		)
		if m.daemonLogPath != "" {
			lines = append(lines, fmt.Sprintf("log path: %s", m.daemonLogPath))
		} else {
			lines = append(lines, "log path: managed by service manager")
		}
	default:
		lines = append(lines, m.styles.muted.Render("daemon status loading..."))
	}

	lines = append(lines,
		"",
		m.styles.muted.Render("CLI equivalents:"),
		"ugh daemon status | start | stop | restart | install | uninstall | logs",
	)

	return strings.Join(lines, "\n")
}

func formatUnixTime(unix int64) string {
	if unix <= 0 {
		return "never"
	}
	return time.Unix(unix, 0).Local().Format(time.RFC3339)
}
