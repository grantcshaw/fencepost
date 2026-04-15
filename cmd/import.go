package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/fencepost/internal/audit"
	"github.com/fencepost/internal/config"
	"github.com/fencepost/internal/keystore"
)

var importCmd = &cobra.Command{
	Use:   "import <service> <key>",
	Short: "Import an existing API key for a service",
	Long:  `Import an existing API key into the keystore without generating a new one. Useful for migrating existing keys.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		key := args[1]

		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("opening keystore: %w", err)
		}

		if existing, _ := ks.Get(service); existing != "" {
			overwrite, _ := cmd.Flags().GetBool("overwrite")
			if !overwrite {
				return fmt.Errorf("service %q already has a key; use --overwrite to replace it", service)
			}
		}

		if err := ks.Set(service, key); err != nil {
			return fmt.Errorf("storing key: %w", err)
		}

		log, err := audit.New(cfg.AuditLogPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not open audit log: %v\n", err)
		} else {
			_ = log.Log(service, "import", fmt.Sprintf("key imported at %s", time.Now().Format(time.RFC3339)))
		}

		fmt.Printf("Imported key for service %q\n", service)
		return nil
	},
}

func init() {
	importCmd.Flags().Bool("overwrite", false, "Overwrite an existing key if one already exists")
	AddCommand(importCmd)
}
