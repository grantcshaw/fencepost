package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func init() {
	var sourceCmd = &cobra.Command{
		Use:   "source",
		Short: "Manage key source metadata",
	}

	var setCmd = &cobra.Command{
		Use:   "set <service> <source>",
		Short: "Set the source for a service key (manual, vault, aws, gcp, azure, env, file)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			if err := s.SetSource(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Source for %q set to %q\n", args[0], args[1])
			return nil
		},
	}

	var getCmd = &cobra.Command{
		Use:   "get <service>",
		Short: "Get the source for a service key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			src, err := s.GetSource(args[0])
			if err != nil {
				return err
			}
			fmt.Println(src)
			return nil
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list <source>",
		Short: "List services by source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			services, err := s.ServicesBySource(args[0])
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SERVICE\tSOURCE")
			for _, svc := range services {
				fmt.Fprintf(w, "%s\t%s\n", svc, args[0])
			}
			return w.Flush()
		},
	}

	sourceCmd.AddCommand(setCmd, getCmd, listCmd)
	AddCommand(sourceCmd)
}
