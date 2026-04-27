package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/densestvoid/fencepost/internal/config"
	"github.com/densestvoid/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var readonlyCmd = &cobra.Command{
	Use:   "readonly",
	Short: "Manage read-only protection for service keys",
}

var readonlySetCmd = &cobra.Command{
	Use:   "set <service>",
	Short: "Mark a service key as read-only",
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
		if err := s.SetReadOnly(args[0], true); err != nil {
			return err
		}
		fmt.Printf("Service %q is now read-only.\n", args[0])
		return nil
	},
}

var readonlyUnsetCmd = &cobra.Command{
	Use:   "unset <service>",
	Short: "Remove read-only protection from a service key",
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
		if err := s.SetReadOnly(args[0], false); err != nil {
			return err
		}
		fmt.Printf("Service %q is no longer read-only.\n", args[0])
		return nil
	},
}

var readonlyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all read-only service keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		keys := s.ReadOnlyKeys()
		if len(keys) == 0 {
			fmt.Println("No read-only keys.")
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

func init() {
	readonlyCmd.AddCommand(readonlySetCmd)
	readonlyCmd.AddCommand(readonlyUnsetCmd)
	readonlyCmd.AddCommand(readonlyListCmd)
	AddCommand(readonlyCmd)
}
