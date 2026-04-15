package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <service>",
	Short: "Archive a service key and remove it from the active store",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(cmd)
		if err != nil {
			return err
		}
		store, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("opening store: %w", err)
		}
		archivePath, _ := cmd.Flags().GetString("archive-file")
		if archivePath == "" {
			archivePath = cfg.StorePath + ".archive.json"
		}
		service := args[0]
		if err := store.Archive(service, archivePath); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "archived service %q → %s\n", service, archivePath)
		return nil
	},
}

var archiveListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all archived service entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(cmd)
		if err != nil {
			return err
		}
		archivePath, _ := cmd.Flags().GetString("archive-file")
		if archivePath == "" {
			archivePath = cfg.StorePath + ".archive.json"
		}
		entries, err := keystore.LoadArchive(archivePath)
		if err != nil {
			return fmt.Errorf("loading archive: %w", err)
		}
		if len(entries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no archived entries found")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tARCHIVED AT")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\n", e.Entry.Service, e.ArchivedAt.Format(time.RFC3339))
		}
		w.Flush()
		return nil
	},
}

func init() {
	archiveCmd.AddCommand(archiveListCmd)
	archiveCmd.PersistentFlags().String("archive-file", "", "path to archive file (default: <store>.archive.json)")
	AddCommand(archiveCmd)
}
