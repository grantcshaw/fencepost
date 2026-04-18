package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/seanenck/fencepost/internal/config"
	"github.com/seanenck/fencepost/internal/keystore"
)

var healthCmd = &cobra.Command{
	Use:   "health [service]",
	Short: "Check health status of one or all keys",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		if err := s.Load(); err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tSTATUS\tREASONS")

		if len(args) == 1 {
			r, err := s.HealthCheck(args[0])
			if err != nil {
				return fmt.Errorf("service %q not found", args[0])
			}
			printHealthRow(w, r)
		} else {
			for _, r := range s.HealthCheckAll() {
				printHealthRow(w, r)
			}
		}
		return w.Flush()
	},
}

func printHealthRow(w *tabwriter.Writer, r keystore.HealthReport) {
	reasons := "-"
	if len(r.Reasons) > 0 {
		reasons = ""
		for i, reason := range r.Reasons {
			if i > 0 {
				reasons += "; "
			}
			reasons += reason
		}
	}
	fmt.Fprintf(w, "%s\t%s\t%s\n", r.Service, r.Status, reasons)
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
