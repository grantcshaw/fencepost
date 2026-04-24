package cmd

import (
	"fmt"

	"github.com/nicholasgasior/fencepost/internal/audit"
	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

func init() {
	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Manage key baselines for drift detection",
	}

	setCmd := &cobra.Command{
		Use:   "set <service>",
		Short: "Record the current key as the baseline for a service",
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
			if err := s.SetBaseline(args[0]); err != nil {
				return err
			}
			log, _ := audit.New(cfg.AuditLogPath)
			_ = log.Log(args[0], "baseline-set", "")
			fmt.Printf("Baseline recorded for %q\n", args[0])
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <service>",
		Short: "Clear the baseline for a service",
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
			if err := s.ClearBaseline(args[0]); err != nil {
				return err
			}
			fmt.Printf("Baseline cleared for %q\n", args[0])
			return nil
		},
	}

	driftCmd := &cobra.Command{
		Use:   "drift <service>",
		Short: "Check whether a service key has drifted from its baseline",
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
			changed, err := s.BaselineChanged(args[0])
			if err != nil {
				return err
			}
			if changed {
				fmt.Printf("%s: DRIFTED (key differs from baseline)\n", args[0])
			} else {
				fmt.Printf("%s: OK (matches baseline)\n", args[0])
			}
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all services that have a baseline recorded",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			names := s.ServicesWithBaseline()
			if len(names) == 0 {
				fmt.Println("No baselines recorded.")
				return nil
			}
			for _, name := range names {
				_, ts, _ := s.GetBaseline(name)
				fmt.Printf("%-24s  recorded at %s\n", name, ts.Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}

	baselineCmd.AddCommand(setCmd, clearCmd, driftCmd, listCmd)
	AddCommand(baselineCmd)
}
