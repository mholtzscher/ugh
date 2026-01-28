package daemon

import (
	"context"
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
	storeOpts store.Options
	logger    *slog.Logger

	startTime time.Time
	mu        sync.Mutex
	running   bool

	// Config with its own mutex for hot-reload
	configMu sync.RWMutex
	config   Config

	// Sync state
	lastSyncTime        time.Time
	lastSyncError       error
	consecutiveFailures int

	// Config reload signal
	reloadCh chan struct{}
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
		reloadCh:  make(chan struct{}, 1),
	}
}

// UpdateConfig updates the daemon configuration at runtime.
// This is called when the config file changes (hot-reload).
func (d *Daemon) UpdateConfig(cfg Config) {
	d.configMu.Lock()
	old := d.config
	d.config = cfg
	d.configMu.Unlock()

	d.logger.Info("config reloaded",
		"periodic_sync", cfg.PeriodicSync,
		"log_level", cfg.LogLevel,
	)

	// Signal the run loop to reload if interval changed
	if old.PeriodicSync != cfg.PeriodicSync {
		select {
		case d.reloadCh <- struct{}{}:
		default:
		}
	}
}

// getConfig returns a copy of the current config.
func (d *Daemon) getConfig() Config {
	d.configMu.RLock()
	defer d.configMu.RUnlock()
	return d.config
}

// Run starts the daemon and blocks until shutdown.
// It handles graceful shutdown on SIGINT/SIGTERM.
func (d *Daemon) Run(ctx context.Context) error {
	d.mu.Lock()
	if d.running {
		d.mu.Unlock()
		return fmt.Errorf("daemon is already running")
	}
	d.running = true
	d.startTime = time.Now()
	d.mu.Unlock()

	// Check if sync is configured
	if d.storeOpts.SyncURL == "" {
		d.logger.Info("sync not configured (db.sync_url not set), nothing to do")
		fmt.Fprintln(os.Stderr, "Daemon exiting: sync not configured. Set db.sync_url in config to enable background sync.")
		return nil
	}

	// Create context that cancels on signals
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	cfg := d.getConfig()
	d.logger.Info("daemon starting",
		"sync_interval", cfg.PeriodicSync,
		"sync_url", d.storeOpts.SyncURL,
	)

	// Perform initial sync
	d.logger.Info("performing initial sync")
	if err := d.doSync(ctx); err != nil {
		d.logger.Warn("initial sync failed", "error", err)
	}

	// Start periodic sync loop
	ticker := time.NewTicker(cfg.PeriodicSync)
	defer ticker.Stop()

	for {
		select {
		case sig := <-sigCh:
			d.logger.Info("received signal, shutting down", "signal", sig)
			return d.shutdown(ctx)
		case <-ctx.Done():
			d.logger.Info("context cancelled, shutting down")
			return d.shutdown(context.Background())
		case <-d.reloadCh:
			// Config was reloaded, update ticker interval
			newCfg := d.getConfig()
			ticker.Reset(newCfg.PeriodicSync)
			d.logger.Info("sync interval updated", "interval", newCfg.PeriodicSync)
		case <-ticker.C:
			if err := d.doSync(ctx); err != nil {
				d.logger.Warn("periodic sync failed", "error", err)
			}
		}
	}
}

// doSync opens the database, syncs, and closes it.
func (d *Daemon) doSync(ctx context.Context) error {
	var err error
	cfg := d.getConfig()
	backoff := cfg.SyncRetryBackoff

	for attempt := 0; attempt <= cfg.SyncRetryMax; attempt++ {
		if attempt > 0 {
			d.logger.Debug("retrying sync", "attempt", attempt, "backoff", backoff)
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

			d.logger.Debug("sync completed successfully")
			return nil
		}

		d.logger.Debug("sync attempt failed", "error", err, "attempt", attempt)
	}

	// All retries exhausted
	d.mu.Lock()
	d.lastSyncError = err
	d.consecutiveFailures++
	failures := d.consecutiveFailures
	d.mu.Unlock()

	d.logger.Error("sync failed after retries", "error", err, "consecutive_failures", failures)
	return err
}

// syncOnce opens the DB, syncs, and closes it.
func (d *Daemon) syncOnce(ctx context.Context) error {
	st, err := store.Open(ctx, d.storeOpts)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}
	defer func() { _ = st.Close() }()

	if err := st.Sync(ctx); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	return nil
}

// shutdown performs graceful shutdown.
func (d *Daemon) shutdown(ctx context.Context) error {
	d.logger.Info("performing final sync before shutdown")

	// Final sync attempt
	if err := d.syncOnce(ctx); err != nil {
		d.logger.Warn("final sync failed", "error", err)
	}

	d.mu.Lock()
	d.running = false
	d.mu.Unlock()

	d.logger.Info("shutdown complete")
	return nil
}

// Uptime returns how long the daemon has been running.
func (d *Daemon) Uptime() time.Duration {
	return time.Since(d.startTime)
}
