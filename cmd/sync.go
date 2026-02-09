package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/output"
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
	return writer.WriteSuccess("synced with remote")
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
	return writer.WriteSuccess("pulled changes from remote")
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
	return writer.WriteSuccess("pushed changes to remote")
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

	rows := []output.KeyValue{
		{Key: "last_pull", Value: formatSyncUnixTime(stats.LastPullUnixTime)},
		{Key: "last_push", Value: formatSyncUnixTime(stats.LastPushUnixTime)},
		{Key: "pending_changes", Value: strconv.FormatInt(stats.CdcOperations, 10)},
		{Key: "network_sent", Value: strconv.FormatInt(stats.NetworkSentBytes, 10) + " bytes"},
		{Key: "network_received", Value: strconv.FormatInt(stats.NetworkReceivedBytes, 10) + " bytes"},
		{Key: "revision", Value: stats.Revision},
	}

	return writer.WriteKeyValues(rows)
}

func formatSyncUnixTime(unixTime int64) string {
	if unixTime <= 0 {
		return "never"
	}
	return strconv.FormatInt(unixTime, 10)
}
