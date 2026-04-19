package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/fencepost/internal/config"
	"github.com/fencepost/internal/keystore"
)

func init() {
	quotaCmd := &cobra.Command{
		Use:   "quota",
		Short: "Manage per-service request quotas",
	}

	setCmd := &cobra.Command{
		Use:   "set <service> <quota>",
		Short: "Set the quota for a service",
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
			v, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("quota must be an integer: %w", err)
			}
			if err := s.SetQuota(args[0], v); err != nil {
				return err
			}
			fmt.Printf("quota for %q set to %d\n", args[0], v)
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <service>",
		Short: "Get the quota for a service",
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
			q, err := s.GetQuota(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("%d\n", q)
			return nil
		},
	}

	var minQuota int
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List services with quota >= min",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			svcs := s.ServicesByQuota(minQuota)
			if len(svcs) == 0 {
				fmt.Println("no services found")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SERVICE\tQUOTA")
			for _, svc := range svcs {
				q, _ := s.GetQuota(svc)
				fmt.Fprintf(w, "%s\t%d\n", svc, q)
			}
			return w.Flush()
		},
	}
	listCmd.Flags().IntVar(&minQuota, "min", 1, "minimum quota value to include")

	quotaCmd.AddCommand(setCmd, getCmd, listCmd)
	AddCommand(quotaCmd)
}
