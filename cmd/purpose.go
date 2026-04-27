package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var purposeCmd = &cobra.Command{
	Use:   "purpose",
	Short: "Manage the purpose field for API keys",
}

var purposeSetCmd = &cobra.Command{
	Use:   "set <service> <purpose>",
	Short: "Set the purpose for a service key",
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
		if err := s.SetPurpose(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("purpose for %q set to %q\n", args[0], args[1])
		return nil
	},
}

var purposeGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get the purpose for a service key",
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
		purpose, err := s.GetPurpose(args[0])
		if err != nil {
			return err
		}
		fmt.Println(purpose)
		return nil
	},
}

var purposeListCmd = &cobra.Command{
	Use:   "list <purpose>",
	Short: "List services with a given purpose",
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
		services, err := s.ServicesByPurpose(args[0])
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tPURPOSE")
		for _, svc := range services {
			fmt.Fprintf(w, "%s\t%s\n", svc, args[0])
		}
		return w.Flush()
	},
}

func init() {
	purposeCmd.AddCommand(purposeSetCmd, purposeGetCmd, purposeListCmd)
	AddCommand(purposeCmd)
}
