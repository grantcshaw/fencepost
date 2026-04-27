package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var revokedCmd = &cobra.Command{
	Use:   "revoked",
	Short: "Manage revoked key status",
}

var revokedSetCmd = &cobra.Command{
	Use:   "set <service> [reason]",
	Short: "Mark a service key as revoked",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		reason := ""
		if len(args) == 2 {
			reason = args[1]
		}
		if err := s.SetRevoked(args[0], reason); err != nil {
			return err
		}
		fmt.Printf("service %q marked as revoked\n", args[0])
		return nil
	},
}

var revokedClearCmd = &cobra.Command{
	Use:   "clear <service>",
	Short: "Remove revoked status from a service key",
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
		if err := s.ClearRevoked(args[0]); err != nil {
			return err
		}
		fmt.Printf("revoked status cleared for %q\n", args[0])
		return nil
	},
}

var revokedListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all revoked service keys",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		keys := s.RevokedKeys()
		if len(keys) == 0 {
			fmt.Println("no revoked keys")
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
	revokedCmd.AddCommand(revokedSetCmd)
	revokedCmd.AddCommand(revokedClearCmd)
	revokedCmd.AddCommand(revokedListCmd)
	AddCommand(revokedCmd)
}
