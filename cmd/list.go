package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"fencepost/internal/config"
	"fencepost/internal/keystore"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed services and their key metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("opening keystore: %w", err)
		}

		services := ks.List()
		if len(services) == 0 {
			fmt.Println("No services registered.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tCREATED\tLAST ROTATED\tDUE FOR ROTATION")
		fmt.Fprintln(w, "-------\t-------\t------------\t----------------")

		policy := keystore.DefaultRotationPolicy()
		for _, svc := range services {
			entry, _ := ks.Get(svc)
			rotated := "never"
			if !entry.RotatedAt.IsZero() {
				rotated = entry.RotatedAt.Format(time.RFC3339)
			}
			due := "no"
			if policy.DueForRotation(entry) {
				due = "YES"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				svc,
				entry.CreatedAt.Format(time.RFC3339),
				rotated,
				due,
			)
		}
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
