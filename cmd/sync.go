package cmd

import (
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

type syncPullResult struct {
	Action     string `json:"action"`
	NewChanges bool   `json:"newChanges"`
	Message    string `json:"message"`
}

type syncPushResult struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

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
	ctx := cmd.Context()
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
	ctx := cmd.Context()
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
		return enc.Encode(syncPullResult{Action: "pull", Message: "pulled changes from remote"})
	}
	_, err = fmt.Fprintln(writer.Out, "pulled changes from remote")
	return err
}

func runSyncPush(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
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
		return enc.Encode(syncPushResult{Action: "push", Message: "pushed changes to remote"})
	}
	_, err = fmt.Fprintln(writer.Out, "pushed changes to remote")
	return err
}

func runSyncStatus(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
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

	if stats.LastPullUnixTime > 0 {
		fmt.Fprintf(tw, "last_pull:\t%s\n", formatUnixTime(stats.LastPullUnixTime))
	} else {
		fmt.Fprintf(tw, "last_pull:\tnever\n")
	}

	if stats.LastPushUnixTime > 0 {
		fmt.Fprintf(tw, "last_push:\t%s\n", formatUnixTime(stats.LastPushUnixTime))
	} else {
		fmt.Fprintf(tw, "last_push:\tnever\n")
	}

	fmt.Fprintf(tw, "pending_changes:\t%d\n", stats.CdcOperations)
	fmt.Fprintf(tw, "network_sent:\t%d bytes\n", stats.NetworkSentBytes)
	fmt.Fprintf(tw, "network_received:\t%d bytes\n", stats.NetworkReceivedBytes)
	fmt.Fprintf(tw, "revision:\t%s\n", stats.Revision)

	return tw.Flush()
}

func formatUnixTime(ts int64) string {
	if ts == 0 {
		return "never"
	}
	return fmt.Sprintf("%d", ts)
}

func init() {
	syncCmd.AddCommand(syncPullCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncStatusCmd)

	syncCmd.Flags().SortFlags = false
}
