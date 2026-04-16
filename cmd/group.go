package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage service groups",
	}

	setCmd := &cobra.Command{
		Use:   "set <service> <group>",
		Short: "Assign a service to a group",
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
			if err := s.SetGroup(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Service %q assigned to group %q\n", args[0], args[1])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [group]",
		Short: "List groups or services in a group",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			if len(args) == 1 {
				for _, svc := range s.ServicesByGroup(args[0]) {
					fmt.Println(svc)
				}
			} else {
				for _, g := range s.Groups() {
					fmt.Println(g)
				}
			}
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <service>",
		Short: "Remove group assignment from a service",
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
			return s.ClearGroup(args[0])
		},
	}

	groupCmd.AddCommand(setCmd, listCmd, clearCmd)
	AddCommand(groupCmd)
}
