package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fencepost/internal/config"
	"github.com/fencepost/internal/keystore"
)

func init() {
	var regionCmd = &cobra.Command{
		Use:   "region",
		Short: "Manage service regions",
	}

	var setRegionCmd = &cobra.Command{
		Use:   "set <service> <region>",
		Short: "Set the region for a service",
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
			return s.SetRegion(args[0], args[1])
		},
	}

	var getRegionCmd = &cobra.Command{
		Use:   "get <service>",
		Short: "Get the region for a service",
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
			region, err := s.GetRegion(args[0])
			if err != nil {
				return err
			}
			if region == "" {
				fmt.Println("(no region set)")
			} else {
				fmt.Println(region)
			}
			return nil
		},
	}

	var listRegionCmd = &cobra.Command{
		Use:   "list <region>",
		Short: "List services in a region",
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
			for _, svc := range s.ServicesByRegion(args[0]) {
				fmt.Println(svc)
			}
			return nil
		},
	}

	regionCmd.AddCommand(setRegionCmd, getRegionCmd, listRegionCmd)
	AddCommand(regionCmd)
}
