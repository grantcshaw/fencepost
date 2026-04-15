package cmd

import (
n	"github.com/spf13/cobra"

	"github.com/smlmbrt/fencepost/internal/audit"
	"github.com/smlmbrt/fencepost/internal/config"
	"github.com/smlmbrt/fencepost/internal/keystore"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup or restore the key store",
}

var backupWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a timestamped backup of the key store to a directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		store, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		meta, err := store.WriteBackup(dir)
		if err != nil {
			return fmt.Errorf("write backup: %w", err)
		}

		logger, err := audit.New(cfg.AuditLogPath)
		if err == nil {
			_ = logger.Log("backup", "all", fmt.Sprintf("backup written to %s", meta.Path))
		}

		fmt.Printf("Backup written: %s (%d services)\n", meta.Path, meta.Services)
		return nil
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore <file>",
	Short: "Restore the key store from a backup file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		store, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		if err := store.RestoreBackup(args[0]); err != nil {
			return fmt.Errorf("restore backup: %w", err)
		}

		logger, err := audit.New(cfg.AuditLogPath)
		if err == nil {
			_ = logger.Log("restore", "all", fmt.Sprintf("store restored from %s", args[0]))
		}

		fmt.Printf("Store restored from %s\n", args[0])
		return nil
	},
}

func init() {
	backupWriteCmd.Flags().String("dir", ".", "Directory to write the backup file into")
	backupCmd.AddCommand(backupWriteCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	AddCommand(backupCmd)
}
