package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/audit"
	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

var expireCmd = &cobra.Command{
	Use:   "expire",
	Short: "List or rotate keys that have exceeded their maximum age",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open keystore: %w", err)
		}

		policy := keystore.DefaultExpiryPolicy()
		expired := keystore.ExpiredKeys(ks, policy)

		if len(expired) == 0 {
			fmt.Println("No expired keys found.")
			return nil
		}

		sort.Strings(expired)

		rotateFlag, _ := cmd.Flags().GetBool("rotate")
		if !rotateFlag {
			fmt.Printf("Expired keys (%d):\n", len(expired))
			for _, svc := range expired {
				fmt.Printf("  %s\n", svc)
			}
			return nil
		}

		logger, err := audit.New(cfg.AuditLogPath)
		if err != nil {
			return fmt.Errorf("open audit log: %w", err)
		}

		for _, svc := range expired {
			newKey := generateKey()
			if err := ks.Rotate(svc, newKey); err != nil {
				fmt.Fprintf(os.Stderr, "rotate %s: %v\n", svc, err)
				continue
			}
			_ = logger.Log(svc, "expire-rotate", "key rotated due to expiry")
			fmt.Printf("Rotated: %s\n", svc)
		}
		return nil
	},
}

func init() {
	expireCmd.Flags().Bool("rotate", false, "automatically rotate all expired keys")
	rootCmd.AddCommand(expireCmd)
}
