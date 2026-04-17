package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Track and query last-accessed timestamps for services",
}

var accessTouchCmd = &cobra.Command{
	Use:   "touch <service>",
	Short: "Record current time as last accessed for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ks, _, err := bootstrap()
		_ = cfg
		if err != nil {
			return err
		}
		if err := ks.SetLastAccessed(args[0]); err != nil {
			return err
		}
		fmt.Printf("Recorded access for %q\n", args[0])
		return nil
	},
}

var accessNeverCmd = &cobra.Command{
	Use:   "never",
	Short: "List services that have never been accessed",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ks, _, err := bootstrap()
		_ = cfg
		if err != nil {
			return err
		}
		services := ks.NeverAccessed()
		if len(services) == 0 {
			fmt.Println("All services have been accessed.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE")
		for _, svc := range services {
			fmt.Fprintln(w, svc)
		}
		return w.Flush()
	},
}

var accessSinceCmd = &cobra.Command{
	Use:   "since <duration>",
	Short: "List services accessed within the given duration (e.g. 24h, 7d)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ks, _, err := bootstrap()
		_ = cfg
		if err != nil {
			return err
		}
		d, err := time.ParseDuration(args[0])
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", args[0], err)
		}
		cutoff := time.Now().UTC().Add(-d)
		services := ks.AccessedSince(cutoff)
		if len(services) == 0 {
			fmt.Println("No services accessed in that period.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE")
		for _, svc := range services {
			fmt.Fprintln(w, svc)
		}
		return w.Flush()
	},
}

func init() {
	accessCmd.AddCommand(accessTouchCmd, accessNeverCmd, accessSinceCmd)
	AddCommand(accessCmd)
}
