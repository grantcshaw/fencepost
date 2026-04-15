package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

var diffCmd = &cobra.Command{
	Use:   "diff <snapshot-file>",
	Short: "Show differences between current store and a snapshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		snapshotPath := args[0]
		snap, err := keystore.LoadSnapshot(snapshotPath)
		if err != nil {
			return fmt.Errorf("load snapshot %q: %w", snapshotPath, err)
		}

		result := s.Diff(snap)

		if len(result.Added)+len(result.Removed)+len(result.Rotated)+len(result.Modified) == 0 {
			fmt.Fprintln(os.Stdout, "No differences found.")
			return nil
		}

		for _, name := range result.Added {
			fmt.Fprintf(os.Stdout, "+ %-30s  added\n", name)
		}
		for _, name := range result.Removed {
			fmt.Fprintf(os.Stdout, "- %-30s  removed\n", name)
		}
		for _, name := range result.Rotated {
			fmt.Fprintf(os.Stdout, "~ %-30s  rotated\n", name)
		}
		for _, name := range result.Modified {
			fmt.Fprintf(os.Stdout, "* %-30s  modified\n", name)
		}

		return nil
	},
}

func init() {
	AddCommand(diffCmd)
}
