package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

func init() {
	trustCmd := &cobra.Command{
		Use:   "trust",
		Short: "Manage trust levels for API keys",
	}

	setCmd := &cobra.Command{
		Use:   "set <service> <level>",
		Short: "Set trust level (none, low, medium, high, full)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s := keystore.New(cfg.StorePath)
			if err := s.SetTrustLevel(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Trust level for %q set to %q\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <service>",
		Short: "Get trust level for a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s := keystore.New(cfg.StorePath)
			level, err := s.GetTrustLevel(args[0])
			if err != nil {
				return err
			}
			fmt.Println(level)
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <service>",
		Short: "Clear trust level for a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s := keystore.New(cfg.StorePath)
			if err := s.ClearTrustLevel(args[0]); err != nil {
				return err
			}
			fmt.Printf("Trust level cleared for %q\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <level>",
		Short: "List services by trust level",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s := keystore.New(cfg.StorePath)
			services, err := s.ServicesByTrustLevel(args[0])
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SERVICE\tTRUST LEVEL")
			for _, svc := range services {
				fmt.Fprintf(w, "%s\t%s\n", svc, args[0])
			}
			return w.Flush()
		},
	}

	trustCmd.AddCommand(setCmd, getCmd, clearCmd, listCmd)
	AddCommand(trustCmd)
}
