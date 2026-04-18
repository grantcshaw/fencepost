package cmd

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var contactCmd = &cobra.Command{
	Use:   "contact",
	Short: "Manage contact information for service keys",
}

var contactSetCmd = &cobra.Command{
	Use:   "set <service> <contact>",
	Short: "Set the contact for a service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		if err := ks.SetContact(args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "contact set for %q\n", args[0])
		return nil
	},
}

var contactGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get the contact for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		c, err := ks.GetContact(args[0])
		if err != nil {
			return err
		}
		if c == "" {
			fmt.Fprintln(os.Stdout, "(no contact set)")
		} else {
			fmt.Fprintln(os.Stdout, c)
		}
		return nil
	},
}

var contactListCmd = &cobra.Command{
	Use:   "list <contact>",
	Short: "List services by contact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		for _, svc := range ks.ServicesByContact(args[0]) {
			fmt.Fprintln(os.Stdout, svc)
		}
		return nil
	},
}

func init() {
	contactCmd.AddCommand(contactSetCmd, contactGetCmd, contactListCmd)
	AddCommand(contactCmd)
}
