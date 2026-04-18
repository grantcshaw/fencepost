package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/janearc/fencepost/internal/config"
	"github.com/janearc/fencepost/internal/keystore"
)

var rotationCountCmd = &cobra.Command{
	Use:   "rotation-count",
	Short: "Manage rotation counts for services",
}

var rotationCountGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get rotation count for a service",
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
		count, err := s.GetRotationCount(args[0])
		if err != nil {
			return err
		}
		fmt.Println(count)
		return nil
	},
}

var rotationCountResetCmd = &cobra.Command{
	Use:   "reset <service>",
	Short: "Reset rotation count for a service to zero",
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
		return s.ResetRotationCount(args[0])
	},
}

var rotationCountListCmd = &cobra.Command{
	Use:   "list <min>",
	Short: "List services with at least <min> rotations",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		min, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid minimum count: %w", err)
		}
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		services := s.ByRotationCount(min)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tROTATIONS")
		for _, svc := range services {
			count, _ := s.GetRotationCount(svc)
			fmt.Fprintf(w, "%s\t%d\n", svc, count)
		}
		return w.Flush()
	},
}

func init() {
	rotationCountCmd.AddCommand(rotationCountGetCmd)
	rotationCountCmd.AddCommand(rotationCountResetCmd)
	rotationCountCmd.AddCommand(rotationCountListCmd)
	RootCmd.AddCommand(rotationCountCmd)
}
