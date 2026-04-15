package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/fencepost/internal/audit"
	"github.com/user/fencepost/internal/config"
)

var (
	auditService string
	auditEvent   string
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "View audit log entries",
	Long:  `Display audit log entries with optional filtering by service or event type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		entries, err := audit.ReadAll(cfg.AuditLogPath)
		if err != nil {
			return fmt.Errorf("reading audit log: %w", err)
		}

		if auditService != "" {
			entries = audit.FilterByService(entries, auditService)
		}
		if auditEvent != "" {
			entries = audit.FilterByEvent(entries, auditEvent)
		}

		if len(entries) == 0 {
			fmt.Println("No audit log entries found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tSERVICE\tEVENT\tMESSAGE")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				e.Timestamp.Format("2006-01-02 15:04:05"),
				e.Service,
				e.Event,
				e.Message,
			)
		}
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.Flags().StringVarP(&auditService, "service", "s", "", "Filter by service name")
	auditCmd.Flags().StringVarP(&auditEvent, "event", "e", "", "Filter by event type (e.g. key.created, key.rotated)")
}
