package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync database with remote server",
		RunE:  runSync,
	}

	syncPullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull changes from remote server",
		RunE:  runSyncPull,
	}

	syncPushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push local changes to remote server",
		RunE:  runSyncPush,
	}

	syncStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show sync status",
		RunE:  runSyncStatus,
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

func runSync(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	st, err := openStore(ctx)
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

	writer := outputWriter()
	if writer.JSON {
		enc := json.NewEncoder(writer.Out)
		return enc.Encode(syncResult{Action: "sync", Message: "synced with remote"})
	}
	_, err = fmt.Fprintln(writer.Out, "synced with remote")
	return err
}

func runSyncPull(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
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

func runSyncPush(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	st, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	if err := st.Push(ctx); err != nil {
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

func runSyncStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
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

	tw := tabwriter.NewWriter(writer.Out, 0, 2, 2, ' ', 0)
	write := func(format string, a ...any) error {
		_, err := fmt.Fprintf(tw, format, a...)
		return err
	}

	if stats.LastPullUnixTime > 0 {
		if err := write("last_pull:\t%d\n", stats.LastPullUnixTime); err != nil {
			return err
		}
	} else {
		if err := write("last_pull:\tnever\n"); err != nil {
			return err
		}
	}

	if stats.LastPushUnixTime > 0 {
		if err := write("last_push:\t%d\n", stats.LastPushUnixTime); err != nil {
			return err
		}
	} else {
		if err := write("last_push:\tnever\n"); err != nil {
			return err
		}
	}

	if err := write("pending_changes:\t%d\n", stats.CdcOperations); err != nil {
		return err
	}
	if err := write("network_sent:\t%d bytes\n", stats.NetworkSentBytes); err != nil {
		return err
	}
	if err := write("network_received:\t%d bytes\n", stats.NetworkReceivedBytes); err != nil {
		return err
	}
	if err := write("revision:\t%s\n", stats.Revision); err != nil {
		return err
	}

	return tw.Flush()
}

func init() {
	syncCmd.AddCommand(syncPullCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncStatusCmd)
}
