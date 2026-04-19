package cmd

import (
	"fmt"
	"os"

	"github.com/seanmorris/fencepost/internal/config"
	"github.com/seanmorris/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

func init() {
	var linkCmd = &cobra.Command{
		Use:   "link",
		Short: "Manage documentation links for services",
	}

	var setLinkCmd = &cobra.Command{
		Use:   "set <service> <url>",
		Short: "Set a documentation link for a service",
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
			return s.SetLink(args[0], args[1])
		},
	}

	var getLinkCmd = &cobra.Command{
		Use:   "get <service>",
		Short: "Get the documentation link for a service",
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
			link, err := s.GetLink(args[0])
			if err != nil {
				return err
			}
			if link == "" {
				fmt.Fprintln(os.Stdout, "(no link set)")
			} else {
				fmt.Fprintln(os.Stdout, link)
			}
			return nil
		},
	}

	var clearLinkCmd = &cobra.Command{
		Use:   "clear <service>",
		Short: "Clear the documentation link for a service",
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
			return s.ClearLink(args[0])
		},
	}

	var listLinkCmd = &cobra.Command{
		Use:   "list",
		Short: "List all services with a documentation link",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			for _, svc := range s.ServicesByLink() {
				link, _ := s.GetLink(svc)
				fmt.Fprintf(os.Stdout, "%-20s %s\n", svc, link)
			}
			return nil
		},
	}

	linkCmd.AddCommand(setLinkCmd, getLinkCmd, clearLinkCmd, listLinkCmd)
	AddCommand(linkCmd)
}
