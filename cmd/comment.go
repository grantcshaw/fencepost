package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/stnley/fencepost/internal/config"
	"github.com/stnley/fencepost/internal/keystore"
)

func init() {
	var clearFlag bool

	commentCmd := &cobra.Command{
		Use:   "comment <service> [text]",
		Short: "Set or clear an inline comment on a service",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}

			service := args[0]

			if clearFlag {
				if err := s.ClearComment(service); err != nil {
					return err
				}
				fmt.Fprintf(os.Stdout, "comment cleared for %q\n", service)
				return nil
			}

			if len(args) < 2 {
				text, err := s.GetComment(service)
				if err != nil {
					return err
				}
				if text == "" {
					fmt.Fprintf(os.Stdout, "no comment set for %q\n", service)
				} else {
					fmt.Fprintf(os.Stdout, "%s\n", text)
				}
				return nil
			}

			if err := s.SetComment(service, args[1]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "comment set for %q\n", service)
			return nil
		},
	}

	commentCmd.Flags().BoolVar(&clearFlag, "clear", false, "clear the comment for the service")
	AddCommand(commentCmd)
}
