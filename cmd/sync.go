package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
)

const (
	syncStatusTabWidth = 2
	syncStatusPadding  = 2
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var (
	syncPullCmd = &cli.Command{
		Name:  "pull",
		Usage: "Pull changes from remote server",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return runSyncPull(ctx)
		},
	}

	syncPushCmd = &cli.Command{
		Name:  "push",
		Usage: "Push local changes to remote server",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return runSyncPush(ctx)
		},
	}

	syncStatusCmd = &cli.Command{
		Name:  "status",
		Usage: "Show sync status",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return runSyncStatus(ctx)
		},
	}

	syncCmd = &cli.Command{
		Name:     "sync",
		Usage:    "Sync database with remote server",
		Category: "Sync",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return runSync(ctx)
		},
		Commands: []*cli.Command{syncPullCmd, syncPushCmd, syncStatusCmd},
	}
)

type syncStatusResult struct {
	Action          string `json:"action"`
	LastPullTime    int64  `json:"lastPullTime"`
	LastPushTime    int64  `json:"lastPushTime"`
	PendingChanges  int64  `json:"pendingChanges"`
	NetworkSent     int64  `json:"networkSentBytes"`
	NetworkReceived int64  `json:"networkReceivedBytes"`
	Revision        string `json:"revision"`
}

type syncResult struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func runSync(ctx context.Context) error {
	st, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	err = st.Sync(ctx)
	if err != nil {
		return fmt.Errorf("sync: %w", err)
	}
	err = st.Push(ctx)
	if err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	writer := outputWriter()
	if writer.JSON {
		enc := json.NewEncoder(writer.Out)
		return enc.Encode(syncResult{Action: "sync", Message: "synced with remote"})
	}
	_, err = fmt.Fprintln(writer.Out, "synced with remote")
	return err
}

func runSyncPull(ctx context.Context) error {
	st, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	err = st.Sync(ctx)
	if err != nil {
		return fmt.Errorf("sync pull: %w", err)
	}

	writer := outputWriter()
	if writer.JSON {
		enc := json.NewEncoder(writer.Out)
		return enc.Encode(syncResult{Action: "pull", Message: "pulled changes from remote"})
	}
	_, err = fmt.Fprintln(writer.Out, "pulled changes from remote")
	return err
}

func runSyncPush(ctx context.Context) error {
	st, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	err = st.Push(ctx)
	if err != nil {
		return fmt.Errorf("sync push: %w", err)
	}

	writer := outputWriter()
	if writer.JSON {
		enc := json.NewEncoder(writer.Out)
		return enc.Encode(syncResult{Action: "push", Message: "pushed changes to remote"})
	}
	_, err = fmt.Fprintln(writer.Out, "pushed changes to remote")
	return err
}

func runSyncStatus(ctx context.Context) error {
	st, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	stats, err := st.SyncStats(ctx)
	if err != nil {
		return fmt.Errorf("sync status: %w", err)
	}

	writer := outputWriter()
	if writer.JSON {
		enc := json.NewEncoder(writer.Out)
		return enc.Encode(syncStatusResult{
			Action:          "status",
			LastPullTime:    stats.LastPullUnixTime,
			LastPushTime:    stats.LastPushUnixTime,
			PendingChanges:  stats.CdcOperations,
			NetworkSent:     stats.NetworkSentBytes,
			NetworkReceived: stats.NetworkReceivedBytes,
			Revision:        stats.Revision,
		})
	}

	tw := tabwriter.NewWriter(writer.Out, 0, syncStatusTabWidth, syncStatusPadding, ' ', 0)
	write := func(format string, a ...any) error {
		_, writeErr := fmt.Fprintf(tw, format, a...)
		return writeErr
	}

	if stats.LastPullUnixTime > 0 {
		err = write("last_pull:\t%d\n", stats.LastPullUnixTime)
		if err != nil {
			return err
		}
	} else {
		err = write("last_pull:\tnever\n")
		if err != nil {
			return err
		}
	}

	if stats.LastPushUnixTime > 0 {
		err = write("last_push:\t%d\n", stats.LastPushUnixTime)
		if err != nil {
			return err
		}
	} else {
		err = write("last_push:\tnever\n")
		if err != nil {
			return err
		}
	}

	err = write("pending_changes:\t%d\n", stats.CdcOperations)
	if err != nil {
		return err
	}
	err = write("network_sent:\t%d bytes\n", stats.NetworkSentBytes)
	if err != nil {
		return err
	}
	err = write("network_received:\t%d bytes\n", stats.NetworkReceivedBytes)
	if err != nil {
		return err
	}
	err = write("revision:\t%s\n", stats.Revision)
	if err != nil {
		return err
	}

	return tw.Flush()
}
