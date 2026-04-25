package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/bxrne/fencepost/internal/config"
	"github.com/bxrne/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var sensitivityCmd = &cobra.Command{
	Use:   "sensitivity",
	Short: "Manage sensitivity levels for service keys",
}

var sensitivitySetCmd = &cobra.Command{
	Use:   "set <service> <level>",
	Short: "Set sensitivity level (public|internal|confidential|restricted)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		if err := s.SetSensitivity(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("sensitivity for %q set to %q\n", args[0], args[1])
		return nil
	},
}

var sensitivityGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get the sensitivity level of a service key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		level, err := s.GetSensitivity(args[0])
		if err != nil {
			return err
		}
		fmt.Println(level)
		return nil
	},
}

var sensitivityListCmd = &cobra.Command{
	Use:   "list <level>",
	Short: "List services with a given sensitivity level",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		services := s.ServicesBySensitivity(args[0])
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tSENSITIVITY")
		for _, svc := range services {
			fmt.Fprintf(w, "%s\t%s\n", svc, args[0])
		}
		return w.Flush()
	},
}

func init() {
	sensitivityCmd.AddCommand(sensitivitySetCmd, sensitivityGetCmd, sensitivityListCmd)
	AddCommand(sensitivityCmd)
}
