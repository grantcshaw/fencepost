package cmd

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

func init() {
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage rotation policies for services",
	}

	setCmd := &cobra.Command{
		Use:   "set <service> <policy>",
		Short: "Set rotation policy (strict, moderate, relaxed, none)",
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
			if err := s.SetPolicy(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "policy for %s set to %s\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <service>",
		Short: "Get the rotation policy for a service",
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
			pol, err := s.GetPolicy(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, pol)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <policy>",
		Short: "List services with a given rotation policy",
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
			for _, svc := range s.ServicesByPolicy(args[0]) {
				fmt.Fprintln(os.Stdout, svc)
			}
			return nil
		},
	}

	policyCmd.AddCommand(setCmd, getCmd, listCmd)
	AddCommand(policyCmd)
}
