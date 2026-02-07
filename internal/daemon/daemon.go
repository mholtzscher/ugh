package daemon

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
)

// Daemon handles periodic background sync to Turso.
// It only opens the database when syncing, avoiding lock contention with the CLI.
type Daemon struct {
	config    Config
	storeOpts store.Options
	logger    *slog.Logger

	startTime time.Time
	mu        sync.Mutex
	running   bool

	// Sync state
	lastSyncTime        time.Time
	lastSyncError       error
	consecutiveFailures int
}

// New creates a new Daemon instance.
func New(opts store.Options, cfg Config, logger *slog.Logger) *Daemon {
	if logger == nil {
		logger = slog.Default()
	}

	return &Daemon{
		config:    cfg,
		storeOpts: opts,
		logger:    logger,
	}
}

// Run starts the daemon and blocks until shutdown.
// It handles graceful shutdown on SIGINT/SIGTERM.
func (d *Daemon) Run(ctx context.Context) error {
	d.mu.Lock()
	if d.running {
		d.mu.Unlock()
		return errors.New("daemon is already running")
	}
	d.running = true
	d.startTime = time.Now()
	d.mu.Unlock()

	// Check if sync is configured
	if d.storeOpts.SyncURL == "" {
		d.logger.InfoContext(ctx, "sync not configured (db.sync_url not set), nothing to do")
		fmt.Fprintln(
			os.Stderr,
			"Daemon exiting: sync not configured. Set db.sync_url in config to enable background sync.",
		)
		return nil
	}

	// Create context that cancels on signals
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	d.logger.InfoContext(ctx, "daemon starting",
		"sync_interval", d.config.PeriodicSync,
		"sync_url", d.storeOpts.SyncURL,
	)

	// Perform initial sync
	d.logger.InfoContext(ctx, "performing initial sync")
	if err := d.doSync(ctx); err != nil {
		d.logger.WarnContext(ctx, "initial sync failed", "error", err)
	}

	// Start periodic sync loop
	ticker := time.NewTicker(d.config.PeriodicSync)
	defer ticker.Stop()

	for {
		select {
		case sig := <-sigCh:
			d.logger.InfoContext(ctx, "received signal, shutting down", "signal", sig)
			d.shutdown(ctx)
			return nil
		case <-ctx.Done():
			d.logger.InfoContext(ctx, "context cancelled, shutting down")
			d.shutdown(context.Background())
			return nil
		case <-ticker.C:
			if err := d.doSync(ctx); err != nil {
				d.logger.WarnContext(ctx, "periodic sync failed", "error", err)
			}
		}
	}
}

// doSync opens the database, syncs, and closes it.
func (d *Daemon) doSync(ctx context.Context) error {
	var err error
	backoff := d.config.SyncRetryBackoff

	for attempt := 0; attempt <= d.config.SyncRetryMax; attempt++ {
		if attempt > 0 {
			d.logger.DebugContext(ctx, "retrying sync", "attempt", attempt, "backoff", backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
		}

		err = d.syncOnce(ctx)
		if err == nil {
			d.mu.Lock()
			d.lastSyncTime = time.Now()
			d.lastSyncError = nil
			d.consecutiveFailures = 0
			d.mu.Unlock()

			d.logger.DebugContext(ctx, "sync completed successfully")
			return nil
		}

		d.logger.DebugContext(ctx, "sync attempt failed", "error", err, "attempt", attempt)
	}

	// All retries exhausted
	d.mu.Lock()
	d.lastSyncError = err
	d.consecutiveFailures++
	failures := d.consecutiveFailures
	d.mu.Unlock()

	d.logger.ErrorContext(ctx, "sync failed after retries", "error", err, "consecutive_failures", failures)
	return err
}

// syncOnce opens the DB, syncs, and closes it.
func (d *Daemon) syncOnce(ctx context.Context) error {
	st, err := store.Open(ctx, d.storeOpts)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}
	defer func() { _ = st.Close() }()

	err = st.Sync(ctx)
	if err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	return nil
}

// shutdown performs graceful shutdown.
func (d *Daemon) shutdown(ctx context.Context) {
	d.logger.InfoContext(ctx, "performing final sync before shutdown")

	// Final sync attempt
	if err := d.syncOnce(ctx); err != nil {
		d.logger.WarnContext(ctx, "final sync failed", "error", err)
	}

	d.mu.Lock()
	d.running = false
	d.mu.Unlock()

	d.logger.InfoContext(ctx, "shutdown complete")
}

// Uptime returns how long the daemon has been running.
func (d *Daemon) Uptime() time.Duration {
	return time.Since(d.startTime)
}
