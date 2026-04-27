package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/smlx/fencepost/internal/config"
	"github.com/smlx/fencepost/internal/keystore"
)

var snapshotDir string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Write a point-in-time snapshot of the keystore",
	Long: `Saves all current key records to a timestamped JSON file in the
specified directory (defaults to ~/.fencepost/snapshots).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		store := keystore.New(cfg.StorePath)

		dir, err := resolveSnapshotDir(snapshotDir)
		if err != nil {
			return err
		}

		path, err := store.WriteSnapshot(dir)
		if err != nil {
			return fmt.Errorf("write snapshot: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Snapshot written to %s\n", path)
		return nil
	},
}

// resolveSnapshotDir returns the snapshot directory to use. If dir is empty,
// it falls back to ~/.fencepost/snapshots.
func resolveSnapshotDir(dir string) (string, error) {
	if dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".fencepost", "snapshots"), nil
}

func init() {
	snapshotCmd.Flags().StringVarP(
		&snapshotDir, "dir", "d", "",
		"directory to write snapshot (default ~/.fencepost/snapshots)",
	)
	AddCommand(snapshotCmd)
}
