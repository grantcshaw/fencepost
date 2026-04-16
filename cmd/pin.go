package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/janearc/fencepost/internal/config"
	"github.com/janearc/fencepost/internal/keystore"
)

func init() {
	var unpin bool

	pinCmd := &cobra.Command{
		Use:   "pin <service>",
		Short: "Pin or unpin a service key to prevent automatic rotation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service := args[0]

			cfg, err := config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			ks, err := keystore.New(cfg.StorePath)
			if err != nil {
				return fmt.Errorf("open keystore: %w", err)
			}

			if unpin {
				if err := ks.UnpinEntry(service); err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("unpinned %q\n", service)
			} else {
				if err := ks.PinEntry(service); err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("pinned %q — key will not be auto-rotated\n", service)
			}
			return nil
		},
	}

	pinCmd.Flags().BoolVar(&unpin, "unpin", false, "Remove pin from service key")
	AddCommand(pinCmd)
}
