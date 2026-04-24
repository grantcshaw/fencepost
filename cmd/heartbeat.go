package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	var silentThreshold string

	heartbeatCmd := &cobra.Command{
		Use:   "heartbeat",
		Short: "Manage heartbeat timestamps for services",
	}

	touchCmd := &cobra.Command{
		Use:   "touch <service>",
		Short: "Record current time as the heartbeat for a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, store, err := loadConfigAndStore(cmd)
			if err != nil {
				return err
			}
			_ = cfg
			if err := store.SetHeartbeat(args[0], time.Now().UTC()); err != nil {
				return err
			}
			fmt.Printf("heartbeat recorded for %q\n", args[0])
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <service>",
		Short: "Clear the heartbeat timestamp for a service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, store, err := loadConfigAndStore(cmd)
			if err != nil {
				return err
			}
			_ = cfg
			return store.ClearHeartbeat(args[0])
		},
	}

	silentCmd := &cobra.Command{
		Use:   "silent",
		Short: "List services with no recent heartbeat",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, store, err := loadConfigAndStore(cmd)
			if err != nil {
				return err
			}
			_ = cfg
			d, err := time.ParseDuration(silentThreshold)
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", silentThreshold, err)
			}
			services := store.SilentServices(d)
			if len(services) == 0 {
				fmt.Println("no silent services")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SERVICE\tLAST HEARTBEAT")
			for _, svc := range services {
				hb, _ := store.GetHeartbeat(svc)
				last := "never"
				if !hb.IsZero() {
					last = hb.Format(time.RFC3339)
				}
				fmt.Fprintf(w, "%s\t%s\n", svc, last)
			}
			return w.Flush()
		},
	}
	silentCmd.Flags().StringVarP(&silentThreshold, "since", "s", "1h", "silence threshold duration (e.g. 30m, 2h)")

	heartbeatCmd.AddCommand(touchCmd, clearCmd, silentCmd)
	AddCommand(heartbeatCmd)
}
