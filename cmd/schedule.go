package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage rotation schedules for services",
}

var scheduleSetCmd = &cobra.Command{
	Use:   "set <service> <duration>",
	Short: "Set a rotation schedule (e.g. 24h, 168h)",
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
		if err := s.SetSchedule(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Schedule for %q set to %s\n", args[0], args[1])
		return nil
	},
}

var scheduleClearCmd = &cobra.Command{
	Use:   "clear <service>",
	Short: "Clear the rotation schedule for a service",
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
		return s.ClearSchedule(args[0])
	},
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services with a rotation schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		services := s.ServicesBySchedule()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tSCHEDULE\tNEXT ROTATION")
		for _, svc := range services {
			sched, _ := s.GetSchedule(svc)
			next, err := s.NextScheduledRotation(svc)
			nextStr := next.Format("2006-01-02 15:04:05")
			if err != nil {
				nextStr = "n/a"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", svc, sched, nextStr)
		}
		return w.Flush()
	},
}

func init() {
	scheduleCmd.AddCommand(scheduleSetCmd)
	scheduleCmd.AddCommand(scheduleClearCmd)
	scheduleCmd.AddCommand(scheduleListCmd)
	AddCommand(scheduleCmd)
}
