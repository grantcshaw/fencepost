package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/audit"
	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

var noteClearFlag bool

var noteCmd = &cobra.Command{
	Use:   "note <service> [text]",
	Short: "Set or clear a note on a service key entry",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		store, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}
		logger, err := audit.New(cfg.AuditLogPath)
		if err != nil {
			return fmt.Errorf("open audit log: %w", err)
		}

		service := args[0]

		if noteClearFlag {
			if err := store.ClearNote(service); err != nil {
				return err
			}
			_ = logger.Log(service, "note_cleared", "note removed")
			fmt.Printf("Note cleared for %q.\n", service)
			return nil
		}

		if len(args) < 2 {
			note, err := store.GetNote(service)
			if err != nil {
				return err
			}
			if note == "" {
				fmt.Printf("No note set for %q.\n", service)
			} else {
				fmt.Printf("Note for %q: %s\n", service, note)
			}
			return nil
		}

		note := args[1]
		if err := store.SetNote(service, note); err != nil {
			return err
		}
		_ = logger.Log(service, "note_set", fmt.Sprintf("note updated: %s", note))
		fmt.Printf("Note set for %q.\n", service)
		return nil
	},
}

func init() {
	noteCmd.Flags().BoolVar(&noteClearFlag, "clear", false, "Clear the note for the service")
	AddCommand(noteCmd)
}
