package daemon

import (
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the daemon as a system service",
	Long: `Install the daemon as a user-level system service.

On Linux (systemd):
  Creates ~/.config/systemd/user/ughd.service and enables it.

On macOS (launchd):
  Creates ~/Library/LaunchAgents/com.ugh.daemon.plist and loads it.

After installation, use 'ugh daemon start' to start the service.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		binaryPath, err := getBinaryPath()
		if err != nil {
			return fmt.Errorf("get binary path: %w", err)
		}

		cfg := service.InstallConfig{
			BinaryPath: binaryPath,
			ConfigPath: getConfigPath(),
		}

		if err := mgr.Install(cfg); err != nil {
			if errors.Is(err, service.ErrAlreadyInstalled) {
				return errors.New("service already installed - use 'ugh daemon uninstall' first to reinstall")
			}
			return fmt.Errorf("install service: %w", err)
		}

		w := deps.OutputWriter()
		status, _ := mgr.Status()
		_, _ = fmt.Fprintln(w.Out, "Service installed at", status.ServicePath)
		_, _ = fmt.Fprintln(w.Out, "Run 'ugh daemon start' to start the daemon")
		return nil
	},
}
