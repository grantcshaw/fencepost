package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/audit"
	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

var checkOnly bool

var rotateCmd = &cobra.Command{
	Use:   "rotate [service]",
	Short: "Rotate the API key for a service, or check which keys are stale",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		store, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		policy := keystore.DefaultRotationPolicy()
		if cfg.MaxKeyAgeDays > 0 {
			policy.MaxAgeDays = cfg.MaxKeyAgeDays
		}

		if checkOnly {
			stale, err := store.StaleKeys(policy)
			if err != nil {
				return err
			}
			if len(stale) == 0 {
				fmt.Println("All keys are up to date.")
				return nil
			}
			fmt.Println("Stale keys:")
			for _, svc := range stale {
				fmt.Printf("  - %s\n", svc)
			}
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("service name required (or use --check)")
		}
		service := args[0]

		newKey := generateKey()
		if err := store.Rotate(service, newKey); err != nil {
			return fmt.Errorf("rotate: %w", err)
		}

		logger, err := audit.New(cfg.AuditLogPath)
		if err != nil {
			return fmt.Errorf("audit logger: %w", err)
		}
		if err := logger.Log(service, "rotate", fmt.Sprintf("rotated at %s", time.Now().Format(time.RFC3339))); err != nil {
			return fmt.Errorf("audit log: %w", err)
		}

		fmt.Printf("Rotated key for %s.\n", service)
		return nil
	},
}

func init() {
	rotateCmd.Flags().BoolVar(&checkOnly, "check", false, "Only report stale keys without rotating")
	rootCmd.AddCommand(rotateCmd)
}
