package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/clitools/fencepost/internal/config"
	"github.com/clitools/fencepost/internal/keystore"
)

func init() {
	var clearFlag bool
	var listFlag string

	categoryCmd := &cobra.Command{
		Use:   "category [service] [category]",
		Short: "Set or view the category for a service",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}

			if listFlag != "" {
				services := s.ServicesByCategory(listFlag)
				if len(services) == 0 {
					fmt.Printf("No services in category %q\n", listFlag)
					return nil
				}
				for _, svc := range services {
					fmt.Println(svc)
				}
				return nil
			}

			service := args[0]

			if clearFlag {
				if err := s.ClearCategory(service); err != nil {
					return err
				}
				fmt.Printf("Category cleared for %q\n", service)
				return nil
			}

			if len(args) == 1 {
				cat, err := s.GetCategory(service)
				if err != nil {
					return err
				}
				if cat == "" {
					fmt.Printf("%s: (no category)\n", service)
				} else {
					fmt.Printf("%s: %s\n", service, cat)
				}
				return nil
			}

			if err := s.SetCategory(service, args[1]); err != nil {
				return err
			}
			fmt.Printf("Category %q set for %q\n", args[1], service)
			return nil
		},
	}

	categoryCmd.Flags().BoolVar(&clearFlag, "clear", false, "Clear the category for a service")
	categoryCmd.Flags().StringVar(&listFlag, "list", "", "List all services in a given category")
	AddCommand(categoryCmd)
}
