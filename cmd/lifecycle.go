package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/clikd-inc/fencepost/internal/config"
	"github.com/clikd-inc/fencepost/internal/keystore"
)

func init() {
	lifecycleCmd := &cobra.Command{
		Use:   "lifecycle",
		Short: "Manage key lifecycle status (active, deprecated, retired)",
	}

	setCmd := &cobra.Command{
		Use:   "set <service> <status>",
		Short: "Set lifecycle status for a service",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			if err := s.SetLifecycle(args[0], keystore.LifecycleStatus(args[1])); err != nil {
				return fmt.Errorf("set lifecycle: %w", err)
			}
			fmt.Printf("lifecycle for %q set to %q\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <service>",
		Short: "Get lifecycle status for a service",
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
			status, err := s.GetLifecycle(args[0])
			if err != nil {
				return err
			}
			fmt.Println(status)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <status>",
		Short: "List services by lifecycle status",
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
			services, err := s.ServicesByLifecycle(keystore.LifecycleStatus(args[0]))
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SERVICE\tLIFECYCLE")
			for _, svc := range services {
				fmt.Fprintf(w, "%s\t%s\n", svc, args[0])
			}
			return w.Flush()
		},
	}

	lifecycleCmd.AddCommand(setCmd, getCmd, listCmd)
	AddCommand(lifecycleCmd)
}
