package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/audit"
	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <service>",
		Short: "Add a service to the watch list",
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
			if err := s.Watch(args[0]); err != nil {
				return err
			}
			log, _ := audit.New(cfg.AuditLogPath)
			_ = log.Log(args[0], "watch", "added to watch list")
			fmt.Printf("Watching %s\n", args[0])
			return nil
		},
	}

	unwatchCmd := &cobra.Command{
		Use:   "unwatch <service>",
		Short: "Remove a service from the watch list",
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
			if err := s.Unwatch(args[0]); err != nil {
				return err
			}
			log, _ := audit.New(cfg.AuditLogPath)
			_ = log.Log(args[0], "unwatch", "removed from watch list")
			fmt.Printf("Unwatched %s\n", args[0])
			return nil
		},
	}

	watchListCmd := &cobra.Command{
		Use:   "watched",
		Short: "List all watched services",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			s, err := keystore.New(cfg.StorePath)
			if err != nil {
				return err
			}
			entries := s.WatchedKeys()
			if len(entries) == 0 {
				fmt.Println("No watched services.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SERVICE\tWATCHED SINCE")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\n", e.Service, e.AddedAt.Format("2006-01-02 15:04:05"))
			}
			return w.Flush()
		},
	}

	AddCommand(watchCmd)
	AddCommand(unwatchCmd)
	AddCommand(watchListCmd)
}
