package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/audit"
	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

var renameOverwrite bool

var renameCmd = &cobra.Command{
	Use:   "rename <old-name> <new-name>",
	Short: "Rename a service entry in the keystore",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldName := args[0]
		newName := args[1]

		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("opening keystore: %w", err)
		}

		if err := ks.Rename(oldName, newName, renameOverwrite); err != nil {
			return err
		}

		log, err := audit.New(cfg.AuditLogPath)
		if err != nil {
			return fmt.Errorf("opening audit log: %w", err)
		}
		if err := log.Log(newName, "renamed", fmt.Sprintf("from %q", oldName)); err != nil {
			return fmt.Errorf("writing audit log: %w", err)
		}

		fmt.Printf("Renamed %q → %q\n", oldName, newName)
		return nil
	},
}

func init() {
	renameCmd.Flags().BoolVar(&renameOverwrite, "overwrite", false, "overwrite destination if it already exists")
	AddCommand(renameCmd)
}
