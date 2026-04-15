package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"fencepost/internal/config"
	"fencepost/internal/keystore"
)

var statusCmd = &cobra.Command{
	Use:   "status [service]",
	Short: "Show key status for one or all services",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		store := keystore.New(cfg.StorePath)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tCREATED\tROTATED\tEXPIRED\tROTATION DUE\tTAGS")

		if len(args) == 1 {
			st, err := store.Status(args[0])
			if err != nil {
				return fmt.Errorf("status %q: %w", args[0], err)
			}
			printStatusRow(w, st)
		} else {
			for _, st := range store.StatusAll() {
				printStatusRow(w, st)
			}
		}

		return w.Flush()
	},
}

func printStatusRow(w *tabwriter.Writer, st keystore.KeyStatus) {
	rotated := "-"
	if !st.RotatedAt.IsZero() {
		rotated = st.RotatedAt.Format("2006-01-02")
	}
	tags := "-"
	if len(st.Tags) > 0 {
		tags = fmt.Sprintf("%v", st.Tags)
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%v\t%s\n",
		st.Service,
		st.CreatedAt.Format("2006-01-02"),
		rotated,
		st.IsExpired,
		st.DueRotation,
		tags,
	)
}

func init() {
	AddCommand(statusCmd)
}
