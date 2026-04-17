package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fencepost/internal/config"
	"github.com/fencepost/internal/keystore"
)

var protocolCmd = &cobra.Command{
	Use:   "protocol",
	Short: "Manage API protocol for services",
}

var protocolSetCmd = &cobra.Command{
	Use:   "set <service> <protocol>",
	Short: "Set the protocol for a service (rest, grpc, graphql, soap, webhook, other)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load(cfgFile)
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		if err := s.SetProtocol(args[0], args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Printf("protocol for %s set to %s\n", args[0], args[1])
	},
}

var protocolGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get the protocol for a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load(cfgFile)
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		p, err := s.GetProtocol(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(p)
	},
}

var protocolListCmd = &cobra.Command{
	Use:   "list <protocol>",
	Short: "List services by protocol",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load(cfgFile)
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		for _, svc := range s.ServicesByProtocol(args[0]) {
			fmt.Println(svc)
		}
	},
}

func init() {
	protocolCmd.AddCommand(protocolSetCmd, protocolGetCmd, protocolListCmd)
	rootCmd.AddCommand(protocolCmd)
}
