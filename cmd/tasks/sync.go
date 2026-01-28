package tasks

import (
	"fmt"

	"github.com/mholtzscher/ugh/internal/output"

	"github.com/urfave/cli/v2"
)

func syncCommand() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "Sync database with remote server",
		Action: func(c *cli.Context) error {
			return runSync(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:      "pull",
				Usage:     "Pull changes from remote server",
				ArgsUsage: "",
				Action:    runSyncPull,
			},
			{
				Name:      "push",
				Usage:     "Push local changes to remote server",
				ArgsUsage: "",
				Action:    runSyncPush,
			},
			{
				Name:      "status",
				Usage:     "Show sync status",
				ArgsUsage: "",
				Action:    runSyncStatus,
			},
		},
	}
}

func runSync(c *cli.Context) error {
	ctx, cancel := deps.WithTimeout(c.Context)
	defer cancel()
	st, err := deps.OpenStore(ctx, c)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	if err := st.Sync(ctx); err != nil {
		return fmt.Errorf("sync: %w", err)
	}
	if err := st.Push(ctx); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	writer := deps.OutputWriter(c)
	return writer.WriteSummary(output.SyncSummary{Action: "sync", Message: "synced with remote"})
}

func runSyncPull(c *cli.Context) error {
	ctx, cancel := deps.WithTimeout(c.Context)
	defer cancel()
	st, err := deps.OpenStore(ctx, c)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	err = st.Sync(ctx)
	if err != nil {
		return fmt.Errorf("sync pull: %w", err)
	}

	writer := deps.OutputWriter(c)
	return writer.WriteSummary(output.SyncSummary{Action: "pull", Message: "pulled changes from remote"})
}

func runSyncPush(c *cli.Context) error {
	ctx, cancel := deps.WithTimeout(c.Context)
	defer cancel()
	st, err := deps.OpenStore(ctx, c)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	if err := st.Push(ctx); err != nil {
		return fmt.Errorf("sync push: %w", err)
	}

	writer := deps.OutputWriter(c)
	return writer.WriteSummary(output.SyncSummary{Action: "push", Message: "pushed changes to remote"})
}

func runSyncStatus(c *cli.Context) error {
	ctx, cancel := deps.WithTimeout(c.Context)
	defer cancel()
	st, err := deps.OpenStore(ctx, c)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	stats, err := st.SyncStats(ctx)
	if err != nil {
		return fmt.Errorf("sync status: %w", err)
	}

	writer := deps.OutputWriter(c)
	return writer.WriteSummary(output.SyncStatusSummary{
		Action:          "status",
		LastPullTime:    stats.LastPullUnixTime,
		LastPushTime:    stats.LastPushUnixTime,
		PendingChanges:  stats.CdcOperations,
		NetworkSent:     stats.NetworkSentBytes,
		NetworkReceived: stats.NetworkReceivedBytes,
		Revision:        stats.Revision,
	})
}
