package cmd

import (
	"fmt"
	"strconv"

	"github.com/nqzyx/fencepost/internal/config"
	"github.com/nqzyx/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var priorityCmd = &cobra.Command{
	Use:   "priority <service> <1-4>",
	Short: "Set priority for a service (1=low, 2=normal, 3=high, 4=critical)",
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
		n, err := strconv.Atoi(args[1])
		if err != nil || n < 1 || n > 4 {
			return fmt.Errorf("priority must be 1, 2, 3, or 4")
		}
		if err := s.SetPriority(args[0], keystore.Priority(n)); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "priority set for %q\n", args[0])
		return nil
	},
}

var priorityListCmd = &cobra.Command{
	Use:   "list <1-4>",
	Short: "List services at a given priority level",
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
		n, err := strconv.Atoi(args[0])
		if err != nil || n < 1 || n > 4 {
			return fmt.Errorf("priority must be 1, 2, 3, or 4")
		}
		results := s.ByPriority(keystore.Priority(n))
		for _, svc := range results {
			fmt.Fprintln(cmd.OutOrStdout(), svc)
		}
		return nil
	},
}

func init() {
	priorityCmd.AddCommand(priorityListCmd)
	AddCommand(priorityCmd)
}
