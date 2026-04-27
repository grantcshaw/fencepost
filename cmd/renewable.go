package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/seankim658/fencepost/internal/config"
	"github.com/seankim658/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

func init() {
	var renewableCmd = &cobra.Command{
		Use:   "renewable",
		Short: "Manage auto-renewable flags for service keys",
	}

	var setCmd = &cobra.Command{
		Use:   "set <service>",
		Short: "Mark a service key as auto-renewable",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			if err := s.SetRenewable(args[0], true); err != nil {
				return err
			}
			fmt.Printf("Service %q marked as auto-renewable.\n", args[0])
			return nil
		},
	}

	var unsetCmd = &cobra.Command{
		Use:   "unset <service>",
		Short: "Remove auto-renewable flag from a service key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			if err := s.SetRenewable(args[0], false); err != nil {
				return err
			}
			fmt.Printf("Auto-renewable flag removed from %q.\n", args[0])
			return nil
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all auto-renewable service keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			keys := s.RenewableKeys()
			if len(keys) == 0 {
				fmt.Println("No auto-renewable services found.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SERVICE")
			for _, k := range keys {
				fmt.Fprintln(w, k)
			}
			return w.Flush()
		},
	}

	renewableCmd.AddCommand(setCmd, unsetCmd, listCmd)
	AddCommand(renewableCmd)
}
