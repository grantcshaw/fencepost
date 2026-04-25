package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/clikd-inc/fencepost/internal/audit"
	"github.com/clikd-inc/fencepost/internal/config"
	"github.com/clikd-inc/fencepost/internal/keystore"
)

func init() {
	attestationCmd := &cobra.Command{
		Use:   "attestation",
		Short: "Manage key attestation methods",
	}

	setCmd := &cobra.Command{
		Use:   "set <service> <method>",
		Short: "Set attestation method (none, tpm, hsm, software, cloud)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			ks, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			if err := ks.SetAttestation(args[0], args[1]); err != nil {
				return err
			}
			log, _ := audit.New(cfg.AuditLogPath)
			log.Log(args[0], "attestation-set", args[1])
			fmt.Printf("attestation for %s set to %s\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <service>",
		Short: "Get attestation method for a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			ks, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			val, err := ks.GetAttestation(args[0])
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <method>",
		Short: "List services using a given attestation method",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			ks, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			services, err := ks.ServicesByAttestation(args[0])
			if err != nil {
				return err
			}
			if len(services) == 0 {
				fmt.Println("no services found")
				return nil
			}
			for _, svc := range services {
				fmt.Fprintln(os.Stdout, svc)
			}
			return nil
		},
	}

	attestationCmd.AddCommand(setCmd, getCmd, listCmd)
	AddCommand(attestationCmd)
}
