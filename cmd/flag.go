package cmd

import (
	"fmt"
	"os"

	"github.com/danielmichaels/fencepost/internal/config"
	"github.com/danielmichaels/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var flagCmd = &cobra.Command{
	Use:   "flag",
	Short: "Manage flags on service keys",
}

var flagSetCmd = &cobra.Command{
	Use:   "set <service> <flag>",
	Short: "Set a flag on a service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		if err := ks.SetFlag(args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "flag %q set on %q\n", args[1], args[0])
		return nil
	},
}

var flagUnsetCmd = &cobra.Command{
	Use:   "unset <service> <flag>",
	Short: "Remove a flag from a service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		return ks.UnsetFlag(args[0], args[1])
	},
}

var flagListCmd = &cobra.Command{
	Use:   "list <flag>",
	Short: "List services with a given flag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		for _, svc := range ks.ServicesByFlag(args[0]) {
			fmt.Println(svc)
		}
		return nil
	},
}

func init() {
	flagCmd.AddCommand(flagSetCmd, flagUnsetCmd, flagListCmd)
	AddCommand(flagCmd)
}
